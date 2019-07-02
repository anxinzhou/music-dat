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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *Manager) BindWalletHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req BindWalletRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	walletId:= req.WalletId
	nickname:= req.Nickname

	walletInfo:= &models.MarketUserTable{
		WalletId: walletId,
		Count: 0,
		Nickname: nickname,
		UserIconUrl: "",
	}

	o:=orm.NewOrm()
	o.Begin()         //TODO single sql
	err = o.Read(walletInfo)
	if err != nil {
		if err!= orm.ErrNoRows {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		} else {
			_,err:=o.Insert(walletInfo)
			if err!=nil {
				o.Rollback()
				logs.Error(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
		}
	}
	o.Commit()
	logs.Info("insert to market user table","nickname",nickname)
	m.wrapperAndSend(c,bq,&BindWalletResponse{
		RQBaseInfo: *bq,
	})
}

func isNicknameExist(nickname string) (bool,error) {
	type fields struct {
		Nickname string `bson:"nickname"`
	}
	filter:= bson.M {
		"nickname": nickname,
	}
	var queryResult fields
	col := models.MongoDB.Collection("users")
	err:=col.FindOne(context.Background(),filter,options.FindOne().SetProjection(bson.M{
		"nickname": true,
	})).Decode(&queryResult)
	if err!=nil {
		if err == mongo.ErrNoDocuments {
			return false,nil
		} else {
			return false,err
		}
	}
	return true,nil
}

func (m *Manager) SetNicknameHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req SetNicknameRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	uuid:= req.Uuid
	nickname:= req.Nickname
	duplicated,err:=isNicknameExist(nickname)
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	if duplicated {
		err:= errors.New("nick name already exists")
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	o:= orm.NewOrm()
	userBaseInfo:= models.MarketUserTable{
		Uuid:uuid,
		Nickname:nickname,
	}
	_,err=o.InsertOrUpdate(&userBaseInfo,"nickname")
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	col := models.MongoDB.Collection("users")
	filter:= bson.M {
		"uuid": uuid,
	}
	update:= bson.M {
		"$set": bson.M {"nickname":nickname},
	}
	_,err=col.UpdateOne(context.Background(),filter,update)
	if err!=nil {
		o.Rollback()
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	o.Commit()
	logs.Warn("insert into user base info table")
	m.wrapperAndSend(c,bq,&SetNicknameResponse{
		RQBaseInfo: *bq,
	})
}

func (m *Manager) IsNicknameDuplicatedHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req IsNicknameDuplicatedRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	nickname:= req.Nickname
	logs.Debug("nick name to test",nickname)
	duplicated,err:=isNicknameExist(nickname)
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	m.wrapperAndSend(c,bq,&IsNicknameDuplicatedResponse{
		RQBaseInfo: *bq,
		Duplicated: duplicated,
	})
}

func (m *Manager) FollowListHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req FollowListRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	nickname:=req.Nickname
	o:= orm.NewOrm()
	type queryInfo struct {
		FolloweeNickname string `orm:"column(followee_nickname)"`
		UserIconUrl string `orm:"column(user_icon_url)"`
		Intro string `orm:"column(intro)"`
	}
	var queryResult []queryInfo
	num,err:=o.Raw(`
		select followee_nickname,user_icon_url, intro from market_user_table as mk, follow_table as ft
		where mk.nickname = ft.followee_nickname and ft.follower_nickname = ?`, nickname).QueryRows(&queryResult)
	if err!=nil {
		if err!= orm.ErrNoRows {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
	}

	followInfo:=make([]*FollowInfo,num)
	for i,_:= range queryResult {
		if queryResult[i].UserIconUrl == "" {
			queryResult[i].UserIconUrl = "default.jpg"
		}
		followInfo[i] = &FollowInfo{
			Nickname: queryResult[i].FolloweeNickname,
			Thumbnail: PathPrefixOfNFT("", PATH_KIND_USER_ICON)+queryResult[i].UserIconUrl,
			Intro: queryResult[i].Intro,
		}
	}
	logs.Debug("user has",num,"followee")

	m.wrapperAndSend(c,bq,&FollowListResponse{
		RQBaseInfo: *bq,
		FollowInfo: followInfo,
	})
}

func (m*Manager) FollowListOperationHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req FollowListOperationRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	operation:=req.Operation
	followNickname:= req.FollowNickname
	nickname:= req.Nickname
	o:= orm.NewOrm()
	o.Begin()
	queryInfo:= models.FollowTable{
		FollowerNickname: nickname,
		FolloweeNickname: followNickname,
	}
	err=o.Read(&queryInfo,"followee_nickname","follower_nickname")
	if operation == FOLLOW_LIST_ADD {
		if err!=nil {
			if err== orm.ErrNoRows {
				_,err:=o.Insert(&queryInfo)
				if err!=nil {
					o.Rollback()
					logs.Error(err.Error())
					m.errorHandler(c, bq, err)
					return
				}
				logs.Warn("follow",followNickname)
			}  else  {
				o.Rollback()
				logs.Error(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
		} else {
			err:= errors.New("already follow "+followNickname)
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
	} else if operation == FOLLOW_LIST_DELETE {
		if err!=nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		} else {
			_,err:=o.Delete(&queryInfo,"followee_nickname","follower_nickname")
			if err!=nil {
				o.Rollback()
				logs.Error(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
			logs.Warn("delete followee",followNickname)
		}
	} else {
		err:= errors.New("operation not exist")
		m.errorHandler(c, bq, err)
		return
	}
	o.Commit()
	m.wrapperAndSend(c,bq,&FollowListOperationResponse{
		RQBaseInfo: *bq,
		FollowNickname:followNickname,
		Operation: operation,
	})
}

func (m*Manager) IsNicknameSetHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req IsNicknameSetRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	uuid:= req.Uuid
	queryInfo:= models.UserBaseInfo{
		Uuid:uuid,
	}
	set:= true
	o:=orm.NewOrm()
	err=o.Read(&queryInfo)
	if err!=nil {
		if err==orm.ErrNoRows {
			set = false
		} else {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
	} else {
		set = true
	}
	m.wrapperAndSend(c,bq,&IsNicknameSetResponse{
		RQBaseInfo: *bq,
		Set: set,
	})
}

func (m *Manager) MarketUserListHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Uuid string `json:"uuid"`
	}
	var req request
	type userInfo struct {
		Uuid string `json:"uuid"`
		Nickname string `json:"nickname"`
		Count int `json:"count"`
		Thumbnail string `json:"thumnail" orm:"column(avatar_file_name)"`
		Followed bool `json:"followed" orm:"column(followed)"`
	}
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err:= errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}
	o := orm.NewOrm()
	var ui []userInfo
	num,err := o.Raw(`
		select ui.uuid,ui.nickname,umi.count,ui.avatar_file_name,count(ft.followee_uuid) as followed from 
		user_market_info as umi inner join user_info as ui on umi.uuid = ui.uuid
		left join (
			select followee_uuid, follower_uuid from follow_table
			where follow_table.follower_uuid =?) 
			as ft on umi.uuid = ft.followee_uuid 
			`,req.Uuid).QueryRows(&ui)
	if err!=nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:= errors.New("unexpected error when query db")
		m.errorHandler(c, action, err)
		return
	}
	if num == 0 {
		ui=make([]userInfo,0)
	}
	for i,_:=range ui {
		ui[i].Thumbnail = util.PathPrefixOfNFT("",common.PATH_KIND_USER_ICON) + ui[i].Thumbnail
	}
	type response struct {
		Status int `json:"status"`
		Action string `json:"action"`
		WalletIdList []userInfo
	}
	m.wrapperAndSend(c, action, &response{
		Status: common.RESPONSE_STATUS_SUCCESS,
		Action: action,
		WalletIdList: ui,
	})
}

