package mobile

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/transactionQueue"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/util"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
)

func (m *Manager) ItemDetailsHandler(c *client.Client, action string,data []byte) {
	type request struct {
		NftLdefIndex string `json:"nftLdefIndex"`
		SupportedType string `json:"supportedType"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err:=errors.New("wrong data format")
		m.errorHandler(c, action, err)
		return
	}

	nftType:= req.SupportedType
	if err:=util.ValidNftType(nftType); err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, action, err)
		return
	}

	switch nftType {
	case common.TYPE_NFT_AVATAR:
		type nftTranData struct {
			NftLdefIndex string `json:"nftLdefIndex"`
			NftName string `json:"nftName"`
			ShortDesc string `json:"shortDesc"`
			LongDesc string `json:"longDesc"`
			NftLifeIndex int `json:"nftLifeIndex"`
			NftPowerIndex int `json:"nftPowerIndex"`
			NftValue int `json:"nftValue" orm:"column(price)"`
			Qty int `json:"qty"`
			Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
		}
		o:=orm.NewOrm()
		dbEngine := beego.AppConfig.String("dbEngine")
		qb,_:=orm.NewQueryBuilder(dbEngine)
		var avatarInfo nftTranData
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
			From("avatar_nft_info").
			InnerJoin("avatar_nft_market_info").
			On("avatar_nft_info.nft_ldef_index = avatar_nft_market_info.nft_ldef_index").
			InnerJoin("nft_info").
			On("avatar_nft_info.nft_ldef_index = nft_info.nft_ldef_index").
			InnerJoin("nft_market_info").
			On("avatar_nft_market_info.nft_ldef_index = nft_market_info.nft_ldef_index").
			InnerJoin("nft_market_place").
			On("avatar_nft_info.nft_ldef_index = nft_market_place.nft_ldef_index").
			Where("avatar_nft_info.nft_ldef_index = ?").OrderBy("timestamp").Desc()
		sql := qb.String()
		err:=o.Raw(sql,req.NftLdefIndex).QueryRow(&avatarInfo)
		if err!=nil {
			if err== orm.ErrNoRows {
				logs.Error(err.Error())
				err:= errors.New("no such item:"+req.NftLdefIndex+" in marketplace, plase check nft ldef index and nft type")
				m.errorHandler(c,action,err)
				return
			} else {
				logs.Error(err.Error())
				err:= errors.New("unexpected error when query database")
				m.errorHandler(c,action,err)
				return
			}
		}
		avatarInfo.Thumbnail = util.PathPrefixOfNFT(common.TYPE_NFT_AVATAR,common.PATH_KIND_MARKET) + avatarInfo.Thumbnail
		type response struct {
			SupportedType string `json:"supportedType"`
			NftLdefIndex string `json:"nftLdefIndex"`
			Status int `json:"status"`
			Action string `json:"action"`
			NftTranData nftTranData `json:"nftTranData"`
		}
		m.wrapperAndSend(c,action,&response{
			SupportedType: req.SupportedType,
			NftLdefIndex: req.NftLdefIndex,
			Status: common.RESPONSE_STATUS_SUCCESS,
			Action: action,
			NftTranData: avatarInfo,
		})
	case common.TYPE_NFT_OTHER:
		type nftTranData struct {
			NftLdefIndex string `json:"nftLdefIndex"`
			NftName string `json:"nftName"`
			ShortDesc string `json:"shortDesc"`
			LongDesc string `json:"longDesc"`
			NftValue int `json:"nftValue" orm:"column(price)"`
			Qty int `json:"qty"`
			Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
			NftParentLdef string `json:"nftParentLdef"`
		}
		o:=orm.NewOrm()
		dbEngine := beego.AppConfig.String("dbEngine")
		qb,_:=orm.NewQueryBuilder(dbEngine)
		var otherInfo nftTranData
		qb.Select(
			"nft_info.nft_ldef_index",
			"nft_info.nft_name",
			"nft_info.short_desc",
			"nft_info.long_desc",
			"nft_market_info.price",
			"nft_market_info.qty",
			"nft_info.file_name",
			"nft_info.nft_parent_ldef",
		).
			From("other_nft_info").
			InnerJoin("other_nft_market_info").
			On("other_nft_info.nft_ldef_index = other_nft_market_info.nft_ldef_index").
			InnerJoin("nft_info").
			On("other_nft_info.nft_ldef_index = nft_info.nft_ldef_index").
			InnerJoin("nft_market_info").
			On("other_nft_info.nft_ldef_index = nft_market_info.nft_ldef_index").
			InnerJoin("nft_market_place").
			On("other_nft_info.nft_ldef_index = nft_market_place.nft_ldef_index").
			Where("other_nft_info.nft_ldef_index = ?").OrderBy("timestamp").Desc()
		sql := qb.String()
		err:=o.Raw(sql,req.NftLdefIndex).QueryRow(&otherInfo)
		if err!=nil {
			if err== orm.ErrNoRows {
				logs.Error(err.Error())
				err:= errors.New("no such item:"+req.NftLdefIndex+" plase check nft ldef index and nft type")
				m.errorHandler(c,action,err)
				return
			} else {
				logs.Error(err.Error())
				err:= errors.New("unexpected error when query database")
				m.errorHandler(c,action,err)
				return
			}
		}
		otherInfo.Thumbnail = util.PathPrefixOfNFT(common.TYPE_NFT_OTHER,common.PATH_KIND_MARKET) + otherInfo.Thumbnail

		type response struct {
			SupportedType string `json:"supportedType"`
			NftLdefIndex string `json:"nftLdefIndex"`
			Status int `json:"status"`
			Action string `json:"action"`
			NftTranData nftTranData `json:"nftTranData"`
		}
		m.wrapperAndSend(c,action,&response{
			SupportedType: req.SupportedType,
			NftLdefIndex: req.NftLdefIndex,
			Status: common.RESPONSE_STATUS_SUCCESS,
			Action: action,
			NftTranData: otherInfo,
		})
	case common.TYPE_NFT_MUSIC:
		type nftTranData struct {
			NftLdefIndex string `json:"nftLdefIndex"`
			NftName string `json:"nftName"`
			ShortDesc string `json:"shortDesc"`
			LongDesc string `json:"longDesc"`
			NftValue int `json:"nftValue" orm:"column(price)"`
			Qty int `json:"qty"`
			Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
		}
		o:=orm.NewOrm()
		dbEngine := beego.AppConfig.String("dbEngine")
		qb,_:=orm.NewQueryBuilder(dbEngine)
		var datInfo nftTranData
		qb.Select(
			"nft_info.nft_ldef_index",
			"nft_info.nft_name",
			"nft_info.short_desc",
			"nft_info.long_desc",
			"nft_market_info.price",
			"nft_market_info.qty",
			"nft_info.file_name",
		).
			From("dat_nft_info").
			InnerJoin("dat_nft_market_info").
			On("dat_nft_info.nft_ldef_index = dat_nft_market_info.nft_ldef_index").
			InnerJoin("nft_info").
			On("dat_nft_info.nft_ldef_index = nft_info.nft_ldef_index").
			InnerJoin("nft_market_info").
			On("dat_nft_info.nft_ldef_index = nft_market_info.nft_ldef_index").
			InnerJoin("nft_market_place").
			On("dat_nft_info.nft_ldef_index = nft_market_place.nft_ldef_index").
			Where("dat_nft_info.nft_ldef_index = ?").OrderBy("timestamp").Desc()
		sql := qb.String()
		err:=o.Raw(sql,req.NftLdefIndex).QueryRow(&datInfo)
		if err!=nil {
			if err== orm.ErrNoRows {
				logs.Error(err.Error())
				err:= errors.New("no such item:"+req.NftLdefIndex+" plase check nft ldef index and nft type")
				m.errorHandler(c,action,err)
				return
			} else {
				logs.Error(err.Error())
				err:= errors.New("unexpected error when query database")
				m.errorHandler(c,action,err)
				return
			}
		}
		datInfo.Thumbnail = util.PathPrefixOfNFT(common.TYPE_NFT_MUSIC,common.PATH_KIND_MARKET) + datInfo.Thumbnail

		type response struct {
			SupportedType string `json:"supportedType"`
			NftLdefIndex string `json:"nftLdefIndex"`
			Status int `json:"status"`
			Action string `json:"action"`
			NftTranData nftTranData `json:"nftTranData"`
		}
		m.wrapperAndSend(c,action,&response{
			SupportedType: req.SupportedType,
			NftLdefIndex: req.NftLdefIndex,
			Status: common.RESPONSE_STATUS_SUCCESS,
			Action: action,
			NftTranData: datInfo,
		})
	}
}

func (m *Manager) GetMPListOfAvatar(c *client.Client, action string) {
	nftType:= common.TYPE_NFT_AVATAR
	o:=orm.NewOrm()
	dbEngine := beego.AppConfig.String("dbEngine")
	qb,_:=orm.NewQueryBuilder(dbEngine)
	type nftTranData struct {
		NftLdefIndex string `json:"nftLdefIndex"`
		NftName string `json:"nftName"`
		ShortDesc string `json:"shortDesc"`
		LongDesc string `json:"longDesc"`
		NftLifeIndex int `json:"nftLifeIndex"`
		NftPowerIndex int `json:"nftPowerIndex"`
		NftValue int `json:"nftValue" orm:"column(price)"`
		Qty int `json:"qty"`
		Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
	}
	var avatarMKPlaceInfo []nftTranData
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
		From("nft_market_place").
		InnerJoin("nft_market_info").
		On("nft_market_place.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("avatar_nft_market_info").
		On("nft_market_place.nft_ldef_index = avatar_nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_market_place.nft_ldef_index = nft_info.nft_ldef_index").
		InnerJoin("avatar_nft_info").
		On("nft_market_place.nft_ldef_index = avatar_nft_info.nft_ldef_index").
		Where("nft_info.nft_type = ?").OrderBy("timestamp").Desc()
	sql := qb.String()
	num,err:=o.Raw(sql,nftType).QueryRows(&avatarMKPlaceInfo)
	if err != nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:= errors.New("unexpected error when query database")
		m.errorHandler(c,action,err)
		return
	}
	logs.Debug("get",num,"from database")
	type response struct {
		SupportedType string `json:"supportedType"`
		Status int `json:"status"`
		Action string `json:"action"`
		NftTranData []nftTranData `json:"nftTranData"`
	}
	if num == 0 {
		avatarMKPlaceInfo= make([]nftTranData,0)
	}
	for i,_:= range avatarMKPlaceInfo {
		avatarMKPlaceInfo[i].Thumbnail = util.PathPrefixOfNFT(nftType,common.PATH_KIND_MARKET) + avatarMKPlaceInfo[i].Thumbnail
	}
	res:= response{
		SupportedType:nftType,
		Status: common.RESPONSE_STATUS_SUCCESS,
		Action: action,
		NftTranData: avatarMKPlaceInfo,
	}
	m.wrapperAndSend(c,action,&res)
}

func (m *Manager) GetMPListOfDat(c *client.Client, action string) {
	nftType:= common.TYPE_NFT_MUSIC
	o:=orm.NewOrm()
	type nftTranData struct {
		NftLdefIndex string `json:"nftLdefIndex"`
		NftName string `json:"nftName"`
		ShortDesc string `json:"shortDesc"`
		LongDesc string `json:"longDesc"`
		NftValue int `json:"nftValue" orm:"column(price)"`
		Qty int `json:"qty"`
		Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
	}
	dbEngine := beego.AppConfig.String("dbEngine")
	qb,_:=orm.NewQueryBuilder(dbEngine)
	var datMKPlaceInfo []nftTranData
	qb.Select(
		"nft_info.nft_ldef_index",
		"nft_info.nft_name",
		"nft_info.short_desc",
		"nft_info.long_desc",
		"nft_market_info.price",
		"nft_market_info.qty",
		"nft_info.file_name",
		).
		From("nft_market_place").
		InnerJoin("nft_market_info").
		On("nft_market_place.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("dat_nft_market_info").
		On("nft_market_place.nft_ldef_index = dat_nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_market_place.nft_ldef_index = nft_info.nft_ldef_index").
		InnerJoin("dat_nft_info").
		On("nft_market_place.nft_ldef_index = dat_nft_info.nft_ldef_index").
		Where("nft_info.nft_type = ?").OrderBy("timestamp").Desc()
	sql := qb.String()
	num,err:=o.Raw(sql,nftType).QueryRows(&datMKPlaceInfo)
	if err != nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:= errors.New("unexpected error when query database")
		m.errorHandler(c,action,err)
		return
	}
	logs.Debug("get",num,"from database")
	type response struct {
		SupportedType string `json:"supportedType"`
		Status int `json:"status"`
		Action string `json:"action"`
		NftTranData []nftTranData `json:"nftTranData"`
	}
	if num == 0 {
		datMKPlaceInfo= make([]nftTranData,0)
	}
	for i,_:= range datMKPlaceInfo {
		datMKPlaceInfo[i].Thumbnail = util.PathPrefixOfNFT(nftType,common.PATH_KIND_MARKET) + datMKPlaceInfo[i].Thumbnail
	}
	res:= response{
		SupportedType: nftType,
		Status: common.RESPONSE_STATUS_SUCCESS,
		Action: action,
		NftTranData: datMKPlaceInfo,
	}
	m.wrapperAndSend(c,action,&res)
}

func (m *Manager) GetMPListOfOther(c *client.Client, action string) {
	nftType:= common.TYPE_NFT_OTHER
	o:=orm.NewOrm()
	dbEngine := beego.AppConfig.String("dbEngine")
	qb,_:=orm.NewQueryBuilder(dbEngine)
	type nftTranData struct {
		NftLdefIndex string `json:"nftLdefIndex"`
		NftName string `json:"nftName"`
		ShortDesc string `json:"shortDesc"`
		LongDesc string `json:"longDesc"`
		NftValue int `json:"nftValue" orm:"column(price)"`
		Qty int `json:"qty"`
		Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
		NftParentLdef string `json:"nftParentLdef"`
	}
	var otherMKPlaceInfo []nftTranData
	qb.Select(
		"nft_info.nft_ldef_index",
		"nft_info.nft_name",
		"nft_info.short_desc",
		"nft_info.long_desc",
		"nft_market_info.price",
		"nft_market_info.qty",
		"nft_info.file_name",
		"nft_info.nft_parent_ldef",
		).
		From("nft_market_place").
		InnerJoin("nft_market_info").
		On("nft_market_place.nft_ldef_index = nft_market_info.nft_ldef_index").
		InnerJoin("other_nft_market_info").
		On("nft_market_place.nft_ldef_index = other_nft_market_info.nft_ldef_index").
		InnerJoin("nft_info").
		On("nft_market_place.nft_ldef_index = nft_info.nft_ldef_index").
		InnerJoin("other_nft_info").
		On("nft_market_place.nft_ldef_index = other_nft_info.nft_ldef_index").
		Where("nft_info.nft_type = ?").OrderBy("timestamp").Desc()
	sql := qb.String()
	num,err:=o.Raw(sql,nftType).QueryRows(&otherMKPlaceInfo)
	if err != nil && err!=orm.ErrNoRows {
		logs.Error(err.Error())
		err:= errors.New("unexpected error when query database")
		m.errorHandler(c,action,err)
		return
	}
	logs.Debug("get",num,"from database")
	type response struct {
		SupportedType string `json:"supportedType"`
		Status int `json:"status"`
		Action string `json:"action"`
		NftTranData []nftTranData `json:"nftTranData"`
	}
	if num == 0 {
		otherMKPlaceInfo= make([]nftTranData,0)
	}
	for i,_:= range otherMKPlaceInfo {
		otherMKPlaceInfo[i].Thumbnail = util.PathPrefixOfNFT(nftType,common.PATH_KIND_MARKET) + otherMKPlaceInfo[i].Thumbnail
	}
	res:= response{
		SupportedType:nftType,
		Status: common.RESPONSE_STATUS_SUCCESS,
		Action: action,
		NftTranData: otherMKPlaceInfo,
	}
	m.wrapperAndSend(c,action,&res)
}


func (m *Manager) GetMPListHandler(c *client.Client, action string, data[]byte) {
	type request struct {
		SupportedType string `json:"supportedType"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err:=errors.New("wrong data format")
		m.errorHandler(c, action, err)
		return
	}

	nftType:= req.SupportedType
	if err:=util.ValidNftType(nftType); err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, action, err)
		return
	}

	switch nftType {
	case common.TYPE_NFT_AVATAR:
		m.GetMPListOfAvatar(c,action)
	case common.TYPE_NFT_OTHER:
		m.GetMPListOfOther(c,action)
	case common.TYPE_NFT_MUSIC:
		m.GetMPListOfDat(c,action)
	}
}

