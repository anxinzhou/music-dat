package web

import (
	"encoding/json"
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"io/ioutil"
)

type AdminController struct {
	web.Controller
}

func (this *AdminController) Login() {
	type AdminRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type AdminResponse struct {
		Uuid        string `json:"uuid"`
		AccessToken string `json:"accessToken"`
	}
	var req AdminRequest
	data, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("error when parsing request")
		sendError(&this.Controller, err, 500)
		return
	}
	err = json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		err := errors.New("incorrect data format")
		sendError(&this.Controller, err, 400)
		return
	}
	logs.Debug("username", req.Username)

	var res AdminResponse
	username := req.Username
	rawPassword := req.Password
	o := orm.NewOrm()
	queryResult := models.CreatorInfo{
		Username: username,
		Password: rawPassword,
	}
	err = o.Read(&queryResult, "username", "password")
	if err != nil {
		if err == orm.ErrNoRows {
			logs.Error(err.Error())
			sendError(&this.Controller, errors.New("wrong user name or password"), 401)
			return
		} else {
			logs.Error(err.Error())
			sendError(&this.Controller, errors.New("unknown error when query database"), 500)
			return
		}
	}

	res.Uuid = queryResult.Uuid
	res.AccessToken = generateAccessToken()
	this.Data["json"] = &res
	this.ServeJSON()
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	logs.Debug("user login successful")
	return
}
