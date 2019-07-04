package models

import (
	"time"
)

type CreatorInfo struct {
	Uuid string `orm:"pk;unique"`
	Username string `orm:"unique"`
	Password string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	UserInfo *UserInfo `orm:"rel(one);on_delete(cascade);"`
}

func (this *CreatorInfo) TableIndex() [][]string {
	return [][]string {
		[]string{"Username","Password"},
	}
}

type UserInfo struct {
	Uuid string `orm:"pk;unique"`
	Nickname string `orm:"unique"`
	AvatarFileName string
	Intro string
	UserMarketInfo *UserMarketInfo `orm:"reverse(one);"`
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	CreatorInfo *CreatorInfo `orm:"reverse(one)"`
}

func (*UserInfo) TableIndex() [][]string {
	return [][]string {
		[]string{"nickname"},
	}
}

type UserMarketInfo struct {
	Uuid string `orm:"pk;"`
	Wallet string
	Count int
	UserInfo *UserInfo `orm:"rel(one);on_delete(cascade);"`
}

type FollowTable struct {
	Id int `orm:"pk;auto"`
	FollowerUuid string
	FolloweeUuid string
	Timestamp time.Time
}

func (this *FollowTable) TableIndex() [][]string {
	return [][]string {
		[]string {"FolloweeUuid"},
		[]string {"FollowerUuid"},
		[]string {"FollowerUuid","FolloweeUuid"},
	}
}

func (this *FollowTable) TableUnique() [][]string {
	return [][]string {
		[]string {"FolloweeUuid","FollowerUuid"},
	}
}

type BerryPurchaseInfo struct {
	TransactionId string `orm:"pk;unique"`
	NumPurchased int
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	AppTranId string
	AppId string
	Status int
	Uuid string `orm:"rel(one);on_delete(cascade);"`
}

func (this *BerryPurchaseInfo) TableIndex() [][]string {
	return [][]string {
		[]string{"AppTranId"},
		[]string{"Uuid"},
	}
}

type NftPurchaseInfo struct {
	PurchaseId string `orm:"pk;unique"`
	Uuid string
	SellerUuid string
	TransactionAddress string
	ActiveTicker string
	TotalPaid int
	NftLdefIndex string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	Status int
	UserInfo *UserInfo `orm:"rel(fk);on_delete(cascade);"`
}

func (this *NftPurchaseInfo) TableIndex() [][]string {
	return [][]string {
		[]string {"Uuid"},
		[]string {"NftLdefIndex"},
		[]string {"SellerUuid"},
	}
}

type NftShoppingCart struct {
	Id    int    `orm:"auto;pk;"`
	NftLdefIndex string
	Uuid string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	NftMarketPlace *NftMarketPlace `orm:"rel(fk);on_delete(cascade);"`
	UserInfo *UserInfo `orm:"rel(fk);on_delete(cascade);"`
}

func (this *NftShoppingCart) TableIndex() [][]string {
	return [][]string {
		[]string{"Uuid"},
		[]string{"Uuid","NftLdefIndex"},
	}
}

func (this *NftShoppingCart) TableUnique() [][] string {
	return [][]string {
		[]string{"Uuid","NftLdefIndex"},
	}
}