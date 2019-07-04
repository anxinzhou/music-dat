package web

import (
	"bytes"
	"crypto/rand"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/nfnt/resize"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/transactionQueue"
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
)


type NumOfChildrenController struct {
	ContractController
}

func (this *NumOfChildrenController) Get() {
	parentIndex := this.Ctx.Input.Param(":parentIndex")
	o := orm.NewOrm()
	r := o.Raw(`
		select count(a.nft_ldef_index) as num  
		from nft_info as a 
		where nft_parent_ldef = ? `, parentIndex)
	type CountQuery struct {
		Num int
	}
	var queryResult CountQuery
	err := r.QueryRow(&queryResult)
	if err != nil {
		if err == orm.ErrNoRows {
			err:=errors.New("this nft does not have children")
			sendError(&this.Controller, err, 500)
			return
		} else {
			logs.Error(err.Error())
			err:=errors.New("unexpected error when query db")
			sendError(&this.Controller, err, 500)
		}
	}
	type response struct {
		Count int `json:"count"`
	}
	res := &response{
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

	o:=orm.NewOrm()
	type nftTranData struct {
		common.OtherNftInfo
		common.MarketPlaceInfo
	}
	dbEngine := beego.AppConfig.String("dbEngine")
	qb,_:=orm.NewQueryBuilder(dbEngine)
	var avatarMKPlaceInfo []nftTranData
	qb.Select("*").
		From("nft_market_place").
		InnerJoin("nft_market_info").
		On("nft_market_place.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("other_nft_market_info").
		On("nft_market_place.nft_ldef_index = other_nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_market_place.nft_ldef_index = nft_info.nft_ldef_index").
		InnerJoin("other_nft_info").
		On("nft_market_place.nft_ldef_index = other_nft_info.nft_ldef_index").
		Where("nft_parent_ldef = ?")
	sql := qb.String()
	num,err:=o.Raw(sql,parentIndex).QueryRows(&avatarMKPlaceInfo)
	if err != nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:= errors.New("unknown error when query database")
		sendError(&this.Controller, err, 500)
		return
	}
	logs.Debug("get",num,"from database")
	type response struct {
		NftTranData []nftTranData `json:"nftTranData"`
	}

	// in case no row in database
	if num == 0 {
		avatarMKPlaceInfo=make([]nftTranData,0)
	}

	res:= response{
		NftTranData: avatarMKPlaceInfo,
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &res
	this.ServeJSON()
}


type UploadController struct {
	ContractController
	TransactionQueue *transactionQueue.TransactionQueue
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
	cipherSavingPath:= path.Join(common.ENCRYPTION_FILE_PATH, nameOfNftType, fileName)
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
		NftInfo: &models.NftInfo{
			NftLdefIndex:nftLdefIndex,
		},
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
		UserInfo: &models.UserInfo{
			Uuid: reqBaseInfo.Uuid,
		},
	}
	err=o.Read(&userMarketInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("unknown error when query db")
		sendError(&this.Controller,err, 500)
		return
	}
	//  update user nft count
	_,err=o.QueryTable("user_market_info").Update(orm.Params{
		"count": orm.ColValue(orm.ColAdd,1),
	})
	if err!=nil {
		logs.Error(err.Error())
		err:=errors.New("unexpected error when query db")
		sendError(&this.Controller,err, 500)
	}
	nftMarketInfo:= models.NftMarketInfo{
		NftLdefIndex: nftLdefIndex,
		SellerWallet: userMarketInfo.Wallet,
		SellerUuid: reqBaseInfo.Uuid,
		Price: price,
		Qty: qty,
		Active: false,
		NumSold: 0,
		NftInfo: &models.NftInfo{
			NftLdefIndex:nftLdefIndex,
		},
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
		NftMarketInfo:&models.NftMarketInfo{
			NftLdefIndex:nftLdefIndex,
		},
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
		MpId: common.MARKETPLACE_ID,
		ActiveTicker: common.ACTIVE_TICKER,
		NftMarketInfo: &models.NftMarketInfo{
			NftLdefIndex:nftLdefIndex,
		},
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
	// TODO use message queue to instead of go channel deal with transaction.
	this.TransactionQueue.Append(&transactionQueue.UploadNftTransaction{
		Uuid: reqBaseInfo.Uuid,
		NftLdefIndex: nftLdefIndex,
		NftType: typeOfNft,
		NftName: reqBaseInfo.NftName,
		DistIndex: "1",
		NftLifeIndex: big.NewInt(int64(nftLifeIndex)),
		NftPowerIndex: big.NewInt(int64(nftPowerIndex)),
		PublicKey: "345435",
	})
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
	musicFileNamePrefix := util.RandomPathFromFileName("file")
	musicFileName:= musicFileNamePrefix+".mp3"

	typeOfNft:= common.TYPE_NFT_MUSIC
	nameOfNftType:= common.NAME_NFT_MUSIC
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
	cipherSavingPath:= path.Join(common.ENCRYPTION_FILE_PATH, nameOfNftType, musicFileName)
	err = ioutil.WriteFile(cipherSavingPath, cipherText, 0777)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("saving file", musicFileName, "to", cipherSavingPath)

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
		FileName: iconFileName,
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
		MusicFileName:musicFileName,
		NftInfo: &models.NftInfo{
			NftLdefIndex:nftLdefIndex,
		},
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
		UserInfo: &models.UserInfo{
			Uuid: reqBaseInfo.Uuid,
		},
	}
	err=o.Read(&userMarketInfo)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("unknown error when query db")
		sendError(&this.Controller,err, 500)
		return
	}
	//  update user nft count
	_,err=o.QueryTable("user_market_info").Update(orm.Params{
		"count": orm.ColValue(orm.ColAdd,1),
	})
	if err!=nil {
		logs.Error(err.Error())
		err:=errors.New("unexpected error when query db")
		sendError(&this.Controller,err, 500)
	}
	nftMarketInfo:= models.NftMarketInfo{
		NftLdefIndex: nftLdefIndex,
		SellerWallet: userMarketInfo.Wallet,
		SellerUuid: reqBaseInfo.Uuid,
		Price: price,
		Qty: number,
		NumSold: 0,
		Active: false,
		NftInfo: &models.NftInfo{
			NftLdefIndex: nftLdefIndex,
		},
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
		NftMarketInfo: &models.NftMarketInfo{
			NftLdefIndex:nftLdefIndex,
		},
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
		MpId: common.MARKETPLACE_ID,
		ActiveTicker: common.ACTIVE_TICKER,
		NftMarketInfo: &models.NftMarketInfo{
			NftLdefIndex:nftLdefIndex,
		},
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
				FileName: iconFileName,
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
	// TODO use message queue to instead of go channel deal with transaction.
	this.TransactionQueue.Append(&transactionQueue.UploadNftTransaction{
		Uuid: reqBaseInfo.Uuid,
		NftLdefIndex: nftLdefIndex,
		NftType: typeOfNft,
		NftName: reqBaseInfo.NftName,
		DistIndex: "1",
		NftLifeIndex: big.NewInt(1),
		NftPowerIndex: big.NewInt(1),
		PublicKey: "345435",
	})
}

func (this *UploadController) uploadOther(reqBaseInfo *uplodBaseInfo) {
	fileNamePrefix := util.RandomPathFromFileName("file")
	fileName:= fileNamePrefix+".jpg"
	parentNftLdefIndex:= this.GetString("parent")
	logs.Debug("parent index",parentNftLdefIndex)

	typeOfNft:= common.TYPE_NFT_OTHER
	nameOfNftType:= common.NAME_NFT_OTHER
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
	cipherSavingPath:= path.Join(common.ENCRYPTION_FILE_PATH, nameOfNftType, fileName)
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
	// nft market info
	price:= util.SmallRandInt()
	qty:= util.SmallRandInt()
	// ---------------------------------------
	// save nft info to database
	// ---------------------------------------
	o:= orm.NewOrm()
	o.Begin()
	// update nft count

	// set nftInfo
	nftInfo:= models.NftInfo{
		NftLdefIndex: nftLdefIndex,
		NftType: typeOfNft,
		NftName: reqBaseInfo.NftName,
		ShortDescription: reqBaseInfo.ShortDesc,
		LongDescription: reqBaseInfo.LongDesc,
		FileName: fileName,
		NftParentLdef: parentNftLdefIndex,
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
	ohterInfo:= models.OtherNftInfo{
		NftLdefIndex: nftLdefIndex,
		NftInfo: &models.NftInfo{
			NftLdefIndex:nftLdefIndex,
		},
	}
	_,err= o.Insert(&ohterInfo)
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
	//  update user nft count
	_,err=o.QueryTable("user_market_info").Update(orm.Params{
		"count": orm.ColValue(orm.ColAdd,1),
	})
	if err!=nil {
		logs.Error(err.Error())
		err:=errors.New("unexpected error when query db")
		sendError(&this.Controller,err, 500)
	}
	nftMarketInfo:= models.NftMarketInfo{
		NftLdefIndex: nftLdefIndex,
		SellerWallet: userMarketInfo.Wallet,
		SellerUuid: reqBaseInfo.Uuid,
		Price: price,
		Qty: qty,
		NumSold: 0,
		Active: false,
		NftInfo: & models.NftInfo{
			NftLdefIndex:nftLdefIndex,
		},
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
		NftMarketInfo: &models.NftMarketInfo{
			NftLdefIndex: nftLdefIndex,
		},
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
		MpId: common.MARKETPLACE_ID,
		ActiveTicker: common.ACTIVE_TICKER,
		NftMarketInfo: &models.NftMarketInfo{
			NftLdefIndex:nftLdefIndex,
		},
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
	// TODO use message queue to instead of go channel deal with transaction.
	this.TransactionQueue.Append(&transactionQueue.UploadNftTransaction{
		Uuid: reqBaseInfo.Uuid,
		NftLdefIndex: nftLdefIndex,
		NftType: typeOfNft,
		NftName: reqBaseInfo.NftName,
		DistIndex: "1",
		NftLifeIndex: big.NewInt(1),
		NftPowerIndex: big.NewInt(1),
		PublicKey: "345435",
	})
}

func (this *UploadController) Upload() {
	kind := this.Ctx.Input.Param(":kind")
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