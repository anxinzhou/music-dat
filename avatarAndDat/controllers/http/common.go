package http

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/astaxie/beego"
)

const FILE_SAVING_PATH = "./resource/"
const ENCRYPTION_FILE_PATH = "./resource/encryption/"
const DECRYPTION_FILE_PATH = "./resource/public/"
const MARKET_PATH = "./resource/market/"

// NFT TYPE
const (
	TYPE_NFT_AVATAR = "721-02"
	TYPE_NFT_MUSIC = "721-04"
)

// NFT NAME
const (
	NAME_NFT_AVATAR = "avatar"
	NAME_NFT_MUSIC = "dat"
)

var (
	symmetricKey []byte
	aesgcm cipher.AEAD
)

func init() {
	symmetricKey = []byte("passphrasewhichneedstobe32bytes!")
	var err error
	bc, err := aes.NewCipher(symmetricKey)
	if err!=nil {
		panic(err)
	}
	aesgcm,err =cipher.NewGCM(bc)
	if err!=nil {
		panic(err)
	}
}

func sendError(c beego.ControllerInterface,err error, statusCode int) {
	controller:=c.(*beego.Controller)
	controller.Ctx.ResponseWriter.ResponseWriter.WriteHeader(500)
	controller.Data["json"] = &ErrorResponse{
		Reason: err.Error(),
	}
	controller.ServeJSON()
}