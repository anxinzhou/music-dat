package routers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/chainHelper"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/transactionQueue"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/mobile"
	myweb "github.com/xxRanger/music-dat/avatarAndDat/controllers/server/web"
)

func InitRouter() {
	logs.Info("initial router")

	// set chain handler
	contractAddress, _:= web.AppConfig.String("contractAddress")
	port,_:=web.AppConfig.String("chainWS")
	address,_:= web.AppConfig.String("masterAddress")
	privateKey, _:= web.AppConfig.String("masterPrivateKey")

	chainHandler,err:= chainHelper.NewChainHandler(&chainHelper.ChainConfig{
		ContractAddress: contractAddress,
		Port: port,
		Account: chainHelper.AccountConfig{
			Address: address,
			PrivateKey: privateKey,
		},
	})
	if err!=nil {
		logs.Error(err.Error())
		panic(err)
	}

	// start transaction sending queue
	// TODO use message queue instead of go channel to ensure fault tolerance
	transactionQueueSzie,err:=web.AppConfig.Int("transactionQueueSize")
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
	chainHelper:= myweb.NewChainHelper()
	upLoadController:= &myweb.UploadController{}
	upLoadController.C = chainHelper
	upLoadController.TransactionQueue = transactionQueue
	nftListController:= &myweb.NftListController{}
	nftListController.C = chainHelper
	rewardController:= &myweb.RewardController{}
	rewardController.C = chainHelper
	rewardController.TransactionQueue = transactionQueue
	childrenOfNFTController := &myweb.ChildrenOfNFTController{}
	childrenOfNFTController.C = chainHelper
	numOfChildrenController:= &myweb.NumOfChildrenController{}
	numOfChildrenController.C = chainHelper
	marketTransactionHistoryController:= &myweb.MarketTransactionHistoryController{}
	marketTransactionHistoryController.C = chainHelper
	nicknameController:= &myweb.NicknameController{}
	introController:= &myweb.IntroController{}
	avatarController:= &myweb.AvatarController{}
	walletController:= &myweb.WalletController{}

	web.Router("/ws", wsHandler)
	web.Router("/admin",&myweb.AdminController{},"post:Login")
	web.Router("/file/:kind(avatar|dat|other)",upLoadController,"post:Upload")
	web.Router("/nftList/:kind(avatar|dat|other)/:uuid:string",nftListController)
	web.Router("/rewardDat/:uuid:string",rewardController,"get:RewardDat")
	web.Router("/nfts/:parentIndex:string/children", childrenOfNFTController)
	web.Router("/nfts/:parentIndex:string/balance", numOfChildrenController)
	web.Router("/market/transactionHistory/:uuid:string",marketTransactionHistoryController,"get:MarketTransactionHistory")
	web.Router("/profile/:uuid/nickname",nicknameController,"get:GetNickname;post:SetNickname")
	web.Router("/profile/:uuid/avatar",avatarController,"get:GetAvatar;post:SetAvatar")
	web.Router("/profile/:uuid/intro",introController,"get:GetIntro;post:SetIntro")
	web.Router("/profile/:uuid/wallet",walletController,"get:GetWallet;post:SetWallet")
}