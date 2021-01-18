package web

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/beego/beego/v2/server/web"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xxRanger/blockchainUtil/chain"
	"github.com/xxRanger/blockchainUtil/contract"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/blockchainUtil/sender"
	"io"
	"math/rand"
	"strconv"
	"time"
)

type ContractController struct {
	web.Controller
	C *ChainHelper
}

type ChainHelper struct {
	account       *sender.User
	smartContract contract.Contract
}

func NewChainHelper() *ChainHelper {
	addressStr, _ := web.AppConfig.String("masterAddress")
	address := common.HexToAddress(addressStr)
	privateKey, _ := web.AppConfig.String("masterPrivateKey")
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(err)
	}
	account := sender.NewUser(address, pk)
	port, _ := web.AppConfig.String("chainWS")
	c, err := chain.NewEthClient(port)
	if err != nil {
		panic(err)
	}

	runmode, _ := web.AppConfig.String("runmode")
	chainKind := sender.CHAIN_KIND_PUBLIC
	if runmode == "prod" {
		chainKind = sender.CHAIN_KIND_PUBLIC
	} else if runmode == "dev" || runmode == "master" || runmode == "test" {
		chainKind = sender.CHAIN_KIND_PRIVATE
	} else {
		panic("unknown run mode")
	}

	account.BindEthClient(c, chainKind)
	contractAddress, _ := web.AppConfig.String("contractAddress")
	smartContract := nft.NewNFT(common.HexToAddress(contractAddress))
	smartContract.BindClient(c)
	return &ChainHelper{
		account:       account,
		smartContract: smartContract,
	}
}

func sendError(c web.ControllerInterface, err error, statusCode int) {
	type ErrorResponse struct {
		Reason string `json:"reason"`
	}
	controller := c.(*web.Controller)
	controller.Ctx.ResponseWriter.ResponseWriter.WriteHeader(statusCode)
	controller.Data["json"] = &ErrorResponse{
		Reason: err.Error(),
	}
	controller.ServeJSON()
}

func generateAccessToken() string {
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().UnixNano()|rand.Int63(), 10))
	accessToken := hex.EncodeToString(h.Sum(nil))
	return accessToken
}
