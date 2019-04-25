package http

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type AdminController struct {
	beego.Controller
}

type LoginInfo struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (this *AdminController) Login() {
	data:= this.Ctx.Input.RequestBody
	var loginInfo LoginInfo
	err:= json.Unmarshal(data,&loginInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(400)
		this.Data["json"]=&ErrorResponse{
			Reason:err.Error(),
		}
		this.ServeJSON()
		return
	}

	// TODO verify
	//email:= loginInfo.Email
	//password:= loginInfo.Password

	//
	this.Ctx.ResponseWriter.ResponseWriter.WriteHeader(200)
	return
}
