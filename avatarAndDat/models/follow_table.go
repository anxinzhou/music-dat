package models

type FollowTable struct {
	Id int `orm:"auto"`
	FolloweeNickname string
	FollowerNickname string
}

func (this *FollowTable) TableIndex() [][]string {
	return [][]string {
		[]string {"FolloweeNickName"},
		[]string {"FollowerNickName"},
		[]string {"FolloweeNickName","FollowerNickName"},
	}
}

func (this *FollowTable) TableUnique() [][]string {
	return [][]string {
		[]string {"FolloweeNickName","FollowerNickName"},
	}
}