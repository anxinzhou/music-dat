package models

type NftMarketTable struct {
	NftLdefIndex string `orm:"pk;unique"`
	MpId string
	NftAdminId string
	Price float64 `orm:"digits(12);decimals(4)"`
	Qty int
	NumSold int
	Active bool
	ActiveTicker string
}

