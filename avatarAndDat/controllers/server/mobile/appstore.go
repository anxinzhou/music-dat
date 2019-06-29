package mobile

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

// token purchase
func (m *Manager) TokenBuyPaidHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req TokenPurchaseRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	actionStatus := req.ActionStatus
	nickname:= req.Nickname
	o := orm.NewOrm()
	if actionStatus == ACTION_STATUS_FINISH {
		purchaseInfo := models.BerryPurchaseTable{
			TransactionId: req.TransactionId,
		}

		o.Begin()
		err = o.Read(&purchaseInfo)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		if purchaseInfo.Status != ACTION_STATUS_PENDING {
			err := errors.New("action in wrong status")
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		purchaseInfo.AppTranId = req.AppTranId
		purchaseInfo.Status = ACTION_STATUS_FINISH
		_, err = o.Update(&purchaseInfo)
		if err != nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}

		// update coin records
		col := models.MongoDB.Collection("users")
		// update coin records
		type fields struct {
			Coin string `bson:"coin"`
		}

		filter:= bson.M{
			"nickname": nickname,
		}

		var queryResult fields

		err = col.FindOne(context.Background(), filter, options.FindOne().SetProjection(bson.M{
			"coin": true,
		})).Decode(&queryResult)
		if err != nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		logs.Debug("nickname", req.Nickname, "coin number:", queryResult.Coin)

		currentBalance, err := strconv.Atoi(queryResult.Coin)
		if err != nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		amount := req.Amount
		update := bson.M{
			"$set": bson.M{"coin": strconv.Itoa(amount + currentBalance)},
		}
		_, err = col.UpdateOne(context.Background(), filter, update)
		if err != nil {
			o.Rollback()
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		logs.Info("update success", "after update, amount:", amount+currentBalance)

		o.Commit()
		logs.Info("insert one record to purchase table")
		m.wrapperAndSend(c, bq, &TokenPurchaseResponse{
			RQBaseInfo:   *bq,
			ActionStatus: ACTION_STATUS_FINISH,
		})
		return
	} else if actionStatus == ACTION_STATUS_PENDING {
		appTranIdBytes := make([]byte, 32)
		_, err := rand.Read(appTranIdBytes)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		appTranId := hex.EncodeToString(appTranIdBytes)
		purchaseInfo := models.BerryPurchaseTable{
			TransactionId: appTranId,
			BuyerNickname:  nickname,
			NumPurchased:  req.Amount,
			AppId:         req.AppId,
			Status:        ACTION_STATUS_PENDING,
		}
		_, err = o.Insert(&purchaseInfo)
		if err != nil {
			logs.Emergency(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		m.wrapperAndSend(c, bq, &TokenPurchaseResponse{
			RQBaseInfo:    *bq,
			ActionStatus:  ACTION_STATUS_PENDING,
			TransactionId: appTranId,
		})
		return
	} else {
		err := errors.New("unknow action status")
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
}