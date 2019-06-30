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
	// call smart contract to create nft
	// ---------------------------------------
	// TODO use message queue to deal with transaction.
}

func (this *UploadController) uploadDat(reqBaseInfo *uplodBaseInfo) {
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
}

func (this *UploadController) uploadOther(reqBaseInfo *uplodBaseInfo) {
	fileNamePrefix := util.RandomPathFromFileName("file")
	fileName:= fileNamePrefix+".jpg"

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
	o.Commit()
	// ---------------------------------------
	// call smart contract to create nft
	// ---------------------------------------
	// TODO use message queue to deal with transaction.
}

func (this *UploadController) Upload() {
	kind := this.GetString(":kind")
	if kind != common.NAME_NFT_AVATAR && kind!= common.NAME_NFT_OTHER && kind!=common.NAME_NFT_MUSIC {
		err:= errors.New("no such nft kind")
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	uuid:= this.GetString("uuid")
	userMarketInfo:= models.UserMarketInfo{
		Uuid: uuid,
	}

	// check if user bind wallet
	o:=orm.NewOrm()
	err:=o.Read(&userMarketInfo)
	if err!=nil {
		if err == orm.ErrNoRows {
			err:=errors.New("user has not bind wallet")
			logs.Error(err.Error())
			sendError(&this.Controller,err,400)
			return
		} else {
			logs.Error(err.Error())
			err:=errors.New("unknown error when query db")
			sendError(&this.Controller,err,400)
			return
		}
	}
	nftName:= this.GetString("nftName")
	shortDesc:= this.GetString("shortDesc")
	longDesc:= this.GetString("longDesc")
	file, _, err := this.GetFile("file")
	if err != nil {
		logs.Error(err.Error())
		err:=errors.New("can not parse file")
		sendError(&this.Controller,err,400)
		return
	}
	reqBaseInfo:= &uplodBaseInfo{
		File:file,
		Uuid:uuid,
		NftName:nftName,
		ShortDesc:shortDesc,
		LongDesc:longDesc,
	}
	switch kind {
	case common.NAME_NFT_AVATAR:
		this.uploadAvatar(reqBaseInfo)
	case common.NAME_NFT_OTHER:
		this.uploadOther(reqBaseInfo)
	case common.NAME_NFT_MUSIC:
		this.uploadDat(reqBaseInfo)
	}

	// resize image and save
	var marketFileName string = fileName
	if kind == NAME_NFT_AVATAR || kind == NAME_NFT_OTHER {
	} else if kind == NAME_NFT_MUSIC {
		iconFile,iconFileHeader,err:= this.GetFile("icon")
		if err!=nil {
			logs.Error(err.Error())
			sendError(&this.Controller,err, 500)
			return
		}
		iconFileName:= RandomPathFromFileName(iconFileHeader.Filename)+".jpg"
		marketFileName = iconFileName
		data,err:= ReadFileFromRequest(iconFile)
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
		filePath := path.Join(MARKET_PATH, NAME_NFT_MUSIC, iconFileName)

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
	}

	// set other nft metadata info
	logs.Info("address of user", walletAddress, "kind of creating", kind)
	logs.Debug("name",nftName)
	logs.Debug("shortDesc",shortDesc)
	logs.Debug("longDesc",longDesc)
	// rand set power and life of nft
	nftPowerIndex = big.NewInt(int64(smallRandInt()))
	nftLifeIndex = big.NewInt(int64(smallRandInt()))

	// rand set character id
	nftCharacterId = strconv.FormatInt(time.Now().UnixNano()|rand2.Int63(), 10)
	h := md5.New()
	io.WriteString(h, nftLdefIndex)
	nftCharacterId = new(big.Int).SetBytes(h.Sum(nil)[:4]).String()

	// generate random nftLdefIndex
	nftLdefIndex = strconv.FormatInt(time.Now().UnixNano()|rand2.Int63(), 10)
	h = md5.New()
	io.WriteString(h, nftLdefIndex)
	nftLdefIndex = new(big.Int).SetBytes(h.Sum(nil)[:6]).String()
	// random public key TODO
	publicKey = []byte("2213")
	if kind == NAME_NFT_AVATAR {
		nftType = TYPE_NFT_AVATAR
		nftLdefIndex = "A" + nftLdefIndex
		distIndex = "0"   // TODO
		typeId = "01"
	} else if kind == NAME_NFT_MUSIC {
		nftType = TYPE_NFT_MUSIC
		nftLdefIndex = "M" + nftLdefIndex
		distIndex = "0"
		typeId = "02"
		// create nft
	} else if kind == NAME_NFT_OTHER {
		nftType = TYPE_NFT_OTHER
		nftLdefIndex = "O" + nftLdefIndex
		distIndex = "0"
		typeId = "05"
		nftParentLdef = this.GetString("parent")
		logs.Debug("parent ldef index",nftParentLdef)
	}
	logs.Info("nftLdefindex", nftLdefIndex)
	o:= orm.NewOrm()
	o.Begin()   //start transaction
	// store nft info to database

	nftInfo := &models.NftInfoTable{
		NftLdefIndex:  nftLdefIndex,
		NftType:       nftType,
		NftName:       nftName,
		DistIndex:     distIndex,
		NftLifeIndex:  nftLifeIndex.Int64(),
		NftPowerIndex: nftPowerIndex.Int64(),
		NftCharacId:   nftCharacterId,
		PublicKey:     string(publicKey),
	}
	_, err = o.Insert(nftInfo)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Info("insert nft info from", nickname, "to db, number:")

	// store mapping info to database

	nftAdminID := strconv.FormatInt(time.Now().UnixNano()|rand2.Int63(), 10)
	h = md5.New()
	io.WriteString(h, nftLdefIndex)
	nftAdminID = new(big.Int).SetBytes(h.Sum(nil)[:4]).String()

	mappingInfo := &models.NftMappingTable{
		NftLdefIndex: nftLdefIndex,
		TypeId:       typeId,
		FileName:     fileName,
		Key:          "0x01", //TODO use mk key to encrypt file
		NftAdminId:   nftAdminID,
		NftParentLdef: nftParentLdef,
		IconFileName: marketFileName,
	}

	// generate random market place id
	mpId := strconv.FormatInt(time.Now().UnixNano()|rand2.Int63(), 10)
	h = md5.New()
	io.WriteString(h, mpId)
	mpId = new(big.Int).SetBytes(h.Sum(nil)[:8]).String()

	_, err = o.Insert(mappingInfo)
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("insert mapping info from", nickname, "to db, number:")

	// store marketplace info to database
	marketInfo := &models.NftMarketTable{
		MpId:         mpId,
		NftLdefIndex: nftLdefIndex,
		NftAdminId:   nftAdminID,
		Price:        price,   //TODO set price
		Qty:          qty, // TODO All for selling
		NumSold:      0,   // already sold
		Active:       true,
		ActiveTicker: ACTIVE_TICKER,
		SellerWalletId:walletAddress,
		SellerNickname:nickname,
		AllowAirdrop:allowAirdrop,
		CreatorPercent: creatorPercent,
		LyricsWriterPercent: lyricsWriterPercent,
		SongComposerPercent: songComposerPercent,
		PublisherPercent: publisherPercent,
		UserPercent:userPercent,
	}
	_, err = o.Insert(marketInfo)
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("insert marketplace info from", nickname, "to db")

	// store admin table to database
	nftAdminInfo:= &models.NftItemAdmin{
		NftAdminId: nftAdminID,
		ShortDescription:shortDesc, //TODO
		LongDescription: longDesc, //TODO
		NumDistribution: qty, //TODO
	}
	_,err = o.Insert(nftAdminInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller,err,500)
		return
	}
	logs.Debug("insert nftadmin table from", nickname, "to db")

	// if unsuccessful to create nft, delete file
	tokenId, _ := new(big.Int).SetString(nftLdefIndex[1:], 10)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	// insert to wallet address table
	_,err = o.QueryTable("market_user_table").Filter("nickname",nickname).Update(orm.Params{
		"count": orm.ColValue(orm.ColAdd,1),
	})
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller,err,500)
		return
	}
	logs.Debug("insert to wallet address table success")

	logs.Warn("create nft, tokenId", tokenId,"for",walletAddress)
	txErr := this.C.account.SendFunction(this.C.smartContract,
		nil,
		nft.FuncMint,
		common.HexToAddress(walletAddress),
		tokenId,
		nftType,
		nftName,
		nftLdefIndex,
		distIndex,
		nftLifeIndex,
		nftPowerIndex,
		nftCharacterId,
		publicKey,
	)
	err = <-txErr
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("create nft success")

	err = o.Commit()
	if err!=nil {
		o.Rollback()
		err:=errors.New("commit to database fail")
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("insert all success, return")
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	res:= &nftInfoListRes{
		ActiveTicker:ACTIVE_TICKER,
		LongDesc:longDesc,
		ShortDesc:shortDesc,
		NftCharacId: nftCharacterId,
		NftLdefIndex:nftLdefIndex,
		NftLifeIndex: nftLifeIndex.Int64(),
		NftName: nftName, //todo
		NftPowerIndex: nftPowerIndex.Int64(),
		NftValue: price,
		Qty: qty,
		SupportedType: nftType,
		Thumbnail: marketFileName,
		CreatorPercent:creatorPercent,
		LyricsWriterPercent:lyricsWriterPercent,
		SongComposerPercent:songComposerPercent,
		PublisherPercent:publisherPercent,
		UserPercent:userPercent,
	}
	this.Data["json"] = res
	this.ServeJSON()
	return
}
