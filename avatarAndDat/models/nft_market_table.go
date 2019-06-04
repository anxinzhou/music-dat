package models

type NftMarketTable struct {
	NftLdefIndex string `orm:"pk;unique"`
	SellerWalletId string
	SellerNickname string
	MpId string
	NftAdminId string
	Price int
	Qty int
	NumSold int
	Active bool
	ActiveTicker string
	AllowAirdrop bool
}

func (this *NftMarketTable) TableIndex() [][]string {
	return [][]string {
		[]string{"NftAdminId"},
		[]string{"SellerWalletId"},
		[]string{"SellerNickname"},
	}
}


