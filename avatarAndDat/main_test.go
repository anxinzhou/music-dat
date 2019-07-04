package main

import (
	"database/sql"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego/session"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"github.com/xxRanger/music-dat/avatarAndDat/routers"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	logs.SetLogFuncCallDepth(3)
	//
	// initialize test database
	//
	// change to test db
	beego.AppConfig.Set("dbName","alphaslot_test")
	models.InitilizeModel()

	// drop database after test
	dbUser:= beego.AppConfig.String("dbUser")
	dbPassword:= beego.AppConfig.String("dbPassword")
	dbUrls:= beego.AppConfig.String("dbUrls")
	dbPort:=beego.AppConfig.String("dbPort")
	dbName:=beego.AppConfig.String("dbName")
	dbEngine:= beego.AppConfig.String("dbEngine")
	dbPath:= dbUser+":"+dbPassword+"@"+"tcp("+dbUrls+":"+dbPort+")"+"/"
	db,err:=sql.Open(dbEngine,dbPath)
	if err!=nil {
		panic(err)
	}
	_,err =db.Exec("drop database "+dbName)
	if err!=nil {
		panic(err)
	}

	// start server
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
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
	createDir()
	routers.InitRouter()
	beego.SetStaticPath("/resource","resource")
	go beego.Run()
	code:=m.Run()
	os.Exit(code)
}

func TestWebsiteApi(t *testing.T) {
	// necessary to wait for server starting
	<-time.After(1*time.Second)

	hostaddr:= beego.AppConfig.String("hostaddr")
	httpport,err:= beego.AppConfig.Int64("httpport")
	if err!=nil {
		t.Error("httpport should be int")
	}
	u:= url.URL{Scheme:"ws",Host:hostaddr+":"+strconv.FormatInt(httpport,10),Path:"/ws"}
	c,_,err:=websocket.DefaultDialer.Dial(u.String(),nil)
	defer c.Close()
	if err!=nil {
		t.Error("can not dail to ",u.String())
	}


	// start test
	testMobileUserUuid:= "4298349238490234456sa"
	testWebSiteUserUuid:= "4298349238490234456sa1"

	// test set nickname
	for {
		_, data, err:= c.ReadMessage()
		var kvs map[string] interface{}
		if err!=nil {
			logs.Error(err.Error())
			break;
		}
		json.Unmarshal(data, &kvs)
		action, ok := kvs["action"]
		if !ok {
			logs.Error("action not exist")
			continue
		}
		switch action:

	}
	logs.Info("begin database test")
}