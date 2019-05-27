package ws

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"math/big"
	rand2 "math/rand"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

func (m *Manager) errorHandler(c *client.Client, bq *RQBaseInfo, err error) {
	bq.Event = "failed"
	res := &ErrorResponse{
		RQBaseInfo: *bq,
		Reason:     err.Error(),
	}
	resWrapper, err := json.Marshal(res)
	if err != nil {
		panic(err)
		return
	}
	c.Send(resWrapper)
}

func (m *Manager) wrapperAndSend(c *client.Client, bq *RQBaseInfo, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	c.Send(data)
}

// NFT Market
func (m *Manager) GetMPList(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req MpListRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	//page := req.Page
	//offset := req.Offset

	nftType := req.SupportedType
	logs.Debug("nft type", nftType)

	// TODO can use prepare to optimize query
	o := orm.NewOrm()
	r := o.Raw(`
		select ni.nft_type, ni.nft_name,
		mk.price,mk.active_ticker, ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,  mp.file_name, mp.icon_file_name, mk.qty 
		from nft_market_table as mk, nft_mapping_table as mp,
		nft_info_table as ni where mk.nft_ldef_index = mp.nft_ldef_index 
		and mk.nft_ldef_index = ni.nft_ldef_index 
		and ni.nft_type = ? `, nftType)
	var nftInfos []NFTInfo

	_, err = r.QueryRows(&nftInfos)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	length := len(nftInfos)
	nis := make([]*MpListNFTInfo, length)

	prefix := PathPrefixOfNFT(nftType, PATH_KIND_MARKET)
	for i := 0; i < length; i++ {
		nftInfo := &nftInfos[i]
		nftResInfo := &MpListNFTInfo{
			SupportedType: nftInfo.SupportedType,
			NftName:       nftInfo.NftName,
			NftValue:      nftInfo.NftValue,
			ActiveTicker:  nftInfo.ActiveTicker,
			NftLifeIndex:  nftInfo.NftLifeIndex,
			NftPowerIndex: nftInfo.NftPowerIndex,
			NftLdefIndex:  nftInfo.NftLdefIndex,
			Thumbnail:     nftInfo.FileName,
			Qty:           nftInfo.Qty,
		}
		if nftType == TYPE_NFT_AVATAR || nftType == TYPE_NFT_OTHER {
			nftResInfo.Thumbnail = prefix + nftInfo.FileName
		} else if nftType == TYPE_NFT_MUSIC {
			nftResInfo.Thumbnail = prefix + nftInfo.IconFileName
		} else {
			err := errors.New("unexpected type")
			logs.Emergency(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		nis[i] = nftResInfo
	}

	m.wrapperAndSend(c, bq, &MpListResponse{
		RQBaseInfo:  *bq,
		NftTranData: nis,
	})
}

func (m *Manager) PurchaseConfirmHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req PurchaseConfirmRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	// currentBalance must be larger than total price of nft
	needToPay := 0
	nftRequestData := req.NftTranData

	for _, itemDetail := range nftRequestData {
		needToPay += itemDetail.NftValue
	}

	//session,_:=models.MongoClient.StartSession()
	//session.StartTransaction()
	ctx := context.Background()
	col := models.MongoDB.Collection("users")

	type fields struct {
		Coin string `bson:"coin"`
	}

	idType := req.AsUser.Type
	var filter bson.M
	if idType == WeChatId || idType == FBId {
		filter = bson.M{
			"uid": req.AsUser.AsId,
		}
	} else if idType == PhoneOrEmailId {
		filter = bson.M{
			"username": req.AsUser.AsId,
		}
	} else {
		err := errors.New("wrong type")
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	var queryResult fields

	err = col.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{
		"coin": true,
	})).Decode(&queryResult)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	currentBalance, err := strconv.Atoi(queryResult.Coin)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	logs.Debug("as id", req.AsUser.AsId, "current balance:", currentBalance)

	finalBalance := currentBalance - needToPay
	if finalBalance < 0 {
		err := errors.New("Insufficient balance")
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	update := bson.M{
		"$set": bson.M{"coin": strconv.Itoa(finalBalance)},
	}
	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	logs.Warn("update balance of user", req.AsUser.AsId, " to", finalBalance)

	asId := req.AsUser.AsId
	walletAddress := req.AsUser.AsWallet
	logs.Debug("wallet address", walletAddress)

	responseNftTranData := make([]*NftPurchaseResponseInfo, len(nftRequestData))
	res := &PurchaseConfirmResponse{
		RQBaseInfo:  *bq,
		NftTranData: responseNftTranData,
	}

	toBeInsert := make([]*models.StorePurchaseHistroy, len(nftRequestData))
	//toBeDelete:=make([]*models.NftMarketTable,len(nftRequestData))
	// send transaction
	o := orm.NewOrm()
	o.Begin() // begin transaction
	var wg sync.WaitGroup
	for i, itemDetail := range nftRequestData {
		wg.Add(1)
		go func(i int, itemDetail *PurchaseNftInfo) {
			defer wg.Done()
			// generate purchase id
			purchaseId := strconv.FormatInt(time.Now().UnixNano()|rand2.Int63(), 10)
			h := md5.New()
			io.WriteString(h, purchaseId)
			purchaseId = new(big.Int).SetBytes(h.Sum(nil)[:8]).String()
			nftLdefIndex := itemDetail.NftLdefIndex
			tokenId, _ := new(big.Int).SetString(nftLdefIndex[1:], 10)
			totalPaid := itemDetail.NftValue
			activeTicker := itemDetail.ActiveTicker
			nftName := itemDetail.NftName
			ownerAddress, err := m.chainHandler.Contract.(*nft.NFT).OwnerOf(tokenId)
			if err != nil {
				o.Rollback()
				logs.Emergency(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
			logs.Debug("purchase owner address", ownerAddress)
			// delete from market user table if balance is zero
			_, err = o.QueryTable("market_user_table").Filter("wallet_id", ownerAddress).Update(
				orm.Params{
					"count": orm.ColValue(orm.ColAdd, -1),
				},
			)
			if err != nil {
				o.Rollback()
				logs.Emergency(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
			tx, txErr := m.chainHandler.ManagerAccount.SendFunction2(m.chainHandler.Contract,
				nil,
				nft.FuncDelegateTransfer,
				common.HexToAddress(ownerAddress),
				common.HexToAddress(walletAddress),
				tokenId)
			err = <-txErr
			var status int
			if err != nil {
				status = PURCHASE_PENDING
				logs.Debug("transfer token unsuccessful", tokenId, "to", walletAddress, "from", ownerAddress)
			} else {
				status = PURCHASE_CONFIRMED
				logs.Debug("transfer token", tokenId, "to", walletAddress, "from", ownerAddress)
			}
			nftPurchaseResponseInfo := &NftPurchaseResponseInfo{
				NftLdefIndex: nftLdefIndex,
				Status:       status,
			}
			responseNftTranData[i] = nftPurchaseResponseInfo
			storeInfo := &models.StorePurchaseHistroy{
				PurchaseId:         purchaseId,
				AsId:               asId,
				WalletId:           walletAddress,
				TransactionAddress: tx.Hash().Hex(),
				NftName:            nftName,
				TotalPaid:          totalPaid,
				NftLdefIndex:       nftLdefIndex,
				ActiveTicker:       activeTicker,
				Status:             status,
			}
			toBeInsert[i] = storeInfo
			toBeDelete := &models.NftMarketTable{
				NftLdefIndex: nftLdefIndex,
			}
			//delete from marketplace
			num, err := o.Delete(toBeDelete)
			if err != nil {
				o.Rollback()
				logs.Emergency("can not delete nft ldef:", nftLdefIndex)
			} else {
				logs.Warn("delete from marketplace table, nftldef:", nftLdefIndex, "num", num)
			}
		}(i, itemDetail)
	}
	wg.Wait()
	logs.Debug("length to be insert", len(toBeInsert))
	num, err := o.InsertMulti(len(toBeInsert), toBeInsert)
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	logs.Info("insert num", num, "to purchase table")
	o.Commit()
	m.wrapperAndSend(c, bq, res)
}

func (m *Manager) ItemDetailsHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req ItemDetailsRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	nftTranData := req.NftTranData
	nftResponseTranData := make([]*ItemDetailsResponseNftInfo, len(nftTranData))
	itemDetailRes := ItemDetailResponse{
		RQBaseInfo:  *bq,
		NftTranData: nftResponseTranData,
	}
	o := orm.NewOrm()
	for i, itemDetailsRequestNftInfo := range nftTranData {
		nftLdefIndex := itemDetailsRequestNftInfo.NftLdefIndex
		nftType := itemDetailsRequestNftInfo.SupportedType
		r := o.Raw(`
		select ni.nft_type, ni.nft_name,
		mk.price,mk.active_ticker, ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,  na.short_description, na.long_description ,mp.file_name,mk.qty from  
		nft_market_table as mk, 
		nft_mapping_table as mp,
		nft_info_table as ni,
		nft_item_admin as na 
		where mk.nft_ldef_index = mp.nft_ldef_index and mk.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id and  ni.nft_ldef_index = ? `, nftLdefIndex)
		var nftResponseInfo ItemDetailsResponseNftInfo
		err = r.QueryRow(&nftResponseInfo)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		thumbnail := PathPrefixOfNFT(nftType, PATH_KIND_MARKET)
		nftResponseInfo.Thumbnail = thumbnail + nftResponseInfo.Thumbnail // appending file name
		nftResponseTranData[i] = &nftResponseInfo
	}
	m.wrapperAndSend(c, bq, itemDetailRes)
}

// nft show
func (m *Manager) NFTDisplayHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req NftShowRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	nftLdefIndex := req.NftLdefIndex
	mp := models.NftMappingTable{
		NftLdefIndex: nftLdefIndex,
	}
	o := orm.NewOrm()
	err = o.Read(&mp)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	fileName := mp.FileName
	//TODO user symmetric key from client to decrypt file
	var encryptedFilePath string
	var decryptedFilePath string
	logs.Debug("nft type from request,", req.SupportedType)
	if req.SupportedType == TYPE_NFT_AVATAR {
		encryptedFilePath = path.Join(ENCRYPTION_FILE_PATH, NAME_NFT_AVATAR, fileName)
		decryptedFilePath = path.Join(DECRYPTION_FILE_PATH, NAME_NFT_AVATAR, fileName)
	} else if req.SupportedType == TYPE_NFT_MUSIC {
		encryptedFilePath = path.Join(ENCRYPTION_FILE_PATH, NAME_NFT_MUSIC, fileName)
		decryptedFilePath = path.Join(DECRYPTION_FILE_PATH, NAME_NFT_MUSIC, fileName)
	} else if req.SupportedType == TYPE_NFT_OTHER {
		encryptedFilePath = path.Join(ENCRYPTION_FILE_PATH, NAME_NFT_OTHER, fileName)
		decryptedFilePath = path.Join(DECRYPTION_FILE_PATH, NAME_NFT_OTHER, fileName)
	} else {
		err := errors.New("unknown supported type")
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	cipherText, err := ioutil.ReadFile(encryptedFilePath)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	nonce, ct := cipherText[:aesgcm.NonceSize()], cipherText[aesgcm.NonceSize():]
	originalData, err := aesgcm.Open(nil, nonce, ct, nil)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	logs.Debug("length of original data", len(originalData))
	if req.SupportedType == TYPE_NFT_AVATAR || req.SupportedType == TYPE_NFT_OTHER {
		out, err := os.Create(decryptedFilePath)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		defer out.Close()
		originalImage, _, err := image.Decode(bytes.NewBuffer(originalData))
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		err = jpeg.Encode(out, originalImage, nil)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
	} else if req.SupportedType == TYPE_NFT_MUSIC {
		err := ioutil.WriteFile(decryptedFilePath, originalData, 0777)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
	}

	m.wrapperAndSend(c, bq, NftShowResponse{
		RQBaseInfo:   *bq,
		NftLdefIndex: nftLdefIndex,
		DecSource: beego.AppConfig.String("prefix") + beego.AppConfig.String("hostaddr") + ":" +
			beego.AppConfig.String("fileport") + "/" + decryptedFilePath,
	})
}

// token purchase
func (m *Manager) TokenBuyPaidHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req TokenPurchaseRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	actionStatus := req.ActionStatus
	o := orm.NewOrm()
	if actionStatus == ACTION_STATUS_FINISH {
		purchaseInfo := models.BerryPurchaseTable{
			TransactionId: req.TransactionId,
		}

		o.Begin()
		err = o.Read(&purchaseInfo)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		if purchaseInfo.Status != ACTION_STATUS_PENDING {
			err := errors.New("action in wrong status")
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		purchaseInfo.AppTranId = req.AppTranId
		purchaseInfo.Status = ACTION_STATUS_FINISH
		_, err = o.Update(&purchaseInfo)
		if err != nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}

		// update coin records
		col := models.MongoDB.Collection("users")
		// update coin records
		type fields struct {
			Coin string `bson:"coin"`
		}

		idType := req.AsUser.Type
		var filter bson.M
		if idType == WeChatId || idType == FBId {
			filter = bson.M{
				"uid": req.AsUser.AsId,
			}
		} else if idType == PhoneOrEmailId {
			filter = bson.M{
				"username": req.AsUser.AsId,
			}
		} else {
			o.Rollback()
			err := errors.New("wrong type")
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		var queryResult fields

		err = col.FindOne(context.Background(), filter, options.FindOne().SetProjection(bson.M{
			"coin": true,
		})).Decode(&queryResult)
		if err != nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		logs.Debug("as id", req.AsUser.AsId, "coin number:", queryResult.Coin)

		currentBalance, err := strconv.Atoi(queryResult.Coin)
		if err != nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
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
			m.errorHandler(c, bq, err)
			return
		}
		logs.Info("update success", "after update, amount:", amount+currentBalance)

		o.Commit()
		logs.Info("insert one record to purchase table")
		m.wrapperAndSend(c, bq, &TokenPurchaseResponse{
			RQBaseInfo:   *bq,
			ActionStatus: ACTION_STATUS_FINISH,
		})
		return
	} else if actionStatus == ACTION_STATUS_PENDING {
		appTranIdBytes := make([]byte, 32)
		_, err := rand.Read(appTranIdBytes)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		appTranId := hex.EncodeToString(appTranIdBytes)
		purchaseInfo := models.BerryPurchaseTable{
			TransactionId: appTranId,
			RefillAsId:    req.AsUser.AsId,
			NumPurchased:  req.Amount,
			AppId:         req.AppId,
			Status:        ACTION_STATUS_PENDING,
		}
		_, err = o.Insert(&purchaseInfo)
		if err != nil {
			logs.Emergency(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		m.wrapperAndSend(c, bq, &TokenPurchaseResponse{
			RQBaseInfo:    *bq,
			ActionStatus:  ACTION_STATUS_PENDING,
			TransactionId: appTranId,
		})
		return
	} else {
		err := errors.New("unknow action status")
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
}

func (m *Manager) MarketUserListHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req MarketUserListRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	o := orm.NewOrm()
	r := o.Raw(`
		select wallet_id,username,count,user_icon_url from market_user_table where count>0`)
	var walletIdList []MarketUserWallet
	_, err = r.QueryRows(&walletIdList)
	if err != nil {
		if err == orm.ErrNoRows {
			logs.Debug(err.Error())
		} else {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
	}

	wl := make([]*MarketUserWallet, len(walletIdList))
	for i, _ := range wl {
		walletIdList[i].Thumbnail = PathPrefixOfNFT("", PATH_KIND_USER_ICON) + walletIdList[i].Thumbnail
		wl[i] = &walletIdList[i]
	}

	m.wrapperAndSend(c, bq, &MarketUserListResponse{
		RQBaseInfo:   *bq,
		WalletIdList: wl,
	})
}

func (m *Manager) UserMarketInfoHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req UserMarketInfoRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	user := req.WalletId

	logs.Debug("user", user, "query user market info")

	nftContract := m.chainHandler.Contract.(*nft.NFT)
	nftList, err := nftContract.TokensOfUser(common.HexToAddress(user))
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	nftLdefIndexs := make([]string, len(nftList))
	for i, tokenId := range nftList {
		ldef, err := nftContract.LdefIndexOfToken(tokenId)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		//logs.Info("ldefIndex",ldef)
		nftLdefIndexs[i] = ldef
	}
	nftTranResponseData := make([]*nftInfoListRes, 0, len(nftList))
	o := orm.NewOrm()
	// get user nftInfo
	for _, nftLdefIndex := range nftLdefIndexs {

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
				logs.Debug(err.Error())
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
		nftTranResponseData = append(nftTranResponseData, nftResInfo)
	}

	// balance of user
	balance := len(nftTranResponseData)

	res := &UserMarketInfoResponse{
		RQBaseInfo:  *bq,
		TotalNFT:    balance,
		NftTranData: nftTranResponseData,
	}
	m.wrapperAndSend(c, bq, res)
}

func (m *Manager) NFTPurchaseHistoryHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req NFTPurchaseHistoryRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	o := orm.NewOrm()
	userName := req.UserName
	var purchaseHistory []models.StorePurchaseHistroy
	logs.Debug("user", userName, "query nft purchase history")
	_, err = o.QueryTable("store_purchase_histroy").
		Filter("as_id", userName).All(&purchaseHistory, "purchase_id",
		"transaction_address",
		"wallet_id",
		"total_paid",
		"active_ticker",
		"nft_ldef_index",
		"timestamp",
		"status")
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	purchaseRecordRes := make([]*NFTPurchaseRecord, len(purchaseHistory))
	for i, _ := range purchaseHistory {
		nftLdefIndex := purchaseHistory[i].NftLdefIndex
		r := o.Raw(`
		select ni.nft_type, ni.nft_name,
		ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,na.short_description, na.long_description,
		mp.file_name, mp.icon_file_name from
		nft_mapping_table as mp,
		nft_info_table as ni,
		nft_item_admin as na
		where mp.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id and  ni.nft_ldef_index = ? `, nftLdefIndex)
		var nftInfo NFTInfo
		err = r.QueryRow(&nftInfo)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}

		nftInfo.ActiveTicker = purchaseHistory[i].ActiveTicker
		nftInfo.NftValue = purchaseHistory[i].TotalPaid
		nftResInfo, err := nftResInfoFromNftInfo(&nftInfo)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}

		purchaseRecordRes[i] = &NFTPurchaseRecord{
			PurchaseId:         purchaseHistory[i].PurchaseId,
			TransactionAddress: purchaseHistory[i].TransactionAddress,
			NftTranData:        nftResInfo,
			WalletId:           purchaseHistory[i].WalletId,
			Timestamp:          chinaTimeFromTimeStamp(purchaseHistory[i].Timestamp),
			Status:             purchaseHistory[i].Status,
		}
	}

	res := &NFTPurchaseHistoryResponse{
		RQBaseInfo:   *bq,
		PurchaseList: purchaseRecordRes,
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

	username := req.Username
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
			UserName:     username,
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
			_, err := o.Delete(&shoppingCartRecord, "nft_ldef_index", "user_name")
			if err != nil {
				o.Rollback()
				logs.Error(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
			logs.Debug("delete ", nftLdefIndex, "success")
		}
	}
	o.Commit()

	m.wrapperAndSend(c, bq, &ShoppingCartChangeResponse{
		RQBaseInfo: *bq,
	})
}

func (m *Manager) ShoppingCartListHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req ShoppingCartListRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	username := req.Username
	var shoppingCartHistory []models.NftShoppingCart
	o := orm.NewOrm()
	_, err = o.QueryTable("nft_shopping_cart").
		Filter("user_name", username).
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
