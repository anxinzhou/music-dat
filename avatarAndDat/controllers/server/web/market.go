package web

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
)

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
	CreatorPercent int `json:"creatorPercent"`
	LyricsWriterPercent int `json:"lyricsWriterPercent"`
	SongComposerPercent int `json:"songComposerPercent"`
	PublisherPercent int `json:"publisherPercent"`
	UserPercent int `json:"userPercent"`
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
		Filter("seller_nickname",nickname).OrderBy("-timestamp").
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
	for i:=len(mkInfos)-1;i>=0;i-- {
		mkInfo:= mkInfos[i]
		nftLdefIndex:= mkInfo.NftLdefIndex
		r := o.Raw(`
		select mk.creator_percent, mk.lyrics_writer_percent, mk.song_composer_percent,
		mk.publisher_percent, mk.user_percent, mk.price,mk.active_ticker, 
		ni.nft_type, ni.nft_name, ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
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
	nickname:= this.Ctx.Input.Param(":nickname")
	var purchaseHistory [] models.StorePurchaseHistroy
	o:=orm.NewOrm()
	cond:= orm.NewCondition()
	cond = cond.And("seller_nickname",nickname).Or("buyer_nickname",nickname)
	num,err:=o.QueryTable("store_purchase_histroy").
		SetCond(cond).
		All(&purchaseHistory,"buyer_nickname","seller_nickname","transaction_address","nft_ldef_index","timestamp")
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
			Time: chinaTimeFromTimeStamp(v.Timestamp),
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
