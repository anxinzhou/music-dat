package models

import (
	"context"
	"database/sql"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	_  "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database //TODO add logic to guranttee concurrency of mongodb
var MongoClient *mongo.Client

func InitilizeModel(force bool, verbose bool) {
	// initialize mysql handler
	logs.Warn("initialize database")
	dbUser, _ := web.AppConfig.String("dbUser")
	dbPassword, _ := web.AppConfig.String("dbPassword")
	dbUrls, _ := web.AppConfig.String("dbUrls")
	dbPort, _ := web.AppConfig.String("dbPort")
	dbName, _ := web.AppConfig.String("dbName")
	dbEngine, _ := web.AppConfig.String("dbEngine")
	dbPath := dbUser + ":" + dbPassword + "@" + "tcp(" + dbUrls + ":" + dbPort + ")" + "/"
	dataSource := dbPath + dbName + "?charset=utf8"

	// create db if not exist
	db, err := sql.Open(dbEngine, dbPath)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("create database if not exists " + dbName)
	if err != nil {
		panic(err)
	}

	orm.RegisterDriver("mysql", orm.DRMySQL)
	err = orm.RegisterDataBase("default", dbEngine, dataSource)
	orm.RegisterModel(
		new(UserInfo),
		new(CreatorInfo),
		new(NftInfo),
		new(NftPurchaseInfo),
		new(AvatarNftInfo),
		new(DatNftInfo),
		new(OtherNftInfo),
		new(NftMarketInfo),
		new(DatNftMarketInfo),
		new(AvatarNftMarketInfo),
		new(OtherNftMarketInfo),
		new(NftMarketPlace),
		new(NftShoppingCart),
		new(UserMarketInfo),
		new(FollowTable),
		new(BerryPurchaseInfo),
	)

	// auto generate table
	err = orm.RunSyncdb("default", force, verbose)
	if err != nil {
		panic(err)
	}
	// set connection pool
	orm.SetMaxOpenConns("default", 2000)
	orm.SetMaxIdleConns("default", 2000)

	// initialize mongodb handler
	mongoDBURL, _ := web.AppConfig.String("mongodbConnection")
	mongoDatabase, _ := web.AppConfig.String("mongodbDatabase")
	MongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoDBURL))
	if err != nil {
		logs.Error(err.Error())
		panic(err)
	}
	MongoDB = MongoClient.Database(mongoDatabase)

	// set test creator
	GenerateTestCreator(4)
}
