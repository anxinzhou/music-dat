package web

import (
	"bytes"
	"crypto/md5"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/nfnt/resize"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/util"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"math/big"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"
	"crypto/rand"
)

type RewardController struct {
	ContractController
}

type nftInfoQuery struct {
	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
	NftName       string `json:"nftName"`
	NftValue      int    `json:"nftValue" orm:"column(price)"`
	ActiveTicker  string `json:"activeTicker"`
	NftLifeIndex  int64  `json:"nftLifeIndex"`
	NftPowerIndex int64  `json:"nftPowerIndex"`
	NftLdefIndex  string `json:"nftLdefIndex"`
	NftCharacId   string `json:"nftCharacId"`
	ShortDesc     string `json:"shortDesc" orm:"column(short_description)"`
	LongDesc      string `json:"longDesc" orm:"column(long_description)"`
	Thumbnail     string `json:"thumbnail" orm:"column(icon_file_name)"`
	Qty           int    `json:"qty"`
	SellerWalletId string
	SellerNickname string
}

type rewardNFTInfo struct {
	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
	NftName       string `json:"nftName"`
	NftValue      int    `json:"nftValue" orm:"column(price)"`
	ActiveTicker  string `json:"activeTicker"`
	NftLifeIndex  int64  `json:"nftLifeIndex"`
	NftPowerIndex int64  `json:"nftPowerIndex"`
	NftLdefIndex  string `json:"nftLdefIndex"`
	NftCharacId   string `json:"nftCharacId"`
	ShortDesc     string `json:"shortDesc" orm:"column(short_description)"`
	LongDesc      string `json:"longDesc" orm:"column(long_description)"`
	Thumbnail     string `json:"thumbnail" orm:"column(icon_file_name)"`
	Qty           int    `json:"qty"`
}

type RewardResponse struct {
	NftTranData []*rewardNFTInfo `json:"nftTranData"`
}

