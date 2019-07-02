package mobile
//
//import (
//	"encoding/json"
//	"errors"
//	"github.com/astaxie/beego/logs"
//	"github.com/astaxie/beego/orm"
//	"github.com/xxRanger/music-dat/avatarAndDat/models"
//	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
//)
//
//func (m *Manager) ShoppingCartListHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
//	var req ShoppingCartListRequest
//	err := json.Unmarshal(data, &req)
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//	nickname := req.Nickname
//	var shoppingCartHistory []models.NftShoppingCart
//	o := orm.NewOrm()
//	_, err = o.QueryTable("nft_shopping_cart").
//		Filter("nickname", nickname).
//		All(&shoppingCartHistory)
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//
//	shoppingCartRecordRes := make([]*ShoppingCartRecord, len(shoppingCartHistory))
//	for i, _ := range shoppingCartHistory {
//		nftLdefIndex := shoppingCartHistory[i].NftLdefIndex
//		logs.Debug("shopping card ldef index", nftLdefIndex)
//		r := o.Raw(`
//		select ni.nft_type, ni.nft_name,
//		mk.price,mk.active_ticker, mk.qty,
//		ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
//		ni.nft_charac_id,na.short_description, na.long_description,
//		mp.file_name, mp.icon_file_name from
//		nft_market_table as mk,
//		nft_mapping_table as mp,
//		nft_info_table as ni,
//		nft_item_admin as na
//		where mk.nft_ldef_index = mp.nft_ldef_index and mk.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id and  ni.nft_ldef_index = ? `, nftLdefIndex)
//		var nftInfo NFTInfo
//		err = r.QueryRow(&nftInfo)
//		if err != nil {
//			if err == orm.ErrNoRows {
//				logs.Info("item not exist in market")
//				continue
//			} else {
//				logs.Error(err.Error())
//				m.errorHandler(c, bq, err)
//				return
//			}
//		}
//
//		nftResInfo, err := nftResInfoFromNftInfo(&nftInfo)
//		if err != nil {
//			logs.Error(err.Error())
//			m.errorHandler(c, bq, err)
//			return
//		}
//		logs.Debug("origin time", shoppingCartHistory[i].Timestamp)
//		shoppingCartRecordRes[i] = &ShoppingCartRecord{
//			Timestamp:   chinaTimeFromTimeStamp(shoppingCartHistory[i].Timestamp),
//			NftTranData: nftResInfo,
//		}
//	}
//
//	res := &ShoppingCartListResponse{
//		RQBaseInfo: *bq,
//		NftList:    shoppingCartRecordRes,
//	}
//	m.wrapperAndSend(c, bq, res)
//}
//
//
//func (m *Manager) ShoppingCartChangeHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
//	var req ShoppingCartChangeRequest
//	err := json.Unmarshal(data, &req)
//	if err != nil {
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//
//	nickname := req.Nickname
//	operation := req.Operation
//	// check operation
//	if operation != SHOPPING_CART_ADD && operation != SHOPPING_CART_DELETE {
//		err := errors.New("unknown shopping cart operation")
//		logs.Error(err.Error())
//		m.errorHandler(c, bq, err)
//		return
//	}
//
//	nftList := req.NFTList
//	o := orm.NewOrm()
//	o.Begin()
//	for _, nftLdefIndex := range nftList {
//		shoppingCartRecord := models.NftShoppingCart{
//			NftLdefIndex: nftLdefIndex,
//			Nickname:     nickname,
//		}
//		if operation == SHOPPING_CART_ADD {
//			_, err := o.Insert(&shoppingCartRecord)
//			if err != nil {
//				o.Rollback()
//				logs.Error(err.Error())
//				m.errorHandler(c, bq, err)
//				return
//			}
//			logs.Debug("insert ", nftLdefIndex, "success")
//		} else if operation == SHOPPING_CART_DELETE {
//			err:=o.Read(&shoppingCartRecord,"nft_ldef_index", "nickname")
//			if err!=nil {
//				o.Rollback()
//				logs.Error(err.Error())
//				m.errorHandler(c, bq, err)
//				return
//			}
//			_, err = o.Delete(&shoppingCartRecord, "nft_ldef_index", "nickname")
//			if err != nil {
//				o.Rollback()
//				logs.Error(err.Error())
//				m.errorHandler(c, bq, err)
//				return
//			}
//			logs.Debug("delete ", nftLdefIndex, "success")
//		}
//	}
//	err=o.Commit()
//
//	m.wrapperAndSend(c, bq, &ShoppingCartChangeResponse{
//		RQBaseInfo: *bq,
//	})
//}