package ws

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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
	logs.Info("nft type",nftType)
	// TODO can use prepare to optimize query
	r := models.O.Raw(`
		select ni.nft_type, ni.nft_name,
		mk.price,mk.active_ticker, ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,  mp.file_name, mk.qty 
		from nft_market_table as mk, nft_mapping_table as mp,
		nft_info_table as ni where mk.nft_ldef_index = mp.nft_ldef_index 
		and mk.nft_ldef_index = ni.nft_ldef_index 
		and ni.nft_type = ? `, nftType)

	var nftInfos []MpListNFTInfo
	_,err= r.QueryRows(&nftInfos)
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	length:= len(nftInfos)
	nis:=make([]*MpListNFTInfo,length)

	var thumbnail string
	if nftType == TYPE_NFT_MUSIC {   // music
		thumbnail = beego.AppConfig.String("prefix")+beego.AppConfig.String("hostaddr")+ ":"+
			beego.AppConfig.String("fileport") + "/resource/market/dat/"
	} else if nftType == TYPE_NFT_AVATAR {  //avatar
		thumbnail = beego.AppConfig.String("prefix")+beego.AppConfig.String("hostaddr")+ ":"+
			beego.AppConfig.String("fileport") + "/resource/market/avatar/"
	} else {
		err := errors.New("unknown supported type")
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	for i:=0;i<length;i++ {
		logs.Info("thumbnail:",nftInfos[i].Thumbnail)
		nftInfos[i].Thumbnail = thumbnail + nftInfos[i].Thumbnail  //appending file name
		nis[i] = &nftInfos[i]
	}

	m.wrapperAndSend(c,bq,&MpListResponse{
		RQBaseInfo: *bq,
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

	col:=models.MongoDB.Collection("users")
	type fields struct {
		Coin string `bson:"coin"`
	}

	idType:= req.AsUser.Type
	var filter bson.M
	if idType == WeChatId || idType == FBId {
		filter=bson.M {
			"uid": req.AsUser.AsId,
		}
	} else if idType == PhoneOrEmailId {
		filter=bson.M {
			"username": req.AsUser.AsId,
		}
	} else {
		err:= errors.New("wrong type")
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	var queryResult fields

	err=col.FindOne(context.Background(), filter, options.FindOne().SetProjection(bson.M{
		"coin":true,
	})).Decode(&queryResult)
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	currentBalance,err:= strconv.Atoi(queryResult.Coin)
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c,bq,err)
		return
	}
	logs.Debug("as id",req.AsUser.AsId,"current balance:",currentBalance)

	// currentBalance must be larger than total price of nft
	needToPay:=0
	nftRequestData:=req.NftTranData

	for _,itemDetail:= range nftRequestData {
		needToPay+=itemDetail.NftValue
	}
	if currentBalance < needToPay {
		err:=errors.New("insufficient balance!")
		m.errorHandler(c, bq, err)
		return
	}

	asId:= req.AsUser.AsId
	walletAddress:= req.AsUser.AsWallet
	logs.Debug("wallet address",walletAddress)

	responseNftTranData:= make([]*NftPurchaseResponseInfo,len(nftRequestData))
	res:=&PurchaseConfirmResponse {
		RQBaseInfo: *bq,
		NftTranData: responseNftTranData,
	}

	toBeInsert:=make([]*models.StorePurchaseHistroy,len(nftRequestData))
	//toBeDelete:=make([]*models.NftMarketTable,len(nftRequestData))
	// send transaction
	models.O.Begin()  // begin transaction
	var wg sync.WaitGroup
	for i,itemDetail := range nftRequestData {
		wg.Add(1)
		go func(i int, itemDetail *PurchaseNftInfo) {
			defer wg.Done()
			// generate purchase id
			purchaseId := strconv.FormatInt(time.Now().UnixNano()|rand2.Int63(), 10)
			h := md5.New()
			io.WriteString(h, purchaseId)
			purchaseId = new(big.Int).SetBytes(h.Sum(nil)[:8]).String()
			nftLdefIndex:= itemDetail.NftLdefIndex
			tokenId,_:= new(big.Int).SetString(nftLdefIndex[1:],10)
			totalPaid:= itemDetail.NftValue
			nftName:= itemDetail.NftName
			ownerAddress,err:= m.chainHandler.Contract.(*nft.NFT).OwnerOf(tokenId)
			// TODO in case blockchain down here
			if err!=nil {
				logs.Emergency(err.Error())
				m.errorHandler(c,bq,err)
				return
			}
			tx,txErr:=m.chainHandler.ManagerAccount.SendFunction2(m.chainHandler.Contract,
				nil,
				nft.FuncDelegateTransfer,
				common.HexToAddress(ownerAddress),
				common.HexToAddress(walletAddress),
				tokenId)
			err= <-txErr
			var status int
			if err!=nil {
				status = PURCHASE_PENDING
				logs.Debug("transfer token unsuccessful",tokenId,"to",walletAddress,"from",ownerAddress)
			} else {
				status = PURCHASE_CONFIRMED
				logs.Debug("transfer token",tokenId,"to",walletAddress,"from",ownerAddress)
			}
			nftPurchaseResponseInfo:=&NftPurchaseResponseInfo{
				NftLdefIndex: nftLdefIndex,
				Status:status,
			}
			responseNftTranData[i] = nftPurchaseResponseInfo
			storeInfo:=&models.StorePurchaseHistroy{
				PurchaseId:purchaseId,
				AsId:asId,
				TransactionAddress:tx.Hash().Hex(),
				NftName:nftName,
				TotalPaid:totalPaid,
				NftLdefIndex:nftLdefIndex,
				Status: status,
			}
			toBeInsert[i] = storeInfo
			toBeDelete := &models.NftMarketTable{
				NftLdefIndex:nftLdefIndex,
			}
			//delete from marketplace
			num,err:=models.O.Delete(toBeDelete)
			if err!= nil {
				models.O.Rollback()
				logs.Emergency("can not delete nft ldef:",nftLdefIndex)
			} else {
				logs.Warn("delete from marketplace table, nftldef:",nftLdefIndex,"num",num)
			}
		}(i,itemDetail)
	}
	wg.Wait()
	num,err:=models.O.InsertMulti(len(toBeInsert),toBeInsert)
	if err!=nil {
		models.O.Rollback()
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	models.O.Commit()
	// update balance of user   TODO  in case update fail
	finalBalance:= currentBalance - needToPay

	update:=bson.M {
		"$set":bson.M {"coin":strconv.Itoa(finalBalance)},

	}
	_,err =col.UpdateOne(context.Background(),filter,update)
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c,bq,err)
		return
	}
	logs.Warn("update balance of user",req.AsUser.AsId," to",finalBalance)


	logs.Info("insert num",num,"to purchase table")
	m.wrapperAndSend(c,bq,res)
}

