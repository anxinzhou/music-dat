package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego/session"
	_ "github.com/xxRanger/music-dat/avatarAndDat/routers"
)

func init() {
	logs.SetLogFuncCallDepth(3)

}

func main() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:8080","http://52.220.30.192:*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "content-type", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
	sessionconf := &session.ManagerConfig{
		CookieName: "begoosessionID",
		Gclifetime: 3600,
	}
	beego.GlobalSessions, _ = session.NewManager("memory", sessionconf)
	go beego.GlobalSessions.GC()
	beego.SetStaticPath("/resource","resource")
	beego.Run()
}
