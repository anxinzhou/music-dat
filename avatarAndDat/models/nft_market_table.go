package models

type NftMarketTable struct {
	NftLdefIndex string `orm:"pk;unique"`
	MpId string
	NftAdminId string
	Price int
	Qty int
	NumSold int
	Active bool
	ActiveTicker string
}

