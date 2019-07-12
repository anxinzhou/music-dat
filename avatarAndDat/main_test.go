package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego/session"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/util"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"github.com/xxRanger/music-dat/avatarAndDat/routers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	logs.SetLogFuncCallDepth(3)
	//logs.SetLevel( logs.LevelError)
	//
	// initialize test database
	//
	// change to test db
	beego.AppConfig.Set("dbName", "alphaslot_test")
	beego.AppConfig.Set("fileBasePath", "resourceTest")
	// start server
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "content-type", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
	sessionconf := &session.ManagerConfig{
		CookieName: "begoosessionID",
		Gclifetime: 3600,
	}
	beego.GlobalSessions, _ = session.NewManager("memory", sessionconf)
	go beego.GlobalSessions.GC()
	createDir()
	routers.InitRouter()
	go beego.Run()
	models.InitilizeModel(true, false)
	code := m.Run()
	// clear file path
	fileBasePath := beego.AppConfig.String("fileBasePath")
	d, err := os.Open(fileBasePath)
	if err != nil {
		panic(err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	for _, name := range names {
		if name == "test" {
			continue
		}
		err = os.RemoveAll(filepath.Join(fileBasePath, name))
		if err != nil {
			panic(err)
		}
	}
	// clean test database after finishing test
	//dbUser:= beego.AppConfig.String("dbUser")
	//dbPassword:= beego.AppConfig.String("dbPassword")
	//dbUrls:= beego.AppConfig.String("dbUrls")
	//dbPort:=beego.AppConfig.String("dbPort")
	//dbName:=beego.AppConfig.String("dbName")
	//dbEngine:= beego.AppConfig.String("dbEngine")
	//dbPath:= dbUser+":"+dbPassword+"@"+"tcp("+dbUrls+":"+dbPort+")"+"/"
	//db,err:=sql.Open(dbEngine,dbPath)
	//if err!=nil {
	//	panic(err)
	//}
	//_,err =db.Exec("drop database "+dbName)
	//if err!=nil {
	//	panic(err)
	//}
	os.Exit(code)
}

type testFunc func(t *testing.T)

func UploadTest(uuid string, httpPath string) testFunc {
	return func(t *testing.T) {

	}
}

func TestWebsiteApi(t *testing.T) {
	// necessary to wait for server starting
	<-time.After(1 * time.Second)
	fileBasePath := beego.AppConfig.String("fileBasePath")
	hostaddr := beego.AppConfig.String("hostaddr")
	httpport := beego.AppConfig.String("httpport")
	u := url.URL{Scheme: "ws", Host: hostaddr + ":" + httpport, Path: "/ws"}
	// start test
	testMobileUserUuid := "4298349238490234456sa"
	testWebSiteUserUuid := "4298349238490234456sa1"
	httpPath := "http://" + hostaddr + ":" + httpport
	// Test upload dat
	extraParams := map[string]string{
		"uuid":      testWebSiteUserUuid,
		"nftName":   "testNft",
		"shortDesc": "this is test nft short desc",
		"longDesc":  "this is test nft long desc",
	}
	type nftInfo struct {
		NftLdefIndex  string `json:"nftLdefIndex"`
		NftType       string `json:"nftType"`
		NftName       string `json:"nftName"`
		ShortDesc     string `json:"shortDesc"`
		LongDesc      string `json:"longDesc"`
		FileName      string `json:"fileName"`
		NftParentLdef string `json:"nftParentLdef"`
		SellerUuid    string `json:"sellerUuid"`
		SellerWallet  string `json:"sellerWallet"`
		Price         int    `json:"price"`
		Qty           int    `json:"qty"`
		NumSold       int    `json:"numSold"`
	}
	var testDatInfo nftInfo
	var testAvatarInfo nftInfo
	var testOtherInfo nftInfo

	t.Run("upload_dat", func(t *testing.T) {
		datPath := fileBasePath + "/test/dat.mp3"
		datIconPath := fileBasePath + "/test/icon.jpg"
		uri := httpPath + "/file" + "/dat"
		file, err := os.Open(datPath)
		defer file.Close()
		if err != nil {
			t.Error(err)
			return
		}
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(datPath))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = io.Copy(part, file)
		if err != nil {
			t.Error(err)
			return
		}
		iconFile, err := os.Open(datIconPath)
		defer file.Close()
		if err != nil {
			t.Error(err)
			return
		}
		iconPart, err := writer.CreateFormFile("icon", filepath.Base(datPath))
		_, err = io.Copy(iconPart, iconFile)
		if err != nil {
			t.Error(err)
			return
		}

		for key, val := range extraParams {
			_ = writer.WriteField(key, val)
		}
		_ = writer.WriteField("allowAirdrop", "1")
		_ = writer.WriteField("number", "100")
		_ = writer.WriteField("price", "1")
		_ = writer.WriteField("creatorPercent", "10")
		_ = writer.WriteField("lyricsWriterPercent", "20")
		_ = writer.WriteField("songComposerPercent", "20")
		_ = writer.WriteField("publisherPercent", "10")
		_ = writer.WriteField("userPercent", "40")
		err = writer.Close()
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest("POST", uri, body)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		var respInfo nftInfo
		err = json.Unmarshal(data, &respInfo)
		if err != nil {
			t.Error(err.Error())
			return
		}
		testDatInfo = respInfo
	})
	// Test upload avatar
	t.Run("upload_avatar", func(t *testing.T) {
		avatarPath := fileBasePath + "/test/avatar.jpg"
		uri := httpPath + "/file" + "/avatar"
		file, err := os.Open(avatarPath)
		defer file.Close()
		if err != nil {
			t.Error(err)
			return
		}
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(avatarPath))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = io.Copy(part, file)
		if err != nil {
			t.Error(err)
			return
		}
		for key, val := range extraParams {
			_ = writer.WriteField(key, val)
		}
		err = writer.Close()
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest("POST", uri, body)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		var respInfo nftInfo
		err = json.Unmarshal(data, &respInfo)
		if err != nil {
			t.Error(err.Error())
			return
		}
		testAvatarInfo = respInfo
	})
	// Test upoad other
	t.Run("upload_other", func(t *testing.T) {
		otherPath := fileBasePath + "/test/other.jpg"
		uri := httpPath + "/file" + "/other"
		file, err := os.Open(otherPath)
		defer file.Close()
		if err != nil {
			t.Error(err)
			return
		}
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(otherPath))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = io.Copy(part, file)
		if err != nil {
			t.Error(err)
			return
		}
		for key, val := range extraParams {
			_ = writer.WriteField(key, val)
		}
		_ = writer.WriteField("parent", testAvatarInfo.NftLdefIndex)
		err = writer.Close()
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest("POST", uri, body)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		var respInfo nftInfo
		err = json.Unmarshal(data, &respInfo)
		if err != nil {
			t.Error(err.Error())
			return
		}
		testOtherInfo = respInfo
	})

	// mp_list_dat
	type mpListRequest struct {
		Action        string `json:"action"`
		SupportedType string `json:"supportedType"`
	}
	type mpListNftInfo struct {
		NftLdefIndex  string `json:"nftLdefIndex"`
		NftType       string `json:"nftType"`
		NftName       string `json:"nftName"`
		ShortDesc     string `json:"shortDesc"`
		LongDesc      string `json:"longDesc"`
		FileName      string `json:"thumbnail"`
		NftParentLdef string `json:"nftParentLdef"`
		SellerUuid    string `json:"sellerUuid"`
		SellerWallet  string `json:"sellerWallet"`
		Price         int    `json:"price"`
		Qty           int    `json:"qty"`
		NumSold       int    `json:"numSold"`
	}
	type mpListResponse struct {
		Status        int             `json:"status"`
		Action        string          `json:"action"`
		SupportedType string          `json:"supportedType"`
		NftTranData   []mpListNftInfo `json:"nftTranData"`
		Reason        string          `json:"reason"`
	}
	t.Run(common.ACTION_MP_LIST+"_dat", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		reqTestDat := &mpListRequest{
			Action:        common.ACTION_MP_LIST,
			SupportedType: common.TYPE_NFT_MUSIC,
		}
		datData, err := json.Marshal(reqTestDat)
		if err != nil {
			t.Error(err.Error())
		}
		err = c.WriteMessage(websocket.TextMessage, datData)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			if !ok {
				t.Error("supportedType not specify")
				return
			}
			switch action {
			case common.ACTION_MP_LIST:
				var respInfo mpListResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("insert fail", respInfo.Reason)
					return
				}
				if len(respInfo.NftTranData) == 0 {
					t.Error("insert fail")
					return
				}
				insertedDat := respInfo.NftTranData[0]
				if testDatInfo.NftLdefIndex != insertedDat.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_MUSIC, common.PATH_KIND_MARKET) + testDatInfo.FileName
				if fileUri != insertedDat.FileName {
					t.Error("insert fail, wrong file path", insertedDat.FileName)
				}
				logs.Info("dat file uri", fileUri)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	t.Run(common.ACTION_MP_LIST+"_avatar", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		// test marketplace list
		reqTestAvatar := &mpListRequest{
			Action:        common.ACTION_MP_LIST,
			SupportedType: common.TYPE_NFT_AVATAR,
		}
		avatarData, err := json.Marshal(reqTestAvatar)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, avatarData)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_MP_LIST:
				var respInfo mpListResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("insert fail", respInfo.Reason)
					return
				}
				if len(respInfo.NftTranData) == 0 {
					t.Error("insert fail")
					return
				}

				insertedAvatar := respInfo.NftTranData[0]
				if testAvatarInfo.NftLdefIndex != insertedAvatar.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_AVATAR, common.PATH_KIND_MARKET) + testAvatarInfo.FileName
				if fileUri != insertedAvatar.FileName {
					t.Error("insert fail, wrong file path", insertedAvatar.FileName)
				}
				logs.Info("avatar file uri", fileUri)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	t.Run(common.ACTION_MP_LIST+"_other", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}

		reqTestOther := &mpListRequest{
			Action:        common.ACTION_MP_LIST,
			SupportedType: common.TYPE_NFT_OTHER,
		}
		otherData, err := json.Marshal(reqTestOther)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, otherData)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_MP_LIST:
				var respInfo mpListResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("insert fail", respInfo.Reason)
					return
				}
				if len(respInfo.NftTranData) == 0 {
					t.Error("insert fail")
					return
				}
				insertedOther := respInfo.NftTranData[0]
				if testOtherInfo.NftLdefIndex != insertedOther.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_OTHER, common.PATH_KIND_MARKET) + testOtherInfo.FileName
				if fileUri != insertedOther.FileName {
					t.Error("insert fail, wrong file path", insertedOther.FileName)
				}
				if testOtherInfo.NftParentLdef != testAvatarInfo.NftLdefIndex {
					t.Error("wrong parent of other nft")
				}
				logs.Info("other file uri", fileUri)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})

	t.Run(common.ACTION_MP_LIST+"_avatar", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		// test marketplace list
		reqTestAvatar := &mpListRequest{
			Action:        common.ACTION_MP_LIST,
			SupportedType: common.TYPE_NFT_AVATAR,
		}
		avatarData, err := json.Marshal(reqTestAvatar)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, avatarData)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_MP_LIST:
				var respInfo mpListResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				if len(respInfo.NftTranData)==0 {
					t.Error("insert fail")
					return
				}
				insertedAvatar := respInfo.NftTranData[0]
				if testAvatarInfo.NftLdefIndex != insertedAvatar.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_AVATAR, common.PATH_KIND_MARKET) + testAvatarInfo.FileName
				if fileUri != insertedAvatar.FileName {
					t.Error("insert fail, wrong file path", insertedAvatar.FileName)
				}
				logs.Info("avatar file uri", fileUri)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	// test item details
	type itemDetailsRequest struct {
		Action        string `json:"action"`
		NftLdefIndex  string `json:"nftLdefIndex"`
		SupportedType string `json:"supportedType"`
	}
	type itemDetailsNftInfo struct {
		NftLdefIndex  string `json:"nftLdefIndex"`
		NftType       string `json:"nftType"`
		NftName       string `json:"nftName"`
		ShortDesc     string `json:"shortDesc"`
		LongDesc      string `json:"longDesc"`
		FileName      string `json:"thumbnail"`
		NftParentLdef string `json:"nftParentLdef"`
		SellerUuid    string `json:"sellerUuid"`
		SellerWallet  string `json:"sellerWallet"`
		Price         int    `json:"price"`
		Qty           int    `json:"qty"`
		NumSold       int    `json:"numSold"`
	}
	type itemDetailsResponse struct {
		Status        int                 `json:"status"`
		Action        string              `json:"action"`
		NftLdefIndex  string              `json:"nftLdefIndex"`
		SupportedType string              `json:"supportedType"`
		NftTranData   *itemDetailsNftInfo `json:"nftTranData"`
		Reason        string
	}
	t.Run(common.ACTION_ITEM_DETAILS+"_dat", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}

		reqTestDat := &itemDetailsRequest{
			Action:        common.ACTION_ITEM_DETAILS,
			SupportedType: common.TYPE_NFT_MUSIC,
			NftLdefIndex:  testDatInfo.NftLdefIndex,
		}
		datData, err := json.Marshal(reqTestDat)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, datData)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_ITEM_DETAILS:
				var respInfo itemDetailsResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				insertedDat := respInfo.NftTranData
				if testDatInfo.NftLdefIndex != insertedDat.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_MUSIC, common.PATH_KIND_MARKET) + testDatInfo.FileName
				logs.Info("item details dat file uri", fileUri)
				if fileUri != insertedDat.FileName {
					t.Error("insert fail, wrong file path", insertedDat.FileName)
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	t.Run(common.ACTION_ITEM_DETAILS+"_avatar", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}

		reqTestAvatar := &itemDetailsRequest{
			Action:        common.ACTION_ITEM_DETAILS,
			SupportedType: common.TYPE_NFT_AVATAR,
			NftLdefIndex:  testAvatarInfo.NftLdefIndex,
		}
		avatarData, err := json.Marshal(reqTestAvatar)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, avatarData)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_ITEM_DETAILS:
				var respInfo itemDetailsResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				insertedAvatar := respInfo.NftTranData
				if testAvatarInfo.NftLdefIndex != insertedAvatar.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_AVATAR, common.PATH_KIND_MARKET) + testAvatarInfo.FileName
				if fileUri != insertedAvatar.FileName {
					t.Error("insert fail, wrong file path", insertedAvatar.FileName)
				}
				logs.Info("avatar file uri", fileUri)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	t.Run(common.ACTION_ITEM_DETAILS+"_other", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}

		reqOtherAvatar := &itemDetailsRequest{
			Action:        common.ACTION_ITEM_DETAILS,
			SupportedType: common.TYPE_NFT_OTHER,
			NftLdefIndex:  testOtherInfo.NftLdefIndex,
		}
		otherData, err := json.Marshal(reqOtherAvatar)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, otherData)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_ITEM_DETAILS:
				var respInfo itemDetailsResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				insertedOther := respInfo.NftTranData
				if testOtherInfo.NftLdefIndex != insertedOther.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_OTHER, common.PATH_KIND_MARKET) + testOtherInfo.FileName
				if fileUri != insertedOther.FileName {
					t.Error("insert fail, wrong file path", insertedOther.FileName)
				}
				if testOtherInfo.NftParentLdef != testAvatarInfo.NftLdefIndex {
					t.Error("wrong parent of other nft")
				}
				logs.Info("other file uri", fileUri)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})

	// test follow list add
	type followListOperationRequest struct {
		Action       string `json:"action"`
		Uuid         string `json:"uuid"`
		FolloweeUuid string `json:"followeeUuid"`
		Operation    int    `json:"operation"`
	}

	type followListOperationResponse struct {
		Action       string `json:"action"`
		Status       int    `json:"status"`
		Operation    int    `json:"operation"`
		FolloweeUuid string `json:"followeeUuid"`
		Reason       string `json:"reason"`
	}

	// test follow list add
	t.Run(common.ACTION_FOLLOW_LIST_OPERATION+"_add", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &followListOperationRequest{
			Action:       common.ACTION_FOLLOW_LIST_OPERATION,
			Uuid:         testMobileUserUuid,
			FolloweeUuid: testWebSiteUserUuid,
			Operation:    common.FOLLOW_LIST_ADD,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_FOLLOW_LIST_OPERATION:
				var respInfo followListOperationResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				var followInfo models.FollowTable
				o := orm.NewOrm()
				err = o.QueryTable("follow_table").
					Filter("followee_uuid", testWebSiteUserUuid).
					Filter("follower_uuid", testMobileUserUuid).
					One(&followInfo)
				if err != nil {
					if err != nil {
						t.Error("follow fail", err)
					}
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})

	// test follow list
	type followListRequest struct {
		Action string `json:"action"`
		Uuid   string `json:"uuid"`
	}

	type followerInfo struct {
		Uuid      string `json:"uuid"`
		Nickname  string `json:"nickname"`
		Thumbnail string `json:"thumbnail"`
		Intro     string `json:"intro"`
	}

	type followListResponse struct {
		Status     int             `json:"status"`
		Action     string          `json:"action"`
		FollowList []*followerInfo `json:"followList"`
		Reason     string          `json:"reason"`
	}

	t.Run(common.ACTION_FOLLOW_LIST, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &followListRequest{
			Action: common.ACTION_FOLLOW_LIST,
			Uuid:   testMobileUserUuid,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_FOLLOW_LIST:
				var respInfo followListResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				if len(respInfo.FollowList) == 0 {
					t.Error("follow info is not corrected inserted")
					return
				}
				insertedFollowInfo := respInfo.FollowList[0]
				if insertedFollowInfo.Uuid != testWebSiteUserUuid {
					t.Error("follow info is not corrected inserted")
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})

	// test market user list
	type marketUserListRequest struct {
		Action string `json:"action"`
		Uuid   string `json:"uuid"`
	}
	type markerUserInfo struct {
		Uuid      string `json:"uuid"`
		Nickname  string `json:"nickname"`
		Count     int    `json:"count"`
		Thumbnail string `json:"thumbnail"`
		Followed  bool   `json:"followed"`
	}
	type marketUserListResponse struct {
		Status       int               `json:"status"`
		Action       string            `json:"action"`
		WalletIdList []*markerUserInfo `json:"walletIdList"`
		Reason       string            `json:"reason"`
	}
	t.Run(common.ACTION_MARKET_USER_LIST, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}

		req := &marketUserListRequest{
			Action: common.ACTION_MARKET_USER_LIST,
			Uuid:   testMobileUserUuid,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_MARKET_USER_LIST:
				var respInfo marketUserListResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				walletInfo := respInfo.WalletIdList
				for _, w := range walletInfo {
					if w.Uuid == testWebSiteUserUuid {
						o := orm.NewOrm()
						var userMarketInfo models.UserMarketInfo
						err = o.QueryTable("user_market_info").
							Filter("uuid", w.Uuid).RelatedSel("UserInfo").
							One(&userMarketInfo)

						if err != nil {
							t.Error(err)
							return
						}
						if userMarketInfo.Count != 3 {
							t.Error("unexpected count")
						}
						userIconPath := util.PathPrefixOfNFT("", common.PATH_KIND_USER_ICON) + userMarketInfo.UserInfo.AvatarFileName + "default.jpg"
						if userIconPath != w.Thumbnail {
							t.Error("wrong user icon path", w.Thumbnail)
						}
						logs.Info("user icon path from market", userIconPath)
						if w.Followed != true {
							t.Error("follow status should be true at this moment")
						}
						return
					}
				}
				t.Error("no website user in marketplace")
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})

	t.Run(common.ACTION_FOLLOW_LIST_OPERATION+"_delete", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &followListOperationRequest{
			Action:       common.ACTION_FOLLOW_LIST_OPERATION,
			Uuid:         testMobileUserUuid,
			FolloweeUuid: testWebSiteUserUuid,
			Operation:    common.FOLLOW_LIST_DELETE,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_FOLLOW_LIST_OPERATION:
				var respInfo followListOperationResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				var followInfo models.FollowTable
				o := orm.NewOrm()
				err = o.QueryTable("follow_table").
					Filter("followee_uuid", testWebSiteUserUuid).
					Filter("follower_uuid", testMobileUserUuid).
					One(&followInfo)
				if err != orm.ErrNoRows {
					t.Error("now there should be no row")
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	// test duplicate nickname
	type duplicateNicknameRequest struct {
		Action   string `json:"action"`
		Nickname string `json:"nickname"`
	}
	type duplicateNicknameResponse struct {
		Status     int    `json:"status"`
		Action     string `json:"action"`
		Reason     string `json:"reason"`
		Duplicated bool   `json:"duplicated"`
		Nickname   string `json:"nickname"`
	}
	t.Run(common.ACTION_IS_NICKNAME_DUPLICATED, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		o := orm.NewOrm()
		userInfo := models.UserInfo{
			Uuid: testWebSiteUserUuid,
		}
		err = o.Read(&userInfo)
		if err != nil {
			t.Error(err)
			return
		}
		nickname := userInfo.Nickname
		req := &duplicateNicknameRequest{
			Action:   common.ACTION_IS_NICKNAME_DUPLICATED,
			Nickname: nickname,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_IS_NICKNAME_DUPLICATED:
				var respInfo duplicateNicknameResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				if !respInfo.Duplicated {
					t.Error("duplicated should be true at this time")
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	// test set nickname with duplicate insert
	type setNicknameRequest struct {
		Action   string `json:"action"`
		Uuid     string `json:"uuid"`
		Nickname string `json:"nickname"`
	}
	type setNicknameResponse struct {
		Status   int    `json:"status"`
		Action   string `json:"action"`
		Reason   string `json:"reason"`
		Nickname string `json:"nickname"`
	}
	testNewUserUuid := "89043850943860xxxx"
	t.Run(common.ACTION_NFT_SET_NICKNAME, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		o := orm.NewOrm()
		userInfo := models.UserInfo{
			Uuid: testWebSiteUserUuid,
		}
		err = o.Read(&userInfo)
		if err != nil {
			t.Error(err)
			return
		}
		nickname := userInfo.Nickname
		req := &setNicknameRequest{
			Action:   common.ACTION_NFT_SET_NICKNAME,
			Uuid:     testNewUserUuid,
			Nickname: nickname,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_NFT_SET_NICKNAME:
				var respInfo setNicknameResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					if respInfo.Reason != "nickname has been registered" {
						t.Error("fail", respInfo.Reason)
						return
					}
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})

	// test set nickname without duplicate insert
	testNewUserUuid2 := "890438509438fdsffds"
	testNickname := "baobao"
	t.Run(common.ACTION_NFT_SET_NICKNAME, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &setNicknameRequest{
			Action:   common.ACTION_NFT_SET_NICKNAME,
			Uuid:     testNewUserUuid2,
			Nickname: testNickname,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_NFT_SET_NICKNAME:
				var respInfo setNicknameResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				o := orm.NewOrm()
				userInfo := models.UserInfo{
					Uuid: testNewUserUuid2,
				}
				err = o.Read(&userInfo)
				if err != nil {
					t.Error(err.Error())
					return
				}
				if userInfo.Nickname != testNickname {
					t.Error("nickname should be equal")
					return
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})

	// test set nickname with duplicate insert
	type isNicknameSetRequest struct {
		Action string `json:"action"`
		Uuid   string `json:"uuid"`
	}
	type isNicknameSetResponse struct {
		Status int    `json:"status"`
		Action string `json:"action"`
		Reason string `json:"reason"`
		Set    bool   `json:"set"`
	}
	t.Run(common.ACTION_IS_NICKNAME_SET, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		o := orm.NewOrm()
		userInfo := models.UserInfo{
			Uuid: testWebSiteUserUuid,
		}
		err = o.Read(&userInfo)
		if err != nil {
			t.Error(err)
			return
		}
		req := &isNicknameSetRequest{
			Action: common.ACTION_IS_NICKNAME_SET,
			Uuid:   testNewUserUuid2,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_IS_NICKNAME_SET:
				var respInfo isNicknameSetResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				if !respInfo.Set {
					t.Error("nickname set should be true")
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})

	// test bind wallet
	type bindWalletRequest struct {
		Action   string `json:"action"`
		Uuid     string `json:"uuid"`
		WalletId string `json:"walletId"`
	}
	type bindWalletResponse struct {
		Status   int    `json:"status"`
		Action   string `json:"action"`
		Reason   string `json:"reason"`
		WalletId string `json:"walletId"`
		Set      bool   `json:"set"`
	}
	testWallet := "0x3c62aa7913bc303ee4b9c07df87b556b6770e3fc"
	t.Run(common.ACTION_NFT_BIND_WALLET, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &bindWalletRequest{
			Action:   common.ACTION_NFT_BIND_WALLET,
			Uuid:     testNewUserUuid2,
			WalletId: testWallet,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_NFT_BIND_WALLET:
				var respInfo bindWalletResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				userMarketInfo := models.UserMarketInfo{
					Uuid: testNewUserUuid2,
				}
				o := orm.NewOrm()
				err = o.Read(&userMarketInfo)
				if err != nil {
					t.Error(err.Error())
					return
				}
				if userMarketInfo.Wallet != testWallet {
					t.Error("set wallet fail")
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})

	// test bind wallet
	type nftTransferRequest struct {
		Action       string `json:"action"`
		SenderUuid   string `json:"senderUuid"`
		ReceiverUuid string `json:"receiverUuid"`
		NftLdefIndex string `json:"nftLdefIndex"`
	}
	type nftTransferResponse struct {
		Status       int    `json:"status"`
		Action       string `json:"action"`
		Reason       string `json:"reason"`
		SenderUuid   string `json:"senderUuid"`
		ReceiverUuid string `json:"receiverUuid"`
		NftLdefIndex string `json:"nftLdefIndex"`
	}
	t.Run(common.ACTION_NFT_TRANSFER, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &nftTransferRequest{
			Action:       common.ACTION_NFT_TRANSFER,
			SenderUuid:   testWebSiteUserUuid,
			ReceiverUuid: testMobileUserUuid,
			NftLdefIndex: testAvatarInfo.NftLdefIndex,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_NFT_TRANSFER:
				var respInfo nftTransferResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				// check count
				buyer := models.UserMarketInfo{
					Uuid: testMobileUserUuid,
				}
				seller := models.UserMarketInfo{
					Uuid: testWebSiteUserUuid,
				}
				o := orm.NewOrm()
				err = o.Read(&buyer)
				if err != nil {
					t.Error(err)
					return
				}
				err = o.Read(&seller)
				if err != nil {
					t.Error(err)
					return
				}
				if buyer.Count != 1 {
					t.Error("buyer count should be 1")
				}
				if seller.Count != 2 {
					t.Error("seller count should be 2")
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	// test shopping cart change add
	type shoppingCartRequest struct {
		Action    string   `json:"action"`
		Operation int      `json:"operation"`
		Uuid      string   `json:"uuid"`
		NftList   []string `json:"nftList"`
	}
	type shoppingCartChangeResponse struct {
		Status    int      `json:"status"`
		Action    string   `json:"action"`
		Operation int      `json:"operation"`
		NftList   []string `json:"nftList"`
		Reason    string   `json:"reason"`
	}
	t.Run(common.ACTION_NFT_SHOPPING_CART_CHANGE+"_add", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &shoppingCartRequest{
			Action:    common.ACTION_NFT_SHOPPING_CART_CHANGE,
			Operation: common.SHOPPING_CART_ADD,
			Uuid:      testMobileUserUuid,
			NftList: []string{
				testAvatarInfo.NftLdefIndex,
			},
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_NFT_SHOPPING_CART_CHANGE:
				var respInfo shoppingCartChangeResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				o := orm.NewOrm()
				shoppingCartInfo := models.NftShoppingCart{
					NftLdefIndex: testAvatarInfo.NftLdefIndex,
					Uuid:         testMobileUserUuid,
				}
				err = o.Read(&shoppingCartInfo, "nft_ldef_index", "uuid")
				if err != nil {
					t.Error(err)
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	// test shopping cart list
	type shoppingCartListRequest struct {
		Action string `json:"action"`
		Uuid   string `json:"uuid"`
	}
	type shoppingCartInfo struct {
		NftLdefIndex  string `json:"nftLdefIndex"`
		NftType       string `json:"supportedType"`
		NftName       string `json:"nftName"`
		ShortDesc     string `json:"shortDesc"`
		LongDesc      string `json:"longDesc"`
		FileName      string `json:"thumbnail"`
		NftParentLdef string `json:"nftParentLdef"`
		Price         int    `json:"price"`
	}
	type shoppingCartListResponse struct {
		Status  int                 `json:"status"`
		Action  string              `json:"action"`
		NftList []*shoppingCartInfo `json:"nftList"`
		Reason  string              `json:"reason"`
	}
	t.Run(common.ACTION_NFT_SHOPPING_CART_LIST, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &shoppingCartListRequest{
			Action: common.ACTION_NFT_SHOPPING_CART_LIST,
			Uuid:   testMobileUserUuid,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_NFT_SHOPPING_CART_LIST:
				var respInfo shoppingCartListResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				if len(respInfo.NftList) == 0 {
					t.Error("no nft found in shopping cart")
					return
				}
				insertedAvatar := respInfo.NftList[0]
				if testAvatarInfo.NftLdefIndex != insertedAvatar.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_AVATAR, common.PATH_KIND_MARKET) + testAvatarInfo.FileName
				if fileUri != insertedAvatar.FileName {
					t.Error("insert fail, wrong file path", insertedAvatar.FileName)
				}
				logs.Info("avatar file uri", fileUri)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})

	t.Run(common.ACTION_NFT_SHOPPING_CART_CHANGE+"_delete", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &shoppingCartRequest{
			Action:    common.ACTION_NFT_SHOPPING_CART_CHANGE,
			Operation: common.SHOPPING_CART_DELETE,
			Uuid:      testMobileUserUuid,
			NftList: []string{
				testAvatarInfo.NftLdefIndex,
			},
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_NFT_SHOPPING_CART_CHANGE:
				var respInfo shoppingCartChangeResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				o := orm.NewOrm()
				shoppingCartInfo := models.NftShoppingCart{
					NftLdefIndex: testAvatarInfo.NftLdefIndex,
					Uuid:         testMobileUserUuid,
				}
				err = o.Read(&shoppingCartInfo, "nft_ldef_index", "uuid")
				if err != orm.ErrNoRows {
					if err != nil {
						t.Error(err)
					} else {
						t.Error("shopping cart should be empty")
					}
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	// test nft purchase
	type nftPurchaseConfirmRequest struct {
		Action      string   `json:"action"`
		Uuid        string   `json:"uuid"`
		NftTranData []string `json:"nftTranData"`
	}
	type nftPurchaseConfirmResponse struct {
		Status      int      `json:"status"`
		Action      string   `json:"action"`
		NftTranData []string `json:"nftTranData"`
		Reason      string   `json:"reason"`
	}
	t.Run(common.ACTION_NFT_PUCHASE_CONFIRM, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &nftPurchaseConfirmRequest{
			Action: common.ACTION_NFT_PUCHASE_CONFIRM,
			Uuid:   testMobileUserUuid,
			NftTranData: []string{
				testDatInfo.NftLdefIndex,
				testAvatarInfo.NftLdefIndex,
				testOtherInfo.NftLdefIndex,
			},
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_NFT_PUCHASE_CONFIRM:
				var respInfo nftPurchaseConfirmResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	// test purchase history
	type nftPurchaseHistoryRequest struct {
		Action        string `json:"action"`
		Uuid          string `json:"uuid"`
		SupportedType string `json:"supportedType"`
	}
	type purchaseNftInfo struct {
		TransactionAddress string `json:"transactionAddress"`
		NftLdefIndex       string `json:"nftLdefIndex"`
		SupportedType      string `json:"supportedType"`
		NftName            string `json:"nftName"`
		ShortDesc          string `json:"shortDesc"`
		LongDesc           string `json:"longDesc"`
		Thumbnail          string `json:"thumbnail"`
		DecSource          string `json:"decSource"`
		Qty                int    `json:"qty"`
		NftLifeIndex       int    `json:"nftLifeIndex"`
		NftPowerIndex      int    `json:"nftPowerIndex"`
	}
	type nftPurchaseHistoryResponse struct {
		Status        int                `json:"status"`
		Action        string             `json:"action"`
		SupportedType string             `json:"supportedType"`
		Reason        string             `json:"reason"`
		PurchaseList  []*purchaseNftInfo `json:"purchaseList"`
	}
	t.Run(common.ACTION_NFT_PURCHASE_HISTORY+"_avatar", func(t *testing.T) {
		logs.Info("test")
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &nftPurchaseHistoryRequest{
			Action:        common.ACTION_NFT_PURCHASE_HISTORY,
			Uuid:          testMobileUserUuid,
			SupportedType: common.TYPE_NFT_AVATAR,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_NFT_PURCHASE_HISTORY:
				var respInfo nftPurchaseHistoryResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				if len(respInfo.PurchaseList) == 0 {
					t.Error("insert fail")
					return
				}

				insertedAvatar := respInfo.PurchaseList[0]
				if testAvatarInfo.NftLdefIndex != insertedAvatar.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_AVATAR, common.PATH_KIND_PUBLIC) + testAvatarInfo.FileName
				if fileUri != insertedAvatar.Thumbnail {
					t.Error("insert fail, wrong file path", insertedAvatar.Thumbnail)
				}
				logs.Info("avatar file uri", fileUri)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	t.Run(common.ACTION_NFT_PURCHASE_HISTORY+"_dat", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &nftPurchaseHistoryRequest{
			Action:        common.ACTION_NFT_PURCHASE_HISTORY,
			Uuid:          testMobileUserUuid,
			SupportedType: common.TYPE_NFT_MUSIC,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_NFT_PURCHASE_HISTORY:
				var respInfo nftPurchaseHistoryResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				if len(respInfo.PurchaseList) == 0 {
					t.Error("insert fail")
					return
				}
				insertedDat := respInfo.PurchaseList[0]
				if testDatInfo.NftLdefIndex != insertedDat.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_MUSIC, common.PATH_KIND_PUBLIC) + testDatInfo.FileName
				if fileUri != insertedDat.Thumbnail {
					t.Error("insert fail, wrong file path", insertedDat.Thumbnail)
				}
				logs.Info("dat file uri", fileUri)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	t.Run(common.ACTION_NFT_PURCHASE_HISTORY+"_other", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &nftPurchaseHistoryRequest{
			Action:        common.ACTION_NFT_PURCHASE_HISTORY,
			Uuid:          testMobileUserUuid,
			SupportedType: common.TYPE_NFT_OTHER,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_NFT_PURCHASE_HISTORY:
				var respInfo nftPurchaseHistoryResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				if len(respInfo.PurchaseList) == 0 {
					t.Error("insert fail")
					return
				}
				insertedOther := respInfo.PurchaseList[0]
				if testOtherInfo.NftLdefIndex != insertedOther.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_OTHER, common.PATH_KIND_PUBLIC) + testOtherInfo.FileName
				if fileUri != insertedOther.Thumbnail {
					t.Error("insert fail, wrong file path", insertedOther.Thumbnail)
				}
				if testOtherInfo.NftParentLdef != testAvatarInfo.NftLdefIndex {
					t.Error("wrong parent of other nft")
				}
				logs.Info("other file uri", fileUri)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})

	// test token purchase history
	type tokenPurchaseRequest struct {
		Action        string `json:"action"`
		Uuid          string `json:"uuid"`
		AppTranId     string `json:"appTranId"`
		TransactionId string `json:"transactionId"`
		AppId         string `json:"appId"`
		Amount        int    `json:"amount"`
		ActionStatus  int    `json:"actionStatus"`
	}
	type tokenPurchaseResponse struct {
		Status        int    `json:"status"`
		Action        string `json:"action"`
		Amount        int    `json:"amount"`
		ActionStatus  int    `json:"actionStatus"`
		Reason        string `json:"reason"`
		TransactionId string `json:"transactionId"`
	}
	var transactionId string
	t.Run(common.ACTION_TOKENBUY_PAID+"_pending", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &tokenPurchaseRequest{
			Action:        common.ACTION_TOKENBUY_PAID,
			Uuid:          testMobileUserUuid,
			AppTranId:     "4fffc3b55bd8270df6fd85d1f17b8ec53adfa64d882",
			TransactionId: "",
			AppId:         "XXX@appleid.com",
			Amount:        4,
			ActionStatus:  common.BERRY_PURCHASE_PENDING,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_TOKENBUY_PAID:
				var respInfo tokenPurchaseResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}
				transactionId = respInfo.TransactionId
				if transactionId == "" {
					t.Error("empty tansaction id")
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	testBuyAmount:=4
	t.Run(common.ACTION_TOKENBUY_PAID+"_finish", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &tokenPurchaseRequest{
			Action:        common.ACTION_TOKENBUY_PAID,
			Uuid:          testMobileUserUuid,
			AppTranId:     "4fffc3b55bd8270df6fd85d1f17b8ec53adfa64d882",
			TransactionId: transactionId,
			AppId:         "XXX@appleid.com",
			Amount:        testBuyAmount,
			ActionStatus:  common.BERRY_PURCHASE_FINISH,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}
		// before send count amount
		// update coin records
		col := models.MongoDB.Collection("users")
		// before coin amount
		type fields struct {
			Coin string `bson:"coin"`
		}

		filter := bson.M{
			"uuid": req.Uuid,
		}

		var queryResult fields

		err = col.FindOne(context.Background(), filter, options.FindOne().SetProjection(bson.M{
			"coin": true,
		})).Decode(&queryResult)
		if err != nil {
			t.Error(err)
		}

		currentBalance, err := strconv.Atoi(queryResult.Coin)
		if err != nil {
			t.Error(err)
		}

		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_TOKENBUY_PAID:
				var respInfo tokenPurchaseResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}

				err = col.FindOne(context.Background(), filter, options.FindOne().SetProjection(bson.M{
					"coin": true,
				})).Decode(&queryResult)
				if err != nil {
					t.Error(err)
				}

				afterBalance, err := strconv.Atoi(queryResult.Coin)
				if err != nil {
					t.Error(err)
				}
				if afterBalance != currentBalance+testBuyAmount {
					t.Error("uncorrect balance", "before:", currentBalance, "after:", afterBalance)
				}
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	// test token purchase history
	type marketUserRequest struct {
		Action        string `json:"action"`
		Uuid          string `json:"uuid"`
		SupportedType string `json:"supportedType"`
	}
	type marketUserResponse struct {
		Status        int    `json:"status"`
		Action        string `json:"action"`
		NftTranData []*mpListNftInfo`json:"nftTranData"`
		Reason        string `json:"reason"`
	}
	t.Run(common.ACTION_USER_MARKET_INFO+"_dat", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &marketUserRequest{
			Action:        common.ACTION_USER_MARKET_INFO,
			Uuid:          testWebSiteUserUuid,
			SupportedType: common.TYPE_NFT_MUSIC,
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}

		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_USER_MARKET_INFO:
				var respInfo marketUserResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}

				if len(respInfo.NftTranData) == 0 {
					t.Error("insert fail")
					return
				}
				insertedDat := respInfo.NftTranData[0]
				if testDatInfo.NftLdefIndex != insertedDat.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_MUSIC, common.PATH_KIND_MARKET) + testDatInfo.FileName
				if fileUri != insertedDat.FileName {
					t.Error("insert fail, wrong file path", insertedDat.FileName)
				}
				logs.Info("dat file uri", fileUri)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	t.Run(common.ACTION_USER_MARKET_INFO+"_all", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &marketUserRequest{
			Action:        common.ACTION_USER_MARKET_INFO,
			Uuid:          testWebSiteUserUuid,
			SupportedType: "",
		}
		data, err := json.Marshal(req)
		if err != nil {
			t.Error(err.Error())
			return
		}

		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Error(err)
			return
		}
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				logs.Error(err.Error())
				break;
			}
			var kvs map[string]interface{}
			json.Unmarshal(data, &kvs)
			action, ok := kvs["action"]
			if !ok {
				t.Error("action not exist")
				return
			}
			switch action {
			case common.ACTION_USER_MARKET_INFO:
				var respInfo marketUserResponse
				err = json.Unmarshal(data, &respInfo)
				if err != nil {
					t.Error(err)
					return
				}
				if respInfo.Status == common.RESPONSE_STATUS_FAIL {
					t.Error("fail", respInfo.Reason)
					return
				}

				if len(respInfo.NftTranData) == 0 {
					t.Error("insert fail")
					return
				}
				insertedDat := respInfo.NftTranData[0]
				if testDatInfo.NftLdefIndex != insertedDat.NftLdefIndex {
					t.Error("insert fail")
				}
				fileUri := util.PathPrefixOfNFT(common.TYPE_NFT_MUSIC, common.PATH_KIND_MARKET) + testDatInfo.FileName
				if fileUri != insertedDat.FileName {
					t.Error("insert fail, wrong file path", insertedDat.FileName)
				}
				logs.Info("dat file uri", fileUri)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
	// wait for transaction to finish
	<-time.After(5*time.Second)
	logs.Info("test")
}
