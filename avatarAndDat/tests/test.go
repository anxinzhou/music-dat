package main

import (
	"github.com/astaxie/beego/logs"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	a:=[]byte("$2a$10$la75mLwUDCkxwdMNdOBaS.UHdjo3MD2iESfAmNTM1/h2vgHkFTdYm")
	b:=[]byte("123456")
	err:=bcrypt.CompareHashAndPassword(a,b)
	if err==nil {
		logs.Info("ok")
	} else {
		logs.Info("not ok")
	}
}