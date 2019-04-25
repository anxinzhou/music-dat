package models

type NftMappingTable struct {
	NftLdefIndex string `orm:"pk;unique"`
	TypeId string
	FileName string
	Key string
	NftAdminId string
}