func (m *Manager) AvatarPurchaseHistory(c *client.Client, action string,uuid string) {
	o:=orm.NewOrm()
	dbEngine := beego.AppConfig.String("dbEngine")
	qb,_:=orm.NewQueryBuilder(dbEngine)
	type nftTranData struct {
		common.AvatarNftMarketInfo
		common.MarketPlaceInfo
	}
	var avatarMKPlaceInfo []nftTranData
	qb.Select("*").
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
	num,err:=o.Raw(sql,uuid).QueryRows(&avatarMKPlaceInfo)
	if err != nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:= errors.New("unknown error when query database")
		m.errorHandler(c,action,err)
		return
	}
	logs.Debug("get",num,"from database")
	type response struct {
		Status int `json:"status"`
		Action string `json:"action"`
		SupportedType string `json:"supportedType"`
		PurchaseList []nftTranData `json:"purchaseList"`
	}
	if num == 0 {
		avatarMKPlaceInfo= make([]nftTranData,0)
	}
	for i,_:=range avatarMKPlaceInfo {
		avatarMKPlaceInfo[i].Thumbnail = util.PathPrefixOfNFT("",common.PATH_KIND_USER_ICON) + ui[i].Thumbnail
	}
	m.wrapperAndSend(c,action,&response{
		Status: common.RESPONSE_STATUS_SUCCESS,
		Action: action,
		SupportedType: common.TYPE_NFT_AVATAR,
		PurchaseList: avatarMKPlaceInfo,
	})
}

