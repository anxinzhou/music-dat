package ws

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/astaxie/beego"
	"time"
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


var (
	symmetricKey []byte
	aesgcm cipher.AEAD
)

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

func nftResInfoFromNftInfo(nftInfo *NFTInfo) (*nftInfoListRes,error) {
	nftResInfo := &nftInfoListRes{
		SupportedType: nftInfo.SupportedType,
		NftName:       nftInfo.NftName,
		NftValue:      nftInfo.NftValue,
		ActiveTicker:  nftInfo.ActiveTicker,
		NftLifeIndex:  nftInfo.NftLifeIndex,
		NftPowerIndex: nftInfo.NftPowerIndex,
		NftLdefIndex:  nftInfo.NftLdefIndex,
		ShortDesc:     nftInfo.ShortDesc,
		LongDesc:      nftInfo.LongDesc,
		Thumbnail:     nftInfo.FileName,
		Qty:           nftInfo.Qty,
	}
	nftType := nftResInfo.SupportedType
	prefix := PathPrefixOfNFT(nftResInfo.SupportedType, PATH_KIND_MARKET)
	if nftType == TYPE_NFT_AVATAR || nftType == TYPE_NFT_OTHER {
		nftResInfo.Thumbnail = prefix + nftInfo.FileName
	} else if nftType == TYPE_NFT_MUSIC {
		nftResInfo.Thumbnail = prefix + nftInfo.IconFileName
	} else {
		err := errors.New("unexpected type")
		return nil,err
	}
	return nftResInfo,nil
}

func chinaTimeFromTimeStamp(timestamp time.Time) string {
	timeLocaltion,err:= time.LoadLocation("Asia/Shanghai")
	if err!=nil {
		panic(err)
	}
	return timestamp.In(timeLocaltion).Format("2006-01-02T15:04:05")
}

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
