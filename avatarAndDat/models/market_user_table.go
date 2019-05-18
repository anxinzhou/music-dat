package models

type MarketUserTable struct {
	WalletId string `orm:"pk;unique"`
	Count int
	Username string
	UserIconUrl string
}

func (this *MarketUserTable) TableIndex() [][]string {
	return [][]string {
		[]string{"Username"},
	}
}