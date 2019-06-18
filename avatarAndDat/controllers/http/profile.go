package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type NicknameController struct {
	beego.Controller
}

func (this *NicknameController) GetNickname() {

}

func (this *NicknameController) SetNickname() {

}

type IntroController struct {
	beego.Controller
}

func (this *IntroController) GetIntro() {
	nickname:= this.Ctx.Input.Param(":nickname")
	logs.Debug("user",nickname,"query intro")
	userInfo:= models.MarketUserTable{
		Nickname:nickname,
	}
	o:=orm.NewOrm()
	err:=o.Read(&userInfo)
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
	type IntroResponse struct {
		Intro string `json:"intro"`
	}
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	this.Data["json"] = &IntroResponse{
		Intro: userInfo.Intro,
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

}

func (this *AvatarController) SetAvatar() {
	file,_,err:=this.GetFile("avatar")
	defer file.Close()
	if err != nil {
		logs.Error(err.Error())
		sendError(&this.Controller,err,400)
		return
	}
	nickname:= this.Ctx.Input.Param(":nickname")
	fileName:= UserIconPathFromNickname(nickname)
	savingPath:= path.Join(USER_ICON_PATH,fileName)
	f, err := os.OpenFile(savingPath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	io.Copy(f,file)
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

