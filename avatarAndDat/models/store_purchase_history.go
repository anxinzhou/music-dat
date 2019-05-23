package models

import "time"

type StorePurchaseHistroy struct {
	PurchaseId string `orm:"pk;unique"`
	AsId string    // alpha user id
	WalletId string // alpha user wallet address
	TransactionAddress string
	NftName string
	ActiveTicker string
	TotalPaid int
	NftLdefIndex string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	Status int
}

func (this *StorePurchaseHistroy) TableIndex() [][]string {
	return [][]string {
		[]string {"AsId"},
		[]string {"WalletId"},
		[]string {"NftLdefIndex"},
		[]string {"TransactionAddress"},
	}
}

