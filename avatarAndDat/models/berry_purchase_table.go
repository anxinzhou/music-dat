package models

import "time"

type BerryPurchaseTable struct {
	TransactionId string `orm:"pk;unique"`
	NumPurchased int
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	AppTranId string
	AppId string
	Status int
	BuyerNickname string
}

func (this *BerryPurchaseTable) TableIndex() [][]string {
	return [][]string {
		[]string{"AppTranId"},
		[]string{"BuyerNickname"},
	}
}
