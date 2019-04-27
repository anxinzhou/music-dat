package models

import "time"

type BerryPurchaseTable struct {
	TransactionId string `orm:"pk;unique"`
	RefillAsId string
	NumPurchased int
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	AppTranId string
	AppId string
	Status int
}
