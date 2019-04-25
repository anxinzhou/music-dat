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
	chainController:= &http.ChainBalanceController{}
	chainController.C = chainHelper

	beego.Router("/", &http.MainController{})
	beego.Router("/ws", wsHandler)
	beego.Router("/admin",&http.AdminController{},"get:Get;post:Login")
	beego.Router("/file/:kind(avatar|dat)",upLoadController,"get:Get;post:Upload")
	beego.Router("/balance/:kind(avatar|dat)",chainController)
}