func (this *RewardController) RewardDat() {
	// only reward one dat now
	nickname := this.Ctx.Input.Param(":nickname")
	walletAddress,err:= models.WalletIdOfNickname(nickname)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}
	o := orm.NewOrm()
	o.Begin()
	qs := o.QueryTable("nft_market_table").Filter("nft_ldef_index__contains", "M").Filter("allow_airdrop",true).Limit(1)
	var mk models.NftMarketTable
	rewardAccount := 1
	err = qs.Limit(rewardAccount).One(&mk)
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	var res RewardResponse
	nftInfoList := make([]*rewardNFTInfo, rewardAccount)
	res.NftTranData = nftInfoList
	nftLdefIndex := mk.NftLdefIndex
	r := o.Raw(`
		select ni.nft_type, ni.nft_name,
		mk.price,mk.active_ticker, mk.seller_wallet_id,mk.seller_nickname,
		ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,na.short_description, na.long_description,mp.icon_file_name,mk.qty from
		nft_market_table as mk,
		nft_mapping_table as mp,
		nft_info_table as ni,
		nft_item_admin as na
		where mk.nft_ldef_index = mp.nft_ldef_index and mk.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id and  ni.nft_ldef_index = ? `, nftLdefIndex)
	var nftResponseInfo nftInfoQuery
	err = r.QueryRow(&nftResponseInfo)
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}
	nftType := nftResponseInfo.SupportedType
	thumbnail := PathPrefixOfNFT(nftType, PATH_KIND_MARKET)
	nftResponseInfo.Thumbnail = thumbnail + nftResponseInfo.Thumbnail
	//nftResponseInfo.Thumbnail = thumbnail + "music.png"

	nftInfoList[0] = &rewardNFTInfo{
		SupportedType: nftResponseInfo.SupportedType,
		NftName: nftResponseInfo.NftName,
		NftValue: nftResponseInfo.NftValue,
		ActiveTicker: nftResponseInfo.ActiveTicker,
		NftLifeIndex: nftResponseInfo.NftLifeIndex,
		NftPowerIndex: nftResponseInfo.NftPowerIndex,
		NftLdefIndex: nftResponseInfo.NftLdefIndex,
		NftCharacId: nftResponseInfo.NftCharacId,
		ShortDesc: nftResponseInfo.ShortDesc,
		LongDesc: nftResponseInfo.LongDesc,
		Thumbnail: nftResponseInfo.Thumbnail,
		Qty:          nftResponseInfo.Qty,
	}
	_, err = o.Delete(&mk)  //TODO comment for testing
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	sellerNickname:= nftResponseInfo.SellerNickname
	sellerWalletAddress:= nftResponseInfo.SellerWalletId
	// add count for buyer
	_,err = o.QueryTable("market_user_table").Filter("nickname",nickname).Update(orm.Params{
		"count": orm.ColValue(orm.ColAdd,1),
	})
	if err!=nil {
		o.Rollback()
		logs.Emergency("can not add count for nickname:", nickname)
		sendError(&this.Controller, err, 500)
		return
	}
	logs.Warn("add count in market table for",nickname)
	// reduce count for sender
	_,err = o.QueryTable("market_user_table").Filter("nickname",sellerNickname).Update(orm.Params{
		"count": orm.ColValue(orm.ColMinus,1),
	})
	if err!=nil {
		o.Rollback()
		logs.Emergency("can not reduce count for nickname:", sellerNickname)
		sendError(&this.Controller, err, 500)
		return
	}
	logs.Warn("reduce count in market table for",sellerNickname)

	if len(nftLdefIndex)<=1 {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	// insert to purchase history
	purchaseId := strconv.FormatInt(time.Now().UnixNano()|rand.Int63(), 10)
	storeInfo := &models.StorePurchaseHistroy{
		PurchaseId:    purchaseId,
		BuyerNickname: nickname,
		BuyerWalletId: walletAddress,
		SellerNickname: sellerNickname,
		SellerWalletId:     sellerWalletAddress,
		TotalPaid:     nftInfoList[0].NftValue,
		NftLdefIndex:  nftLdefIndex,
		ActiveTicker:  nftInfoList[0].ActiveTicker,
		Status:       PURCHASE_CONFIRMED ,
	}
	_, err = o.Insert(storeInfo)
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	tokenId, _ := new(big.Int).SetString(nftLdefIndex[1:], 10)

	nftContract := this.C.smartContract.(*nft.NFT)
	_, txErr := this.C.account.SendFunction2(nftContract,
		nil,
		nft.FuncDelegateTransfer,
		common.HexToAddress(sellerWalletAddress),
		common.HexToAddress(walletAddress),
		tokenId) // TODO redis to cache unsuccessful transaction
	err = <-txErr
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}
	o.Commit()
	this.Data["json"] = res
	this.ServeJSON()
}


type NumOfChildrenController struct {
	ContractController
}

type NumOfChildrenRes struct {
	Count int `json:"count"`
}

func (this *NumOfChildrenController) Get() {
	parentIndex := this.Ctx.Input.Param(":parentIndex")
	o := orm.NewOrm()
	r := o.Raw(`
		select count(a.nft_ldef_index) as num 
		from nft_mapping_table as a,
		nft_market_table as b 
		where a.nft_parent_ldef = ? and a.nft_ldef_index = b.nft_ldef_index `, parentIndex)
	type CountQuery struct {
		Num int
	}
	var queryResult CountQuery
	err := r.QueryRow(&queryResult)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	res := &NumOfChildrenRes{
		Count: queryResult.Num,
	}

	this.Data["json"] = res
	this.ServeJSON()
}

type ChildrenOfNFTController struct {
	ContractController
}

