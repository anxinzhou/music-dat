package http

import (
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xxRanger/blockchainUtil/chain"
	"github.com/xxRanger/blockchainUtil/contract"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/blockchainUtil/sender"
)

type ErrorResponse struct {
	Reason string `json:"reason"`
}

type ContractController struct {
	beego.Controller
	C *ChainHelper
}

type ChainHelper struct {
	account *sender.User
	smartContract contract.Contract
}

func NewChainHelper() *ChainHelper {
	address := common.HexToAddress(beego.AppConfig.String("masterAddress"))
	privateKey := beego.AppConfig.String("masterPrivateKey")
	pk,err:=crypto.HexToECDSA(privateKey)
	if err!=nil {
		panic(err)
	}
	account:=sender.NewUser(address,pk)
	port:= beego.AppConfig.String("chainWS")
	c,err:=chain.NewEthClient(port)
	if err!=nil {
		panic(err)
	}
	account.BindEthClient(c,sender.CHAIN_KIND_PRIVATE)
	contractAddress:= beego.AppConfig.String("contractAddress")
	smartContract:= nft.NewNFT(common.HexToAddress(contractAddress))
	smartContract.BindClient(c)
	return &ChainHelper{
		account:account,
		smartContract: smartContract,
	}
}
