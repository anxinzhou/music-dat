package mobile

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/util"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func (m *Manager) BindWalletHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Uuid     string `json:"uuid"`
		WalletId string `json:"wallet_id"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}

	userMarketInfo := &models.UserMarketInfo{
		Uuid:   req.Uuid,
		Wallet: req.WalletId,
		Count:  0,
		UserInfo: &models.UserInfo{
			Uuid: req.Uuid,
		},
	}
	o := orm.NewOrm()
	_, err = o.InsertOrUpdate(&userMarketInfo, "wallet")
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("unexpected error when query db")
		m.errorHandler(c, action, err)
		return
	}

	type response struct {
		Status   int    `json:"status"`
		Action   string `json:"action"`
		WalletId string `json:"wallet_id"`
	}
	m.wrapperAndSend(c, action, &response{
		Status:   common.RESPONSE_STATUS_SUCCESS,
		Action:   action,
		WalletId: req.WalletId,
	})
}

func (m *Manager) SetNicknameHandler(c *client.Client, action string, data []byte) {
	// TODO set nickname in mongodb
	type request struct {
		Uuid     string `json:"uuid"`
		Nickname string `json:"nickname"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}

	userInfo := &models.UserInfo{
		Uuid:     req.Uuid,
		Nickname: req.Nickname,
	}
	o := orm.NewOrm()
	o.Begin()
	_, err = o.InsertOrUpdate(&userInfo, "nickname")
	if err != nil {
		o.Rollback()
		if strings.Contains(err.Error(), common.DUPLICATE_ENTRY) {
			logs.Error(err.Error())
			err := errors.New("duplicate nickname")
			m.errorHandler(c, action, err)
			return
		} else {
			logs.Error(err.Error())
			err := errors.New("unexpected error when query db")
			m.errorHandler(c, action, err)
			return
		}
	}
	// update mongodb
	col := models.MongoDB.Collection("users")
	filter:= bson.M {
		"uuid": req.Uuid,
	}
	update:= bson.M {
		"$set": bson.M {"nickname":req.Nickname},
	}
	_,err=col.UpdateOne(context.Background(),filter,update)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		err:=errors.New("can not set nickname in mongodb")
		m.errorHandler(c, action, err)
		return
	}
	o.Commit()
	type response struct {
		Status   int    `json:"status"`
		Action   string `json:"action"`
		Nickname string `json:"nickname"`
	}
	m.wrapperAndSend(c, action, &response{
		Status:   common.RESPONSE_STATUS_SUCCESS,
		Action:   action,
		Nickname: req.Nickname,
	})
}

func (m *Manager) IsNicknameDuplicatedHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Nickname string `json:"nickname"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}

	userInfo := &models.UserInfo{
		Nickname: req.Nickname,
	}
	o := orm.NewOrm()
	err = o.Read(&userInfo, "nickname")
	duplicated := true
	if err != nil {
		if err == orm.ErrNoRows {
			duplicated = false
		} else {
			logs.Error(err.Error())
			err := errors.New("unexpected error when query db")
			m.errorHandler(c, action, err)
			return
		}
	}
	type response struct {
		Status     int    `json:"status"`
		Action     string `json:"action"`
		Duplicated bool   `json:"duplicated"`
		Nickname   string `json:"nickname"`
	}
	m.wrapperAndSend(c, action, &response{
		Status:     common.RESPONSE_STATUS_SUCCESS,
		Action:     action,
		Duplicated: duplicated,
		Nickname:   req.Nickname,
	})
}

func (m *Manager) FollowListHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Uuid string `json:"uuid"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}
	type followeeInfo struct {
		Uuid      string `json:"uuid" orm:"column(followee_uuid)"`
		Nickname  string `json:"nickname"`
		Thumbnail string `json:"thumnail" orm:"column(avatar_file_name)"`
		Intro     string `json:"intro"`
	}
	o := orm.NewOrm()
	dbEngine := beego.AppConfig.String("dbEngine")
	qb, _ := orm.NewQueryBuilder(dbEngine)

	var followInfo []followeeInfo
	qb.Select(
		"nft_info.nft_ldef_index",
		"nft_info.nft_name",
		"nft_info.short_desc",
		"nft_info.long_desc",
		"avatar_nft_info.nft_life_index",
		"avatar_nft_info.nft_power_index",
		"nft_market_info.price",
		"nft_market_info.qty",
		"nft_info.file_name",
	).
		From("follow_table").
		InnerJoin("user_info").
		On("follow_table.followee_uuid = user_info.uuid").
		Where("follow_table.follower_uuid = ?").OrderBy("follow_table.timestamp").Desc()
	sql := qb.String()
	num, err := o.Raw(sql, req.Uuid).QueryRows(&followInfo)
	if err != nil && err != orm.ErrNoRows {
		logs.Error(err.Error())
		err := errors.New("unexpected error when query database")
		m.errorHandler(c, action, err)
		return
	}
	if num == 0 {
		followInfo = make([]followeeInfo, 0)
	}

	for i, _ := range followInfo {
		followInfo[i].Thumbnail = util.PathPrefixOfNFT("", common.PATH_KIND_USER_ICON) + followInfo[i].Thumbnail
	}

	type response struct {
		Status     int            `json:"status"`
		Action     string         `json:"action"`
		FollowList []followeeInfo `json:"followList"`
	}
	m.wrapperAndSend(c, action, &response{
		Status:     common.RESPONSE_STATUS_SUCCESS,
		Action:     action,
		FollowList: followInfo,
	})
}

func (m *Manager) FollowListOperationHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Uuid         string `json:"uuid"`
		FolloweeUuid string `json:"followee_uuid"`
		Operation    int    `json:"operation"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}
	if err:= util.ValidFollowListOperation(req.Operation);err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, action, err)
		return
	}
	followInfo:= models.FollowTable{
		FollowerUuid: req.Uuid,
		FolloweeUuid: req.FolloweeUuid,
	}
	switch req.Operation {
	case common.FOLLOW_LIST_ADD:
		o:=orm.NewOrm()
		_,err:=o.Insert(&followInfo)
		if err!=nil {
			if !strings.Contains(err.Error(),common.DUPLICATE_ENTRY) {
				logs.Error(err.Error())
				err:= errors.New("unexpected error when query database")
				m.errorHandler(c, action, err)
				return
			}
		}
	case common.FOLLOW_LIST_DELETE:
		o:=orm.NewOrm()
		_,err:=o.Delete(&followInfo,req.Uuid,req.FolloweeUuid)
		if err!=nil {
			logs.Error(err.Error())
			err:= errors.New("unexpected error when query database")
			m.errorHandler(c, action, err)
			return
		}
	default:
		panic("unexpected follow list operation")
	}
	type response struct {
		Status int `json:"status"`
		Action string `json:"action"`
		FollowUuid string `json:"followUuid"`
		Operation int `json:"operation"`
	}
	m.wrapperAndSend(c,action,&response{
		Status: common.RESPONSE_STATUS_SUCCESS,
		Action: action,
		FollowUuid: req.FolloweeUuid,
		Operation: req.Operation,
	})
}

func (m *Manager) IsNicknameSetHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Uuid         string `json:"uuid"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}

	set := true
	o := orm.NewOrm()
	userInfo:= models.UserInfo{
		Uuid: req.Uuid,
	}
	err = o.Read(&userInfo)
	if err != nil {
		if err == orm.ErrNoRows {
			set = false
		} else {
			logs.Error(err.Error())
			err:= errors.New("unexpected error when query database")
			m.errorHandler(c, action, err)
			return
		}
	}

	type response struct {
		Status int `json:"status"`
		Action string `json:"action"`
		Set bool `json:"set"`
	}
	m.wrapperAndSend(c, action, &response{
		Status: common.RESPONSE_STATUS_SUCCESS,
		Action: action,
		Set:set,
	})
}