func (this *ChildrenOfNFTController) Get() {
	parentIndex := this.Ctx.Input.Param(":parentIndex")
	o := orm.NewOrm()
	r := o.Raw(`
		select ni.nft_type, ni.nft_name,
		mk.price,mk.active_ticker, ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,na.short_description, na.long_description,mp.file_name,mk.qty from
		nft_market_table as mk,
		nft_mapping_table as mp,
		nft_info_table as ni,
		nft_item_admin as na
		where mp.nft_parent_ldef= ? and mk.nft_ldef_index = mp.nft_ldef_index and mp.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id`, parentIndex)
	nftResponseInfo := []*nftInfoListRes{}
	_, err := r.QueryRows(&nftResponseInfo)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	thumbnail := PathPrefixOfNFT(TYPE_NFT_OTHER, PATH_KIND_MARKET)
	for _, nftInfo := range nftResponseInfo {
		nftInfo.Thumbnail = thumbnail + nftInfo.Thumbnail
	}
	var res nftListResponse
	res.NftTranData = nftResponseInfo
	this.Data["json"] = &res
	this.ServeJSON()
}


type UploadController struct {
	ContractController
}

type uplodBaseInfo struct {
	File multipart.File
	Uuid string
	NftName string
	ShortDesc string
	LongDesc string
}

func (this *UploadController) uploadAvatar(reqBaseInfo *uplodBaseInfo) {
	fileNamePrefix := util.RandomPathFromFileName("file")
	fileName:= fileNamePrefix+".jpg"

	typeOfNft:= common.TYPE_NFT_AVATAR
	nameOfNftType:= common.NAME_NFT_AVATAR
	// ---------------------------------------
	// calculate ciphertext and save
	// ---------------------------------------
	data,err:= util.ReadFile(reqBaseInfo.File)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}

	logs.Debug("len of data", len(data))
	nonce := make([]byte, util.Aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	cipherText := util.Aesgcm.Seal(nonce, nonce, data, nil)

	// saving ciphertext
	logs.Debug("saving ciphertext")
	cipherSavingPath:= path.Join(common.ENCRYPTION_FILE_PATH, typeOfNft, fileName)
	err = ioutil.WriteFile(cipherSavingPath, cipherText, 0777)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("saving file", fileName, "to", cipherSavingPath)

	// ---------------------------------------
	// resize image and save to folder market
	// ---------------------------------------
	originImage, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	newImage := resize.Resize(200, 200, originImage, resize.Lanczos3)
	filePath:= path.Join(common.MARKET_PATH,nameOfNftType,fileName)
	err=util.SaveImage(newImage,filePath)
	if err!=nil {
		logs.Error(err.Error())
		err:= errors.New("can not save resized image")
		sendError(&this.Controller,err, 500)
		return
	}

	// ---------------------------------------
	// initialize nft info
	// ---------------------------------------
	//nft info
	nftLdefIndex:= util.RandomNftLdefIndex(typeOfNft)
	nftLifeIndex:= util.SmallRandInt()
	nftPowerIndex:= util.SmallRandInt()
	// nft market info
	price:= util.SmallRandInt()
	qty:= util.SmallRandInt()
	// ---------------------------------------
	// save nft info to database
	// ---------------------------------------
	o:= orm.NewOrm()
	o.Begin()
	// set nftInfo
	nftInfo:= models.NftInfo{
		NftLdefIndex: nftLdefIndex,
		NftType: typeOfNft,
		NftName: reqBaseInfo.NftName,
		ShortDescription: reqBaseInfo.ShortDesc,
		LongDescription: reqBaseInfo.LongDesc,
		FileName: fileName,
		NftParentLdef: "",
	}
	_,err=o.Insert(&nftInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert nft into database")
		sendError(&this.Controller,err, 500)
		return
	}
	// set avatar info
	avatarInfo:= models.AvatarNftInfo{
		NftLdefIndex: nftLdefIndex,
		NftLifeIndex: nftLifeIndex,
		NftPowerIndex: nftPowerIndex,
	}
	_,err= o.Insert(&avatarInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert avatar into database")
		sendError(&this.Controller,err, 500)
		return
	}
	// ---------------------------------------
	// save nft market info to database
	// ---------------------------------------
	userMarketInfo:= models.UserMarketInfo{
		Uuid: reqBaseInfo.Uuid,
	}
	err=o.Read(&userMarketInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("unknown error when query db")
		sendError(&this.Controller,err, 500)
		return
	}
	nftMarketInfo:= models.NftMarketInfo{
		NftLdefIndex: nftLdefIndex,
		SellerWallet: userMarketInfo.Wallet,
		SellerUuid: reqBaseInfo.Uuid,
		Price: price,
		Qty: qty,
		NumSold: 0,
	}
	_,err=o.Insert(&nftMarketInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert nft market into database")
		sendError(&this.Controller,err, 500)
		return
	}

	avatarMarketInfo:= models.AvatarNftMarketInfo{
		NftLdefIndex: nftLdefIndex,
	}
	_,err=o.Insert(&avatarMarketInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert avatar market into database")
		sendError(&this.Controller,err, 500)
		return
	}
	o.Commit()
	// ---------------------------------------
	// record nft in market Place
	// ---------------------------------------
	mkPlace:= models.NftMarketPlace{
		NftLdefIndex: nftLdefIndex,
		NftType: typeOfNft,
		MpId: common.MARKETPLACE_ID,
		Active: true,
		ActiveTicker: common.ACTIVE_TICKER,
	}
	_,err=o.Insert(&mkPlace)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert avatar into market place")
		sendError(&this.Controller,err, 500)
		return
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	res:= &common.AvatarNftMarketInfo{
		AvatarNftInfo: common.AvatarNftInfo{
			NftInfo: common.NftInfo{
				NftLdefIndex: nftLdefIndex,
				NftType: typeOfNft,
				NftName: reqBaseInfo.NftName,
				ShortDescription: reqBaseInfo.ShortDesc,
				LongDescription: reqBaseInfo.LongDesc,
				FileName: fileName,
				NftParentLdef: "",
			},
		},
		NftMarketInfo:common.NftMarketInfo{
			SellerWallet: userMarketInfo.Wallet,
			SellerUuid: userMarketInfo.Uuid,
			Price: price,
			Qty: qty,
			NumSold: 0,
		} ,
	}
	this.Data["json"] = res
	this.ServeJSON()
	// ---------------------------------------
	// call smart contract to create nft
	// ---------------------------------------
	// TODO use message queue to deal with transaction.
}

