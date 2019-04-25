package models

type NftInfoTable struct {
	NftLdefIndex string `orm:"pk;unique"`
	NftType string
	NftName string
	DistIndex string
	NftLifeIndex int64
	NftPowerIndex int64
	NftCharacId string
	PublicKey string
}
