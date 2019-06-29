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

type SetIntroRequest struct {
	Intro string `json:"intro"`
}

func (this *IntroController) SetIntro() {
	var req SetIntroRequest
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

	nickname:= this.Ctx.Input.Param(":nickname")
	intro:= req.Intro
	logs.Debug("user",nickname,"query intro")
	o:=orm.NewOrm()
	userInfo:= models.MarketUserTable{
		Nickname: nickname,
		Intro: intro,
	}
	_,err=o.Update(&userInfo,"intro")
	if err!=nil {
		if err == orm.ErrNoRows {
			err:= errors.New("no such user")
			if err!=nil {
				logs.Error(err.Error())
				sendError(&this.Controller,err,400)
				return
			}
		} else {
			logs.Error(err.Error())
			sendError(&this.Controller,err,500)
			return
		}
	}
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

type SetWalletRequest struct {
	Address string `json:"Address"`
}

func (this *WalletController) SetWallet() {
	var req SetWalletRequest
	nickname:= this.Ctx.Input.Param(":nickname")
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
	address:= req.Address
	walletInfo:= models.MarketUserTable{
		WalletId: address,
		Nickname: nickname,
	}
	o:=orm.NewOrm()
	_,err=o.Update(&walletInfo,"wallet_id")
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,500)
		return
	}
	logs.Warn("set wallet of nickname",nickname,"to",address)
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	return
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
	file,_,err:=this.GetFile("avatar")
	defer file.Close()
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}

	data,err:= ReadFileFromRequest(file)
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err, 500)
		return
	}
	h:=md5.New()
	h.Write(data)
	fileName:=hex.EncodeToString(h.Sum(nil))+".jpg"


	nickname:= this.Ctx.Input.Param(":nickname")
	savingPath:= path.Join(USER_ICON_PATH,fileName)
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
	userInfo:= models.MarketUserTable{
		Nickname:nickname,
		UserIconUrl: fileName,
	}
	o:= orm.NewOrm()
	_,err=o.Update(&userInfo,"user_icon_url")
	if err!=nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,500)
		return
	}
	logs.Info("save user icon")
	type setAvatarResponse struct {
		AvatarUrl string `json:"avatarUrl"`
	}
	res:= &setAvatarResponse{
		AvatarUrl: PathPrefixOfNFT("",PATH_KIND_USER_ICON)+fileName,
	}
	this.Data["json"]= res
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.ServeJSON()
}
