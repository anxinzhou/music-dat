package models

type NftItemAdmin struct {
	NftAdminID string `orm:"pk;unique"`
	ShortDescription string
	LongDescription string
	NumDistribution int
}
