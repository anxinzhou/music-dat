package mobile

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/util"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
)

func (m *Manager) GetMPListOfAvatar(c *client.Client, action string) {
	nftType:= common.TYPE_NFT_AVATAR
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
		avatarMKPlaceInfo[i].FileName = util.PathPrefixOfNFT(nftType,common.PATH_KIND_MARKET)
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
		datMKPlaceInfo[i].FileName = util.PathPrefixOfNFT(nftType,common.PATH_KIND_MARKET)
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
		otherMKPlaceInfo[i].FileName = util.PathPrefixOfNFT(nftType,common.PATH_KIND_MARKET)
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

func (m *Manager) NFTTransferHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req NftTransferRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	senderNickname := req.SenderNickname
	receiverNickname := req.ReceiverNickname
	nftInfo := req.NftTranData
	nftLdefIndex := nftInfo.NftLdefIndex
	supportedType := nftInfo.SupportedType
	if err := validNftLdefIndex(nftLdefIndex); err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	if err := validSupportedType(supportedType); err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	senderWalletId, err := models.WalletIdOfNickname(senderNickname)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	receiverWalletId, err := models.WalletIdOfNickname(receiverNickname)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	o:=orm.NewOrm()
	o.Begin()
	// add count for sender
	_,err = o.QueryTable("market_user_table").Filter("nickname",senderNickname).Update(orm.Params{
		"count": orm.ColValue(orm.ColAdd,1),
	})
	if err!=nil {
		o.Rollback()
		logs.Emergency("can not add count for nickname:", senderNickname)
		m.errorHandler(c, bq, err)
		return
	}
	logs.Warn("add count in market table for",senderNickname)
	// reduce count for receiver
	_,err = o.QueryTable("market_user_table").Filter("nickname",receiverNickname).Update(orm.Params{
		"count": orm.ColValue(orm.ColMinus,1),
	})
	if err!=nil {
		o.Rollback()
		logs.Emergency("can not add count for nickname:", receiverNickname)
		m.errorHandler(c, bq, err)
		return
	}
	logs.Warn("reduce count in market table for",receiverNickname)


	tokenId := TokenIdFromNftLdefIndex(nftLdefIndex)
	txErr := m.chainHandler.ManagerAccount.SendFunction(m.chainHandler.Contract,
		nil,
		nft.FuncDelegateTransfer,
		common.HexToAddress(senderWalletId),
		common.HexToAddress(receiverWalletId),
		tokenId)
	err = <-txErr
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	o.Commit()
	logs.Debug("send", nftLdefIndex, "from", senderWalletId, "to", receiverWalletId)
	m.wrapperAndSend(c, bq, &NftTransferResponse{
		RQBaseInfo: *bq,
		Status:     NFT_TRANSFER_SUCCESS,
	})
}
