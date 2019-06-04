package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/http"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/ws"
)

func init() {
	logs.Info("initial router")

	m:=ws.NewManager()
	m.Init()
	// set chain handler
	chainHandler,err:= ws.NewChainHandler(&ws.ChainConfig{
		ContractAddress: beego.AppConfig.String("contractAddress"),
		Port: beego.AppConfig.String("chainWS"),
		Account: ws.AccountConfig{
			Address: beego.AppConfig.String("masterAddress"),
			PrivateKey:beego.AppConfig.String("masterPrivateKey"),
		},
	})
	if err!=nil {
		logs.Error(err.Error())
		panic(err)
	}
	m.SetChainHandler(chainHandler)
	wsHandler:= &ws.WebSocketHandler{
		M: m,
	}
	chainHelper:= http.NewChainHelper()
	upLoadController:= &http.UploadController{}
	upLoadController.C = chainHelper
	nftBalanceController:= &http.NftBalanceController{}
	nftBalanceController.C = chainHelper
	nftListController:= &http.NftListController{}
	nftListController.C = chainHelper
	rewardController:= &http.RewardController{}
	rewardController.C = chainHelper
	childrenOfNFTController := &http.ChildrenOfNFTController{}
	childrenOfNFTController.C = chainHelper
	numOfChildrenController:= &http.NumOfChildrenController{}
	numOfChildrenController.C = chainHelper
	marketTransactionHistoryController:= &http.MarketTransactionHistoryController{}
	marketTransactionHistoryController.C = chainHelper

	beego.Router("/", &http.MainController{})
	beego.Router("/ws", wsHandler)
	beego.Router("/admin",&http.AdminController{},"post:Login")
	beego.Router("/file/:kind(avatar|dat|other)",upLoadController,"get:Get;post:Upload")
	beego.Router("/balance/:nickname:string",nftBalanceController)
	beego.Router("/nftList/:nickname:string",nftListController)
	beego.Router("/rewardDat/:nickname:string",rewardController,"get:RewardDat")
	beego.Router("/nfts/:parentIndex:string/children", childrenOfNFTController)
	beego.Router("/nfts/:parentIndex:string/balance", numOfChildrenController)
	beego.Router("/wallet",&http.ImportWalletController{},"post:ImportWallet")
	beego.Router("/market/transactionHistory/:nickname:string",marketTransactionHistoryController,"get:MarketTransactionHistory")
}


