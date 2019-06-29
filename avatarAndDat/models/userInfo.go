package models

import (
	"time"
)

type CreatorInfo struct {
	Uuid string `orm:"pk;unique"`
	Username string `orm:"unique"`
	Password string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	UserInfo *UserInfo `orm:"reverse(one)"`
}

func (this *CreatorInfo) TableIndex() [][]string {
	return [][]string {
		[]string{"Username","Password"},
		[]string{"Nickname"},
	}
}

type UserInfo struct {
	Uuid string `orm:"pk;"`
	Nickname string `orm:"unique"`
	AvatarFileName string
	Intro string
	UserMarketInfo *UserMarketInfo `orm:"reverse(one);"`
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	CreatorInfo *CreatorInfo `orm:"rel(one);on_delete(cascade);"`
}

type UserMarketInfo struct {
	Uuid string `orm:"pk;"`
	Wallet string
	Count int
	UserInfo *UserInfo `orm:"rel(one);on_delete(cascade);"`
}

type FollowTable struct {
	Id int `orm:"auto"`
	FolloweeUuid string
	FollowerUuid string
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
	NftType string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	Status int
	UserInfo *UserInfo `orm:"rel(fk);on_delete(cascade);"`
}

func (this *NftPurchaseInfo) TableIndex() [][]string {
	return [][]string {
		[]string {"Uuid"},
		[]string {"NftLdefIndex"},
		[]string {"TransactionAddress"},
		[]string {"SellerNickname"},
		[]string {"NftType"},
	}
}

