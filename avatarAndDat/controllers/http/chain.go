package http

import (
	"github.com/astaxie/beego/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xxRanger/blockchainUtil/contract/nft"
)

type NftBalanceController struct {
	ContractController
}

type nftBalanceResponse struct {
	Count int `json:"count"`
}

func (this *NftBalanceController) Get() {
	user := this.Ctx.Input.Param(":user")
	nftContract:= this.C.smartContract.(*nft.NFT)
	count,err:= nftContract.BalanceOf(common.HexToAddress(user))
	if err!=nil {
		logs.Error(err.Error())
		sendError(this,err,500)
		return
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"]= &nftBalanceResponse{
		Count: int(count.Int64()),
	}
}

type NftListController struct {
	ContractController
}

type nftInfoListRes struct {
	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
	NftName string `json:"nftName"`
	NftValue int `json:"nftValue" orm:"column(price)"`
	ActiveTicker string `json:"activeTicker"`
	NftLifeIndex int64 `json:"nftLifeIndex"`
	NftPowerIndex int64 `json:"nftPowerIndex"`
	NftLdefIndex string `json:"nftLdefIndex"`
	NftCharacId string `json:"nftCharacId"`
	Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
	Qty int `json:"qty"`
}

type nftListResponse struct {
	NftTranData []*nftInfoListRes `json:"nftTranData"`
}

func (this *NftListController) Get() {
	//user := this.Ctx.Input.Param(":user")
	//kind := this.Ctx.Input.Param(":kind")
	//logs.Debug("user",user,"query nft list")
	//nftContract:=this.C.smartContract.(*nft.NFT)
	//nftList,err:= nftContract.TokensOfUser(common.HexToAddress(user))
	//if err!=nil {
	//	logs.Error(err.Error())
	//	sendError(this,err,500)
	//	return
	//}
	//
	//nftLdefIndexs:=make([]string,len(nftList))
	//for i,tokenId:= range nftList {
	//	ldef,err:= nftContract.LdefIndexOfToken(tokenId)
	//	if err!=nil {
	//		logs.Error(err.Error())
	//		sendError(this,err,500)
	//		return
	//	}
	//	nftLdefIndexs[i] = ldef
	//}
	//nftTranResponseData:= make([]*nftInfoListRes,len(nftList))
	//
	//r := models.O.Raw(`
	//	select ni.nft_type, ni.nft_name,
	//	mk.price,mk.active_ticker, ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
	//	ni.nft_charac_id,mp.file_name,mk.qty from
	//	nft_market_table as mk,
	//	nft_mapping_table as mp,
	//	nft_info_table as ni,
	//	nft_item_admin as na
	//	where mk.nft_ldef_index = mp.nft_ldef_index and mk.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id and  ni.nft_ldef_index = ? `, nftLdefIndex)
	//var nftResponseInfo ItemDetailsResponseNftInfo
	//err = r.QueryRow(&nftResponseInfo)
	//if err!=nil {
	//	logs.Error(err.Error())
	//	m.errorHandler(c, bq, err)
	//	return
	//}
	//
	//res:= &nftListResponse{
	//	NftTranData: nftTranResponseData,
	//}
	//this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	//this.Data["json"]= &res
}