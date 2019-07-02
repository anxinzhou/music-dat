package mobile
//
//import (
//	"context"
//	"crypto/md5"
//	"encoding/json"
//	"errors"
//	"github.com/astaxie/beego"
//	"github.com/astaxie/beego/logs"
//	"github.com/astaxie/beego/orm"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/xxRanger/blockchainUtil/contract/nft"
//	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
//	"github.com/xxRanger/music-dat/avatarAndDat/models"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"io"
//	"math/big"
//	rand2 "math/rand"
//	"strconv"
//	"sync"
//	"time"
//)
//
//func (m *Manager) NFTPurchaseHistoryHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
//	var req NFTPurchaseHistoryRequest
//	err := json.Unmarshal(data, &req)
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//	nftType:=req.SupportedType
//	if err:=validSupportedType(nftType);err!=nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//
//	o := orm.NewOrm()
//	nickname := req.Nickname
//	var purchaseHistory []models.StorePurchaseHistroy
//	logs.Debug("user", nickname, "query nft purchase history")
//	_, err = o.QueryTable("store_purchase_histroy").
//		Filter("buyer_nickname", nickname).Filter("nft_type",nftType).All(&purchaseHistory, "purchase_id",
//		"transaction_address",
//		"buyer_wallet_id",
//		"total_paid",
//		"active_ticker",
//		"nft_ldef_index",
//		"timestamp",
//		"status")
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//	purchaseRecordRes := make([]*NFTPurchaseRecord, len(purchaseHistory))
//	for i, _ := range purchaseHistory {
//		nftLdefIndex := purchaseHistory[i].NftLdefIndex
//		r := o.Raw(`
//		select ni.nft_type, ni.nft_name,
//		ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
//		ni.nft_charac_id,na.short_description, na.long_description,
//		mp.file_name, mp.icon_file_name from
//		nft_mapping_table as mp,
//		nft_info_table as ni,
//		nft_item_admin as na
//		where mp.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id and  ni.nft_ldef_index = ? `, nftLdefIndex)
//		var nftInfo NFTInfo
//		err = r.QueryRow(&nftInfo)
//		if err != nil {
//			logs.Error(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//
//
//		nftInfo.ActiveTicker = purchaseHistory[i].ActiveTicker
//		nftInfo.NftValue = purchaseHistory[i].TotalPaid
//		fileName := nftInfo.FileName
//		decryptedFilePath,err:=DecryptFile(fileName,req.SupportedType)
//		if err!=nil {
//			logs.Error(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//
//		nftResInfo, err := nftResInfoFromNftInfo(&nftInfo)
//		nftResInfo.DecSource = beego.AppConfig.String("prefix") + beego.AppConfig.String("hostaddr") + ":" +
//			beego.AppConfig.String("fileport") + "/" + decryptedFilePath
//		if err != nil {
//			logs.Error(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//
//
//		purchaseRecordRes[i] = &NFTPurchaseRecord{
//			PurchaseId:         purchaseHistory[i].PurchaseId,
//			TransactionAddress: purchaseHistory[i].TransactionAddress,
//			NftTranData:        nftResInfo,
//			WalletId:           purchaseHistory[i].BuyerWalletId,
//			Timestamp:          chinaTimeFromTimeStamp(purchaseHistory[i].Timestamp),
//			Status:             purchaseHistory[i].Status,
//		}
//	}
//
//	res := &NFTPurchaseHistoryResponse{
//		RQBaseInfo:   *bq,
//		PurchaseList: purchaseRecordRes,
//	}
//	m.wrapperAndSend(c, bq, res)
//}
//
//func (m *Manager) UserMarketInfoHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
//	var req UserMarketInfoRequest
//	err := json.Unmarshal(data, &req)
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//	nickname := req.Nickname
//	var mkInfos []models.NftMarketTable
//	o:=orm.NewOrm()
//	num,err:=o.QueryTable("nft_market_table").
//		Filter("seller_nickname",nickname).
//		All(&mkInfos,"nft_ldef_index")
//
//	if err!=nil {
//		if err == orm.ErrNoRows {
//			logs.Info("no row in marketplace now")
//			mkInfos = make([]models.NftMarketTable,0)
//		} else {
//			logs.Error(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//	}
//	nftTranResponseData := make([]*nftInfoListRes, 0, num)
//
//	// get user nftInfo
//	for _, mkInfo := range mkInfos {
//		nftLdefIndex:= mkInfo.NftLdefIndex
//		r := o.Raw(`
//		select ni.nft_type, ni.nft_name,
//		mk.price,mk.active_ticker, mk.qty,
//		ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
//		ni.nft_charac_id,na.short_description, na.long_description,
//		mp.file_name, mp.icon_file_name from
//		nft_market_table as mk,
//		nft_mapping_table as mp,
//		nft_info_table as ni,
//		nft_item_admin as na
//		where mk.nft_ldef_index = mp.nft_ldef_index and mk.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id and  ni.nft_ldef_index = ? `, nftLdefIndex)
//		var nftInfo NFTInfo
//		err = r.QueryRow(&nftInfo)
//		if err != nil {
//			if err == orm.ErrNoRows {
//				logs.Debug(err.Error())
//				continue
//			} else {
//				logs.Error(err.Error())
//				m.errorHandler(c, bq, err)
//				return
//			}
//		}
//		nftResInfo, err := nftResInfoFromNftInfo(&nftInfo)
//		if err != nil {
//			logs.Error(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//		nftTranResponseData = append(nftTranResponseData, nftResInfo)
//	}
//
//	// balance of user
//	balance := len(nftTranResponseData)
//
//	res := &UserMarketInfoResponse{
//		RQBaseInfo:  *bq,
//		TotalNFT:    balance,
//		NftTranData: nftTranResponseData,
//	}
//	m.wrapperAndSend(c, bq, res)
//}
//
//// NFT Market
//func (m *Manager) GetMPList(c *client.Client, bq *RQBaseInfo, data []byte) {
//	var req MpListRequest
//	err := json.Unmarshal(data, &req)
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//
//	//page := req.Page
//	//offset := req.Offset
//
//	nftType := req.SupportedType
//	logs.Debug("nft type", nftType)
//
//	// TODO can use prepare to optimize query
//	o := orm.NewOrm()
//	r := o.Raw(`
//		select ni.nft_type, ni.nft_name,
//		mk.price,mk.active_ticker, ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
//		ni.nft_charac_id,  mp.file_name, mp.icon_file_name, mk.qty
//		from nft_market_table as mk, nft_mapping_table as mp,
//		nft_info_table as ni where mk.nft_ldef_index = mp.nft_ldef_index
//		and mk.nft_ldef_index = ni.nft_ldef_index
//		and ni.nft_type = ? `, nftType)
//	var nftInfos []NFTInfo
//
//	_, err = r.QueryRows(&nftInfos)
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//	length := len(nftInfos)
//	nis := make([]*MpListNFTInfo, length)
//
//	prefix := PathPrefixOfNFT(nftType, PATH_KIND_MARKET)
//	for i := 0; i < length; i++ {
//		nftInfo := &nftInfos[i]
//		nftResInfo := &MpListNFTInfo{
//			SupportedType: nftInfo.SupportedType,
//			NftName:       nftInfo.NftName,
//			NftValue:      nftInfo.NftValue,
//			ActiveTicker:  nftInfo.ActiveTicker,
//			NftLifeIndex:  nftInfo.NftLifeIndex,
//			NftPowerIndex: nftInfo.NftPowerIndex,
//			NftLdefIndex:  nftInfo.NftLdefIndex,
//			Thumbnail:     nftInfo.FileName,
//			Qty:           nftInfo.Qty,
//		}
//		if nftType == TYPE_NFT_AVATAR || nftType == TYPE_NFT_OTHER {
//			nftResInfo.Thumbnail = prefix + nftInfo.FileName
//		} else if nftType == TYPE_NFT_MUSIC {
//			nftResInfo.Thumbnail = prefix + nftInfo.IconFileName
//		} else {
//			err := errors.New("unexpected type")
//			logs.Emergency(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//		nis[i] = nftResInfo
//	}
//
//	m.wrapperAndSend(c, bq, &MpListResponse{
//		RQBaseInfo:  *bq,
//		NftTranData: nis,
//	})
//}
//
//func (m *Manager) PurchaseConfirmHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
//	var req PurchaseConfirmRequest
//	err := json.Unmarshal(data, &req)
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//
//	// currentBalance must be larger than total price of nft
//	needToPay := 0
//	nftRequestData := req.NftTranData
//	o := orm.NewOrm()
//	for _, itemDetail := range nftRequestData {
//		nftLdefIndex := itemDetail.NftLdefIndex
//		var nftMKInfo models.NftMarketTable
//		err := o.QueryTable("nft_market_table").
//			Filter("nft_ldef_index", nftLdefIndex).
//			One(&nftMKInfo, "price")
//		if err != nil {
//			logs.Error(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//		price := nftMKInfo.Price
//		needToPay += price
//	}
//	logs.Debug("need to pay", needToPay)
//	//session,_:=models.MongoClient.StartSession()
//	//session.StartTransaction()
//	ctx := context.Background()
//	col := models.MongoDB.Collection("users")
//
//	type fields struct {
//		Coin string `bson:"coin"`
//	}
//
//	uuid:= req.Uuid
//	filter := bson.M{
//		"uuid": uuid,
//	}
//
//	var queryResult fields
//
//	err = col.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{
//		"coin": true,
//	})).Decode(&queryResult)
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//	currentBalance, err := strconv.Atoi(queryResult.Coin)
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//	logs.Debug("uuid", uuid, "current balance:", currentBalance)
//
//	finalBalance := currentBalance - needToPay
//	if finalBalance < 0 {
//		err := errors.New("Insufficient balance")
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//	update := bson.M{
//		"$set": bson.M{"coin": strconv.Itoa(finalBalance)},
//	}
//	_, err = col.UpdateOne(ctx, filter, update)
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//
//	logs.Warn("update balance of user", uuid, " to", finalBalance)
//
//
//	walletAddress,err := models.WalletIdOfNickname(nickname)
//	if err!=nil {
//		err:= errors.New("nickname "+nickname+" doest bind wallet")
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//
//	logs.Debug("wallet address", walletAddress)
//
//	responseNftTranData := make([]*NftPurchaseResponseInfo, len(nftRequestData))
//	res := &PurchaseConfirmResponse{
//		RQBaseInfo:  *bq,
//		NftTranData: responseNftTranData,
//	}
//
//	type transferPayLoad struct {
//		TokenId    *big.Int
//		PurchaseId string
//	}
//	toBeTransfer := make([]*transferPayLoad, len(nftRequestData))
//	o.Begin() // begin transaction
//
//	nftOwners := make([]string, len(nftRequestData))
//
//	for i, itemDetail := range nftRequestData {
//		purchaseId := strconv.FormatInt(time.Now().UnixNano()|rand2.Int63(), 10)
//		h := md5.New()
//		io.WriteString(h, purchaseId)
//		purchaseId = new(big.Int).SetBytes(h.Sum(nil)[:8]).String()
//		nftLdefIndex := itemDetail.NftLdefIndex
//		nftType:= itemDetail.SupportedType
//		if err:=validNftLdefIndex(nftLdefIndex);err!=nil {
//			logs.Emergency(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//		tokenId:=TokenIdFromNftLdefIndex(nftLdefIndex)
//		toBeTransfer[i] = &transferPayLoad{
//			TokenId:    tokenId,
//			PurchaseId: purchaseId,
//		}
//
//		var nftMKInfo models.NftMarketTable
//		err := o.QueryTable("nft_market_table").
//			Filter("nft_ldef_index", nftLdefIndex).
//			One(&nftMKInfo, "price")
//		if err != nil {
//			logs.Error(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//		totalPaid:= nftMKInfo.Price
//		activeTicker := nftMKInfo.ActiveTicker
//		if err != nil {
//			o.Rollback()
//			logs.Emergency(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//
//		// query owner info
//		var mkInfo models.NftMarketTable
//		err = o.QueryTable("nft_market_table").
//			Filter("nft_ldef_index", nftLdefIndex).
//			One(&mkInfo, "seller_wallet_id", "seller_nickname")
//		if err != nil {
//			o.Rollback()
//			logs.Emergency(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//		sellerWalletId := mkInfo.SellerWalletId
//		sellerNickname := mkInfo.SellerNickname
//		logs.Debug("purchase seller address", sellerWalletId)
//		nftOwners[i] = sellerWalletId
//
//		status := PURCHASE_PENDING
//		nftPurchaseResponseInfo := &NftPurchaseResponseInfo{
//			NftLdefIndex: nftLdefIndex,
//			Status:       status,
//		}
//		responseNftTranData[i] = nftPurchaseResponseInfo
//		storeInfo := &models.StorePurchaseHistroy{
//			PurchaseId:    purchaseId,
//			BuyerNickname: nickname,
//			BuyerWalletId: walletAddress,
//			SellerNickname: sellerNickname,
//			SellerWalletId:     sellerWalletId,
//			TotalPaid:     totalPaid,
//			NftLdefIndex:  nftLdefIndex,
//			ActiveTicker:  activeTicker,
//			Status:        status,
//			NftType: nftType,
//		}
//		_, err = o.Insert(storeInfo)
//		if err != nil {
//			o.Rollback()
//			logs.Emergency(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//		toBeDelete := &models.NftMarketTable{
//			NftLdefIndex: nftLdefIndex,
//		}
//		//delete from marketplace
//		num, err := o.Delete(toBeDelete)
//		if err != nil {
//			o.Rollback()
//			logs.Emergency("can not delete nft ldef:", nftLdefIndex)
//			m.errorHandler(c, bq, err)
//			return
//		}
//		logs.Warn("delete from marketplace table, nftldef:", nftLdefIndex, "num", num)
//
//		// add count for buyer
//		_,err = o.QueryTable("market_user_table").Filter("nickname",nickname).Update(orm.Params{
//			"count": orm.ColValue(orm.ColAdd,1),
//		})
//		if err!=nil {
//			o.Rollback()
//			logs.Emergency("can not add count for nickname:", nickname)
//			m.errorHandler(c, bq, err)
//			return
//		}
//		logs.Warn("add count in market table for",nickname)
//		// reduce count for sender
//		_,err = o.QueryTable("market_user_table").Filter("nickname",sellerNickname).Update(orm.Params{
//			"count": orm.ColValue(orm.ColMinus,1),
//		})
//		if err!=nil {
//			o.Rollback()
//			logs.Emergency("can not add count for nickname:", sellerNickname)
//			m.errorHandler(c, bq, err)
//			return
//		}
//		logs.Warn("reduce count in market table for",sellerNickname)
//	}
//	o.Commit()
//
//	//TODO using message queue to deal with pending transaction
//	wg := sync.WaitGroup{}
//	for i, itemDetail := range nftRequestData {
//		wg.Add(1)
//		go func(i int, itemDetail *NftBaseInfo) {
//			defer wg.Done()
//			o := orm.NewOrm()
//			o.Begin()
//			payload := toBeTransfer[i]
//			tokenId := payload.TokenId
//			purchaseId := payload.PurchaseId
//			ownerAddress := nftOwners[i]
//			tx, err := m.chainHandler.ManagerAccount.PackTransaction(
//				m.chainHandler.Contract,
//				nil,
//				nft.FuncDelegateTransfer,
//				common.HexToAddress(ownerAddress),
//				common.HexToAddress(walletAddress),
//				tokenId)
//			if err != nil {
//				logs.Emergency(err.Error())
//				return
//			}
//			tx, err = m.chainHandler.ManagerAccount.SignTransaction(tx)
//			status := PURCHASE_CONFIRMED
//			storeInfo := &models.StorePurchaseHistroy{
//				PurchaseId:         purchaseId,
//				TransactionAddress: tx.Hash().Hex(),
//				Status:             status,
//			}
//			_, err = o.Update(storeInfo, "transaction_address", "status")
//			if err != nil {
//				o.Rollback()
//				logs.Emergency(err.Error())
//				return
//			}
//			logs.Warn("transfer token", tokenId, "to", walletAddress, "from", ownerAddress)
//			txErr := m.chainHandler.ManagerAccount.SendTransaction(tx)
//			err = <-txErr
//			if err != nil {
//				o.Rollback()
//				logs.Emergency(err.Error())
//				return
//			}
//			o.Commit()
//			responseNftTranData[i].Status = status
//		}(i, itemDetail)
//	}
//	wg.Wait()
//	m.wrapperAndSend(c, bq, res)
//}
//
//func (m *Manager) MarketUserListHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
//	var req MarketUserListRequest
//	err := json.Unmarshal(data, &req)
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//	o := orm.NewOrm()
//	r := o.Raw(`
//		select mk.wallet_id,mk.nickname,count,user_icon_url,intro
//		 from market_user_table as mk, creator_info as ci
//		 where mk.nickname = ci.nickname
//			`)
//	var walletIdList []MarketUserWallet
//	_, err = r.QueryRows(&walletIdList)
//	if err != nil {
//		if err == orm.ErrNoRows {
//			logs.Debug(err.Error())
//		} else {
//			logs.Error(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//	}
//	nickname:= req.Nickname
//	wl := make([]*MarketUserWallet, len(walletIdList))
//	for i, _ := range wl {
//		if walletIdList[i].Thumbnail == "" {
//			walletIdList[i].Thumbnail = "default.jpg"
//		}
//		walletIdList[i].Thumbnail = PathPrefixOfNFT("", PATH_KIND_USER_ICON) + walletIdList[i].Thumbnail
//		wl[i] = &walletIdList[i]
//		followNickname:= wl[i].Nickname
//		queryInfo:= models.FollowTable{
//			FolloweeNickname:followNickname,
//			FollowerNickname: nickname,
//		}
//		err:=o.Read(&queryInfo,"followee_nickname","follower_nickname")
//		if err!=nil {
//			if err== orm.ErrNoRows {
//				wl[i].Followed = false
//			} else {
//				logs.Error(err.Error())
//				m.errorHandler(c, bq, err)
//				return
//			}
//		} else {
//			wl[i].Followed = true
//		}
//	}
//
//	m.wrapperAndSend(c, bq, &MarketUserListResponse{
//		RQBaseInfo:   *bq,
//		WalletIdList: wl,
//	})
//}
