package models

import "time"

type CreatorInfo struct {
	Username string `orm:"pk;unique"`
	Password string
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	Nickname string
}

func (this *CreatorInfo) TableIndex() [][]string {
	return [][]string {
		[]string{"Username","Password"},
		[]string{"Nickname"},
	}
}
