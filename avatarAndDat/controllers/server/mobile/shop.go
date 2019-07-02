package mobile

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/util"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"math/big"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func (m *Manager) PurchaseConfirmHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Uuid string `json:"uuid"`
		NftTranData []string `json:"nftTranData"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err:=errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}


	// get buyer info
	var userMarketInfo models.UserMarketInfo
	o:=orm.NewOrm()
	err=o.QueryTable("user_market_info").Filter("uuid",req.Uuid).
		One(&userMarketInfo, "nickname")
	if err!=nil {
		if err == orm.ErrNoRows {
			err:=errors.New("user has not binded wallet")
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
	buyerUuid:= req.Uuid
	buyerWallet:= userMarketInfo.Wallet

	// get user wallet, user must bind wallet before purchase nft

	// currentBalance must be larger than total price of nft
	needToPay := 0
	nftRequestData := req.NftTranData
	o := orm.NewOrm()
	for _, nftLdefIndex := range nftRequestData {
		var nftMKInfo models.NftMarketInfo
		err := o.QueryTable("nft_market_info").
			Filter("nft_ldef_index", nftLdefIndex).
			One(&nftMKInfo, "price")
		if err != nil {
			if err == orm.ErrNoRows {
				err:=errors.New("nft "+nftLdefIndex+" not exist in marketplace")
				logs.Error(err.Error())
				m.errorHandler(c, action, err)
				return
			} else {
				logs.Error(err.Error())
				err:=errors.New("unexpected error when query db")
				logs.Error(err.Error())
				m.errorHandler(c, action, err)
				return
			}
		}
		needToPay += nftMKInfo.Price
	}
	logs.Debug("need to pay", needToPay)
	//session,_:=models.MongoClient.StartSession()
	//session.StartTransaction()
	ctx := context.Background()
	col := models.MongoDB.Collection("users")

	type fields struct {
		Coin string `bson:"coin"`
	}

	uuid:= req.Uuid
	filter := bson.M{
		"uuid": uuid,
	}

	var queryResult fields

	err = col.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{
		"coin": true,
	})).Decode(&queryResult)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, action, err)
		return
	}
	currentBalance, err := strconv.Atoi(queryResult.Coin)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, action, err)
		return
	}
	logs.Debug("uuid", uuid, "current balance:", currentBalance)

	finalBalance := currentBalance - needToPay
	if finalBalance < 0 {
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
		logs.Error(err.Error())
		m.errorHandler(c, action, err)
		return
	}

	logs.Warn("update balance of user", uuid, " to", finalBalance)


	o.Begin() // begin transaction
	for i, nftLdefIndex := range nftRequestData {
		purchaseId := util.RandomPurchaseId()
		var mkInfo models.NftMarketInfo
		err:=o.QueryTable("nft_market_place").
			RelatedSel("nft_market_info").
			Filter("nft_ldef_index",nftLdefIndex).
			One(&mkInfo)
		if err!=nil {
			if err == orm.ErrNoRows {
				err:=errors.New("nft "+nftLdefIndex+" not exist in marketplace")
				logs.Error(err.Error())
				m.errorHandler(c, action, err)
				return
			} else {
				logs.Error(err.Error())
				err:=errors.New("unexpected error when query db")
				logs.Error(err.Error())
				m.errorHandler(c, action, err)
				return
			}
		}

		sellerWallet:= mkInfo.SellerWallet
		sellerUuid:= mkInfo.SellerUuid

		h := md5.New()
		io.WriteString(h, purchaseId)
		purchaseId = new(big.Int).SetBytes(h.Sum(nil)[:8]).String()
		nftLdefIndex := itemDetail.NftLdefIndex
		nftType:= itemDetail.SupportedType
		if err:=validNftLdefIndex(nftLdefIndex);err!=nil {
			logs.Emergency(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		tokenId:=TokenIdFromNftLdefIndex(nftLdefIndex)
		toBeTransfer[i] = &transferPayLoad{
			TokenId:    tokenId,
			PurchaseId: purchaseId,
		}

		var nftMKInfo models.NftMarketTable
		err := o.QueryTable("nft_market_table").
			Filter("nft_ldef_index", nftLdefIndex).
			One(&nftMKInfo, "price")
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		totalPaid:= nftMKInfo.Price
		activeTicker := nftMKInfo.ActiveTicker
		if err != nil {
			o.Rollback()
			logs.Emergency(err.Error())
			m.errorHandler(c, bq, err)
			return
		}

		// query owner info
		var mkInfo models.NftMarketTable
		err = o.QueryTable("nft_market_table").
			Filter("nft_ldef_index", nftLdefIndex).
			One(&mkInfo, "seller_wallet_id", "seller_nickname")
		if err != nil {
			o.Rollback()
			logs.Emergency(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		sellerWalletId := mkInfo.SellerWalletId
		sellerNickname := mkInfo.SellerNickname
		logs.Debug("purchase seller address", sellerWalletId)
		nftOwners[i] = sellerWalletId

		status := PURCHASE_PENDING
		nftPurchaseResponseInfo := &NftPurchaseResponseInfo{
			NftLdefIndex: nftLdefIndex,
			Status:       status,
		}
		responseNftTranData[i] = nftPurchaseResponseInfo
		storeInfo := &models.StorePurchaseHistroy{
			PurchaseId:    purchaseId,
			BuyerNickname: nickname,
			BuyerWalletId: walletAddress,
			SellerNickname: sellerNickname,
			SellerWalletId:     sellerWalletId,
			TotalPaid:     totalPaid,
			NftLdefIndex:  nftLdefIndex,
			ActiveTicker:  activeTicker,
			Status:        status,
			NftType: nftType,
		}
		_, err = o.Insert(storeInfo)
		if err != nil {
			o.Rollback()
			logs.Emergency(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		toBeDelete := &models.NftMarketTable{
			NftLdefIndex: nftLdefIndex,
		}
		//delete from marketplace
		num, err := o.Delete(toBeDelete)
		if err != nil {
			o.Rollback()
			logs.Emergency("can not delete nft ldef:", nftLdefIndex)
			m.errorHandler(c, bq, err)
			return
		}
		logs.Warn("delete from marketplace table, nftldef:", nftLdefIndex, "num", num)

		// add count for buyer
		_,err = o.QueryTable("market_user_table").Filter("nickname",nickname).Update(orm.Params{
			"count": orm.ColValue(orm.ColAdd,1),
		})
		if err!=nil {
			o.Rollback()
			logs.Emergency("can not add count for nickname:", nickname)
			m.errorHandler(c, bq, err)
			return
		}
		logs.Warn("add count in market table for",nickname)
		// reduce count for sender
		_,err = o.QueryTable("market_user_table").Filter("nickname",sellerNickname).Update(orm.Params{
			"count": orm.ColValue(orm.ColMinus,1),
		})
		if err!=nil {
			o.Rollback()
			logs.Emergency("can not add count for nickname:", sellerNickname)
			m.errorHandler(c, bq, err)
			return
		}
		logs.Warn("reduce count in market table for",sellerNickname)
	}
	o.Commit()

	// TODO send transfer transaction
}

func (m *Manager) ShoppingCartListHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req ShoppingCartListRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	nickname := req.Nickname
	var shoppingCartHistory []models.NftShoppingCart
	o := orm.NewOrm()
	_, err = o.QueryTable("nft_shopping_cart").
		Filter("nickname", nickname).
		All(&shoppingCartHistory)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	shoppingCartRecordRes := make([]*ShoppingCartRecord, len(shoppingCartHistory))
	for i, _ := range shoppingCartHistory {
		nftLdefIndex := shoppingCartHistory[i].NftLdefIndex
		logs.Debug("shopping card ldef index", nftLdefIndex)
		r := o.Raw(`
		select ni.nft_type, ni.nft_name,
		mk.price,mk.active_ticker, mk.qty,
		ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,na.short_description, na.long_description,
		mp.file_name, mp.icon_file_name from
		nft_market_table as mk,
		nft_mapping_table as mp,
		nft_info_table as ni,
		nft_item_admin as na
		where mk.nft_ldef_index = mp.nft_ldef_index and mk.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id and  ni.nft_ldef_index = ? `, nftLdefIndex)
		var nftInfo NFTInfo
		err = r.QueryRow(&nftInfo)
		if err != nil {
			if err == orm.ErrNoRows {
				logs.Info("item not exist in market")
				continue
			} else {
				logs.Error(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
		}

		nftResInfo, err := nftResInfoFromNftInfo(&nftInfo)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		logs.Debug("origin time", shoppingCartHistory[i].Timestamp)
		shoppingCartRecordRes[i] = &ShoppingCartRecord{
			Timestamp:   chinaTimeFromTimeStamp(shoppingCartHistory[i].Timestamp),
			NftTranData: nftResInfo,
		}
	}

	res := &ShoppingCartListResponse{
		RQBaseInfo: *bq,
		NftList:    shoppingCartRecordRes,
	}
	m.wrapperAndSend(c, bq, res)
}


func (m *Manager) ShoppingCartChangeHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req ShoppingCartChangeRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	nickname := req.Nickname
	operation := req.Operation
	// check operation
	if operation != SHOPPING_CART_ADD && operation != SHOPPING_CART_DELETE {
		err := errors.New("unknown shopping cart operation")
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	nftList := req.NFTList
	o := orm.NewOrm()
	o.Begin()
	for _, nftLdefIndex := range nftList {
		shoppingCartRecord := models.NftShoppingCart{
			NftLdefIndex: nftLdefIndex,
			Nickname:     nickname,
		}
		if operation == SHOPPING_CART_ADD {
			_, err := o.Insert(&shoppingCartRecord)
			if err != nil {
				o.Rollback()
				logs.Error(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
			logs.Debug("insert ", nftLdefIndex, "success")
		} else if operation == SHOPPING_CART_DELETE {
			err:=o.Read(&shoppingCartRecord,"nft_ldef_index", "nickname")
			if err!=nil {
				o.Rollback()
				logs.Error(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
			_, err = o.Delete(&shoppingCartRecord, "nft_ldef_index", "nickname")
			if err != nil {
				o.Rollback()
				logs.Error(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
			logs.Debug("delete ", nftLdefIndex, "success")
		}
	}
	err=o.Commit()

	m.wrapperAndSend(c, bq, &ShoppingCartChangeResponse{
		RQBaseInfo: *bq,
	})
}

// token purchase
func (m *Manager) TokenBuyPaidHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Uuid string `json:"uuid"`
		AppTranId string `json:"appTranId"`
		TransactionId string `json:"transactionId"`
		AppId string `json:"appId"`
		Amount int `json:"amout"`
		ActionStatus int `json:"actionStatus"`
	}
	type response struct {
		Status int `json:"status"`
		Action string `json:"action"`
		Amount int `json:"amount"`
		ActionStatus int `json:"actionStatus"`
		TransactionId string `json:"transactionId"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err:= errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}

	actionStatus := req.ActionStatus
	if actionStatus == common.BERRY_PURCHASE_FINISH {
		purchaseInfo := models.BerryPurchaseInfo{
			TransactionId: req.TransactionId,
		}
		o:=orm.NewOrm()
		o.Begin()
		err = o.Read(&purchaseInfo)
		if err != nil {
			o.Rollback()
			if err == orm.ErrNoRows {
				err:=errors.New("can not find "+req.TransactionId+" maybe need to do pending first")
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
		_, err = o.Update(&purchaseInfo,"status")
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

		filter:= bson.M{
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
		amount := req.Amount
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
			Status: common.RESPONSE_STATUS_SUCCESS,
			Action: action,
			Amount: req.Amount,
			ActionStatus: common.BERRY_PURCHASE_FINISH,
			TransactionId: req.TransactionId,
		})
		return
	} else if actionStatus == common.BERRY_PURCHASE_PENDING {
		transactionId:= util.RandomPurchaseId()
		purchaseInfo := models.BerryPurchaseInfo{
			TransactionId: transactionId,
			NumPurchased:  req.Amount,
			AppId:         req.AppId,
			Status:        common.BERRY_PURCHASE_PENDING,
			Uuid: req.Uuid,
		}
		o:= orm.NewOrm()
		_, err = o.Insert(&purchaseInfo)
		if err != nil {
			logs.Emergency(err.Error())
			m.errorHandler(c, action, err)
			return
		}
		m.wrapperAndSend(c, action, &response{
			Status: common.RESPONSE_STATUS_SUCCESS,
			Action: action,
			Amount: req.Amount,
			ActionStatus: common.BERRY_PURCHASE_PENDING,
			TransactionId: req.TransactionId,
		})
		return
	} else {
		err := errors.New("unknow berry purchase action status")
		logs.Error(err.Error())
		m.errorHandler(c, action, err)
		return
	}
}

