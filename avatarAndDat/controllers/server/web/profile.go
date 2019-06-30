package web

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/util"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"io/ioutil"
	"os"
	"path"
)

type NicknameController struct {
	beego.Controller
}

func (this *NicknameController) GetNickname() {
	uuid:= this.Ctx.Input.Param(":uuid")
	var userInfo models.UserInfo
	o:=orm.NewOrm()
	err:=o.QueryTable("user_info").Filter("uuid",uuid).
		One(&userInfo,"nickname")
	if err!=nil {
		if err==orm.ErrNoRows {
			err := errors.New("no such user")
			sendError(&this.Controller,err,404)
			return
		} else {
			sendError(&this.Controller,err,500)
			return
		}
	}
	nickname:= userInfo.Nickname
	type response struct {
		Nickname string `json:"nickname"`
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &response{
		Nickname: nickname,
	}
	this.ServeJSON()
}

func (this *NicknameController) SetNickname() {
	uuid:= this.Ctx.Input.Param(":uuid")
	type request struct {
		Nickname string `json:"nickname"`
	}
	var req request
	data, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err!=nil {
		logs.Error(err.Error())
		err:= errors.New("can not parse request")
		sendError(&this.Controller,err,400)
		return
	}
	err=json.Unmarshal(data,&req)
	if err!=nil {
		logs.Error(err.Error())
		err:= errors.New("data format error")
		sendError(&this.Controller,err,400)
		return
	}

	nickname:= req.Nickname
	userInfo:=models.UserInfo {
		Uuid: uuid,
		Nickname: nickname,
	}
	o:=orm.NewOrm()
	o.Begin()
	tmpInfo:= models.UserInfo{
		Nickname: nickname,
	}
	err=o.Read(&tmpInfo,"nickname")
	if err!=nil {
		if err!=orm.ErrNoRows {
			logs.Error("unknown error when query db")
			sendError(&this.Controller,err,500)
			return
		}
	}
	if err == nil && tmpInfo.Uuid != uuid {
		err:=errors.New("duplicate nickname")
		logs.Error(err.Error())
		sendError(&this.Controller,err,500)
		return
	}
	_,err=o.Update(&userInfo,"nickname")
	if err!=nil {
		if err==orm.ErrNoRows {
			err := errors.New("no such user")
			sendError(&this.Controller,err,404)
			return
		} else {
			sendError(&this.Controller,err,500)
			return
		}
	}
	logs.Info("update",uuid,"intro")
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
}

type IntroController struct {
	beego.Controller
}

func (this *IntroController) GetIntro() {
	uuid:= this.Ctx.Input.Param(":uuid")
	var userInfo models.UserInfo
	o:=orm.NewOrm()
	err:=o.QueryTable("user_info").Filter("uuid",uuid).
		One(&userInfo,"intro")
	if err!=nil {
		if err==orm.ErrNoRows {
			err := errors.New("no such user")
			sendError(&this.Controller,err,404)
			return
		} else {
			sendError(&this.Controller,err,500)
			return
		}
	}
	intro:= userInfo.Intro
	type response struct {
		Intro string `json:"intro"`
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &response{
		Intro: intro,
	}
	this.ServeJSON()
}

func (this *IntroController) SetIntro() {
	uuid:= this.Ctx.Input.Param(":uuid")
	type request struct {
		Intro string `json:"intro"`
	}
	var req request
	data, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err!=nil {
		logs.Error(err.Error())
		err:= errors.New("can not parse request")
		sendError(&this.Controller,err,400)
		return
	}
	err=json.Unmarshal(data,&req)
	if err!=nil {
		logs.Error(err.Error())
		err:= errors.New("data format error")
		sendError(&this.Controller,err,400)
		return
	}

	intro:= req.Intro
	userInfo:=models.UserInfo {
		Uuid: uuid,
		Intro: intro,
	}
	o:=orm.NewOrm()
	_,err=o.Update(&userInfo,"intro")
	if err!=nil {
		if err==orm.ErrNoRows {
			err := errors.New("no such user")
			sendError(&this.Controller,err,404)
			return
		} else {
			sendError(&this.Controller,err,500)
			return
		}
	}
	logs.Info("update",uuid,"intro")
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
}

type WalletController struct {
	beego.Controller
}

func (this *WalletController) GetWallet() {
	uuid:= this.Ctx.Input.Param(":uuid")
	userInfo:=models.UserInfo {
		Uuid:uuid,
	}
	o:=orm.NewOrm()
	err:=o.Read(&userInfo)
	if err!=nil {
		if err==orm.ErrNoRows {
			err := errors.New("no such user")
			sendError(&this.Controller,err,404)
			return
		} else {
			sendError(&this.Controller,err,500)
			return
		}
	}
	if userInfo.UserMarketInfo!=nil {
		err:=o.Read(userInfo.UserMarketInfo)
		if err!=nil {
			sendError(&this.Controller,err,500)
			return
		}
	} else {
		err:= errors.New("user have not set wallet")
		sendError(&this.Controller,err,404)
		return
	}

	wallet:=userInfo.UserMarketInfo.Wallet
	count:= userInfo.UserMarketInfo.Count
	type response struct {
		Wallet string `json:"wallet"`
		Count int 	`json:"count"`
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &response{
		Wallet: wallet,
		Count: count,
	}
	this.ServeJSON()
}

func (this *WalletController) SetWallet() {
	type request struct {
		Wallet string `json:"wallet"`
	}
	uuid:= this.Ctx.Input.Param(":uuid")
	var req request
	data, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err!=nil {
		logs.Error(err.Error())
		err:= errors.New("can not parse request")
		sendError(&this.Controller,err,400)
		return
	}
	err=json.Unmarshal(data,&req)
	if err!=nil {
		logs.Error(err.Error())
		err:= errors.New("data format error")
		sendError(&this.Controller,err,400)
		return
	}

	o:=orm.NewOrm()
	userInfo:= models.UserInfo{
		Uuid: uuid,
	}
	err=o.Read(&userInfo)
	if err!=nil {
		if err== orm.ErrNoRows {
			err:= errors.New("user not exist")
			sendError(&this.Controller,err,400)
			return
		} else {
			logs.Error(err.Error())
			err:= errors.New("unknown error when query databse")
			sendError(&this.Controller,err,500)
			return
		}
	}

	wallet:= req.Wallet
	userMKInfo:=models.UserMarketInfo {
		Uuid: uuid,
		Wallet: wallet,
		Count: 0,
	}
	_,err=o.InsertOrUpdate(&userMKInfo,"wallet")
	if err!=nil {
		sendError(&this.Controller,err,500)
		return
	}
	logs.Info("inser or update",uuid,"wallet")
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
}

type AvatarController struct {
	beego.Controller
}

func (this *AvatarController) GetAvatar() {
	uuid:= this.Ctx.Input.Param(":uuid")
	var userInfo models.UserInfo
	o:=orm.NewOrm()
	err:=o.QueryTable("user_info").Filter("uuid",uuid).
		One(&userInfo,"avatar_file_name")
	if err!=nil {
		if err==orm.ErrNoRows {
			err := errors.New("no such user")
			sendError(&this.Controller,err,404)
			return
		} else {
			sendError(&this.Controller,err,500)
			return
		}
	}

	avatarFileName:= userInfo.AvatarFileName
	avatarUrl:=util.PathPrefixOfNFT("",common.PATH_KIND_USER_ICON)+ avatarFileName

	type response struct {
		AvatarUrl string `json:"avatarUrl"`
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &response{
		AvatarUrl: avatarUrl,
	}
	this.ServeJSON()
}

func (this *AvatarController) SetAvatar() {
	uuid:= this.Ctx.Input.Param(":uuid")
	file,_,err:=this.GetFile("avatar")
	defer file.Close()
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}

	data,err:= util.ReadFile(file)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	h:=md5.New()
	h.Write(data)
	fileName:=hex.EncodeToString(h.Sum(nil))+".jpg"


	savingPath:= path.Join(common.USER_ICON_PATH,fileName)
	f, err := os.OpenFile(savingPath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	_,err=f.Write(data)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	defer f.Close()
	userInfo:= models.UserInfo{
		Uuid:uuid,
		AvatarFileName:fileName,
	}
	o:= orm.NewOrm()
	_,err=o.Update(&userInfo,"avatar_file_name")
	if err!=nil {
		if err== orm.ErrNoRows {
			err:= errors.New("user not exist")
			sendError(&this.Controller,err,400)
			return
		} else {
			logs.Error(err.Error())
			err:= errors.New("unknown error when query databse")
			sendError(&this.Controller,err,500)
			return
		}
	}
	logs.Info("update",uuid,"file path")
	type response struct {
		AvatarUrl string `json:"avatarUrl"`
	}
	res:= &response{
		AvatarUrl: util.PathPrefixOfNFT("",common.PATH_KIND_USER_ICON)+fileName,
	}
	this.Data["json"]= res
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.ServeJSON()
}
