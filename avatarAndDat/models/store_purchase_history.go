package models

import "time"

type StorePurchaseHistroy struct {
	PurchaseId string `orm:"pk;unique"`
	AsId string
	TransactionAddress string
	NftName string
	TotalPaid int
	NftLdefIndex string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	Status int
}