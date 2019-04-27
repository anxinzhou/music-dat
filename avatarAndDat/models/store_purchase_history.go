package models

import "time"

type StorePurchaseHistroy struct {
	PurchaseID string `orm:"pk;unique"`
	ASID string
	TransactionAddress string
	NFTName string
	TotalPaid int
	NFTLdefIndex string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	Status int
}