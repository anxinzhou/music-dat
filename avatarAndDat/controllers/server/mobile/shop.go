package mobile

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/transactionQueue"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/util"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
)

func (m *Manager) PurchaseConfirmHandler(c *client.Client, action string, data []byte) {
	// TODO this api need to design carefully
	type request struct {
		Uuid        string   `json:"uuid"`
		NftTranData []string `json:"nftTranData"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}

	// get buyer info
	var userMarketInfo models.UserMarketInfo
	o := orm.NewOrm()
	err = o.QueryTable("user_market_info").Filter("uuid", req.Uuid).
		One(&userMarketInfo)
	if err != nil {
		if err == orm.ErrNoRows {
			err := errors.New("user has not binded wallet")
			logs.Error(err.Error())
			m.errorHandler(c, action, err)
			return
		} else {
			logs.Error(err.Error())
			err := errors.New("unexpected error when query db")
			m.errorHandler(c, action, err)
			return
		}
	}

	// currentBalance must be larger than total price of nft
	needToPay := 0
	nftRequestData := req.NftTranData
	transactionList:= make([]*transactionQueue.NftPurchaseTransaction, len(nftRequestData))
	o.Begin() // begin transaction
	for i, nftLdefIndex := range nftRequestData {
		var nftMarketPlaceInfo models.NftMarketPlace
		err:=o.QueryTable("nft_market_place").
			RelatedSel("NftMarketInfo").Filter("nft_ldef_index",nftLdefIndex).
			One(&nftMarketPlaceInfo)
		if err != nil {
			o.Rollback()
			if err == orm.ErrNoRows {
				err := errors.New("nft " + nftLdefIndex + " not exist in marketplace")
				logs.Error(err.Error())
				m.errorHandler(c, action, err)
				return
			} else {
				logs.Error(err.Error())
				err := errors.New("unexpected error when query db")
				logs.Error(err.Error())
				m.errorHandler(c, action, err)
				return
			}
		}
		needToPay += nftMarketPlaceInfo.NftMarketInfo.Price

		// insert into purchase info
		purchaseId:= util.RandomPurchaseId()
		nftPuchaseInfo:= models.NftPurchaseInfo{
			PurchaseId: purchaseId,
			Uuid: req.Uuid,
			SellerUuid: nftMarketPlaceInfo.NftMarketInfo.SellerUuid,
			TransactionAddress: "", // determined after send transaction
			ActiveTicker: nftMarketPlaceInfo.ActiveTicker,
			NftLdefIndex: nftMarketPlaceInfo.NftLdefIndex,
			Status: common.PURCHASE_PENDING, // change to finish after send transaction
			UserInfo: &models.UserInfo{
				Uuid: req.Uuid,
			},
		}
		_,err=o.Insert(&nftPuchaseInfo)
		if err!=nil {
			o.Rollback()
			logs.Error(err.Error())
			err:= errors.New("unexpected error when query databas")
			m.errorHandler(c, action, err)
			return
		}
		// reduce count for seller
		_,err = o.QueryTable("user_market_info").Filter("uuid",nftMarketPlaceInfo.NftMarketInfo.SellerUuid).Update(orm.Params{
			"count": orm.ColValue(orm.ColMinus,1),
		})
		if err!=nil {
			o.Rollback()
			if err== orm.ErrNoRows {
				err:=errors.New("no such user "+req.Uuid)
				logs.Error(err.Error())
				m.errorHandler(c, action, err)
				return
			} else {
				logs.Error(err.Error())
				err:=errors.New("unexpected error when query db")
				m.errorHandler(c, action, err)
				return
			}
		}
		transactionList[i] = &transactionQueue.NftPurchaseTransaction{
			Uuid: req.Uuid,
			SellerUuid: nftMarketPlaceInfo.NftMarketInfo.SellerUuid,
			NftLdefIndex:nftLdefIndex,
			PurchaseId: purchaseId,
		}
	}
	count:= len(nftRequestData)
	// add count for buyer
	_,err = o.QueryTable("user_market_info").Filter("uuid",req.Uuid).Update(orm.Params{
		"count": orm.ColValue(orm.ColAdd,count),
	})
	if err!=nil {
		o.Rollback()
		if err== orm.ErrNoRows {
			err:=errors.New("no such user "+req.Uuid)
			logs.Error(err.Error())
			m.errorHandler(c, action, err)
			return
		} else {
			logs.Error(err.Error())
			err:=errors.New("unexpected error when query db")
			m.errorHandler(c, action, err)
			return
		}
	}
	logs.Debug("need to pay", needToPay)
	ctx := context.Background()
	col := models.MongoDB.Collection("users")

	type fields struct {
		Coin string `bson:"coin"`
	}

	uuid := req.Uuid
	filter := bson.M{
		"uuid": uuid,
	}

	var queryResult fields

	err = col.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{
		"coin": true,
	})).Decode(&queryResult)
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("no such user")
		m.errorHandler(c, action, err)
		return
	}
	currentBalance, err := strconv.Atoi(queryResult.Coin)
	if err != nil {
		panic("wrong coin type")
	}
	logs.Debug("uuid", uuid, "current balance:", currentBalance)

	finalBalance := currentBalance - needToPay
	if finalBalance < 0 {
		o.Rollback()
		err := errors.New("Insufficient balance")
		logs.Error(err.Error())
		m.errorHandler(c, action, err)
		return
	}
	update := bson.M{
		"$set": bson.M{"coin": strconv.Itoa(finalBalance)},
	}
	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("update balance of user fail")
		m.errorHandler(c, action, err)
		return
	}

	logs.Warn("update balance of user", uuid, " to", finalBalance)
	o.Commit()
	type response struct {
		Status int `json:"status"`
		Action string `json:"action"`
		NftTranData []string `json:"nftTranData"`
	}
	m.wrapperAndSend(c,action,&response{
		Status: common.RESPONSE_STATUS_SUCCESS,
		Action: action,
		NftTranData: req.NftTranData,
	})

	// TODO send transfer transaction
	for _,transaction:=range transactionList {
		m.TransactionQueue.Append(transaction)
	}
}

func (m *Manager) ShoppingCartListHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Uuid string `json:"uuid"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}
	type nftTranData struct {
		SupportedType string `json:"supportedType" orm:"column(nft_type)"`
		NftName       string `json:"nftName"`
		Price         int    `json:"price"`
		ActiveTicker  string `json:"activeTicker"`
		NftLdefIndex  string `json:"nftLdefIndex"`
		ShortDesc     string `json:"shortDesc"`
		LongDesc      string `json:"longDesc"`
		Thumbnail     string `json:"thumbnail" orm:"column(file_name)"`
		Timestamp     string `json:"timestamp"`
	}
	var shoppingInfo []nftTranData
	dbEngine := beego.AppConfig.String("dbEngine")
	qb, _ := orm.NewQueryBuilder(dbEngine)
	qb.Select(
		"nft_info.nft_ldef_index",
		"nft_info.nft_name",
		"nft_info.nft_type",
		"nft_info.short_desc",
		"nft_info.long_desc",
		"nft_market_info.price",
		"nft_market_info.qty",
		"nft_info.file_name",
	).
		From("nft_shopping_cart").
		InnerJoin("nft_market_info").
		On("nft_shopping_cart.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_shopping_cart.nft_ldef_index = nft_info.nft_ldef_index").
		Where("nft_shopping_cart.uuid = ?").OrderBy("nft_shopping_cart.timestamp").Desc()
	sql := qb.String()
	o := orm.NewOrm()
	num, err := o.Raw(sql, req.Uuid).QueryRows(&shoppingInfo)
	if err != nil && err != orm.ErrNoRows {
		logs.Error(err.Error())
		err := errors.New("unexpected error when query database")
		m.errorHandler(c, action, err)
		return
	}
	if num == 0 {
		shoppingInfo = make([]nftTranData, 0)
	}
	for i, v := range shoppingInfo {
		var nftType string
		if v.SupportedType == common.TYPE_NFT_AVATAR {
			nftType = common.TYPE_NFT_AVATAR
		} else if v.SupportedType == common.TYPE_NFT_MUSIC {
			nftType = common.TYPE_NFT_MUSIC
		} else if v.SupportedType == common.TYPE_NFT_OTHER {
			nftType = common.TYPE_NFT_OTHER
		} else {
			panic("unexpected type when query database")
		}
		shoppingInfo[i].Thumbnail = util.PathPrefixOfNFT(nftType, common.PATH_KIND_MARKET) + shoppingInfo[i].Thumbnail
	}
	type response struct {
		Status  int           `json:"status"`
		Action  string        `json:"action"`
		NftList []nftTranData `json:"nftList"`
	}
	m.wrapperAndSend(c, action, &response{
		Status:  common.RESPONSE_STATUS_SUCCESS,
		Action:  action,
		NftList: shoppingInfo,
	})
}