// nft show
func (m *Manager) NFTDisplayHandler(c *client.Client, action string, data []byte) {
	type request struct {
		NftLdefIndex string `json:"nftLdefIndex"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err:=errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}

	nftInfo:= models.NftInfo{
		NftLdefIndex: req.NftLdefIndex,
	}
	o:=orm.NewOrm()
	err=o.Read(&nftInfo)
	if err!=nil {
		if err == orm.ErrNoRows {
			err:=errors.New("nft "+req.NftLdefIndex+" not exists")
			logs.Error(err.Error())
			m.errorHandler(c, action, err)
			return
		} else {
			logs.Error(err.Error())
			err := errors.New("unexpected error when query db")
			m.errorHandler(c, action, err)
			return
		}
	}

	nftType:= nftInfo.NftType
	fileName := nftInfo.FileName
	decryptedFilePath,err:=util.DecryptFile(fileName,nftType)
	if err!=nil {
		logs.Error(err.Error())
		err:=errors.New("can not decrypt file for the nft")
		m.errorHandler(c, action, err)
		return
	}
	type response struct {
		Status int `json:"status"`
		Action string `json:"action"`
		NftLdefIndex string `json:"nftLdefIndex"`
		DecSource string `json:"decSource"`
	}
	m.wrapperAndSend(c, action,&response{
		Status: common.RESPONSE_STATUS_SUCCESS,
		Action: action,
		NftLdefIndex: req.NftLdefIndex,
		DecSource: util.PathPrefixOfNFT(nftType,common.PATH_KIND_PUBLIC)+decryptedFilePath,
	})
}

func (m *Manager) NFTTransferHandler(c *client.Client, action string, data []byte) {
	type request struct {
		SenderUuid string `json:"senderUuid"`
		ReceiverUuid string `json:"receiverUuid"`
		NftLdefIndex string `json:"nftLdefIndex"`
	}
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err:=errors.New("wrong request data format")
		m.errorHandler(c, action, err)
		return
	}

	// change count of user
	o:=orm.NewOrm()
	o.Begin()
	// add count for buyer
	_,err = o.QueryTable("user_market_info").Filter("uuid",req.ReceiverUuid).Update(orm.Params{
		"count": orm.ColValue(orm.ColAdd,1),
	})
	if err!=nil {
		o.Rollback()
		if err== orm.ErrNoRows {
			err:=errors.New("no such user "+req.ReceiverUuid)
			logs.Error(err.Error())
			m.errorHandler(c, action, err)
			return
		} else {
			logs.Error(err.Error())
			err:=errors.New("unexpected error when query db")
			m.errorHandler(c, action, err)
			return
		}
	}
	// reduce count for seller
	_,err = o.QueryTable("user_market_info").Filter("uuid",req.SenderUuid).Update(orm.Params{
		"count": orm.ColValue(orm.ColMinus,1),
	})
	if err!=nil {
		o.Rollback()
		if err== orm.ErrNoRows {
			err:=errors.New("no such user "+req.SenderUuid)
			logs.Error(err.Error())
			m.errorHandler(c, action, err)
			return
		} else {
			logs.Error(err.Error())
			err:=errors.New("unexpected error when query db")
			m.errorHandler(c, action, err)
			return
		}
	}
	o.Commit()
	type response struct {
		Status int  `json:"status"`
		Action string  `json:"action"`
		ReceiverUuid string `json:"receiverUuid"`
		NftLdefIndex string `json:"nftLdefIndex"`
	}
	m.wrapperAndSend(c,action,&response{
		Status: common.RESPONSE_STATUS_SUCCESS,
		Action: action,
		ReceiverUuid: req.ReceiverUuid,
		NftLdefIndex: req.NftLdefIndex,
	})
	// TODO use message queue instead of go channel
	m.TransactionQueue.Append(&transactionQueue.TransferNftTransaction{
		Uuid: req.ReceiverUuid,
		SellerUuid: req.SenderUuid,
		NftLdefIndex: req.NftLdefIndex,
	})
}
