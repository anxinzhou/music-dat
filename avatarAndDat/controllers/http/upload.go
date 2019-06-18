package http

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nfnt/resize"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"math/big"
	rand2 "math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

type UploadController struct {
	ContractController
}

func (this *UploadController) sendError(err error, statusCode int) {
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(500)
	this.Data["json"] = &ErrorResponse{
		Reason: err.Error(),
	}
	this.ServeJSON()
}

func (this *UploadController) Get() {
	fileName := this.GetString("name")
	logs.Debug("download", fileName)
	this.Ctx.Output.Download(FILE_SAVING_PATH + fileName)
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
}

func (this *UploadController) Upload() {
	kind := this.Ctx.Input.Param(":kind")
	file, header, err := this.GetFile("file")
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}

	// check parameters from website

	// nft metadata
	var (
		nftType        string
		nftName        string
		nftLdefIndex   string
		nftLifeIndex   *big.Int
		distIndex      string
		nftPowerIndex  *big.Int
		nftCharacterId string
		publicKey      []byte
		shortDesc string
		longDesc string
		nftParentLdef string
		typeId string
		qty int
		price int
		allowAirdrop bool
	)

	// set nftmetadata from website
	// TODO
	// now can set price and qty from website
	if kind == NAME_NFT_MUSIC {
		price,err = this.GetInt("price")
		if err!=nil {
			logs.Error(err.Error())
			sendError(&this.Controller,err,400)
			return
		}
		qty,err =this.GetInt("number")
		if err!=nil {
			logs.Error(err.Error())
			sendError(&this.Controller,err,400)
			return
		}
	} else {
		price=1
		qty= 100
	}

	allowAirdrop,err = this.GetBool("allowAirdrop")
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}

	nickname:= this.GetString("nickname")
	walletAddress,err := models.WalletIdOfNickname(nickname)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}

	//// get input from user
	nftName = this.GetString("nftName")
	shortDesc = this.GetString("shortDesc")
	longDesc = this.GetString("longDesc")
	fileNamePrefix := RandomPathFromFileName(header.Filename)
	var fileName string
	if kind == NAME_NFT_AVATAR {
		fileName = fileNamePrefix + ".jpg"
	} else if kind == NAME_NFT_MUSIC {
		fileName = fileNamePrefix + ".mp3"
	} else if kind == NAME_NFT_OTHER {
		fileName = fileNamePrefix + ".jpg"
	} else {
		panic("unexpected kind name")
	}

	// calculate ciphertext
	data,err:= ReadFileFromRequest(file)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}

	logs.Debug("len of data", len(data))
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	cipherText := aesgcm.Seal(nonce, nonce, data, nil)

	// saving ciphertext
	logs.Debug("saving ciphertext")
	var cipherSavingPath string
	if kind == NAME_NFT_AVATAR {
		cipherSavingPath = path.Join(ENCRYPTION_FILE_PATH, NAME_NFT_AVATAR, fileName)
	} else if kind == NAME_NFT_MUSIC {
		cipherSavingPath = path.Join(ENCRYPTION_FILE_PATH, NAME_NFT_MUSIC, fileName)
	} else if kind == NAME_NFT_OTHER {
		cipherSavingPath = path.Join(ENCRYPTION_FILE_PATH,NAME_NFT_OTHER,fileName)
	}

	err = ioutil.WriteFile(cipherSavingPath, cipherText, 0777)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("saving file", fileName, "to", cipherSavingPath)

	// resize image and save
	var marketFileName string = fileName
	if kind == NAME_NFT_AVATAR || kind == NAME_NFT_OTHER {
		originImage, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			logs.Error(err.Error())
			sendError(&this.Controller,err, 500)
			return
		}
		newImage := resize.Resize(200, 200, originImage, resize.Lanczos3)
		var filePath string
		if kind == NAME_NFT_AVATAR {
			filePath = path.Join(MARKET_PATH, NAME_NFT_AVATAR, fileName)
		} else if kind == NAME_NFT_OTHER {
			filePath = path.Join(MARKET_PATH, NAME_NFT_OTHER, fileName)
		}

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
	}
	this.Data["json"] = res
	this.ServeJSON()
	return
}