func (this *UploadController) uploadDat(reqBaseInfo *uplodBaseInfo) {
	iconFile,iconFileHeader,err:= this.GetFile("icon")
	if err!=nil {
		logs.Error(err.Error())
		err:= errors.New("no icon file name specified")
		sendError(&this.Controller,err, 500)
		return
	}
	allowAirdrop,err:= this.GetBool("allowAirdrop")
	if err!=nil {
		err:=errors.New("value of allowAirDrop should be bool")
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	number,err:= this.GetInt("number")
	if err!=nil {
		err:=errors.New("value of number should be int")
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	price,err:= this.GetInt("price")
	if err!=nil {
		err:=errors.New("value of price should be int")
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	creatorPercent,err:= this.GetFloat("creatorPercent")
	if err!=nil {
		err:=errors.New("value of creator percent should be float")
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	lyricsWriterPercent,err:= this.GetFloat("lyricsWriterPercent")
	if err!=nil {
		err:=errors.New("value of lyricsWriterPercent should be float")
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	songComposerPercent,err:= this.GetFloat("songComposerPercent")
	if err!=nil {
		err:=errors.New("value of songComposerPercent should be float")
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	publisherPercent,err:= this.GetFloat("publisherPercent")
	if err!=nil {
		err:=errors.New("value of publisherPercent should be float")
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	userPercent,err:= this.GetFloat("userPercent")
	if err!=nil {
		err:=errors.New("value of user percent should be float")
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	fileNamePrefix := util.RandomPathFromFileName("file")
	fileName:= fileNamePrefix+".mp3"

	typeOfNft:= common.TYPE_NFT_MUSIC
	nameOfNftType:= common.TYPE_NFT_MUSIC
	// ---------------------------------------
	// calculate ciphertext and save
	// ---------------------------------------
	data,err:= util.ReadFile(reqBaseInfo.File)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}

	logs.Debug("len of data", len(data))
	nonce := make([]byte, util.Aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	cipherText := util.Aesgcm.Seal(nonce, nonce, data, nil)

	// saving ciphertext
	logs.Debug("saving ciphertext")
	cipherSavingPath:= path.Join(common.ENCRYPTION_FILE_PATH, typeOfNft, fileName)
	err = ioutil.WriteFile(cipherSavingPath, cipherText, 0777)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("saving file", fileName, "to", cipherSavingPath)

	// ---------------------------------------
	// resize image and save to folder market
	// ---------------------------------------
	iconFileName:= util.RandomPathFromFileName(iconFileHeader.Filename)

	data,err= util.ReadFile(iconFile)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}

	originImage, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	newImage := resize.Resize(80, 80, originImage, resize.Lanczos3)
	filePath := path.Join(common.MARKET_PATH, nameOfNftType, iconFileName)

	out, err := os.Create(filePath)
	defer out.Close()
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	err = jpeg.Encode(out, newImage, nil)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}

	// ---------------------------------------
	// initialize nft info
	// ---------------------------------------
	//nft info
	nftLdefIndex:= util.RandomNftLdefIndex(typeOfNft)

	// ---------------------------------------
	// save nft info to database
	// ---------------------------------------
	o:= orm.NewOrm()
	o.Begin()
	// set nftInfo
	nftInfo:= models.NftInfo{
		NftLdefIndex: nftLdefIndex,
		NftType: typeOfNft,
		NftName: reqBaseInfo.NftName,
		ShortDescription: reqBaseInfo.ShortDesc,
		LongDescription: reqBaseInfo.LongDesc,
		FileName: fileName,
		NftParentLdef: "",
	}
	_,err=o.Insert(&nftInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert nft into database")
		sendError(&this.Controller,err, 500)
		return
	}
	// set dat info
	DatInfo:= models.DatNftInfo{
		NftLdefIndex: nftLdefIndex,
		IconFileName:iconFileName,
	}
	_,err= o.Insert(&DatInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert dat into database")
		sendError(&this.Controller,err, 500)
		return
	}
	// ---------------------------------------
	// save nft market info to database
	// ---------------------------------------
	userMarketInfo:= models.UserMarketInfo{
		Uuid: reqBaseInfo.Uuid,
	}
	err=o.Read(&userMarketInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("unknown error when query db")
		sendError(&this.Controller,err, 500)
		return
	}
	nftMarketInfo:= models.NftMarketInfo{
		NftLdefIndex: nftLdefIndex,
		SellerWallet: userMarketInfo.Wallet,
		SellerUuid: reqBaseInfo.Uuid,
		Price: price,
		Qty: number,
		NumSold: 0,
	}
	_,err=o.Insert(&nftMarketInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert nft market into database")
		sendError(&this.Controller,err, 500)
		return
	}

	datMarketInfo:= models.DatNftMarketInfo{
		NftLdefIndex: nftLdefIndex,
		AllowAirdrop: allowAirdrop,
		CreatorPercent: creatorPercent,
		LyricsWriterPercent: lyricsWriterPercent,
		SongComposerPercent: songComposerPercent,
		PublisherPercent: publisherPercent,
		UserPercent: userPercent,
	}
	_,err=o.Insert(&datMarketInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert dat market into database")
		sendError(&this.Controller,err, 500)
		return
	}
	// ---------------------------------------
	// record nft in market Place
	// ---------------------------------------
	mkPlace:= models.NftMarketPlace{
		NftLdefIndex: nftLdefIndex,
		NftType: typeOfNft,
		MpId: common.MARKETPLACE_ID,
		Active: true,
		ActiveTicker: common.ACTIVE_TICKER,
	}
	_,err=o.Insert(&mkPlace)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert dat into market place")
		sendError(&this.Controller,err, 500)
		return
	}
	o.Commit()
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	res:= &common.DatNftMarketInfo{
		DatNftInfo: common.DatNftInfo{
			NftInfo: common.NftInfo{
				NftLdefIndex: nftLdefIndex,
				NftType: typeOfNft,
				NftName: reqBaseInfo.NftName,
				ShortDescription: reqBaseInfo.ShortDesc,
				LongDescription: reqBaseInfo.LongDesc,
				FileName: fileName,
				NftParentLdef: "",
			},
		},
		NftMarketInfo:common.NftMarketInfo{
			SellerWallet: userMarketInfo.Wallet,
			SellerUuid: userMarketInfo.Uuid,
			Price: price,
			Qty: number,
			NumSold: 0,
		} ,
	}
	this.Data["json"] = res
	this.ServeJSON()
	// ---------------------------------------
	// call smart contract to create nft
	// ---------------------------------------
	// TODO use message queue to deal with transaction.
}

func (this *UploadController) uploadOther(reqBaseInfo *uplodBaseInfo) {
	fileNamePrefix := util.RandomPathFromFileName("file")
	fileName:= fileNamePrefix+".jpg"
	parentNftLdefIndex:= this.GetString("parent")

	typeOfNft:= common.TYPE_NFT_OTHER
	nameOfNftType:= common.TYPE_NFT_OTHER
	// ---------------------------------------
	// calculate ciphertext and save
	// ---------------------------------------
	data,err:= util.ReadFile(reqBaseInfo.File)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}

	logs.Debug("len of data", len(data))
	nonce := make([]byte, util.Aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	cipherText := util.Aesgcm.Seal(nonce, nonce, data, nil)

	// saving ciphertext
	logs.Debug("saving ciphertext")
	cipherSavingPath:= path.Join(common.ENCRYPTION_FILE_PATH, typeOfNft, fileName)
	err = ioutil.WriteFile(cipherSavingPath, cipherText, 0777)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("saving file", fileName, "to", cipherSavingPath)

	// ---------------------------------------
	// resize image and save to folder market
	// ---------------------------------------
	originImage, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	newImage := resize.Resize(200, 200, originImage, resize.Lanczos3)
	filePath:= path.Join(common.MARKET_PATH,nameOfNftType,fileName)
	err=util.SaveImage(newImage,filePath)
	if err!=nil {
		logs.Error(err.Error())
		err:= errors.New("can not save resized image")
		sendError(&this.Controller,err, 500)
		return
	}

	// ---------------------------------------
	// initialize nft info
	// ---------------------------------------
	//nft info
	nftLdefIndex:= util.RandomNftLdefIndex(typeOfNft)
	nftLifeIndex:= util.SmallRandInt()
	nftPowerIndex:= util.SmallRandInt()
	// nft market info
	price:= util.SmallRandInt()
	qty:= util.SmallRandInt()
	// ---------------------------------------
	// save nft info to database
	// ---------------------------------------
	o:= orm.NewOrm()
	o.Begin()
	// set nftInfo
	nftInfo:= models.NftInfo{
		NftLdefIndex: nftLdefIndex,
		NftType: typeOfNft,
		NftName: reqBaseInfo.NftName,
		ShortDescription: reqBaseInfo.ShortDesc,
		LongDescription: reqBaseInfo.LongDesc,
		FileName: fileName,
		NftParentLdef: "",
	}
	_,err=o.Insert(&nftInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert nft into database")
		sendError(&this.Controller,err, 500)
		return
	}
	// set avatar info
	avatarInfo:= models.AvatarNftInfo{
		NftLdefIndex: nftLdefIndex,
		NftLifeIndex: nftLifeIndex,
		NftPowerIndex: nftPowerIndex,
	}
	_,err= o.Insert(&avatarInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert other into database")
		sendError(&this.Controller,err, 500)
		return
	}
	// ---------------------------------------
	// save nft market info to database
	// ---------------------------------------
	userMarketInfo:= models.UserMarketInfo{
		Uuid: reqBaseInfo.Uuid,
	}
	err=o.Read(&userMarketInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("unknown error when query db")
		sendError(&this.Controller,err, 500)
		return
	}
	nftMarketInfo:= models.NftMarketInfo{
		NftLdefIndex: nftLdefIndex,
		SellerWallet: userMarketInfo.Wallet,
		SellerUuid: reqBaseInfo.Uuid,
		Price: price,
		Qty: qty,
		NumSold: 0,
	}
	_,err=o.Insert(&nftMarketInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert nft market into database")
		sendError(&this.Controller,err, 500)
		return
	}

	otherMarketInfo:= models.OtherNftMarketInfo{
		NftLdefIndex: nftLdefIndex,
	}
	_,err=o.Insert(&otherMarketInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert other market into database")
		sendError(&this.Controller,err, 500)
		return
	}
	// ---------------------------------------
	// record nft in market Place
	// ---------------------------------------
	mkPlace:= models.NftMarketPlace{
		NftLdefIndex: nftLdefIndex,
		NftType: typeOfNft,
		MpId: common.MARKETPLACE_ID,
		Active: true,
		ActiveTicker: common.ACTIVE_TICKER,
	}
	_,err=o.Insert(&mkPlace)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not insert other into market place")
		sendError(&this.Controller,err, 500)
		return
	}
	o.Commit()
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	res:= &common.OtherNftMarketInfo{
		OtherNftInfo: common.OtherNftInfo{
			NftInfo: common.NftInfo{
				NftLdefIndex: nftLdefIndex,
				NftType: typeOfNft,
				NftName: reqBaseInfo.NftName,
				ShortDescription: reqBaseInfo.ShortDesc,
				LongDescription: reqBaseInfo.LongDesc,
				FileName: fileName,
				NftParentLdef: parentNftLdefIndex,
			},
		},
		NftMarketInfo:common.NftMarketInfo{
			SellerWallet: userMarketInfo.Wallet,
			SellerUuid: userMarketInfo.Uuid,
			Price: price,
			Qty: qty,
			NumSold: 0,
		} ,
	}
	this.Data["json"] = res
	this.ServeJSON()
	// ---------------------------------------
	// call smart contract to create nft
	// ---------------------------------------
	// TODO use message queue to deal with transaction.
}

func (this *UploadController) Upload() {
	kind := this.GetString(":kind")
	if err:= util.ValidNftName(kind); err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller, err, 400)
		return
	}
	uuid := this.GetString("uuid")
	userMarketInfo := models.UserMarketInfo{
		Uuid: uuid,
	}

	// check if user bind wallet
	o := orm.NewOrm()
	err := o.Read(&userMarketInfo)
	if err != nil {
		if err == orm.ErrNoRows {
			err := errors.New("user has not bind wallet")
			logs.Error(err.Error())
			sendError(&this.Controller, err, 400)
			return
		} else {
			logs.Error(err.Error())
			err := errors.New("unknown error when query db")
			sendError(&this.Controller, err, 400)
			return
		}
	}
	nftName := this.GetString("nftName")
	shortDesc := this.GetString("shortDesc")
	longDesc := this.GetString("longDesc")
	file, _, err := this.GetFile("file")
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("can not parse file")
		sendError(&this.Controller, err, 400)
		return
	}
	reqBaseInfo := &uplodBaseInfo{
		File:      file,
		Uuid:      uuid,
		NftName:   nftName,
		ShortDesc: shortDesc,
		LongDesc:  longDesc,
	}
	switch kind {
	case common.NAME_NFT_AVATAR:
		this.uploadAvatar(reqBaseInfo)
	case common.NAME_NFT_OTHER:
		this.uploadOther(reqBaseInfo)
	case common.NAME_NFT_MUSIC:
		this.uploadDat(reqBaseInfo)
	}
}