func (m *Manager) MarketUserListHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Uuid string `json:"uuid"`
	}
	var req request
	type userInfo struct {
		Uuid      string `json:"uuid"`
		Nickname  string `json:"nickname"`
		Count     int    `json:"count"`
		Thumbnail string `json:"thumnail" orm:"column(avatar_file_name)"`
		Followed  bool   `json:"followed" orm:"column(followed)"`
	}
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}
	o := orm.NewOrm()
	var ui []userInfo
	num, err := o.Raw(`
		select ui.uuid,ui.nickname,umi.count,ui.avatar_file_name,count(ft.followee_uuid) as followed from 
		user_market_info as umi inner join user_info as ui on umi.uuid = ui.uuid
		left join (
			select followee_uuid, follower_uuid from follow_table
			where follow_table.follower_uuid =?) 
			as ft on umi.uuid = ft.followee_uuid 
			`, req.Uuid).QueryRows(&ui)
	if err != nil && err != orm.ErrNoRows {
		logs.Error(err.Error())
		err := errors.New("unexpected error when query db")
		m.errorHandler(c, action, err)
		return
	}
	if num == 0 {
		ui = make([]userInfo, 0)
	}
	for i, _ := range ui {
		ui[i].Thumbnail = util.PathPrefixOfNFT("", common.PATH_KIND_USER_ICON) + ui[i].Thumbnail
	}
	type response struct {
		Status       int    `json:"status"`
		Action       string `json:"action"`
		WalletIdList []userInfo
	}
	m.wrapperAndSend(c, action, &response{
		Status:       common.RESPONSE_STATUS_SUCCESS,
		Action:       action,
		WalletIdList: ui,
	})
}

