package models

import "time"

type StorePurchaseHistroy struct {
	PurchaseID string `orm:"pk;unique"`
	ASID string
	TransactionAddress string
	NFTName string
	TotalPaid float64 `orm:"digits(12);decimals(4)"`
	NFTLdefIndex string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
}