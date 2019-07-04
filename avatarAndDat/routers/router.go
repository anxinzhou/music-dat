package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/chainHelper"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/transactionQueue"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/mobile"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/web"
)

func init() {
	logs.Info("initial router")

	// set chain handler
	chainHandler,err:= chainHelper.NewChainHandler(&chainHelper.ChainConfig{
		ContractAddress: beego.AppConfig.String("contractAddress"),
		Port: beego.AppConfig.String("chainWS"),
		Account: chainHelper.AccountConfig{
			Address: beego.AppConfig.String("masterAddress"),
			PrivateKey:beego.AppConfig.String("masterPrivateKey"),
		},
	})
	if err!=nil {
		logs.Error(err.Error())
		panic(err)
	}

	// start transaction sending queue
	// TODO use message queue instead of go channel to ensure fault tolerance
	transactionQueueSzie,err:=beego.AppConfig.Int("transactionQueueSize")
	if err!=nil {
		panic("transaction queue size should be int")
	}
	transactionPool:= make(chan interface{},transactionQueueSzie)
	transactionQueue:= transactionQueue.NewTransactionQueue(transactionPool,chainHandler)
	logs.Info("start transaction queue")
	go transactionQueue.Start()

	// set router
	m:= mobile.NewManager()
	m.Init()
	m.SetChainHandler(chainHandler)
	m.SetTransactionQueue(transactionQueue)
	wsHandler:= &mobile.WebSocketHandler{
		M: m,
	}
	chainHelper:= web.NewChainHelper()
	upLoadController:= &web.UploadController{}
	upLoadController.C = chainHelper
	upLoadController.TransactionQueue = transactionQueue
	nftListController:= &web.NftListController{}
	nftListController.C = chainHelper
	rewardController:= &web.RewardController{}
	rewardController.C = chainHelper
	rewardController.TransactionQueue = transactionQueue
	childrenOfNFTController := &web.ChildrenOfNFTController{}
	childrenOfNFTController.C = chainHelper
	numOfChildrenController:= &web.NumOfChildrenController{}
	numOfChildrenController.C = chainHelper
	marketTransactionHistoryController:= &web.MarketTransactionHistoryController{}
	marketTransactionHistoryController.C = chainHelper
	nicknameController:= &web.NicknameController{}
	introController:= &web.IntroController{}
	avatarController:= &web.AvatarController{}
	walletController:= &web.WalletController{}

	beego.Router("/ws", wsHandler)
	beego.Router("/admin",&web.AdminController{},"post:Login")
	beego.Router("/file/:kind(avatar|dat|other)",upLoadController,"post:Upload")
	beego.Router("/nftList/:kind(avatar|dat|other)/:uuid:string",nftListController)
	beego.Router("/rewardDat/:uuid:string",rewardController,"get:RewardDat")
	beego.Router("/nfts/:parentIndex:string/children", childrenOfNFTController)
	beego.Router("/nfts/:parentIndex:string/balance", numOfChildrenController)
	beego.Router("/market/transactionHistory/:uuid:string",marketTransactionHistoryController,"get:MarketTransactionHistory")
	beego.Router("/profile/:uuid/nickname",nicknameController,"get:GetNickname;post:SetNickname")
	beego.Router("/profile/:uuid/avatar",avatarController,"get:GetAvatar;post:SetAvatar")
	beego.Router("/profile/:uuid/intro",introController,"get:GetIntro;post:SetIntro")
	beego.Router("/profile/:uuid/wallet",walletController,"get:GetWallet;post:SetWallet")
}


