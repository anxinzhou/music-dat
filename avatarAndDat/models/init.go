package models

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var O orm.Ormer   			//TODO how to efficiently gurantee concurrentcy of mysql
var MongoDB *mongo.Database   //TODO add logic to guranttee concurrency of mongodb

func init() {
	// initialize mysql handler
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
		new(MarketUserTable),
		new(NftShoppingCart),
		//new(CoinRecords),
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

	// initialize mongodb handler
	mongoDBURL:= beego.AppConfig.String("mongodbConnection")
	mongoDatabase:= beego.AppConfig.String("mongodbDatabase")
	client,err:=mongo.Connect(context.Background(), options.Client().ApplyURI(mongoDBURL))
	if err!=nil {
		logs.Error(err.Error())
		panic(err)
	}
	MongoDB=client.Database(mongoDatabase)
}