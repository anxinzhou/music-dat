package http

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/jameskeane/bcrypt"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"io/ioutil"
	"math/rand"
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
	Nickname string `json:"nickname"`
	Uuid string `json:"uuid"`
	AccessToken string `json:"accessToken"`
	Address string `json:"address"`
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
		Username string `bson:"username"`
		Password string `bson:"password"`
		AvatarUrl string `bson:"avatar_url"`
		Nickname string `bson:"nickname"`
		Uuid string `bson:"uuid"`
	}
	col:=models.MongoDB.Collection("users")
	var res AdminResponse
	if loginType == LOGIN_TYPE_USERNAME {
		username := req.Username
		rawPassword := req.Password
		filter:= bson.M {
			"username": username,
		}
		var queryResult usersTablefields
		err:= col.FindOne(context.Background(),filter,options.FindOne().SetProjection(bson.M{
			"username": true,
			"password": true,
			"avatar_url":true,
			"nickname":true,
			"uuid": true,
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
		hashedPassword:= queryResult.Password

		// raw password to hased password
		h:= sha256.New()
		h.Write([]byte(rawPassword))
		password:=hex.EncodeToString(h.Sum(nil))
		logs.Debug("raw password",rawPassword)
		logs.Debug("password",password)
		logs.Debug("hashed password",hashedPassword)
		match:=bcrypt.Match(password,hashedPassword)
		if !match {
			err:=errors.New("wrong password")
			logs.Error(err.Error())
			sendError(&this.Controller, err, 401)
			return
		}

		nickname:= queryResult.Nickname
		res.Nickname = queryResult.Nickname
		userInfo:= models.MarketUserTable{
			Nickname:res.Nickname,
		}
		o:=orm.NewOrm()
		o.Begin()
		err=o.Read(&userInfo,"nickname")
		if err!=nil && err!=orm.ErrNoRows {
			o.Rollback()
			logs.Error(err.Error())
			sendError(&this.Controller, err, 500)
			return
		}
		var iconPath string
		if err == orm.ErrNoRows {
			walletInfo:= &models.MarketUserTable{
				WalletId: "",
				Count: 0,
				Nickname: nickname,
				UserIconUrl: "",
			}
			_,err:=o.Insert(walletInfo)
			if err!=nil {
				o.Rollback()
				logs.Error(err.Error())
				sendError(&this.Controller, err, 500)
				return
			}
			iconPath = ""
			res.Address= ""
		} else {
			iconPath = userInfo.UserIconUrl
			walletAddress:=userInfo.WalletId
			res.Address = walletAddress
		}
		if iconPath == "" {
			iconPath = "default.jpg"
		}
 		pathPrefix:=PathPrefixOfNFT("",PATH_KIND_USER_ICON)
		iconUrl:= pathPrefix+iconPath
		res.AvatarUrl = iconUrl
		o.Commit()
		//Get avatar from origin

		//res.AvatarUrl = queryResult.AvatarUrl
		//fileName:= UserIconPathFromNickname(res.Nickname)
		//if _,err:=os.Stat(fileName); os.IsNotExist(err) {
		//	logs.Info("creating user icon file in local")
		//	logs.Debug("avatar url",res.AvatarUrl)
		//	response,err:=http.Get(res.AvatarUrl)
		//	defer response.Body.Close()
		//	if err!=nil {
		//		logs.Error(err.Error())
		//		sendError(&this.Controller,err,500)
		//		return
		//	}
		//
		//	originImage,err:= png.Decode(response.Body)
		//	if err!=nil {
		//		logs.Error(err.Error())
		//		sendError(&this.Controller,err,500)
		//		return
		//	}
		//	newImage:= resize.Resize(100,100,originImage,resize.Lanczos3)
		//	var filePath string
		//	filePath = path.Join(BASE_FILE_PATH,PATH_KIND_USER_ICON,fileName)
		//	out, err:= os.Create(filePath)
		//	defer out.Close()
		//	if err!=nil {
		//		logs.Error(err.Error())
		//		sendError(&this.Controller,err,500)
		//		return
		//	}
		//	err = jpeg.Encode(out,newImage,nil)
		//	if err!=nil {
		//		logs.Error(err.Error())
		//		sendError(&this.Controller,err,500)
		//		return
		//	}
		//}
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
//
//type ImportWalletController struct {
//	beego.Controller
//}
//
//type ImportWalletRequest struct {
//	Nickname string `json:"nickname"`
//	WalletId string `json:"walletId"`
//}
//
//func (this *ImportWalletController) ImportWallet() {
//	var req ImportWalletRequest
//	data, err := ioutil.ReadAll(this.Ctx.Request.Body)
//	if err!=nil {
//		logs.Error(err.Error())
//		sendError(&this.Controller,err,400)
//		return
//	}
//	err=json.Unmarshal(data,&req)
//	if err!=nil {
//		logs.Error(err.Error())
//		sendError(&this.Controller,err,400)
//		return
//	}
//
//	nickname:= req.Nickname
//	walletId:= req.WalletId
//	iconFileName:= UserIconPathFromNickname(nickname)
//	// user file should have been saved at local this time
//	//if _,err:=os.Stat(path.Join(BASE_FILE_PATH,PATH_KIND_USER_ICON,iconFileName)); os.IsNotExist(err) {
//	//	logs.Error(err.Error())
//	//	sendError(&this.Controller,err,500)
//	//	return
//	//}
//	logs.Debug("icon file path",iconFileName)
//	walletInfo:= &models.MarketUserTable{
//		WalletId: walletId,
//		Count: 0,
//		Nickname: nickname,
//		UserIconUrl: iconFileName,
//	}
//	o:=orm.NewOrm()
//	o.Begin()         //TODO single sql
//	err = o.Read(walletInfo)
//	if err != nil {
//		if err!= orm.ErrNoRows {
//			o.Rollback()
//			logs.Error(err.Error())
//			sendError(&this.Controller,err,500)
//			return
//		} else {
//			_,err:=o.Insert(walletInfo)
//			if err!=nil {
//				o.Rollback()
//				logs.Error(err.Error())
//				sendError(&this.Controller,err,500)
//				return
//			}
//		}
//	}
//	o.Commit()
//	logs.Info("insert to market user table","nickname",nickname)
//	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
//	return
//}
