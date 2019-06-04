package models

import "time"

type NftShoppingCart struct {
	Id    int    `orm:"auto"`
	NftLdefIndex string
	Nickname string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
}

func (this *NftShoppingCart) TableIndex() [][]string {
	return [][]string {
		[]string{"Nickname","NftLdefIndex"},
		[]string{"Nickname"},
		[]string{"NftLdefIndex",},
	}
}

func (this *NftShoppingCart) TableUnique() [][] string {
	return [][]string {
		[]string{"Nickname","NftLdefIndex"},
	}
}

