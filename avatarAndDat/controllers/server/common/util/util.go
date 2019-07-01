package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"math/big"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	symmetricKey []byte
	Aesgcm cipher.AEAD
)

func init() {
	symmetricKey = []byte("passphrasewhichneedstobe32bytes!")
	var err error
	bc, err := aes.NewCipher(symmetricKey)
	if err!=nil {
		panic(err)
	}
	Aesgcm,err =cipher.NewGCM(bc)
	if err!=nil {
		panic(err)
	}
}

func PathPrefixOfNFT(nftType string, pathKind string) string {
	pathPrefix := beego.AppConfig.String("prefix") + beego.AppConfig.String("hostaddr") + ":" +
		beego.AppConfig.String("fileport") +"/resource/"
	switch pathKind {
	case common.PATH_KIND_MARKET:
		pathPrefix = pathPrefix+common.PATH_KIND_MARKET + "/"
	case common.PATH_KIND_ENCRYPT:
		pathPrefix = pathPrefix+common.PATH_KIND_ENCRYPT + "/"
	case common.PATH_KIND_PUBLIC:
		pathPrefix= pathPrefix+ common.PATH_KIND_PUBLIC + "/"
	case common.PATH_KIND_DEFAULT:
		pathPrefix= pathPrefix+ common.PATH_KIND_DEFAULT + "/"
		return pathPrefix
	case common.PATH_KIND_USER_ICON:
		pathPrefix = pathPrefix+ common.PATH_KIND_USER_ICON + "/"
		return pathPrefix
	default:
		panic("wrong path kind")
	}
	switch nftType {
	case common.TYPE_NFT_AVATAR:
		pathPrefix = pathPrefix + common.NAME_NFT_AVATAR + "/"
	case common.TYPE_NFT_MUSIC:
		pathPrefix = pathPrefix+ common.NAME_NFT_MUSIC + "/"
	case common.TYPE_NFT_OTHER:
		pathPrefix = pathPrefix+ common.NAME_NFT_OTHER + "/"
	default:
		panic("wrong nft type")
	}
	return pathPrefix
}

func ValidNftType(supportedType string) error {
	if supportedType!= common.TYPE_NFT_MUSIC && supportedType!= common.TYPE_NFT_OTHER && supportedType!=common.TYPE_NFT_AVATAR {
		err:= errors.New("unsupported nft type")
		return err
	}
	return nil
}

func ValidNftName(nftName string) error {
	if nftName != common.NAME_NFT_AVATAR && nftName != common.NAME_NFT_OTHER && nftName != common.NAME_NFT_MUSIC {
		err := errors.New("no such nft name")
		return err
	}
	return nil
}

func ValidNftLdefIndex(nftLdefIndex string) error {
	if len(nftLdefIndex)<=1 {
		err:=errors.New("insufficient length of nft ldef index")
		return err
	}
	return nil
}

func TokenIdFromNftLdefIndex(nftLdefIndex string) (*big.Int,error) {
	if err:=ValidNftLdefIndex(nftLdefIndex); err!=nil {
		return nil,err
	}
	tokenId,_:=new(big.Int).SetString(nftLdefIndex[1:],10)
	return tokenId,nil
}


func DecryptFile(fileName string, nftType string) (string,error) {
	if err:=ValidNftType(nftType);err!=nil {
		err:= errors.New("unknown decryption type")
		return "",err
	}
	//TODO user symmetric key from client to decrypt file
	var encryptedFilePath string
	var decryptedFilePath string
	logs.Debug("nft type from request,", nftType)
	if nftType == common.TYPE_NFT_AVATAR {
		encryptedFilePath = path.Join(common.ENCRYPTION_FILE_PATH, common.NAME_NFT_AVATAR, fileName)
		decryptedFilePath = path.Join(common.DECRYPTION_FILE_PATH, common.NAME_NFT_AVATAR, fileName)
	} else if nftType == common.TYPE_NFT_MUSIC {
		encryptedFilePath = path.Join(common.ENCRYPTION_FILE_PATH, common.NAME_NFT_MUSIC, fileName)
		decryptedFilePath = path.Join(common.DECRYPTION_FILE_PATH, common.NAME_NFT_MUSIC, fileName)
	} else if nftType == common.TYPE_NFT_OTHER {
		encryptedFilePath = path.Join(common.ENCRYPTION_FILE_PATH, common.NAME_NFT_OTHER, fileName)
		decryptedFilePath = path.Join(common.DECRYPTION_FILE_PATH, common.NAME_NFT_OTHER, fileName)
	}
	cipherText, err := ioutil.ReadFile(encryptedFilePath)
	if err!=nil {
		return "",err
	}

	nonce, ct := cipherText[:Aesgcm.NonceSize()], cipherText[Aesgcm.NonceSize():]
	originalData, err := Aesgcm.Open(nil, nonce, ct, nil)
	if err!=nil {
		return "",err
	}

	logs.Debug("length of original data", len(originalData))
	if nftType == common.TYPE_NFT_AVATAR || nftType == common.TYPE_NFT_OTHER {
		out, err := os.Create(decryptedFilePath)
		if err != nil {
			return "",err
		}
		defer out.Close()
		originalImage, _, err := image.Decode(bytes.NewBuffer(originalData))
		if err != nil {
			return "",err
		}
		err = jpeg.Encode(out, originalImage, nil)
		if err != nil {
			return "",err
		}
	} else if nftType == common.TYPE_NFT_MUSIC {
		err := ioutil.WriteFile(decryptedFilePath, originalData, 0777)
		if err != nil {
			return "",err
		}
	}
	return decryptedFilePath,nil
}

func ChinaTimeFromTimeStamp(timestamp time.Time) string {
	timeLocaltion,err:= time.LoadLocation("Asia/Shanghai")
	if err!=nil {
		panic(err)
	}
	return timestamp.In(timeLocaltion).Format("2006-01-02T15:04:05")
}

func ReadFile(file multipart.File) ([]byte, error){
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

func RandomPathFromFileName(fileName string) string {
	h := md5.New()
	io.WriteString(h, fileName)
	io.WriteString(h, strconv.FormatInt(time.Now().UnixNano()|rand.Int63(), 10))
	return new(big.Int).SetBytes(h.Sum(nil)[:10]).String()
}

func SmallRandInt() int {
	return int(rand.Int31())%75 + 26    //26 to 100
}

func SaveImage(img image.Image, filePath string) error {
	out, err := os.Create(filePath)
	defer out.Close()
	if err != nil {
		return err
	}
	err = jpeg.Encode(out, img, nil)
	if err!=nil {
		return err
	}
	return nil
}

func RandomNftLdefIndex(nftType string) string {
	postPrefix:= strconv.FormatInt(time.Now().UnixNano()|rand.Int63(), 10)
	prefix:=""
	switch nftType {
	case common.TYPE_NFT_AVATAR:
		prefix = "A"
	case common.TYPE_NFT_MUSIC:
		prefix = "M"
	case common.TYPE_NFT_OTHER:
		prefix = "O"
	}
	return prefix+postPrefix
}