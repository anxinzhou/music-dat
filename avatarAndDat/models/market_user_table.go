package models

type MarketUserTable struct {
	WalletId string `orm:"pk;unique"`
	Count int
}

