package models

import "time"

type NftMarketInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	SellerWallet string
	SellerUuid string
	Price int
	Qty int
	NumSold int
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	NftInfoTable *NftInfo `orm:"rel(one);on_delete(cascade);"`
}

func (this *NftMarketInfo) TableIndex() [][]string {
	return [][]string {
		[]string{"SellerWalletId"},
		[]string{"SellerNickname"},
	}
}

type DatNftMarketInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	AllowAirdrop bool
	CreatorPercent float64
	LyricsWriterPercent float64
	SongComposerPercent float64
	PublisherPercent float64
	UserPercent float64
	NftMarketInfo *NftMarketInfo `orm:"rel(one);on_delete(cascade);"`
}

type AvatarNftMarketInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftMarketInfo *NftMarketInfo `orm:"rel(one);on_delete(cascade);"`
}

type OtherNftMarketInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftMarketInfo *NftMarketInfo `orm:"rel(one);on_delete(cascade);"`
}

type NftMarketPlace struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftType string
	MpId string
	Active bool
	ActiveTicker string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	NftMarketInfo *NftMarketInfo `orm:"rel(one);on_delete(cascade);"`
}

func (NftMarketPlace *NftMarketPlace) TableIndex() [][]string {
	return [][]string {
		[]string {"NftType"},
	}
}

type NftShoppingCart struct {
	Id    int    `orm:"auto;pk;"`
	NftLdefIndex string
	Uuid string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	NftInfo *NftInfo `orm:"rel(fk);"`
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
		[]string{"Nickname","NftLdefIndex"},
	}
}
