package models

import (
	"github.com/astaxie/beego/orm"
)

type MarketUserTable struct {
	WalletId string
	Count int
	Nickname string `orm:"pk;unique"`
	UserIconUrl string
}

func (this *MarketUserTable) TableIndex() [][]string {
	return [][]string {
		[]string{"WalletId"},
	}
}

func WalletIdOfNickname(nickname string) (string,error) {
	o:= orm.NewOrm()
	userMKInfo := MarketUserTable{
		Nickname: nickname,
	}
	err := o.Read(&userMKInfo)
	if err != nil {
		return "",err
	}
	return userMKInfo.WalletId,nil
}