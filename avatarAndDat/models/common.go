package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
)

func GetNftFullInfo(nftLdefIndex string, nftInfo interface{}) error {
	o:= orm.NewOrm()
	r := o.Raw(`
		select ni.nft_type, ni.nft_name, 
		mk.price,mk.active_ticker, mk.qty,
		ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,na.short_description, na.long_description,
		mp.file_name, mp.icon_file_name from 
		nft_market_table as mk,
		nft_mapping_table as mp,
		nft_info_table as ni,
		nft_item_admin as na
		where mk.nft_ldef_index = mp.nft_ldef_index and mk.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id and  ni.nft_ldef_index = ? `, nftLdefIndex)
	err := r.QueryRow(nftInfo)
	return err
}

func GenerateTestCreator(num int)  {
	uuid:= "429834923849023895sdfsdf08430594"
	nickname:= "AlphaBrain"
	password:= "123456"
	username:="alphaslottest"
	o:=orm.NewOrm()
	for i:=0;i<num;i++ {
		postPrefix:=""
		if i!=0 {
			postPrefix= uuid + strconv.FormatInt(int64(i), 10)
		}
		uuid+=postPrefix
		nickname+=postPrefix
		username+=postPrefix
		userInfo:= CreatorInfo{
			Uuid: uuid,
			Username: username,
			Password: password,
			Nickname: nickname,
		}
		err:= o.Read(&userInfo,"username")
		if err!=nil && err!=orm.ErrNoRows {
			panic(err)
		}
		if err == orm.ErrNoRows {
			_,err:=o.Insert(&userInfo)
			if err!=nil {
				panic(err)
			}
			logs.Info("uuid",uuid,"insert into creator info table")
		} else {
			logs.Info("uuid",uuid,"already created in creator info table")
		}

		mkInfo:= MarketUserTable{
			Uuid: uuid,
			WalletId: "0xaC39b311DCEb2A4b2f5d8461c1cdaF756F4F7Ae9",
			Count: 0,
			Nickname: nickname,
			UserIconUrl: "",
		}

		err = o.Read(&mkInfo,"uuid")
		if err!=nil && err!=orm.ErrNoRows {
			panic(err)
		}
		logs.Info("uuid",uuid,"already created")
		if err == orm.ErrNoRows {
			_,err:=o.Insert(&userInfo)
			if err!=nil {
				panic(err)
			}
			logs.Info("uuid",uuid,"insert into market user table")
		} else {
			logs.Info("uuid",uuid,"already created in market user table")
		}
	}
}