func (m *Manager) DatPurchaseHistory(c *client.Client,action string, uuid string) {
	o:=orm.NewOrm()
	dbEngine := beego.AppConfig.String("dbEngine")
	qb,_:=orm.NewQueryBuilder(dbEngine)
	type nftTranData struct {
		common.DatNftMarketInfo
		common.MarketPlaceInfo
	}
	var datMKPlaceInfo []nftTranData
	qb.Select("*").
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
	num,err:=o.Raw(sql,uuid).QueryRows(&datMKPlaceInfo)
	if err != nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:= errors.New("unknown error when query database")
		m.errorHandler(c,action,err)
		return
	}
	logs.Debug("get",num,"from database")
	type response struct {
		Status int `json:"status"`
		Action string `json:"action"`
		SupportedType string `json:"supportedType"`
		PurchaseList []nftTranData `json:"purchaseList"`
	}
	if num == 0 {
		datMKPlaceInfo= make([]nftTranData,0)
	}
	for i,_:=range avatarMKPlaceInfo {
		avatarMKPlaceInfo[i].Thumbnail = util.PathPrefixOfNFT("",common.PATH_KIND_USER_ICON) + ui[i].Thumbnail
	}
	m.wrapperAndSend(c,action,&response{
		Status: common.RESPONSE_STATUS_SUCCESS,
		Action: action,
		SupportedType: common.TYPE_NFT_MUSIC,
		PurchaseList: datMKPlaceInfo,
	})
}

func (m *Manager) OtherPurchaseHistory(c *client.Client, action string ,uuid string) {
	o:=orm.NewOrm()
	dbEngine := beego.AppConfig.String("dbEngine")
	qb,_:=orm.NewQueryBuilder(dbEngine)
	type nftTranData struct {
		common.OtherNftMarketInfo
		common.MarketPlaceInfo
	}
	var otherMKPlaceInfo []nftTranData
	qb.Select("*").
		From("nft_purchase_info").
		InnerJoin("nft_market_info").
		On("nft_purchase_info.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("other_nft_market_info").
		On("nft_purchase_info.nft_ldef_index = dat_nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_purchase_info.nft_ldef_index = nft_info.nft_ldef_index").
		InnerJoin("other_nft_info").
		On("nft_purchase_info.nft_ldef_index = other_nft_info.nft_ldef_index").
		Where("nft_purchase_info.uuid = ?").OrderBy("timestamp").Desc()
	sql := qb.String()
	num,err:=o.Raw(sql,uuid).QueryRows(&otherMKPlaceInfo)
	if err != nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:= errors.New("unknown error when query database")
		m.errorHandler(c,action,err)
		return
	}
	logs.Debug("get",num,"from database")
	type response struct {
		Status int `json:"status"`
		Action string `json:"action"`
		SupportedType string `json:"supportedType"`
		PurchaseList []nftTranData `json:"purchaseList"`
	}
	if num == 0 {
		otherMKPlaceInfo= make([]nftTranData,0)
	}
	for i,_:=range avatarMKPlaceInfo {
		avatarMKPlaceInfo[i].Thumbnail = util.PathPrefixOfNFT("",common.PATH_KIND_USER_ICON) + ui[i].Thumbnail
	}
	m.wrapperAndSend(c,action,&response{
		Status: common.RESPONSE_STATUS_SUCCESS,
		Action: action,
		SupportedType: common.TYPE_NFT_OTHER,
		PurchaseList: otherMKPlaceInfo,
	})
}

func (m *Manager) NFTPurchaseHistoryHandler(c *client.Client, action string, data []byte) {
	type request struct {
		Action string `json:"action"`
		Uuid string `json:"uuid"`
		SupportedType string `json:"supportedType"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err:=errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}
	nftType:=req.SupportedType
	if err:=util.ValidNftType(nftType);err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, action, err)
		return
	}
	switch req.SupportedType {
	case common.TYPE_NFT_MUSIC:
		m.DatPurchaseHistory(c,req.Uuid)
	case common.TYPE_NFT_OTHER:
		m.OtherPurchaseHistory(c,req.Uuid)
	case common.TYPE_NFT_AVATAR:
		m.AvatarPurchaseHistory(c,req.Uuid)
	}
}