func (m *Manager) AvatarPurchaseHistory(c *client.Client, action string, uuid string) {
	o := orm.NewOrm()
	dbEngine := beego.AppConfig.String("dbEngine")
	qb, _ := orm.NewQueryBuilder(dbEngine)
	type nftTranData struct {
		NftLdefIndex       string `json:"nftLdefIndex"`
		SupportedType      string `json:"supportedType" orm:"column(nft_type)"`
		NftName            string `json:"nftName"`
		ShortDesc          string `json:"shortDesc"`
		LongDesc           string `json:"longDesc"`
		NftLifeIndex       int    `json:"nftLifeIndex"`
		NftPowerIndex      int    `json:"nftPowerIndex"`
		NftValue           int    `json:"nftValue" orm:"column(price)"`
		Qty                int    `json:"qty"`
		Thumbnail          string `json:"thumnail" orm:"column(file_name)"`
		TransactionAddress string `json:"transactionAddress"`
	}
	var avatarPurchaseInfo []nftTranData
	qb.Select(
		"nft_info.nft_ldef_index",
		"nft_info.nft_type",
		"nft_info.nft_name",
		"nft_info.short_desc",
		"nft_info.long_desc",
		"avatar_nft_info.nft_life_index",
		"avatar_nft_info.nft_power_index",
		"nft_market_info.price",
		"nft_market_info.qty",
		"nft_info.file_name",
		"nft_purchase_info.transaction_address",
	).
		From("nft_purchase_info").
		InnerJoin("nft_market_info").
		On("nft_purchase_info.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("avatar_nft_market_info").
		On("nft_purchase_info.nft_ldef_index = avatar_nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_purchase_info.nft_ldef_index = nft_info.nft_ldef_index").
		InnerJoin("avatar_nft_info").
		On("nft_purchase_info.nft_ldef_index = avatar_nft_info.nft_ldef_index").
		Where("nft_purchase_info.uuid = ?").OrderBy("timestamp").Desc()
	sql := qb.String()
	num, err := o.Raw(sql, uuid).QueryRows(&avatarPurchaseInfo)
	if err != nil && err != orm.ErrNoRows {
		logs.Error(err.Error())
		err := errors.New("unknown error when query database")
		m.errorHandler(c, action, err)
		return
	}
	logs.Debug("get", num, "from database")
	type response struct {
		Status        int           `json:"status"`
		Action        string        `json:"action"`
		SupportedType string        `json:"supportedType"`
		PurchaseList  []nftTranData `json:"purchaseList"`
	}
	if num == 0 {
		avatarPurchaseInfo = make([]nftTranData, 0)
	}
	for i, _ := range avatarPurchaseInfo {
		avatarPurchaseInfo[i].Thumbnail = util.PathPrefixOfNFT("", common.PATH_KIND_USER_ICON) + avatarPurchaseInfo[i].Thumbnail
	}
	m.wrapperAndSend(c, action, &response{
		Status:        common.RESPONSE_STATUS_SUCCESS,
		Action:        action,
		SupportedType: common.TYPE_NFT_AVATAR,
		PurchaseList:  avatarPurchaseInfo,
	})
}

func (m *Manager) DatPurchaseHistory(c *client.Client, action string, uuid string) {
	o := orm.NewOrm()
	dbEngine := beego.AppConfig.String("dbEngine")
	qb, _ := orm.NewQueryBuilder(dbEngine)
	type nftTranData struct {
		NftLdefIndex       string `json:"nftLdefIndex"`
		SupportedType      string `json:"supportedType" orm:"column(nft_type)"`
		NftName            string `json:"nftName"`
		ShortDesc          string `json:"shortDesc"`
		LongDesc           string `json:"longDesc"`
		NftLifeIndex       int    `json:"nftLifeIndex"`
		NftPowerIndex      int    `json:"nftPowerIndex"`
		NftValue           int    `json:"nftValue" orm:"column(price)"`
		Qty                int    `json:"qty"`
		Thumbnail          string `json:"thumnail" orm:"column(file_name)"`
		DecSource          string `json:"decSource" orm:"column(music_file_name)"`
		TransactionAddress string `json:"transactionAddress"`
	}
	var datPurchaseInfo []nftTranData
	qb.Select(
		"nft_info.nft_ldef_index",
		"nft_info.nft_type",
		"nft_info.nft_name",
		"nft_info.short_desc",
		"nft_info.long_desc",
		"nft_info.file_name",
		"nft_market_info.price",
		"nft_market_info.qty",
		"dat_nft_info.music_file_name",
		"nft_purchase_info.transaction_address",
	).
		From("nft_purchase_info").
		InnerJoin("nft_market_info").
		On("nft_purchase_info.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("dat_nft_market_info").
		On("nft_purchase_info.nft_ldef_index = dat_nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_purchase_info.nft_ldef_index = nft_info.nft_ldef_index").
		InnerJoin("dat_nft_info").
		On("nft_purchase_info.nft_ldef_index = dat_nft_info.nft_ldef_index").
		Where("nft_purchase_info.uuid = ?").OrderBy("timestamp").Desc()
	sql := qb.String()
	num, err := o.Raw(sql, uuid).QueryRows(&datPurchaseInfo)
	if err != nil && err != orm.ErrNoRows {
		logs.Error(err.Error())
		err := errors.New("unknown error when query database")
		m.errorHandler(c, action, err)
		return
	}
	logs.Debug("get", num, "from database")
	type response struct {
		Status        int           `json:"status"`
		Action        string        `json:"action"`
		SupportedType string        `json:"supportedType"`
		PurchaseList  []nftTranData `json:"purchaseList"`
	}
	if num == 0 {
		datPurchaseInfo = make([]nftTranData, 0)
	}
	for i, _ := range datPurchaseInfo {
		datPurchaseInfo[i].Thumbnail = util.PathPrefixOfNFT("", common.PATH_KIND_USER_ICON) + datPurchaseInfo[i].Thumbnail
		decryptedFilePath, err := util.DecryptFile(datPurchaseInfo[i].DecSource, common.TYPE_NFT_MUSIC)
		if err != nil {
			logs.Error(err.Error())
			err := errors.New("can not decrypt file for " + datPurchaseInfo[i].NftLdefIndex)
			m.errorHandler(c, action, err)
			return
		}
		datPurchaseInfo[i].DecSource = util.PathPrefixOfNFT(common.TYPE_NFT_MUSIC, common.PATH_KIND_PUBLIC) + decryptedFilePath
	}
	m.wrapperAndSend(c, action, &response{
		Status:        common.RESPONSE_STATUS_SUCCESS,
		Action:        action,
		SupportedType: common.TYPE_NFT_MUSIC,
		PurchaseList:  datPurchaseInfo,
	})
}

func (m *Manager) OtherPurchaseHistory(c *client.Client, action string, uuid string) {
	o := orm.NewOrm()
	dbEngine := beego.AppConfig.String("dbEngine")
	qb, _ := orm.NewQueryBuilder(dbEngine)
	type nftTranData struct {
		NftLdefIndex       string `json:"nftLdefIndex"`
		SupportedType      string `json:"supportedType" orm:"column(nft_type)"`
		NftName            string `json:"nftName"`
		ShortDesc          string `json:"shortDesc"`
		LongDesc           string `json:"longDesc"`
		NftLifeIndex       int    `json:"nftLifeIndex"`
		NftPowerIndex      int    `json:"nftPowerIndex"`
		NftValue           int    `json:"nftValue" orm:"column(price)"`
		Qty                int    `json:"qty"`
		Thumbnail          string `json:"thumnail" orm:"column(file_name)"`
		TransactionAddress string `json:"transactionAddress"`
	}
	var otherPurchaseInfo []nftTranData
	qb.Select(
		"nft_info.nft_ldef_index",
		"nft_info.nft_type",
		"nft_info.nft_name",
		"nft_info.short_desc",
		"nft_info.long_desc",
		"nft_market_info.price",
		"nft_market_info.qty",
		"nft_info.file_name",
		"nft_purchase_info.transaction_address",
	).
		From("nft_purchase_info").
		InnerJoin("nft_market_info").
		On("nft_purchase_info.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("other_nft_market_info").
		On("nft_purchase_info.nft_ldef_index = other_nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_purchase_info.nft_ldef_index = nft_info.nft_ldef_index").
		InnerJoin("other_nft_info").
		On("nft_purchase_info.nft_ldef_index = other_nft_info.nft_ldef_index").
		Where("nft_purchase_info.uuid = ?").OrderBy("timestamp").Desc()
	sql := qb.String()
	num, err := o.Raw(sql, uuid).QueryRows(&otherPurchaseInfo)
	if err != nil && err != orm.ErrNoRows {
		logs.Error(err.Error())
		err := errors.New("unknown error when query database")
		m.errorHandler(c, action, err)
		return
	}
	logs.Debug("get", num, "from database")
	type response struct {
		Status        int           `json:"status"`
		Action        string        `json:"action"`
		SupportedType string        `json:"supportedType"`
		PurchaseList  []nftTranData `json:"purchaseList"`
	}
	if num == 0 {
		otherPurchaseInfo = make([]nftTranData, 0)
	}
	for i, _ := range otherPurchaseInfo {
		otherPurchaseInfo[i].Thumbnail = util.PathPrefixOfNFT("", common.PATH_KIND_USER_ICON) + otherPurchaseInfo[i].Thumbnail
	}
	m.wrapperAndSend(c, action, &response{
		Status:        common.RESPONSE_STATUS_SUCCESS,
		Action:        action,
		SupportedType: common.TYPE_NFT_OTHER,
		PurchaseList:  otherPurchaseInfo,
	})
}

func (m *Manager) NFTPurchaseHistoryHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Action        string `json:"action"`
		Uuid          string `json:"uuid"`
		SupportedType string `json:"supportedType"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}
	nftType := req.SupportedType
	if err := util.ValidNftType(nftType); err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, action, err)
		return
	}
	switch req.SupportedType {
	case common.TYPE_NFT_MUSIC:
		m.DatPurchaseHistory(c, action, req.Uuid)
	case common.TYPE_NFT_OTHER:
		m.OtherPurchaseHistory(c, action, req.Uuid)
	case common.TYPE_NFT_AVATAR:
		m.AvatarPurchaseHistory(c, action, req.Uuid)
	}
}
