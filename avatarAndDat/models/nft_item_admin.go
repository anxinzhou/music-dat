package models

type NftItemAdmin struct {
	NftAdminId string `orm:"pk;unique"`
	ShortDescription string
	LongDescription string
	NumDistribution int
}
