package http

import (
	"crypto/md5"
	"github.com/astaxie/beego/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"io"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

const FILE_SAVING_PATH = "./resource/"

type UploadController struct {
	ContractController
}

func (this *UploadController) Get() {
	fileName:=this.GetString("name")
	logs.Debug("download",fileName)
	this.Ctx.Output.Download(FILE_SAVING_PATH+fileName)
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
}

func (this *UploadController) Upload() {
	kind := this.Ctx.Input.Param(":kind")
	file, header, err := this.GetFile("file")
	user := this.GetString("address")
	logs.Info("address of user", user, "kind of creating", kind)
	if err != nil {
		logs.Error(err.Error())
		this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(400)
		this.Data["json"] = &ErrorResponse{
			Reason: err.Error(),
		}
		this.ServeJSON()
		return
	}

	//save file
	var filePath string
	var fileName string
	if file != nil {
		fileName = header.Filename
		filePath = fileName   // TODO md5 for file name
		err = this.SaveToFile("file", filePath)
		if err != nil {
			logs.Error(err.Error())
			this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(400)
			this.Data["json"] = &ErrorResponse{
				Reason: err.Error(),
			}
			this.ServeJSON()
			return
		} else {
			logs.Info("save", fileName, "to", filePath)
		}
	}

	logs.Warn("create nft")
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
	nftLdefIndex = strconv.FormatInt(time.Now().UnixNano()|rand.Int63(), 10)
	h := md5.New()
	io.WriteString(h, nftLdefIndex)
	nftLdefIndex = new(big.Int).SetBytes(h.Sum(nil)[:4]).String()
	// random public key TODO
	publicKey = []byte("2213")
	if kind == "avatar" {
		nftType = "721-02"
		nftName = "Turtle yellow"
		nftLdefIndex = "A" + nftLdefIndex
		distIndex = "0"
		nftLifeIndex = big.NewInt(100)
		nftPowerIndex = big.NewInt(100)
		nftCharacterId = "MGM09-G73673"
	} else {
		nftType = "721-04"
		nftName = "IP693794957349"
		nftLdefIndex = "M" + nftLdefIndex
		distIndex = "0"
		nftLifeIndex = big.NewInt(0)
		nftPowerIndex = big.NewInt(0)
		nftCharacterId = ""
		// create nft
	}
	logs.Info("nftLdefindex",nftLdefIndex)
	// if unsuccessful to create nft, delete file
	tokenId, _ := new(big.Int).SetString(nftLdefIndex[1:], 10)
	if err != nil {
		logs.Error(err.Error())
		this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(400)
		this.Data["json"] = &ErrorResponse{
			Reason: err.Error(),
		}
		this.ServeJSON()
		return
	}

	logs.Info("tokenId",tokenId)
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
		this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(400)
		this.Data["json"] = &ErrorResponse{
			Reason: err.Error(),
		}
		this.ServeJSON()
		return
	}
	logs.Debug("create nft success")

	type NFTMappingTable struct {
		NFTLdefIndex string `orm:"pk;unique"`
		TypeID string
		FileName string
		Key string
		NFTAdminID string
	}

	type NFTMarketTable struct {
		MPId string `orm:"pk;unique"`
		NFTLdefIndex string
		NFTAdminId string
		Price float64 `orm:"digits(12);decimals(4)"`
		QTY float64 `orm:"digits(12);decimals(4)"`
		NumSold float64 `orm:"digits(12);decimals(4)"`
		Active bool
	}

	// store nft info to database

	nftInfo:=&models.NftInfoTable{
		NftLdefIndex:nftLdefIndex,
		NftType:nftType,
		NftName:nftName,
		DistIndex:distIndex,
		NftLifeIndex:nftLifeIndex.Int64(),
		NftPowerIndex:nftPowerIndex.Int64(),
		NftCharacId:nftCharacterId,
		PublicKey: string(publicKey),
	}
	_,err = models.O.Insert(nftInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(500)
		this.Data["json"] = &ErrorResponse{
			Reason: err.Error(),
		}
		this.ServeJSON()
		return
	}
	logs.Info("insert nft info from",user,"to db, number:")

	// store mapping info to database
	var typeId string
	if kind == "avatar" {
		typeId = "01"
	} else {
		typeId = "02"
	}

	nftAdminID:="12" //TODO admin table

	mappingInfo:= &models.NftMappingTable{
		NftLdefIndex: nftLdefIndex,
		TypeId: typeId,
		FileName: fileName,
		Key: "0x01",   //TODO use mk key to encrypt file
		NftAdminId: nftAdminID,
	}

	// generate random market place id
	mpId := strconv.FormatInt(time.Now().UnixNano()|rand.Int63(), 10)
	h = md5.New()
	io.WriteString(h, mpId)
	mpId = new(big.Int).SetBytes(h.Sum(nil)[:8]).String()

	_,err= models.O.Insert(mappingInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(500)
		this.Data["json"] = &ErrorResponse{
			Reason: err.Error(),
		}
		this.ServeJSON()
		return
	}
	logs.Info("insert mapping info from",user,"to db, number:")

	// store marketplace info to database
	marketInfo:= &models.NftMarketTable{
		MpId:mpId,
		NftLdefIndex:nftLdefIndex,
		NftAdminId: nftAdminID,
		Price: 1,  //TODO set price
		Qty: 100,  // TODO All for selling
		NumSold: 0, // already sold
		Active: true,
		ActiveTicker:"berry",
	}

	_,err= models.O.Insert(marketInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(500)
		this.Data["json"] = &ErrorResponse{
			Reason: err.Error(),
		}
		this.ServeJSON()
		return
	}
	logs.Info("insert marketplace info from",user,"to db")

	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	return
}
