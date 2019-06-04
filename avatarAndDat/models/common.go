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