package web

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/util"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
)

type NftListController struct {
	ContractController
}

func (this *NftListController) GetAvatar(uuid string) {
	o:=orm.NewOrm()
	dbEngine := beego.AppConfig.String("dbEngine")
	qb,_:=orm.NewQueryBuilder(dbEngine)
	type nftTranData struct {
		common.AvatarNftMarketInfo
		common.MarketPlaceInfo
	}
	var avatarMKPlaceInfo []nftTranData
	qb.Select("*").
		From("nft_market_place").
		InnerJoin("nft_market_info").
		On("nft_market_place.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("avatar_nft_market_info").
		On("nft_market_place.nft_ldef_index = avatar_nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_market_place.nft_ldef_index = nft_info.nft_ldef_index").
		InnerJoin("avatar_nft_info").
		On("nft_market_place.nft_ldef_index = avatar_nft_info.nft_ldef_index").
		Where("nft_market_info.seller_uuid = ?")
	sql := qb.String()
	num,err:=o.Raw(sql,uuid).QueryRows(&avatarMKPlaceInfo)
	if err != nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:= errors.New("unknown error when query database")
		sendError(&this.Controller, err, 500)
		return
	}
	logs.Debug("get",num,"from database")
	type response struct {
		NftTranData []nftTranData `json:"nftTranData"`
	}

	res:= response{
		NftTranData: avatarMKPlaceInfo,
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *NftListController) GetOther(uuid string) {
	o:=orm.NewOrm()
	dbEngine := beego.AppConfig.String("dbEngine")
	qb,_:=orm.NewQueryBuilder(dbEngine)
	type nftTranData struct {
		common.OtherNftInfo
		common.MarketPlaceInfo
	}
	var avatarMKPlaceInfo []nftTranData
	qb.Select("*").
		From("nft_market_place").
		InnerJoin("nft_market_info").
		On("nft_market_place.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("other_nft_market_info").
		On("nft_market_place.nft_ldef_index = other_nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_market_place.nft_ldef_index = nft_info.nft_ldef_index").
		InnerJoin("other_nft_info").
		On("nft_market_place.nft_ldef_index = other_nft_info.nft_ldef_index").
		Where("nft_market_info.seller_uuid = ?")
	sql := qb.String()
	num,err:=o.Raw(sql,uuid).QueryRows(&avatarMKPlaceInfo)
	if err != nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:= errors.New("unknown error when query database")
		sendError(&this.Controller, err, 500)
		return
	}
	logs.Debug("get",num,"from database")
	type response struct {
		NftTranData []nftTranData `json:"nftTranData"`
	}

	res:= response{
		NftTranData: avatarMKPlaceInfo,
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *NftListController) GetDat(uuid string) {
	o:=orm.NewOrm()
	type nftTranData struct {
		common.DatNftInfo
		common.MarketPlaceInfo
	}
	dbEngine := beego.AppConfig.String("dbEngine")
	qb,_:=orm.NewQueryBuilder(dbEngine)
	var avatarMKPlaceInfo []nftTranData
	qb.Select("*").
		From("nft_market_place").
		InnerJoin("nft_market_info").
		On("nft_market_place.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("dat_nft_market_info").
		On("nft_market_place.nft_ldef_index = dat_nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_market_place.nft_ldef_index = nft_info.nft_ldef_index").
		InnerJoin("dat_nft_info").
		On("nft_market_place.nft_ldef_index = dat_nft_info.nft_ldef_index").
		Where("nft_market_info.seller_uuid = ?")
	sql := qb.String()
	num,err:=o.Raw(sql,uuid).QueryRows(&avatarMKPlaceInfo)
	if err != nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:= errors.New("unknown error when query database")
		sendError(&this.Controller, err, 500)
		return
	}
	logs.Debug("get",num,"from database")
	type response struct {
		NftTranData []nftTranData `json:"nftTranData"`
	}

	res:= response{
		NftTranData: avatarMKPlaceInfo,
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *NftListController) Get() {
	kind:= this.Ctx.Input.Param(":kind")
	if err:= util.ValidNftName(kind); err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller, err, 400)
		return
	}
	uuid:= this.Ctx.Input.Param("uuid")
	userMarketInfo := models.UserMarketInfo{
		Uuid: uuid,
	}

	// check if user exist
	o := orm.NewOrm()
	err := o.Read(&userMarketInfo)
	if err != nil {
		if err == orm.ErrNoRows {
			err := errors.New("no such user in marketplace")
			logs.Error(err.Error())
			sendError(&this.Controller, err, 400)
			return
		} else {
			logs.Error(err.Error())
			err := errors.New("unknown error when query db")
			sendError(&this.Controller, err, 500)
			return
		}
	}

	switch kind {
	case common.NAME_NFT_AVATAR:
		this.GetAvatar(uuid)
	case common.NAME_NFT_OTHER:
		this.GetOther(uuid)
	case common.NAME_NFT_MUSIC:
		this.GetDat(uuid)
	}
}

type MarketTransactionHistoryController struct {
	ContractController
}

type NftPurchaseInfo struct {
	NftLdefIndex string `json:"nftLdefIndex"`
	Buyer string `json:"buyer"`
	Seller string `json:"seller"`
	TransactionAddress string `json:"transactionAddress"`
	Time string `json:"time"`
}

type MarketHistoryResponse struct {
	NftPurchaseInfo []*NftPurchaseInfo `json:"nftPurchaseInfo"`
}

func (this *MarketTransactionHistoryController) MarketTransactionHistory() {
	uuid:= this.Ctx.Input.Param(":uuid")
	type purchaseNftInfo struct {
		Uuid string `json:"uuid"`
		SellerUuid string `json:"seller_uuid"`
		TransactionAddress string `json:"transactionAddress"`
		NftLdefIndex string `json:"nftLdefIndex"`
		Timestamp string `json:"timestamp"`
	}
	var nftTranData []purchaseNftInfo
	o:=orm.NewOrm()
	num,err:=o.QueryTable("nft_purchase_info").
		Filter("uuid",uuid).
		All(&nftTranData,"transaction_address","nft_ldef_index","uuid","seller_uuid","timestamp")
	if err!=nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:=errors.New("unexpected error when query db")
		sendError(&this.Controller, err, 500)
		return
	}
	type response struct {
		NftTranData []purchaseNftInfo `json:"nftTranData"`
	}
	if num == 0 {
		nftTranData = make([]purchaseNftInfo,num)
	}
	res:= response{
		NftTranData: nftTranData,
	}
	this.Data["json"] = &res
	this.ServeJSON()
}
