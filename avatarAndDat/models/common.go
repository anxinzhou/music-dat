package models

import (
	"github.com/astaxie/beego/orm"
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
	//uuid:= "4298349238490234456sa"
	//nickname:= "AlphaBrain"
	//password:= "123456"
	//username:="alphaslottest"
	//intro:= "this is alphabrain "
	//for i:=0;i<num;i++ {
	//	postPrefix:=""
	//	if i!=0 {
	//		postPrefix= strconv.FormatInt(int64(i), 10)
	//	}
	//	uuid:=uuid+postPrefix
	//	nickname:=nickname+postPrefix
	//	username:=username+postPrefix
	//	intro:=intro+postPrefix
	//	createUser(uuid,nickname,username,password,intro,"0xaC39b311DCEb2A4b2f5d8461c1cdaF756F4F7Ae9")
	//}
	createUser("48320958456gfdgz","YulieSu","YulieSu","YulieSu2019","","0xaC39b311DCEb2A4b2f5d8461c1cdaF756F4F7Ae9")
	createUser("4832095812z6gfdc","Cassette","Cassette","Cassette2019","","0xaC39b311DCEb2A4b2f5d8461c1cdaF756F4F7Ae9")
	createUser("48320z581256gfdc","Kaze.P.C","Kaze.P.C","Kaze.P.C2019","","0xaC39b311DCEb2A4b2f5d8461c1cdaF756F4F7Ae9")
	createUser("58320z581256gfdc","LazyCat","LazyCat","12345678","","0xaC39b311DCEb2A4b2f5d8461c1cdaF756F4F7Ae9")
}

func createUser(uuid , nickname ,  username,password, intro, wallet string)  {
	o:=orm.NewOrm()
	userInfo:= UserInfo{
		Uuid:uuid,
		Nickname:nickname,
		AvatarFileName:"",
		Intro: intro,
	}
	err:=o.Read(&userInfo,"uuid")
	if err!=nil && err!=orm.ErrNoRows {
		o.Rollback()
		panic(err)
	}
	if err == orm.ErrNoRows {
		_,err:=o.Insert(&userInfo)
		if err!=nil {
			o.Rollback()
			panic(err)
		}
	}
	creatorInfo:= CreatorInfo{
		Uuid: uuid,
		Username: username,
		Password: password,
	}
	err=o.Read(&creatorInfo,"uuid")
	if err!=nil && err!=orm.ErrNoRows {
		o.Rollback()
		panic(err)
	}
	if err == orm.ErrNoRows {
		creatorInfo.UserInfo = &UserInfo {
			Uuid:uuid,
		}
		_,err:=o.Insert(&creatorInfo)
		if err!=nil {
			o.Rollback()
			panic(err)
		}
	}
	mkInfo:= UserMarketInfo{
		Uuid: uuid,
		Wallet: "0xaC39b311DCEb2A4b2f5d8461c1cdaF756F4F7Ae9",
		Count: 0,
	}
	err=o.Read(&mkInfo,"uuid")
	if err!=nil && err!=orm.ErrNoRows {
		o.Rollback()
		panic(err)
	}
	if err == orm.ErrNoRows {
		mkInfo.UserInfo = &UserInfo {
			Uuid:uuid,
		}
		_,err:=o.Insert(&mkInfo)
		if err!=nil {
			o.Rollback()
			panic(err)
		}
	}
}