package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var O orm.Ormer

func init() {
	logs.Warn("initialize database")
	dbUser:= beego.AppConfig.String("dbUser")
	dbPassword:= beego.AppConfig.String("dbPassword")
	dbUrls:= beego.AppConfig.String("dbUrls")
	dbPort:=beego.AppConfig.String("dbPort")
	dbName:=beego.AppConfig.String("dbName")
	dbEngine:= beego.AppConfig.String("dbEngine")

	dataSource:=dbUser+":"+dbPassword+"@"+"tcp("+dbUrls+":"+dbPort+")"+"/"+dbName

	orm.RegisterDriver("mysql",orm.DRMySQL)
	err:=orm.RegisterDataBase("default",dbEngine,dataSource)
	if err!=nil {
		panic(err)
	}

	orm.RegisterModel(
		new(BerryPurchaseTable),
		new(NftItemAdmin),
		new(NftMappingTable),
		new(StorePurchaseHistroy),
		new(NftMarketTable),
		new(NftInfoTable),
	)

	// auto generate table
	verbose:=true
	force:=false
	err= orm.RunSyncdb("default",force,verbose)
	if err!=nil {
		panic(err)
	}

	// set oramer object
	O = orm.NewOrm()
	O.Using("default")
}