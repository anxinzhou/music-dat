package http

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/nfnt/resize"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

type AdminController struct {
	beego.Controller
}

type LoginInfo struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type AdminRequest struct {
	LoginType int `json:"loginType"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminResponse struct {
	AvatarUrl string `json:"avatarUrl"`
	NickName string `json:"nickName"`
	AccessToken string `json:"accessToken"`
}

func (this *AdminController) Login() {
	var req AdminRequest
	data, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,500)
		return
	}
	err=json.Unmarshal(data,&req)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	logs.Debug("username",req.Username)
	loginType:=req.LoginType

	logs.Debug("loginType",loginType)

	// verify if user is in db
	type usersTablefields struct {
		AvatarUrl string `bson:"avatar_url"`
		NickName string `bson:"nickname"`
	}
	col:=models.MongoDB.Collection("users")
	var res AdminResponse
	if loginType == LOGIN_TYPE_USERNAME {
		username := req.Username
		password := req.Password
		filter:= bson.M {
			"username": username,
			"password": password,
		}
		var queryResult usersTablefields
		err:= col.FindOne(context.Background(),filter,options.FindOne().SetProjection(bson.M{
			"username": true,
			"avatar_url":true,
			"nickname":true,
		})).Decode(&queryResult)
		if err!=nil {
			if err == mongo.ErrNoDocuments {
				// no such user
				logs.Error(err.Error())
				sendError(&this.Controller,errors.New("no such user"),401)
				return
			} else {
				logs.Error(err.Error())
				sendError(&this.Controller, errors.New("no such user"), 500)
				return
			}
		}

		res.AvatarUrl = queryResult.AvatarUrl
		res.NickName = queryResult.NickName
		fileName:= UserIconPathFromUserName(username)
		if _,err:=os.Stat(fileName); os.IsNotExist(err) {
			logs.Info("creating user icon file in local")
			logs.Debug("avatar url",res.AvatarUrl)
			response,err:=http.Get(res.AvatarUrl)
			defer response.Body.Close()
			if err!=nil {
				logs.Error(err.Error())
				sendError(&this.Controller,err,500)
				return
			}

			originImage,err:= png.Decode(response.Body)
			if err!=nil {
				logs.Error(err.Error())
				sendError(&this.Controller,err,500)
				return
			}
			newImage:= resize.Resize(100,100,originImage,resize.Lanczos3)
			var filePath string
			filePath = path.Join(BASE_FILE_PATH,PATH_KIND_USER_ICON,fileName)
			out, err:= os.Create(filePath)
			defer out.Close()
			if err!=nil {
				logs.Error(err.Error())
				sendError(&this.Controller,err,500)
				return
			}
			err = jpeg.Encode(out,newImage,nil)
			if err!=nil {
				logs.Error(err.Error())
				sendError(&this.Controller,err,500)
				return
			}
		}
	} else {
		err:= errors.New("unsupported login type")
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}

	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().UnixNano()|rand.Int63(), 10))
	accessToken:= hex.EncodeToString(h.Sum(nil))
	res.AccessToken = accessToken
	this.Data["json"] = &res
	this.ServeJSON()
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	logs.Debug("user login successful")
	return
}

type ImportWalletController struct {
	beego.Controller
}

type ImportWalletRequest struct {
	Username string `json:"username"`
	WalletId string `json:"walletId"`
}

func (this *ImportWalletController) ImportWallet() {
	var req ImportWalletRequest
	data, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	err=json.Unmarshal(data,&req)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}

	username:= req.Username
	walletId:= req.WalletId
	iconFileName:= UserIconPathFromUserName(username)
	// user file should have been saved at local this time
	if _,err:=os.Stat(path.Join(BASE_FILE_PATH,PATH_KIND_USER_ICON,iconFileName)); os.IsNotExist(err) {
		logs.Error(err.Error())
		sendError(&this.Controller,err,500)
		return
	}
	logs.Debug("icon file path",iconFileName)
	walletInfo:= &models.MarketUserTable{
		WalletId: walletId,
		Count: 0,
		Username: username,
		UserIconUrl: iconFileName,
	}

	models.O.Begin()         //TODO single sql
	err = models.O.Read(walletInfo)
	if err != nil {
		if err!= orm.ErrNoRows {
			models.O.Rollback()
			logs.Error(err.Error())
			sendError(&this.Controller,err,500)
			return
		} else {
			_,err:=models.O.Insert(walletInfo)
			if err!=nil {
				models.O.Rollback()
				logs.Error(err.Error())
				sendError(&this.Controller,err,500)
				return
			}
		}
	}
	models.O.Commit()
	logs.Info("insert to market user table","username",username)
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	return
}
