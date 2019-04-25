package ws

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
)

func (m *Manager) errorHandler(c *client.Client, bq *RQBaseInfo, err error) {
	bq.Event = "failed"
	res := &ErrorResponse{
		RQBaseInfo: *bq,
		Reason:     err.Error(),
	}
	resWrapper, err := json.Marshal(res)
	if err != nil {
		panic(err)
		return
	}
	c.Send(resWrapper)
}

func (m *Manager) wrapperAndSend(c *client.Client, bq *RQBaseInfo, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	c.Send(data)
}

func (m *Manager) GetMPList(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req MpListRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	//page := req.Page
	//offset := req.Offset

	nftType := req.SupportedType
	logs.Info("nft type",nftType)
	// TODO can use prepare to optimize query
	r := models.O.Raw(`
		select ni.nft_type, ni.nft_name,
		mk.price,mk.active_ticker, ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,  mp.file_name, mk.qty 
		from nft_market_table as mk, nft_mapping_table as mp,
		nft_info_table as ni where mk.nft_ldef_index = mp.nft_ldef_index 
		and mk.nft_ldef_index = ni.nft_ldef_index 
		and ni.nft_type = ? `, nftType)

	var nftInfos []NFTInfo
	r.QueryRows(&nftInfos)
	length:= len(nftInfos)
	nis:=make([]*NFTInfo,length)

	var thumbnail string
	if nftType == "721-04" {   // music
		thumbnail = beego.AppConfig.String("httpaddr")+ ":"+
			beego.AppConfig.String("httpport") + "/resource/"
	} else if nftType == "721-02" {  //avatar
		thumbnail = beego.AppConfig.String("httpaddr")+ ":"+
			beego.AppConfig.String("httpport") + "/resource/"
	} else {
		err := errors.New("unknown supported type")
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	for i:=0;i<length;i++ {
		logs.Info("thumbnail:",nftInfos[i].Thumbnail)
		nftInfos[i].Thumbnail = thumbnail + nftInfos[i].Thumbnail
		nis[i] = &nftInfos[i]
	}

	m.wrapperAndSend(c,bq,&MpListResponse{
		RQBaseInfo: *bq,
		NftData: nis,
	})
}
