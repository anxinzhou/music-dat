package http

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"math/big"
	"sync"
)

type NftBalanceController struct {
	ContractController
}

type nftBalanceResponse struct {
	Count int `json:"count"`
}

func (this *NftBalanceController) Get() {
	user := this.Ctx.Input.Param(":user")
	nftContract, ok := this.C.smartContract.(*nft.NFT)
	logs.Debug("contract address", nftContract.Address())
	if !ok {
		err := errors.New("can not convert smart contract")
		sendError(this, err, 500)
		return
	}
	logs.Debug("user", user, "query balance")
	count, err := nftContract.BalanceOf(common.HexToAddress(user))
	if err != nil {
		logs.Error(err.Error())
		sendError(this, err, 500)
		return
	}
	logs.Debug("balance of user", count)
	this.Data["json"] = &nftBalanceResponse{
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
	user := this.Ctx.Input.Param(":user")
	logs.Debug("user", user, "query nft list")
	nftContract := this.C.smartContract.(*nft.NFT)
	nftList, err := nftContract.TokensOfUser(common.HexToAddress(user))
	if err != nil {
		logs.Error(err.Error())
		sendError(this, err, 500)
		return
	}

	nftLdefIndexs := make([]string, len(nftList))
	for i, tokenId := range nftList {
		ldef, err := nftContract.LdefIndexOfToken(tokenId)
		if err != nil {
			logs.Error(err.Error())
			sendError(this, err, 500)
			return
		}
		//logs.Info("ldefIndex",ldef)
		nftLdefIndexs[i] = ldef
	}
	nftTranResponseData := make([]*nftInfoListRes, 0, len(nftList))

	for _, nftLdefIndex := range nftLdefIndexs {
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
		if err != nil {
			if err == orm.ErrNoRows {
				logs.Debug(err.Error())
				continue
			} else {
				logs.Error(err.Error())
				sendError(this, err, 500)
				return
			}
		}

		var thumbnail string
		if nftResponseInfo.SupportedType == TYPE_NFT_MUSIC { // music
			thumbnail = beego.AppConfig.String("prefix") + beego.AppConfig.String("hostaddr") + ":" +
				beego.AppConfig.String("fileport") + "/resource/market/dat/"
		} else if nftResponseInfo.SupportedType == TYPE_NFT_AVATAR { //avatar
			thumbnail = beego.AppConfig.String("prefix") + beego.AppConfig.String("hostaddr") + ":" +
				beego.AppConfig.String("fileport") + "/resource/market/avatar/"
		} else {
			err := errors.New("unknown supported type")
			logs.Error(err.Error())
			sendError(this, err, 400)
			return
		}
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

func (this *RewardController) RewardDat() {
	user := this.Ctx.Input.Param(":user")
	models.O.Begin()
	qs := models.O.QueryTable("nft_market_table").Filter("nft_ldef_index__contains","M").Limit(1)
	var mks []models.NftMarketTable
	rewardAccount := 1
	_, err := qs.Limit(rewardAccount).All(&mks)
	if err != nil {
		models.O.Rollback()
		logs.Error(err.Error())
		sendError(this, err, 500)
		return
	}
	var res nftListResponse
	nftInfoList := make([]*nftInfoListRes, len(mks))
	res.NftTranData = nftInfoList
	wg := sync.WaitGroup{}
	for i, mk := range mks {
		wg.Add(1)
		go func(i int, mk models.NftMarketTable) {
			defer wg.Done()
			nftLdefIndex := mk.NftLdefIndex

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
			if err != nil {
				models.O.Rollback()
				logs.Error(err.Error())
				sendError(this, err, 500)
				return
			}

			nftType:= nftResponseInfo.SupportedType
			var thumbnail string
			if nftType == TYPE_NFT_MUSIC { // music
				thumbnail = beego.AppConfig.String("prefix") + beego.AppConfig.String("hostaddr") + ":" +
					beego.AppConfig.String("fileport") + "/resource/default/"
			} else {
				models.O.Rollback()
				err:= errors.New("Unknown type")
				logs.Error(err.Error())
				sendError(this, err, 500)
				return
			}
			//nftResponseInfo.Thumbnail = thumbnail + nftResponseInfo.Thumbnail // TODO appending file name
			nftResponseInfo.Thumbnail = thumbnail + "music.png"

			nftInfoList[i] = &nftResponseInfo
			//_, err = models.O.Delete(&mk)  //TODO comment for testing
			//if err != nil {
			//	models.O.Rollback()
			//	logs.Error(err.Error())
			//	sendError(this, err, 500)
			//	return
			//}

			tokenId, _ := new(big.Int).SetString(nftLdefIndex[1:], 10)

			nftContract := this.C.smartContract.(*nft.NFT)
			ownerAddress, err := nftContract.OwnerOf(tokenId)
			if err != nil {
				models.O.Rollback()
				logs.Error(err.Error())
				sendError(this, err, 500)
				return
			}
			_, txErr := this.C.account.SendFunction2(nftContract,
				nil,
				nft.FuncDelegateTransfer,
				common.HexToAddress(ownerAddress),
				common.HexToAddress(user),
				tokenId)    // TODO redis to cache unsuccessful transaction
			err = <-txErr
			if err != nil {
				models.O.Rollback()
				logs.Error(err.Error())
				sendError(this, err, 500)
				return
			}
		}(i, mk)
	}
	wg.Wait()
	models.O.Commit()
	this.Data["json"] = res
	this.ServeJSON()
}
