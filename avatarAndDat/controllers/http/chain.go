package http

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

type NftBalanceController struct {
	ContractController
}

type nftBalanceResponse struct {
	Count int `json:"count"`
}

func (this *NftBalanceController) Get() {
	nickname := this.Ctx.Input.Param(":nickname")
	o:=orm.NewOrm()
	userMkInfo:=models.MarketUserTable {
		Nickname:nickname,
	}
	count:=0
	logs.Debug("nickname",nickname,"query balance")
	err:=o.Read(&userMkInfo)
	if err!=nil && err!= orm.ErrNoRows {
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	count = userMkInfo.Count

	this.Data["json"] = &nftBalanceResponse{
		Count: count,
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.ServeJSON()
}

type NftListController struct {
	ContractController
}

type nftInfoListRes struct {
	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
	NftName       string `json:"nftName"`
	NftValue      int    `json:"nftValue" orm:"column(price)"`
	ActiveTicker  string `json:"activeTicker"`
	NftLifeIndex  int64  `json:"nftLifeIndex"`
	NftPowerIndex int64  `json:"nftPowerIndex"`
	NftLdefIndex  string `json:"nftLdefIndex"`
	NftCharacId   string `json:"nftCharacId"`
	ShortDesc     string `json:"shortDesc" orm:"column(short_description)"`
	LongDesc      string `json:"longDesc" orm:"column(long_description)"`
	Thumbnail     string `json:"thumbnail" orm:"column(file_name)"`
	Qty           int    `json:"qty"`
}

type nftListResponse struct {
	NftTranData []*nftInfoListRes `json:"nftTranData"`
}

func (this *NftListController) Get() {
	nickname := this.Ctx.Input.Param(":nickname")
	logs.Debug("user", nickname, "query nft list")
	o:=orm.NewOrm()
	var mkInfos []models.NftMarketTable
	num,err:=o.QueryTable("nft_market_table").
		Filter("seller_nickname",nickname).
		All(&mkInfos,"nft_ldef_index")
	if err!=nil {
		if err == orm.ErrNoRows {
			logs.Info("no row in marketplace now")
			mkInfos = make([]models.NftMarketTable,0)
		} else {
			logs.Error(err.Error())
			sendError(&this.Controller, err, 500)
			return
		}
	}
	logs.Debug("number of list",num)

	nftTranResponseData := make([]*nftInfoListRes, 0, num)
	for _, mkInfo := range mkInfos {
		nftLdefIndex:= mkInfo.NftLdefIndex
		r := o.Raw(`
		select ni.nft_type, ni.nft_name,
		mk.price,mk.active_ticker, ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,na.short_description, na.long_description,mp.file_name,mk.qty from
		nft_market_table as mk,
		nft_mapping_table as mp,
		nft_info_table as ni,
		nft_item_admin as na
		where mk.nft_ldef_index = mp.nft_ldef_index and mk.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id and  ni.nft_ldef_index = ? `, nftLdefIndex)
		var nftResponseInfo nftInfoListRes
		err = r.QueryRow(&nftResponseInfo)
		if err != nil {
			if err == orm.ErrNoRows {
				logs.Debug(err.Error())
				continue
			} else {
				logs.Error(err.Error())
				sendError(&this.Controller, err, 500)
				return
			}
		}

		thumbnail := PathPrefixOfNFT(nftResponseInfo.SupportedType, PATH_KIND_MARKET)
		nftResponseInfo.Thumbnail = thumbnail + nftResponseInfo.Thumbnail
		nftTranResponseData = append(nftTranResponseData, &nftResponseInfo)
	}

	res := &nftListResponse{
		NftTranData: nftTranResponseData,
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &res
	this.ServeJSON()
}

type RewardController struct {
	ContractController
}

type nftInfoQuery struct {
	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
	NftName       string `json:"nftName"`
	NftValue      int    `json:"nftValue" orm:"column(price)"`
	ActiveTicker  string `json:"activeTicker"`
	NftLifeIndex  int64  `json:"nftLifeIndex"`
	NftPowerIndex int64  `json:"nftPowerIndex"`
	NftLdefIndex  string `json:"nftLdefIndex"`
	NftCharacId   string `json:"nftCharacId"`
	ShortDesc     string `json:"shortDesc" orm:"column(short_description)"`
	LongDesc      string `json:"longDesc" orm:"column(long_description)"`
	Thumbnail     string `json:"thumbnail" orm:"column(icon_file_name)"`
	Qty           int    `json:"qty"`
	SellerWalletId string
	SellerNickname string
}

type rewardNFTInfo struct {
	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
	NftName       string `json:"nftName"`
	NftValue      int    `json:"nftValue" orm:"column(price)"`
	ActiveTicker  string `json:"activeTicker"`
	NftLifeIndex  int64  `json:"nftLifeIndex"`
	NftPowerIndex int64  `json:"nftPowerIndex"`
	NftLdefIndex  string `json:"nftLdefIndex"`
	NftCharacId   string `json:"nftCharacId"`
	ShortDesc     string `json:"shortDesc" orm:"column(short_description)"`
	LongDesc      string `json:"longDesc" orm:"column(long_description)"`
	Thumbnail     string `json:"thumbnail" orm:"column(icon_file_name)"`
	Qty           int    `json:"qty"`
}

type RewardResponse struct {
	NftTranData []*rewardNFTInfo `json:"nftTranData"`
}

func (this *RewardController) RewardDat() {
	// only reward one dat now
	nickname := this.Ctx.Input.Param(":nickname")
	walletAddress,err:= models.WalletIdOfNickname(nickname)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}
	o := orm.NewOrm()
	o.Begin()
	qs := o.QueryTable("nft_market_table").Filter("nft_ldef_index__contains", "M").Filter("allow_airdrop",true).Limit(1)
	var mk models.NftMarketTable
	rewardAccount := 1
	err = qs.Limit(rewardAccount).One(&mk)
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	var res RewardResponse
	nftInfoList := make([]*rewardNFTInfo, rewardAccount)
	res.NftTranData = nftInfoList
	nftLdefIndex := mk.NftLdefIndex
	r := o.Raw(`
		select ni.nft_type, ni.nft_name,
		mk.price,mk.active_ticker, ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,na.short_description, na.long_description,mp.icon_file_name,mk.qty from
		nft_market_table as mk,
		nft_mapping_table as mp,
		nft_info_table as ni,
		nft_item_admin as na
		where mk.nft_ldef_index = mp.nft_ldef_index and mk.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id and  ni.nft_ldef_index = ? `, nftLdefIndex)
	var nftResponseInfo nftInfoQuery
	err = r.QueryRow(&nftResponseInfo)
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}
	nftType := nftResponseInfo.SupportedType
	thumbnail := PathPrefixOfNFT(nftType, PATH_KIND_MARKET)
	nftResponseInfo.Thumbnail = thumbnail + nftResponseInfo.Thumbnail
	//nftResponseInfo.Thumbnail = thumbnail + "music.png"

	nftInfoList[0] = &rewardNFTInfo{
		SupportedType: nftResponseInfo.SupportedType,
		NftName: nftResponseInfo.NftName,
		NftValue: nftResponseInfo.NftValue,
		ActiveTicker: nftResponseInfo.ActiveTicker,
		NftLifeIndex: nftResponseInfo.NftLifeIndex,
		NftPowerIndex: nftResponseInfo.NftPowerIndex,
		NftLdefIndex: nftResponseInfo.NftLdefIndex,
		NftCharacId: nftResponseInfo.NftCharacId,
		ShortDesc: nftResponseInfo.ShortDesc,
		LongDesc: nftResponseInfo.LongDesc,
		Thumbnail: nftResponseInfo.Thumbnail,
		Qty:          nftResponseInfo.Qty,
	}
	_, err = o.Delete(&mk)  //TODO comment for testing
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	sellerNickname:= nftResponseInfo.SellerNickname
	sellerWalletAddress:= nftResponseInfo.SellerWalletId
	// add count for buyer
	_,err = o.QueryTable("market_user_table").Filter("nickname",nickname).Update(orm.Params{
		"count": orm.ColValue(orm.ColAdd,1),
	})
	if err!=nil {
		o.Rollback()
		logs.Emergency("can not add count for nickname:", nickname)
		sendError(&this.Controller, err, 500)
		return
	}
	logs.Warn("add count in market table for",nickname)
	// reduce count for sender
	_,err = o.QueryTable("market_user_table").Filter("nickname",sellerNickname).Update(orm.Params{
		"count": orm.ColValue(orm.ColMinus,1),
	})
	if err!=nil {
		o.Rollback()
		logs.Emergency("can not reduce count for nickname:", sellerNickname)
		sendError(&this.Controller, err, 500)
		return
	}
	logs.Warn("reduce count in market table for",sellerNickname)

	if len(nftLdefIndex)<=1 {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	// insert to purchase history
	purchaseId := strconv.FormatInt(time.Now().UnixNano()|rand.Int63(), 10)
	storeInfo := &models.StorePurchaseHistroy{
		PurchaseId:    purchaseId,
		BuyerNickname: nickname,
		BuyerWalletId: walletAddress,
		SellerNickname: sellerNickname,
		SellerWalletId:     sellerWalletAddress,
		TotalPaid:     nftInfoList[0].NftValue,
		NftLdefIndex:  nftLdefIndex,
		ActiveTicker:  nftInfoList[0].ActiveTicker,
		Status:       PURCHASE_CONFIRMED ,
	}
	_, err = o.Insert(storeInfo)
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	tokenId, _ := new(big.Int).SetString(nftLdefIndex[1:], 10)

	nftContract := this.C.smartContract.(*nft.NFT)
	_, txErr := this.C.account.SendFunction2(nftContract,
		nil,
		nft.FuncDelegateTransfer,
		common.HexToAddress(sellerWalletAddress),
		common.HexToAddress(walletAddress),
		tokenId) // TODO redis to cache unsuccessful transaction
	err = <-txErr
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}
	o.Commit()
	this.Data["json"] = res
	this.ServeJSON()
}

type NumOfChildrenController struct {
	ContractController
}

type NumOfChildrenRes struct {
	Count int `json:"count"`
}

func (this *NumOfChildrenController) Get() {
	parentIndex := this.Ctx.Input.Param(":parentIndex")
	o := orm.NewOrm()
	r := o.Raw(`
		select count(a.nft_ldef_index) as num 
		from nft_mapping_table as a,
		nft_market_table as b 
		where a.nft_parent_ldef = ? and a.nft_ldef_index = b.nft_ldef_index `, parentIndex)
	type CountQuery struct {
		Num int
	}
	var queryResult CountQuery
	err := r.QueryRow(&queryResult)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	res := &NumOfChildrenRes{
		Count: queryResult.Num,
	}

	this.Data["json"] = res
	this.ServeJSON()
}

type ChildrenOfNFTController struct {
	ContractController
}

func (this *ChildrenOfNFTController) Get() {
	parentIndex := this.Ctx.Input.Param(":parentIndex")
	o := orm.NewOrm()
	r := o.Raw(`
		select ni.nft_type, ni.nft_name,
		mk.price,mk.active_ticker, ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,na.short_description, na.long_description,mp.file_name,mk.qty from
		nft_market_table as mk,
		nft_mapping_table as mp,
		nft_info_table as ni,
		nft_item_admin as na
		where mp.nft_parent_ldef= ? and mk.nft_ldef_index = mp.nft_ldef_index and mp.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id`, parentIndex)
	nftResponseInfo := []*nftInfoListRes{}
	_, err := r.QueryRows(&nftResponseInfo)
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller, err, 500)
		return
	}

	thumbnail := PathPrefixOfNFT(TYPE_NFT_OTHER, PATH_KIND_MARKET)
	for _, nftInfo := range nftResponseInfo {
		nftInfo.Thumbnail = thumbnail + nftInfo.Thumbnail
	}
	var res nftListResponse
	res.NftTranData = nftResponseInfo
	this.Data["json"] = &res
	this.ServeJSON()
}

