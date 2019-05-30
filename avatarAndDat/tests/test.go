package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/astaxie/beego/logs"
	"github.com/jameskeane/bcrypt"
)

func main() {
	a:="$2a$10$la75mLwUDCkxwdMNdOBaS.UHdjo3MD2iESfAmNTM1/h2vgHkFTdYm"
	h:= sha256.New()
	h.Write([]byte("123456"))
	b:=hex.EncodeToString(h.Sum(nil))
	logs.Info(b)
	ok:=bcrypt.Match(b,a)
	if ok {
		logs.Info("ok")
	} else {
		logs.Info("not ok")
	}
}