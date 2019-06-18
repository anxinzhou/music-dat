package http

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"io"
	"math/big"
	"math/rand"
	"mime/multipart"
	"strconv"
	"time"
)

const FILE_SAVING_PATH = "./resource/"
const ENCRYPTION_FILE_PATH = "./resource/encryption/"
const DECRYPTION_FILE_PATH = "./resource/public/"
const MARKET_PATH = "./resource/market/"
const USER_ICON_PATH = "./resource/userIcon/"

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

// base file path
const (
	BASE_FILE_PATH = "resource"
)

// purchase nft status
const (
	PURCHASE_CONFIRMED = 1
	PURCHASE_PENDING = 2
)

// path kind
const (
	PATH_KIND_MARKET = "market"
	PATH_KIND_PUBLIC = "public"
	PATH_KIND_ENCRYPT = "encrypt"
	PATH_KIND_DEFAULT = "default"
	PATH_KIND_USER_ICON = "userIcon"
)

// login type
const (
	LOGIN_TYPE_USERNAME = 3
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

func ReadFileFromRequest(file multipart.File) ([]byte, error){
	var dataBuffer bytes.Buffer
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	for {
		if err == io.EOF {
			dataBuffer.Write(buffer[:n])
			break;
		} else if err != nil {
			return nil,err
		} else {
			dataBuffer.Write(buffer[:n])
			n, err = file.Read(buffer)
		}
	}
	data := dataBuffer.Bytes()
	return data,nil
}

// one to one mapping
func UserIconPathFromNickname(nickname string) string {
	h:=md5.New()
	io.WriteString(h,nickname)
	return hex.EncodeToString(h.Sum(nil))+".jpg"
}

func chinaTimeFromTimeStamp(timestamp time.Time) string {
	timeLocaltion,err:= time.LoadLocation("Asia/Shanghai")
	if err!=nil {
		panic(err)
	}
	return timestamp.In(timeLocaltion).Format("2006-01-02T15:04:05")
}

func RandomPathFromFileName(fileName string) string {
	h := md5.New()
	io.WriteString(h, fileName)
	io.WriteString(h, strconv.FormatInt(time.Now().UnixNano()|rand.Int63(), 10))
	return new(big.Int).SetBytes(h.Sum(nil)[:10]).String()
}

func PathPrefixOfNFT(nftType string, pathKind string) string {
	pathPrefix := beego.AppConfig.String("prefix") + beego.AppConfig.String("hostaddr") + ":" +
		beego.AppConfig.String("fileport") +"/resource/"
	switch pathKind {
	case PATH_KIND_MARKET:
		pathPrefix = pathPrefix+PATH_KIND_MARKET + "/"
	case PATH_KIND_ENCRYPT:
		pathPrefix = pathPrefix+PATH_KIND_ENCRYPT + "/"
	case PATH_KIND_PUBLIC:
		pathPrefix= pathPrefix+PATH_KIND_PUBLIC + "/"
	case PATH_KIND_DEFAULT:
		pathPrefix= pathPrefix+PATH_KIND_DEFAULT + "/"
		return pathPrefix
	case PATH_KIND_USER_ICON:
		pathPrefix = pathPrefix+PATH_KIND_USER_ICON + "/"
		return pathPrefix
	default:
		panic("wrong path kind")
	}
	switch nftType {
	case TYPE_NFT_AVATAR:
		pathPrefix = pathPrefix + NAME_NFT_AVATAR + "/"
	case TYPE_NFT_MUSIC:
		pathPrefix = pathPrefix+NAME_NFT_MUSIC + "/"
	case TYPE_NFT_OTHER:
		pathPrefix = pathPrefix+NAME_NFT_OTHER + "/"
	default:
		panic("wrong nft type")
	}
	return pathPrefix
}

func sendError(c beego.ControllerInterface,err error, statusCode int) {
	controller:=c.(*beego.Controller)
	controller.Ctx.ResponseWriter.ResponseWriter.WriteHeader(statusCode)
	controller.Data["json"] = &ErrorResponse{
		Reason: err.Error(),
	}
	controller.ServeJSON()
}

func smallRandInt() int {
	return int(rand.Int31())%75 + 26    //26 to 100
}