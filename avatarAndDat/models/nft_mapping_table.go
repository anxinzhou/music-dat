package models

type NftMappingTable struct {
	NftLdefIndex string `orm:"pk;unique"`
	TypeId string
	FileName string
	Key string
	NftAdminId string
	NftParentLdef string
}

func (this *NftMappingTable) TableIndex() [][]string {
	return [][]string {
		[]string{"NftAdminId"},
		[]string{"NftParentLdef"},
	}
}

