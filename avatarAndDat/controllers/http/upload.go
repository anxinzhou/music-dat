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
		this.sendError(err, 400)
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
			this.sendError(err, 500)
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
		this.sendError(err, 500)
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
		this.sendError(err, 500)
		return
	}
	logs.Debug("saving file", fileName, "to", cipherSavingPath)

	// resize image and save
	if kind == NAME_NFT_AVATAR {
		originImage, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			logs.Error(err.Error())
			this.sendError(err, 500)
			return
		}
		newImage := resize.Resize(200, 200, originImage, resize.Lanczos3)
		filePath := path.Join(MARKET_PATH, NAME_NFT_AVATAR, fileName)
		out, err := os.Create(filePath)
		defer out.Close()
		if err != nil {
			logs.Error(err.Error())
			this.sendError(err, 500)
			return
		}
		err = jpeg.Encode(out, newImage, nil)
		if err != nil {
			logs.Error(err.Error())
			this.sendError(err, 500)
			return
		}
	}

	user := this.GetString("address")
	logs.Info("address of user", user, "kind of creating", kind)

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
	)

	// generate random nftLdefIndex
	nftLdefIndex = strconv.FormatInt(time.Now().UnixNano()|rand2.Int63(), 10)
	h = md5.New()
	io.WriteString(h, nftLdefIndex)
	nftLdefIndex = new(big.Int).SetBytes(h.Sum(nil)[:4]).String()
	// random public key TODO
	publicKey = []byte("2213")
	if kind == NAME_NFT_AVATAR {
		nftType = TYPE_NFT_AVATAR
		nftName = "Turtle yellow"
		nftLdefIndex = "A" + nftLdefIndex
		distIndex = "0"
		nftLifeIndex = big.NewInt(100)
		nftPowerIndex = big.NewInt(100)
		nftCharacterId = "MGM09-G73673"
	} else if kind == NAME_NFT_MUSIC {
		nftType = TYPE_NFT_MUSIC
		nftName = "IP693794957349"
		nftLdefIndex = "M" + nftLdefIndex
		distIndex = "0"
		nftLifeIndex = big.NewInt(0)
		nftPowerIndex = big.NewInt(0)
		nftCharacterId = ""
		// create nft
	}
	logs.Info("nftLdefindex", nftLdefIndex)
	// if unsuccessful to create nft, delete file
	tokenId, _ := new(big.Int).SetString(nftLdefIndex[1:], 10)
	if err != nil {
		logs.Error(err.Error())
		this.sendError(err, 500)
		return
	}

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
		logs.Error(err.Error())
		this.sendError(err, 500)
		return
	}
	logs.Debug("create nft success")
	// TODO in case update database fail while create nft success
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
		this.sendError(err, 500)
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
		this.sendError(err, 500)
		return
	}
	logs.Debug("insert mapping info from", user, "to db, number:")

	// store marketplace info to database
	marketInfo := &models.NftMarketTable{
		MpId:         mpId,
		NftLdefIndex: nftLdefIndex,
		NftAdminId:   nftAdminID,
		Price:        1,   //TODO set price
		Qty:          100, // TODO All for selling
		NumSold:      0,   // already sold
		Active:       true,
		ActiveTicker: "berry",
	}
	_, err = models.O.Insert(marketInfo)
	if err != nil {
		models.O.Rollback()
		logs.Error(err.Error())
		this.sendError(err, 500)
		return
	}
	logs.Debug("insert marketplace info from", user, "to db")

	// store admin table to database
	nftAdminInfo:= &models.NftItemAdmin{
		NftAdminId: nftAdminID,
		ShortDescription:"todo", //TODO
		LongDescription: "todo", //TODO
		NumDistribution: 100, //TODO
	}
	_,err = models.O.Insert(nftAdminInfo)
	if err!=nil {
		models.O.Rollback()
		logs.Error(err.Error())
		this.sendError(err,500)
		return
	}
	logs.Debug("insert nftadmin table from", user, "to db")
	models.O.Commit()
	logs.Debug("insert all success, return")
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	return
}