func (m *Manager) ShoppingCartChangeHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Operation int      `json:"operation"`
		Uuid      string   `json:"uuid"`
		NftList   []string `json:"nftList"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}
	if err := util.ValidShoppingCartOperation(req.Operation); err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, action, err)
		return
	}

	switch req.Operation {
	case common.SHOPPING_CART_ADD:
		o := orm.NewOrm()
		o.Begin()
		for _, nftLdefIndex := range req.NftList {
			shoppingCartInfo := models.NftShoppingCart{
				NftLdefIndex: nftLdefIndex,
				Uuid:         req.Uuid,
				NftMarketPlace: &models.NftMarketPlace{
					NftLdefIndex: nftLdefIndex,
				},
				UserInfo: &models.UserInfo{
					Uuid: req.Uuid,
				},
			}
			_, err := o.Insert(&shoppingCartInfo)
			if err != nil {
				if !strings.Contains(err.Error(), common.DUPLICATE_ENTRY) {
					logs.Error(err.Error())
					err := errors.New("unexpected error when query db")
					m.errorHandler(c, action, err)
					return
				}
			}
		}
		o.Commit()
	case common.SHOPPING_CART_DELETE:
		o := orm.NewOrm()
		o.Begin()
		for _, nftLdefIndex := range req.NftList {
			shoppingCartInfo := models.NftShoppingCart{
				NftLdefIndex: nftLdefIndex,
				Uuid:         req.Uuid,
				NftMarketPlace: &models.NftMarketPlace{
					NftLdefIndex: nftLdefIndex,
				},
				UserInfo: &models.UserInfo{
					Uuid: req.Uuid,
				},
			}
			_, err := o.Delete(&shoppingCartInfo,"nft_ldef_index","uuid")
			if err != nil {
				logs.Error(err.Error())
				err := errors.New("unexpected error when query db")
				m.errorHandler(c, action, err)
				return
			}
		}
		o.Commit()
	}
	type response struct {
		Status    int      `json:"status"`
		Action    string   `json:"action"`
		Operation int      `json:"operation"`
		NftList   []string `json:"nftList"`
	}

	m.wrapperAndSend(c, action, &response{
		Status:    common.RESPONSE_STATUS_SUCCESS,
		Action:    action,
		Operation: req.Operation,
		NftList:   req.NftList,
	})
}

// token purchase
func (m *Manager) TokenBuyPaidHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Uuid          string `json:"uuid"`
		AppTranId     string `json:"appTranId"`
		TransactionId string `json:"transactionId"`
		AppId         string `json:"appId"`
		Amount        int    `json:"amount"`
		ActionStatus  int    `json:"actionStatus"`
	}
	type response struct {
		Status        int    `json:"status"`
		Action        string `json:"action"`
		Amount        int    `json:"amount"`
		ActionStatus  int    `json:"actionStatus"`
		TransactionId string `json:"transactionId"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}

	actionStatus := req.ActionStatus
	if actionStatus == common.BERRY_PURCHASE_FINISH {
		purchaseInfo := models.BerryPurchaseInfo{
			TransactionId: req.TransactionId,
		}
		o := orm.NewOrm()
		o.Begin()
		err = o.Read(&purchaseInfo)
		if err != nil {
			o.Rollback()
			if err == orm.ErrNoRows {
				err := errors.New("can not find " + req.TransactionId + " maybe need to do pending first")
				logs.Error(err.Error())
				m.errorHandler(c, action, err)
				return
			} else {
				logs.Error(err.Error())
				err := errors.New("unexpected error when query db")
				m.errorHandler(c, action, err)
				return
			}
		}
		if purchaseInfo.Status != common.BERRY_PURCHASE_PENDING {
			o.Rollback()
			err := errors.New("berry purchase has already finished")
			logs.Error(err.Error())
			m.errorHandler(c, action, err)
			return
		}
		purchaseInfo.AppTranId = req.AppTranId
		purchaseInfo.Status = common.BERRY_PURCHASE_FINISH
		_, err = o.Update(&purchaseInfo, "status")
		if err != nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, action, err)
			return
		}

		// update coin records
		col := models.MongoDB.Collection("users")
		// update coin records
		type fields struct {
			Coin string `bson:"coin"`
		}

		filter := bson.M{
			"uuid": req.Uuid,
		}

		var queryResult fields

		err = col.FindOne(context.Background(), filter, options.FindOne().SetProjection(bson.M{
			"coin": true,
		})).Decode(&queryResult)
		if err != nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, action, err)
			return
		}
		logs.Debug("uuid", req.Uuid, "coin number:", queryResult.Coin)

		currentBalance, err := strconv.Atoi(queryResult.Coin)
		if err != nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, action, err)
			return
		}
		amount := purchaseInfo.NumPurchased
		update := bson.M{
			"$set": bson.M{"coin": strconv.Itoa(amount + currentBalance)},
		}
		_, err = col.UpdateOne(context.Background(), filter, update)
		if err != nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, action, err)
			return
		}
		logs.Info("update success", "after update, amount:", amount+currentBalance)

		o.Commit()
		m.wrapperAndSend(c, action, &response{
			Status:        common.RESPONSE_STATUS_SUCCESS,
			Action:        action,
			Amount:        req.Amount,
			ActionStatus:  common.BERRY_PURCHASE_FINISH,
			TransactionId: req.TransactionId,
		})
		return
	} else if actionStatus == common.BERRY_PURCHASE_PENDING {
		transactionId := util.RandomPurchaseId()
		purchaseInfo := models.BerryPurchaseInfo{
			TransactionId: transactionId,
			NumPurchased:  req.Amount,
			AppId:         req.AppId,
			Status:        common.BERRY_PURCHASE_PENDING,
			Uuid:          req.Uuid,
			UserInfo: &models.UserInfo{
				Uuid:req.Uuid,
			},
		}
		o := orm.NewOrm()
		_, err = o.Insert(&purchaseInfo)
		if err != nil {
			logs.Emergency(err.Error())
			m.errorHandler(c, action, err)
			return
		}
		m.wrapperAndSend(c, action, &response{
			Status:        common.RESPONSE_STATUS_SUCCESS,
			Action:        action,
			Amount:        req.Amount,
			ActionStatus:  common.BERRY_PURCHASE_PENDING,
			TransactionId: transactionId,
		})
		return
	} else {
		err := errors.New("unknow berry purchase action status")
		logs.Error(err.Error())
		m.errorHandler(c, action, err)
		return
	}
}
