package http

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/astaxie/beego"
	"math/rand"
	"path"
)

const FILE_SAVING_PATH = "./resource/"
const ENCRYPTION_FILE_PATH = "./resource/encryption/"
const DECRYPTION_FILE_PATH = "./resource/public/"
const MARKET_PATH = "./resource/market/"

// NFT TYPE
const (
	TYPE_NFT_AVATAR = "721-02"
	TYPE_NFT_MUSIC = "721-04"
	TYPE_NFT_OTHER = "721-05"
)

// NFT NAME
const (
	NAME_NFT_AVATAR = "avatar"
	NAME_NFT_MUSIC = "dat"
	NAME_NFT_OTHER = "other"
)

// ACTIVE_TICKER
const (
	ACTIVE_TICKER = "berry"
)

// path kind
const (
	PATH_KIND_MARKET = "market"
	PATH_KIND_PUBLIC = "public"
	PATH_KIND_ENCRYPT = "encrypt"
	PATH_KIND_DEFAULT = "default"
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

func PathPrefixOfNFT(nftType string, pathKind string) string {
	pathPrefix := beego.AppConfig.String("prefix") + beego.AppConfig.String("hostaddr") + ":" +
		beego.AppConfig.String("fileport") +"/resource"
	switch pathKind {
	case PATH_KIND_MARKET:
		pathPrefix = path.Join(pathPrefix,PATH_KIND_MARKET)
	case PATH_KIND_ENCRYPT:
		pathPrefix = path.Join(pathPrefix,PATH_KIND_ENCRYPT)
	case PATH_KIND_PUBLIC:
		pathPrefix= path.Join(pathPrefix,PATH_KIND_PUBLIC)
	case PATH_KIND_DEFAULT:
		pathPrefix= path.Join(pathPrefix,PATH_KIND_DEFAULT)
		return pathPrefix
	default:
		panic("wrong path kind")
	}
	switch nftType {
	case TYPE_NFT_AVATAR:
		pathPrefix = path.Join(pathPrefix,NAME_NFT_AVATAR)
	case TYPE_NFT_MUSIC:
		pathPrefix = path.Join(pathPrefix,NAME_NFT_MUSIC)
	case TYPE_NFT_OTHER:
		pathPrefix = path.Join(pathPrefix,NAME_NFT_OTHER)
	default:
		panic("wrong nft type")
	}
	return pathPrefix
}

func sendError(c beego.ControllerInterface,err error, statusCode int) {
	controller:=c.(*beego.Controller)
	controller.Ctx.ResponseWriter.ResponseWriter.WriteHeader(500)
	controller.Data["json"] = &ErrorResponse{
		Reason: err.Error(),
	}
	controller.ServeJSON()
}

func smallRandInt() int {
	return int(rand.Int31())%75 + 26    //26 to 100
}