package web

import (
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/transactionQueue"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/util"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"time"
)

type NftListController struct {
	ContractController
}

func (this *NftListController) GetAvatar(uuid string) {
	o := orm.NewOrm()
	dbEngine, _ := web.AppConfig.String("dbEngine")
	qb, _ := orm.NewQueryBuilder(dbEngine)
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
	num, err := o.Raw(sql, uuid).QueryRows(&avatarMKPlaceInfo)
	if err != nil && err != orm.ErrNoRows {
		logs.Error(err.Error())
		err := errors.New("unknown error when query database")
		sendError(&this.Controller, err, 500)
		return
	}
	logs.Debug("get", num, "from database")
	type response struct {
		NftTranData []nftTranData `json:"nftTranData"`
	}
	if num == 0 {
		avatarMKPlaceInfo = make([]nftTranData, 0)
	}
	res := response{
		NftTranData: avatarMKPlaceInfo,
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *NftListController) GetOther(uuid string) {
	o := orm.NewOrm()
	dbEngine, _ := web.AppConfig.String("dbEngine")
	qb, _ := orm.NewQueryBuilder(dbEngine)
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
	num, err := o.Raw(sql, uuid).QueryRows(&otherMKPlaceInfo)
	if err != nil && err != orm.ErrNoRows {
		logs.Error(err.Error())
		err := errors.New("unknown error when query database")
		sendError(&this.Controller, err, 500)
		return
	}
	logs.Debug("get", num, "from database")
	type response struct {
		NftTranData []nftTranData `json:"nftTranData"`
	}
	if num == 0 {
		otherMKPlaceInfo = make([]nftTranData, 0)
	}
	res := response{
		NftTranData: otherMKPlaceInfo,
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *NftListController) GetDat(uuid string) {
	o := orm.NewOrm()
	type nftTranData struct {
		common.DatNftMarketInfo
		common.MarketPlaceInfo
	}
	dbEngine, _ := web.AppConfig.String("dbEngine")
	qb, _ := orm.NewQueryBuilder(dbEngine)
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
	num, err := o.Raw(sql, uuid).QueryRows(&datMKPlaceInfo)
	if err != nil && err != orm.ErrNoRows {
		logs.Error(err.Error())
		err := errors.New("unknown error when query database")
		sendError(&this.Controller, err, 500)
		return
	}
	logs.Debug("get", num, "from database")
	type response struct {
		NftTranData []nftTranData `json:"nftTranData"`
	}
	if num == 0 {
		datMKPlaceInfo = make([]nftTranData, 0)
	}
	res := response{
		NftTranData: datMKPlaceInfo,
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *NftListController) Get() {
	kind := this.Ctx.Input.Param(":kind")
	if err := util.ValidNftName(kind); err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller, err, 400)
		return
	}
	uuid := this.Ctx.Input.Param(":uuid")
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
	uuid := this.Ctx.Input.Param(":uuid")
	type purchaseNftInfo struct {
		BuyerNickname      string `json:"buyerNickname"`
		SellerNickname     string `json:"sellerNickname" `
		TransactionAddress string `json:"transactionAddress" `
		NftLdefIndex       string `json:"nftLdefIndex" orm:"nft_ldef_index"`
		Timestamp          string `json:"timestamp"`
	}
	var nftTranData []*purchaseNftInfo
	o := orm.NewOrm()
	num, err := o.Raw(
		`
			select ni.transaction_address, ni.nft_ldef_index, ni.timestamp,
			buyer.nickname as buyer_nickname, seller.nickname as seller_nickname 
			from nft_purchase_info as ni
			inner join user_info as buyer on ni.uuid = buyer.uuid 
			inner join user_info as seller on ni.seller_uuid = seller.uuid
			where ni.uuid = ? or ni.seller_uuid = ? order by ni.timestamp desc
			`, uuid, uuid).QueryRows(&nftTranData)
	if err != nil && err != orm.ErrNoRows {
		logs.Error(err.Error())
		err := errors.New("unexpected error when query db")
		sendError(&this.Controller, err, 500)
		return
	}
	type response struct {
		NftTranData []*purchaseNftInfo `json:"nftTranData"`
	}
	if num == 0 {
		nftTranData = make([]*purchaseNftInfo, num)
	}
	res := response{
		NftTranData: nftTranData,
	}
	this.Data["json"] = &res
	this.ServeJSON()
}

type RewardController struct {
	ContractController
	TransactionQueue *transactionQueue.TransactionQueue
}

func (this *RewardController) RewardDat() {
	// only reward one dat now
	uuid := this.Ctx.Input.Param(":uuid")
	//
	// get nft info from database
	//
	o := orm.NewOrm()
	dbEngine, _ := web.AppConfig.String("dbEngine")
	qb, _ := orm.NewQueryBuilder(dbEngine)
	type nftRewardInfo struct {
		NftLdefIndex  string `json:"nftLdefIndex"`
		SupportedType string `json:"supportedType";orm:"nft_type"`
		NftName       string `json:"nftName"`
		Thumbnail     string `json:"thumbnail";orm:"music_file_name"`
	}
	type response struct {
		NftTranData []nftRewardInfo `json:"nftTranData"`
	}

	type queryInfo struct {
		NftLdefIndex string
		NftType      string
		NftName      string
		IconFileName string
		SellerUuid   string
		SellerWallet string
		ActiveTicker string
		Price        int
		Timestamp    time.Time
	}
	var nftMarketInfo queryInfo
	qb.Select("nft_info.nft_ldef_index",
		"nft_info.nft_type",
		"nft_info.nft_name",
		"dat_nft_info.music_file_name",
		"nft_market_info.seller_uuid",
		"nft_market_info.seller_wallet",
		"nft_market_info.price",
		"nft_market_place.active_ticker",
		"nft_market_place.timestamp").
		From("nft_market_place").
		InnerJoin("nft_market_info").
		On("nft_market_place.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("dat_nft_market_info").
		On("nft_market_place.nft_ldef_index = dat_nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_market_place.nft_ldef_index = nft_info.nft_ldef_index").
		InnerJoin("dat_nft_info").
		On("nft_market_place.nft_ldef_index = dat_nft_info.nft_ldef_index").
		Where("dat_nft_market_info.allow_airdrop = true").
		Limit(1)
	sql := qb.String()
	err := o.Raw(sql).QueryRow(&nftMarketInfo)
	if err != nil {
		if err == orm.ErrNoRows {
			logs.Debug("no dat in marketplace now")
			this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
			res := &response{
				NftTranData: make([]nftRewardInfo, 0),
			}
			this.Data["json"] = res
			this.ServeJSON()
			return
		} else {
			logs.Error(err.Error())
			err := errors.New("unexpected error when query database")
			sendError(&this.Controller, err, 500)
			return
		}
	}

	to, err := o.Begin()
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("start the transaction failed")
		sendError(&this.Controller, err, 500)
		return
	}

	//
	// change nft count in database
	//

	buyerMarketInfo := models.UserMarketInfo{
		Uuid: uuid,
	}
	err = o.ReadForUpdate(&buyerMarketInfo)
	if err != nil {
		to.Rollback()
		if err == orm.ErrNoRows {
			err := errors.New("User " + uuid + " has not binded wallet")
			logs.Error(err.Error())
			sendError(&this.Controller, err, 500)
			return
		} else {
			logs.Error(err.Error())
			err := errors.New("unexpected error when query database")
			sendError(&this.Controller, err, 500)
			return
		}
	}
	sellerMarketInfo := models.UserMarketInfo{
		Uuid: nftMarketInfo.SellerUuid,
	}
	//
	err = o.ReadForUpdate(&sellerMarketInfo)
	if err != nil {
		to.Rollback()
		logs.Error(err.Error())
		err := errors.New("unexpected error when query database")
		sendError(&this.Controller, err, 500)
		return
	}
	buyerMarketInfo.Count += 1
	sellerMarketInfo.Count -= 1
	_, err = o.Update(&buyerMarketInfo, "count")
	if err != nil {
		to.Rollback()
		logs.Error(err.Error())
		err := errors.New("unexpected error when query database")
		sendError(&this.Controller, err, 500)
		return
	}
	o.Update(&sellerMarketInfo, "count")
	if err != nil {
		to.Rollback()
		logs.Error(err.Error())
		err := errors.New("unexpected error when query database")
		sendError(&this.Controller, err, 500)
		return
	}

	nftMarketPlaceInfo := models.NftMarketPlace{
		NftLdefIndex: nftMarketInfo.NftLdefIndex,
	}
	err = o.Read(&nftMarketPlaceInfo)
	if err != nil {
		to.Rollback()
		logs.Error(err.Error())
		err := errors.New("unexpected error when query database")
		sendError(&this.Controller, err, 500)
		return
	}

	nftInfo := models.NftMarketInfo{
		NftLdefIndex: nftMarketInfo.NftLdefIndex,
	}
	err = o.Read(&nftInfo)
	if err != nil {
		to.Rollback()
		logs.Error(err.Error())
		err := errors.New("unexpected error when query database")
		sendError(&this.Controller, err, 500)
		return
	}
	err = o.Read(nftInfo.NftInfo)
	if err != nil {
		to.Rollback()
		logs.Error(err.Error())
		err := errors.New("unexpected error when query database")
		sendError(&this.Controller, err, 500)
		return
	}
	// add numsold
	nftInfo.NumSold += 1
	_, err = o.Update(&nftInfo, "num_sold")
	if err != nil {
		to.Rollback()
		logs.Error(err.Error())
		err := errors.New("unexpected error when query database")
		sendError(&this.Controller, err, 500)
		return
	}
	// delete from market
	if nftInfo.NumSold == nftInfo.Qty {
		_, err := o.Delete(&nftMarketPlaceInfo)
		if err != nil {
			to.Rollback()
			logs.Error(err.Error())
			err := errors.New("unexpected error when query database")
			sendError(&this.Controller, err, 500)
			return
		}
	} else {

	}
	//
	// insert into purchase history
	//
	distributionLdefIndex := util.RandomNftLdefIndex(common.TYPE_NFT_MUSIC)
	purchaseId := util.RandomPurchaseId()
	nftPuchaseInfo := models.NftPurchaseInfo{
		TotalPaid:          nftMarketInfo.Price,
		PurchaseId:         purchaseId,
		Uuid:               uuid,
		SellerUuid:         nftMarketInfo.SellerUuid,
		TransactionAddress: "", // determined after send transaction
		ActiveTicker:       nftMarketInfo.ActiveTicker,
		DistributionIndex:  distributionLdefIndex,
		NftLdefIndex:       nftMarketInfo.NftLdefIndex,
		Status:             common.PURCHASE_PENDING, // change to finish after send transaction
		UserInfo: &models.UserInfo{
			Uuid: uuid,
		},
	}
	_, err = o.Insert(&nftPuchaseInfo)
	if err != nil {
		to.Rollback()
		logs.Error(err.Error())
		err := errors.New("unexpected error when query databas")
		sendError(&this.Controller, err, 500)
		return
	}
	rewardNftInfos := make([]nftRewardInfo, 1)
	rewardNftInfos[0] = nftRewardInfo{
		NftLdefIndex:  nftMarketInfo.NftLdefIndex,
		SupportedType: nftMarketInfo.NftType,
		NftName:       nftMarketInfo.NftName,
		Thumbnail:     util.PathPrefixOfNFT(common.TYPE_NFT_MUSIC, common.PATH_KIND_MARKET) + nftMarketInfo.IconFileName,
	}
	res := &response{
		NftTranData: rewardNftInfos,
	}

	to.Commit()
	this.Data["json"] = res
	this.ServeJSON()

	// todo use message queue instead go channel
	this.TransactionQueue.Append(&transactionQueue.RewardNftTransaction{
		Uuid:         uuid,
		SellerUuid:   nftMarketInfo.SellerUuid,
		NftLdefIndex: nftMarketInfo.NftLdefIndex,
		DistIndex:    distributionLdefIndex,
		PurchaseId:   purchaseId,
	})
}
