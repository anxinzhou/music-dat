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