type MarketTransactionHistoryController struct {
	ContractController
}

type NftPurchaseInfo struct {
	NftLdefIndex string `json:"nftLdefIndex"`
	Buyer string `json:"buyer"`
	Seller string `json:"seller"`
	TransactionAddress string `json:"transactionAddress"`
}

type MarketHistoryResponse struct {
	NftPurchaseInfo []*NftPurchaseInfo `json:"nftPurchaseInfo"`
}

func (this *MarketTransactionHistoryController) MarketTransactionHistory() {
	nickname:= this.Ctx.Input.Param(":nickname")
	var purchaseHistory [] models.StorePurchaseHistroy
	o:=orm.NewOrm()
	cond:= orm.NewCondition()
	cond = cond.And("seller_nickname",nickname).Or("buyer_nickname",nickname)
	num,err:=o.QueryTable("store_purchase_histroy").
		SetCond(cond).
		All(&purchaseHistory,"buyer_nickname","seller_nickname","transaction_address","nft_ldef_index")
	if err!=nil {
		if err==orm.ErrNoRows {
			purchaseHistory = make([]models.StorePurchaseHistroy,0)
			logs.Error(err.Error())
		} else {
			logs.Error(err.Error())
			sendError(&this.Controller, err, 500)
			return
		}
	}

	nftPurchaseInfo:=make([]*NftPurchaseInfo,num)
	for i,v:=range purchaseHistory {
		ni:= &NftPurchaseInfo{
			NftLdefIndex: v.NftLdefIndex,
			Buyer:v.BuyerNickname,
			Seller:v.SellerNickname,
			TransactionAddress:v.TransactionAddress,
		}
		nftPurchaseInfo[i] = ni
	}
	logs.Debug("purchase history record of",nickname,"has",num,"record")
	res:=&MarketHistoryResponse{
		NftPurchaseInfo: nftPurchaseInfo,
	}
	this.Data["json"] = &res
	this.ServeJSON()
}
