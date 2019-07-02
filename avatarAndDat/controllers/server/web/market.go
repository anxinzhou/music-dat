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
		Where("nft_market_info.seller_uuid = ?").OrderBy("timestamp").Desc()
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
	if num == 0 {
		avatarMKPlaceInfo= make([]nftTranData,0)
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
		common.OtherNftMarketInfo
		common.MarketPlaceInfo
	}
	var otherMKPlaceInfo []nftTranData
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
		Where("nft_market_info.seller_uuid = ?").OrderBy("timestamp").Desc()
	sql := qb.String()
	num,err:=o.Raw(sql,uuid).QueryRows(&otherMKPlaceInfo)
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
	if num == 0 {
		otherMKPlaceInfo= make([]nftTranData,0)
	}
	res:= response{
		NftTranData: otherMKPlaceInfo,
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *NftListController) GetDat(uuid string) {
	o:=orm.NewOrm()
	type nftTranData struct {
		common.DatNftMarketInfo
		common.MarketPlaceInfo
	}
	dbEngine := beego.AppConfig.String("dbEngine")
	qb,_:=orm.NewQueryBuilder(dbEngine)
	var datMKPlaceInfo []nftTranData
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
		Where("nft_market_info.seller_uuid = ?").OrderBy("timestamp").Desc()
	sql := qb.String()
	num,err:=o.Raw(sql,uuid).QueryRows(&datMKPlaceInfo)
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
	if num == 0 {
		datMKPlaceInfo= make([]nftTranData,0)
	}
	res:= response{
		NftTranData: datMKPlaceInfo,
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
	uuid:= this.Ctx.Input.Param(":uuid")
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

func (this *MarketTransactionHistoryController) MarketTransactionHistory() {
	uuid:= this.Ctx.Input.Param(":uuid")
	type purchaseNftInfo struct {
		BuyerNickname string  `json:"buyerNickname"`
		SellerNickname string `json:"sellerNickname" `
		TransactionAddress string `json:"transactionAddress" `
		NftLdefIndex string `json:"nftLdefIndex" orm:"nft_ldef_index"`
		Timestamp string `json:"timestamp"`
	}
	var nftTranData []*purchaseNftInfo
	o:=orm.NewOrm()
	num,err:=o.Raw(
		`
			select ni.transaction_address, ni.nft_ldef_index, ni.timestamp,
			buyer.nickname as buyer_nickname, seller.nickname as seller_nickname 
			from nft_purchase_info as ni
			inner join user_info as buyer on ni.uuid = buyer.uuid 
			inner join user_info as seller on ni.seller_uuid = seller.uuid
			where ni.uuid = ? or ni.seller_uuid = ?
			`,uuid,uuid).QueryRows(&nftTranData)
	if err!=nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:=errors.New("unexpected error when query db")
		sendError(&this.Controller, err, 500)
		return
	}
	type response struct {
		NftTranData []*purchaseNftInfo `json:"nftTranData"`
	}
	if num == 0 {
		nftTranData = make([]*purchaseNftInfo,num)
	}
	res:= response{
		NftTranData: nftTranData,
	}
	this.Data["json"] = &res
	this.ServeJSON()
}
