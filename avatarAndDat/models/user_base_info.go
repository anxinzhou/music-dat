package models

type UserBaseInfo struct {
	Uuid string `orm:"pk;unique"`
	Nickname string
}