func (m *Manager) ItemDetailsHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req ItemDetailsRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	nftTranData:=req.NftTranData
	nftResponseTranData:= make([]*ItemDetailsResponseNftInfo,len(nftTranData))
	itemDetailRes:=ItemDetailResponse{
		RQBaseInfo:*bq,
		NftTranData:nftResponseTranData,
	}

	for i,itemDetailsRequestNftInfo:= range nftTranData {
		nftLdefIndex:= itemDetailsRequestNftInfo.NftLdefIndex
		nftType:= itemDetailsRequestNftInfo.SupportedType
		r := models.O.Raw(`
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
		if err!=nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		var thumbnail string
		if nftType == TYPE_NFT_MUSIC {   // music
			thumbnail = beego.AppConfig.String("prefix")+beego.AppConfig.String("hostaddr")+ ":"+
				beego.AppConfig.String("fileport") + "/resource/market/dat/"
		} else if nftType == TYPE_NFT_AVATAR {  //avatar
			thumbnail = beego.AppConfig.String("prefix")+beego.AppConfig.String("hostaddr")+ ":"+
				beego.AppConfig.String("fileport") + "/resource/market/avatar/"
		} else {
			err := errors.New("unknown supported type")
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		nftResponseInfo.Thumbnail = thumbnail + nftResponseInfo.Thumbnail  // appending file name
		nftResponseTranData[i] = &nftResponseInfo
	}
	m.wrapperAndSend(c,bq,itemDetailRes)
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
	nftLdefIndex:= req.NftLdefIndex
	mp:=models.NftMappingTable{
		NftLdefIndex: nftLdefIndex,
	}
	err =models.O.Read(&mp)
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c,bq,err)
		return
	}
	fileName:=mp.FileName
	//TODO user symmetric key from client to decrypt file
	var encryptedFilePath string
	var decryptedFilePath string
	logs.Debug("nft type from request,",req.SupportedType)
	if req.SupportedType == TYPE_NFT_AVATAR {
		encryptedFilePath = path.Join(ENCRYPTION_FILE_PATH,NAME_NFT_AVATAR,fileName)
		decryptedFilePath = path.Join(DECRYPTION_FILE_PATH,NAME_NFT_AVATAR,fileName)
	} else if req.SupportedType == TYPE_NFT_MUSIC {
		encryptedFilePath = path.Join(ENCRYPTION_FILE_PATH,NAME_NFT_MUSIC,fileName)
		decryptedFilePath = path.Join(DECRYPTION_FILE_PATH,NAME_NFT_MUSIC,fileName)
	} else {
		err := errors.New("unknown supported type")
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	cipherText,err:=ioutil.ReadFile(encryptedFilePath)
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	nonce,ct:= cipherText[:aesgcm.NonceSize()],cipherText[aesgcm.NonceSize():]
	originalData ,err:= aesgcm.Open(nil,nonce,ct,nil)
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c,bq,err)
		return
	}

	logs.Debug("length of original data",len(originalData))
	if req.SupportedType == TYPE_NFT_AVATAR {
		out,err:= os.Create(decryptedFilePath)
		if err!=nil {
			logs.Error(err.Error())
			m.errorHandler(c,bq,err)
			return
		}
		defer out.Close()
		originalImage,_,err:=image.Decode(bytes.NewBuffer(originalData))
		if err!=nil {
			logs.Error(err.Error())
			m.errorHandler(c,bq,err)
			return
		}
		err=jpeg.Encode(out,originalImage,nil)
		if err!=nil {
			logs.Error(err.Error())
			m.errorHandler(c,bq,err)
			return
		}
	} else if req.SupportedType == TYPE_NFT_MUSIC {
		 err:=ioutil.WriteFile(decryptedFilePath,originalData,0777)
		 if err!=nil {
		 	logs.Error(err.Error())
		 	m.errorHandler(c,bq,err)
		 	return
		 }
	}

	m.wrapperAndSend(c,bq,NftShowResponse{
		RQBaseInfo:*bq,
		NftLdefIndex: nftLdefIndex,
		DecSource: beego.AppConfig.String("prefix")+beego.AppConfig.String("hostaddr")+ ":"+
			beego.AppConfig.String("fileport")+"/"+decryptedFilePath,
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

	actionStatus:= req.ActionStatus
	if actionStatus == ACTION_STATUS_FINISH {
		purchaseInfo:=models.BerryPurchaseTable {
			TransactionId:req.AppTranId,
		}

		// update coin records
		col:=models.MongoDB.Collection("users")
		// update coin records
		type fields struct {
			Coin string `bson:"coin"`
		}

		idType:= req.AsUser.Type
		var filter bson.M
		if idType == WeChatId || idType == FBId {
			filter=bson.M {
				"uid": req.AsUser.AsId,
			}
		} else if idType == PhoneOrEmailId {
			filter=bson.M {
				"username": req.AsUser.AsId,
			}
		} else {
			err:= errors.New("wrong type")
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		var queryResult fields

		err=col.FindOne(context.Background(), filter, options.FindOne().SetProjection(bson.M{
			"coin":true,
		})).Decode(&queryResult)
		if err!=nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		logs.Debug("as id",req.AsUser.AsId,"coin number:",queryResult.Coin)


		if err!=nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		currentBalance,err:= strconv.Atoi(queryResult.Coin)
		if err!=nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		amount:= req.Amount
		update:=bson.M {
			"$set":bson.M {"coin":strconv.Itoa(amount+currentBalance)},

		}
		_,err =col.UpdateOne(context.Background(),filter,update)
		if err!=nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		logs.Info("update success","after update, amount:",amount+currentBalance)

		models.O.Begin()
		err=models.O.Read(&purchaseInfo)
		if err!=nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		purchaseInfo.Status = ACTION_STATUS_FINISH
		_,err=models.O.Update(&purchaseInfo)
		if err!=nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		models.O.Commit()
		logs.Info("insert one record to purchase table")
		m.wrapperAndSend(c,bq,&TokenPurchaseResponse{
			RQBaseInfo: *bq,
			ActionStatus: ACTION_STATUS_FINISH,
		})
		return
	} else if actionStatus == ACTION_STATUS_PENDING {
		purchaseInfo:= models.BerryPurchaseTable{
			TransactionId: req.AppTranId,
			RefillAsId: req.AsUser.AsId,
			NumPurchased: req.Amount,
			AppTranId: req.AppTranId,
			AppId: req.AppId,
			Status: ACTION_STATUS_PENDING,
		}
		_,err=models.O.Insert(&purchaseInfo)
		if err!=nil {
			logs.Emergency(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		m.wrapperAndSend(c,bq,&TokenPurchaseResponse{
			RQBaseInfo: *bq,
			ActionStatus: ACTION_STATUS_PENDING,
		})
		return
	} else {
		err:=errors.New("unknow action status")
		logs.Error(err.Error())
		m.errorHandler(c,bq,err)
		return
	}
}