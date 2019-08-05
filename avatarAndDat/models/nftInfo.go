package models

import "time"

type NftInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftType string
	NftName string
	ShortDesc string `orm:"short_desc;type(text)"`
	LongDesc string `orm:"long_desc;type(text)"`
	FileName string
	NftParentLdef string
	AvatarNftInfoNftLdefIndex *AvatarNftInfo`orm:"reverse(one)"`
	DatNftInfoNftLdefIndex *DatNftInfo `orm:"reverse(one)"`
	OtherNftInfoNftLdefIndex *OtherNftInfo `orm:"reverse(one)"`
	NftMarketInfoNftLdefIndex *NftMarketInfo `orm:"reverse(one)"`
}

type AvatarNftInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftLifeIndex int
	NftPowerIndex int
	NftInfo * NftInfo `orm:"rel(one);"`
}

type DatNftInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	MusicFileName string
	NftInfo * NftInfo `orm:"rel(one);"`
}

type OtherNftInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftInfo * NftInfo `orm:"rel(one);"`
}

type NftMarketInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	SellerWallet string
	SellerUuid string
	Price int
	Qty int
	NumSold int
	Active bool
	NftInfo *NftInfo `orm:"rel(one);"`
	AvatarNftMarketInfo *AvatarNftMarketInfo `orm:"reverse(one)"`
	DatNftMarketInfo *DatNftMarketInfo `orm:"reverse(one)"`
	OtherNftMarketInfo *OtherNftMarketInfo `orm:"reverse(one)"`
	NftMarketPlace *NftMarketPlace `orm:"reverse(one)"`
}

func (this *NftMarketInfo) TableIndex() [][]string {
	return [][]string {
		[]string{"SellerWallet"},
		[]string{"SellerUuid"},
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
	NftMarketInfo *NftMarketInfo `orm:"rel(one);"`
}

type AvatarNftMarketInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftMarketInfo *NftMarketInfo `orm:"rel(one);"`
}

type OtherNftMarketInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftMarketInfo *NftMarketInfo `orm:"rel(one);"`
}

type NftMarketPlace struct {
	NftLdefIndex string `orm:"pk;unique"`
	MpId string
	ActiveTicker string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	NftMarketInfo *NftMarketInfo `orm:"rel(one);"`
}
