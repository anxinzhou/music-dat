package models

type NftInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftType string
	NftName string
	ShortDescription string
	LongDescription string
	FileName string
	NftParentLdef string
}

type AvatarNftInfo struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftLifeIndex int64
	NftPowerIndex int64
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

