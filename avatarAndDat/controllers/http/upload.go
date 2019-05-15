package http

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"github.com/astaxie/beego/logs"
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
	fileName := header.Filename
	h := md5.New()
	io.WriteString(h, fileName)
	fileName = new(big.Int).SetBytes(h.Sum(nil)[:8]).String()
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}

	if kind == NAME_NFT_AVATAR {
		fileName = fileName + ".jpg"
	} else if kind == NAME_NFT_MUSIC {
		fileName = fileName + ".mp3"
	} else {
		panic("unexpected kind name")
	}

	// calculate ciphertext
	var dataBuffer bytes.Buffer
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	for {
		if err == io.EOF {
			dataBuffer.Write(buffer[:n])
			break;
		} else if err != nil {
			logs.Error(err.Error())
			sendError(&this.Controller,err, 500)
			return
		} else {
			dataBuffer.Write(buffer[:n])
			n, err = file.Read(buffer)
		}
	}
	data := dataBuffer.Bytes()

	logs.Debug("len of data", len(data))
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	cipherText := aesgcm.Seal(nonce, nonce, data, nil)

	//test
	//nonce,ct := cipherText[:aesgcm.NonceSize()],cipherText[aesgcm.NonceSize():]
	//plainText,err:=aesgcm.Open(nil,nonce,ct,nil)
	//if err!=nil{
	//	panic(err)
	//}
	//if bytes.Compare(plainText,data) !=0 {
	//	logs.Emergency("can not pass test")
	//	return
	//} else {
	//	logs.Emergency("pass test")
	//	return
	//}

	// saving ciphertext
	logs.Debug("saving ciphertext")
	var cipherSavingPath string
	if kind == NAME_NFT_AVATAR {
		cipherSavingPath = path.Join(ENCRYPTION_FILE_PATH, NAME_NFT_AVATAR, fileName)
	} else if kind == NAME_NFT_MUSIC {
		cipherSavingPath = path.Join(ENCRYPTION_FILE_PATH, NAME_NFT_MUSIC, fileName)
	} else {
		logs.Error("unexpected type")
		panic(err)
	}
	err = ioutil.WriteFile(cipherSavingPath, cipherText, 0777)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("saving file", fileName, "to", cipherSavingPath)

	// resize image and save
	var marketFilePath string = ""
	if kind == NAME_NFT_AVATAR {
		originImage, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			logs.Error(err.Error())
			sendError(&this.Controller,err, 500)
			return
		}
		newImage := resize.Resize(200, 200, originImage, resize.Lanczos3)
		filePath := path.Join(MARKET_PATH, NAME_NFT_AVATAR, fileName)
		marketFilePath = filePath
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

	// create nft
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
	)

	// TODO
	price:= 1
	qty:=100
	//
	//type payload struct {
	//	Address string `json:"address"`
	//	ShortDesc string `json:"shortDesc"`
	//	LongDesc string `json:"longDesc"`
	//	NftName string `json:"nftName"`
	//}

	//var p payload
	//pData,err:=ioutil.ReadAll(this.Ctx.Request.Body)
	//if err!=nil {
	//	logs.Error(err.Error())
	//	sendError(&this.Controller,err,400)
	//	return
	//}
	//err=json.Unmarshal(pData,&p)
	//if err!=nil {
	//	logs.Error(err.Error())
	//	sendError(&this.Controller,err,400)
	//	return
	//}
	//user:= p.Address
	//nftName = p.NftName
	//shortDesc = p.ShortDesc
	//longDesc = p.LongDesc

	user := this.GetString("address")
	//// get input from user
	nftName = this.GetString("nftName")
	shortDesc = this.GetString("shortDesc")
	longDesc = this.GetString("longDesc")
	logs.Info("address of user", user, "kind of creating", kind)
	logs.Debug("name",nftName)
	logs.Debug("shortDesc",shortDesc)
	logs.Debug("longDesc",longDesc)
	// rand set power and life of nft
	nftPowerIndex = big.NewInt(int64(smallRandInt()))
	nftLifeIndex = big.NewInt(int64(smallRandInt()))

	// rand set character id
	nftCharacterId = strconv.FormatInt(time.Now().UnixNano()|rand2.Int63(), 10)
	h = md5.New()
	io.WriteString(h, nftLdefIndex)
	nftCharacterId = new(big.Int).SetBytes(h.Sum(nil)[:4]).String()

	// generate random nftLdefIndex
	nftLdefIndex = strconv.FormatInt(time.Now().UnixNano()|rand2.Int63(), 10)
	h = md5.New()
	io.WriteString(h, nftLdefIndex)
	nftLdefIndex = new(big.Int).SetBytes(h.Sum(nil)[:4]).String()
	// random public key TODO
	publicKey = []byte("2213")
	if kind == NAME_NFT_AVATAR {
		nftType = TYPE_NFT_AVATAR
		nftLdefIndex = "A" + nftLdefIndex
		distIndex = "0"   // TODO
	} else if kind == NAME_NFT_MUSIC {
		nftType = TYPE_NFT_MUSIC
		nftLdefIndex = "M" + nftLdefIndex
		distIndex = "0"
		// create nft
	}
	logs.Info("nftLdefindex", nftLdefIndex)

	models.O.Begin()   //start transaction
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
	_, err = models.O.Insert(nftInfo)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Info("insert nft info from", user, "to db, number:")

	// store mapping info to database
	var typeId string
	if kind == "avatar" {
		typeId = "01"
	} else {
		typeId = "02"
	}

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
	}

	// generate random market place id
	mpId := strconv.FormatInt(time.Now().UnixNano()|rand2.Int63(), 10)
	h = md5.New()
	io.WriteString(h, mpId)
	mpId = new(big.Int).SetBytes(h.Sum(nil)[:8]).String()

	_, err = models.O.Insert(mappingInfo)
	if err != nil {
		models.O.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("insert mapping info from", user, "to db, number:")

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
	}
	_, err = models.O.Insert(marketInfo)
	if err != nil {
		models.O.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("insert marketplace info from", user, "to db")

	// store admin table to database
	nftAdminInfo:= &models.NftItemAdmin{
		NftAdminId: nftAdminID,
		ShortDescription:shortDesc, //TODO
		LongDescription: longDesc, //TODO
		NumDistribution: qty, //TODO
	}
	_,err = models.O.Insert(nftAdminInfo)
	if err!=nil {
		models.O.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller,err,500)
		return
	}
	logs.Debug("insert nftadmin table from", user, "to db")

	// if unsuccessful to create nft, delete file
	tokenId, _ := new(big.Int).SetString(nftLdefIndex[1:], 10)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	// insert to wallet address table
	walletInfo:= &models.MarketUserTable{
		WalletId: user,
		Count: 1,
	}
	_,err=models.O.InsertOrUpdate(walletInfo,"count=count+1")
	if err!=nil {
		models.O.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller,err,500)
		return
	}
	logs.Debug("insert to wallet address table success")

	logs.Warn("create nft, tokenId", tokenId)
	txErr := this.C.account.SendFunction(this.C.smartContract,
		nil,
		nft.FuncMint,
		common.HexToAddress(user),
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
		models.O.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	logs.Debug("create nft success")

	models.O.Commit()
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
		Thumbnail: marketFilePath,
	}
	this.Data["json"] = res
	this.ServeJSON()
	return
}
