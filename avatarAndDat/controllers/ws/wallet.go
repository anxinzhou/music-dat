package ws

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *Manager) BindWalletHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req BindWalletRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	walletId:= req.WalletId
	nickname:= req.Nickname

	walletInfo:= &models.MarketUserTable{
		WalletId: walletId,
		Count: 0,
		Nickname: nickname,
		UserIconUrl: "",
	}

	o:=orm.NewOrm()
	o.Begin()         //TODO single sql
	err = o.Read(walletInfo)
	if err != nil {
		if err!= orm.ErrNoRows {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		} else {
			_,err:=o.Insert(walletInfo)
			if err!=nil {
				o.Rollback()
				logs.Error(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
		}
	}
	o.Commit()
	logs.Info("insert to market user table","nickname",nickname)
	m.wrapperAndSend(c,bq,&BindWalletResponse{
		RQBaseInfo: *bq,
	})
}

func isNicknameExist(nickname string) (bool,error) {
	type fields struct {
		Nickname string `bson:"nickname"`
	}
	filter:= bson.M {
		"nickname": nickname,
	}
	var queryResult fields
	col := models.MongoDB.Collection("users")
	err:=col.FindOne(context.Background(),filter,options.FindOne().SetProjection(bson.M{
		"nickname": true,
	})).Decode(&queryResult)
	if err!=nil {
		if err == mongo.ErrNoDocuments {
			return false,nil
		} else {
			return false,err
		}
	}
	return true,nil
}

func (m *Manager) SetNicknameHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req SetNicknameRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	uuid:= req.Uuid
	nickname:= req.Nickname
	duplicated,err:=isNicknameExist(nickname)
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	if duplicated {
		err:= errors.New("nick name already exists")
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	col := models.MongoDB.Collection("users")
	filter:= bson.M {
		"uuid": uuid,
	}
	update:= bson.M {
		"$set": bson.M {"nickname":nickname},
	}
	_,err=col.UpdateOne(context.Background(),filter,update)
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	m.wrapperAndSend(c,bq,&SetNicknameResponse{
		RQBaseInfo: *bq,
	})
}

func (m *Manager) IsNicknameDuplicatedHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req IsNicknameDuplicatedRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	nickname:= req.Nickname
	logs.Debug("nick name to test",nickname)
	duplicated,err:=isNicknameExist(nickname)
	if err!=nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	m.wrapperAndSend(c,bq,&IsNicknameDuplicatedResponse{
		RQBaseInfo: *bq,
		Duplicated: duplicated,
	})
}

func (m *Manager) FollowListHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req FollowListRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	nickname:=req.Nickname
	o:= orm.NewOrm()
	type queryInfo struct {
		FolloweeNickname string `orm:"column(followee_nickname)"`
		UserIconUrl string `orm:"column(user_icon_url)"`
	}
	var queryResult []queryInfo
	num,err:=o.Raw(`
		select followee_nickname,user_icon_url from market_user_table as mk, follow_table as ft
		where mk.nickname = ft.followee_nickname and ft.follower_nickname = ?`, nickname).QueryRows(&queryResult)
	if err!=nil {
		if err!= orm.ErrNoRows {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
	}

	followInfo:=make([]*FollowInfo,num)
	for i,_:= range queryResult {
		followInfo[i] = &FollowInfo{
			Nickname: queryResult[i].FolloweeNickname,
			Thumbnail: PathPrefixOfNFT("", PATH_KIND_USER_ICON)+queryResult[i].UserIconUrl,
		}
	}
	logs.Debug("user has",num,"followee")

	m.wrapperAndSend(c,bq,&FollowListResponse{
		RQBaseInfo: *bq,
		FollowInfo: followInfo,
	})
}

func (m*Manager) FollowListOperationHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req FollowListOperationRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	operation:=req.Operation
	followNickname:= req.FollowNickname
	nickname:= req.Nickname
	o:= orm.NewOrm()
	o.Begin()
	queryInfo:= models.FollowTable{
		FollowerNickname: nickname,
		FolloweeNickname: followNickname,
	}
	err=o.Read(&queryInfo,"followee_nickname","follower_nickname")
	if operation == FOLLOW_LIST_ADD {
		if err!=nil {
			if err== orm.ErrNoRows {
				_,err:=o.Insert(&queryInfo)
				if err!=nil {
					o.Rollback()
					logs.Error(err.Error())
					m.errorHandler(c, bq, err)
					return
				}
				logs.Warn("follow",followNickname)
			}  else  {
				o.Rollback()
				logs.Error(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
		} else {
			err:= errors.New("already follow "+followNickname)
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
	} else if operation == FOLLOW_LIST_DELETE {
		if err!=nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		} else {
			_,err:=o.Delete(&queryInfo,"followee_nickname","follower_nickname")
			if err!=nil {
				o.Rollback()
				logs.Error(err.Error())
				m.errorHandler(c, bq, err)
				return
			}
			logs.Warn("delete followee",followNickname)
		}
	} else {
		err:= errors.New("operation not exist")
		m.errorHandler(c, bq, err)
		return
	}
	o.Commit()
	m.wrapperAndSend(c,bq,&FollowListOperationResponse{
		RQBaseInfo: *bq,
		FollowNickname:followNickname,
		Operation: operation,
	})
}