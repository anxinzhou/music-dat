package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/mobile"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/web"
)

func init() {
	logs.Info("initial router")

	m:= mobile.NewManager()
	m.Init()
	// set chain handler
	chainHandler,err:= mobile.NewChainHandler(&mobile.ChainConfig{
		ContractAddress: beego.AppConfig.String("contractAddress"),
		Port: beego.AppConfig.String("chainWS"),
		Account: mobile.AccountConfig{
			Address: beego.AppConfig.String("masterAddress"),
			PrivateKey:beego.AppConfig.String("masterPrivateKey"),
		},
	})
	if err!=nil {
		logs.Error(err.Error())
		panic(err)
	}
	m.SetChainHandler(chainHandler)
	wsHandler:= &mobile.WebSocketHandler{
		M: m,
	}
	chainHelper:= web.NewChainHelper()
	upLoadController:= &web.UploadController{}
	upLoadController.C = chainHelper
	nftListController:= &web.NftListController{}
	nftListController.C = chainHelper
	rewardController:= &web.RewardController{}
	rewardController.C = chainHelper
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
	beego.Router("/nftList/:uuid:string",nftListController)
	beego.Router("/rewardDat/:uuid:string",rewardController,"get:RewardDat")
	beego.Router("/nfts/:parentIndex:string/children", childrenOfNFTController)
	beego.Router("/nfts/:parentIndex:string/balance", numOfChildrenController)
	//beego.Router("/wallet",&http.ImportWalletController{},"post:ImportWallet")
	beego.Router("/market/transactionHistory/:uuid:string",marketTransactionHistoryController,"get:MarketTransactionHistory")
	beego.Router("/profile/:uuid/nickname",nicknameController,"get:GetNickname;post:SetNickname")
	beego.Router("/profile/:uuid/avatar",avatarController,"get:GetAvatar;post:SetAvatar")
	beego.Router("/profile/:uuid/intro",introController,"get:GetIntro;post:SetIntro")
	beego.Router("/profile/:uuid/wallet",walletController,"get:GetWallet;post:SetWallet")
}


