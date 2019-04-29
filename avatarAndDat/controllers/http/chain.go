package http

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
)

type NftBalanceController struct {
	ContractController
}

type nftBalanceResponse struct {
	Count int `json:"count"`
}

func (this *NftBalanceController) Get() {
	user := this.Ctx.Input.Param(":user")
	nftContract,ok := this.C.smartContract.(*nft.NFT)
	logs.Debug("contract address",nftContract.Address())
	if !ok {
		err:= errors.New("can not convert smart contract")
		sendError(this,err,500)
		return
	}
	logs.Debug("user",user,"query balance")
	count,err:= nftContract.BalanceOf(common.HexToAddress(user))
	if err!=nil {
		logs.Error(err.Error())
		sendError(this,err,500)
		return
	}
	logs.Debug("balance of user",count)
	this.Data["json"]= &nftBalanceResponse{
		Count: int(count.Int64()),
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.ServeJSON()
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
	ShortDesc string `json:"shortDesc" orm:"column(short_description)"`
	LongDesc string `json:"longDesc" orm:"column(long_description)"`
	Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
	Qty int `json:"qty"`
}

type nftListResponse struct {
	NftTranData []*nftInfoListRes `json:"nftTranData"`
}

func (this *NftListController) Get() {
	user := this.Ctx.Input.Param(":user")
	logs.Debug("user",user,"query nft list")
	nftContract:=this.C.smartContract.(*nft.NFT)
	nftList,err:= nftContract.TokensOfUser(common.HexToAddress(user))
	if err!=nil {
		logs.Error(err.Error())
		sendError(this,err,500)
		return
	}

	nftLdefIndexs:=make([]string,len(nftList))
	for i,tokenId:= range nftList {
		ldef,err:= nftContract.LdefIndexOfToken(tokenId)
		if err!=nil {
			logs.Error(err.Error())
			sendError(this,err,500)
			return
		}
		//logs.Info("ldefIndex",ldef)
		nftLdefIndexs[i] = ldef
	}
	nftTranResponseData:= make([]*nftInfoListRes,0,len(nftList))

	for _,nftLdefIndex:= range nftLdefIndexs {
		r := models.O.Raw(`
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
		if err!=nil {
			if err==orm.ErrNoRows {
				logs.Debug(err.Error())
				continue
			} else {
				logs.Error(err.Error())
				sendError(this,err,500)
				return
			}
		}

		var thumbnail string
		if nftResponseInfo.SupportedType == TYPE_NFT_MUSIC {   // music
			thumbnail = beego.AppConfig.String("hostaddr")+ ":"+
				beego.AppConfig.String("fileport") + "/resource/market/dat/"
		} else if nftResponseInfo.SupportedType == TYPE_NFT_AVATAR {  //avatar
			thumbnail = beego.AppConfig.String("hostaddr")+ ":"+
				beego.AppConfig.String("fileport") + "/resource/market/avatar/"
		} else {
			err := errors.New("unknown supported type")
			logs.Error(err.Error())
			sendError(this,err,400)
			return
		}
		nftResponseInfo.Thumbnail = thumbnail + nftResponseInfo.Thumbnail
		nftTranResponseData = append(nftTranResponseData,&nftResponseInfo)
	}

	res:= &nftListResponse{
		NftTranData: nftTranResponseData,
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"]= &res
	this.ServeJSON()
}