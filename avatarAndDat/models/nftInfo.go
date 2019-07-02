package models

type NftInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftType string
	NftName string
	ShortDescription string
	LongDescription string
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
	NftInfo * NftInfo `orm:"rel(one);on_delete(cascade);"`
}

type DatNftInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	IconFileName string
	NftInfo * NftInfo `orm:"rel(one);on_delete(cascade);"`
}

type OtherNftInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftInfo * NftInfo `orm:"rel(one);on_delete(cascade);"`
}

