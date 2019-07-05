package main

import (
	"bytes"
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
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	logs.SetLogFuncCallDepth(3)
	//
	// initialize test database
	//
	// change to test db
	beego.AppConfig.Set("dbName", "alphaslot_test")
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
	beego.SetStaticPath("/resource", "resource")
	go beego.Run()
	models.InitilizeModel(true, false)
	code := m.Run()

	//// clean test database after finishing test
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

func TestWebsiteApi(t *testing.T) {
	// necessary to wait for server starting
	<-time.After(1 * time.Second)

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
		datPath := "./resource/test/dat.mp3"
		datIconPath := "./resource/test/icon.jpg"
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
		_ = writer.WriteField("number", "1")
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
		avatarPath := "./resource/test/avatar.jpg"
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
		otherPath := "./resource/test/other.jpg"
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
		Status        int       `json:"status"`
		Action        string    `json:"action"`
		SupportedType string    `json:"supportedType"`
		NftTranData   []mpListNftInfo `json:"nftTranData"`
		Reason        string    `json:"reason"`
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
		Action string `json:"action"`
		NftLdefIndex string `json:"nftLdefIndex"`
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
		Status int `json:"status"`
		Action string `json:"action"`
		NftLdefIndex string `json:"nftLdefIndex"`
		SupportedType string `json:"supportedType"`
		NftTranData *itemDetailsNftInfo `json:"nftTranData"`
		Reason string
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
			NftLdefIndex: testDatInfo.NftLdefIndex,
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
			NftLdefIndex: testAvatarInfo.NftLdefIndex,
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
			NftLdefIndex: testOtherInfo.NftLdefIndex,
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
		Action string `json:"action"`
		Uuid string `json:"uuid"`
		FolloweeUuid string `json:"followeeUuid"`
		Operation int `json:"operation"`
	}

	type followListOperationResponse struct {
		Action string `json:"action"`
		Status int `json:"status"`
		Operation int `json:"operation"`
		FolloweeUuid string `json:"followeeUuid"`
		Reason string `json:"reason"`
	}

	// test follow list add
	t.Run(common.ACTION_FOLLOW_LIST_OPERATION+"_add", func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &followListOperationRequest{
			Action:        common.ACTION_FOLLOW_LIST_OPERATION,
			Uuid: testMobileUserUuid,
			FolloweeUuid: testWebSiteUserUuid,
			Operation: common.FOLLOW_LIST_ADD,
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
				o:=orm.NewOrm()
				err=o.QueryTable("follow_table").
					Filter("followee_uuid",testWebSiteUserUuid).
					Filter("follower_uuid",testMobileUserUuid).
					One(&followInfo)
				if err!=nil {
					if err!=nil {
						t.Error("follow fail",err)
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
		Uuid string  `json:"uuid"`
	}


	type followerInfo struct {
		Uuid string `json:"uuid"`
		Nickname string `json:"nickname"`
		Thumbnail string `json:"thumbnail"`
		Intro string `json:"intro"`
	}

	type followListResponse struct {
		Status int `json:"status"`
		Action string `json:"action"`
		FollowList []*followerInfo `json:"followList"`
		Reason string `json:"reason"`
	}

	t.Run(common.ACTION_FOLLOW_LIST, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &followListRequest{
			Action:        common.ACTION_FOLLOW_LIST,
			Uuid: testMobileUserUuid,
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
				if len(respInfo.FollowList)==0 {
					t.Error("follow info is not corrected inserted")
					return
				}
				insertedFollowInfo:=  respInfo.FollowList[0]
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
		Uuid string  `json:"uuid"`
	}
	type markerUserInfo struct {
		Uuid string  `json:"uuid"`
		Nickname string `json:"nickname"`
		Count int `json:"count"`
		Thumbnail string `json:"thumbnail"`
		Followed bool `json:"followed"`
	}
	type marketUserListResponse struct {
		Status int `json:"status"`
		Action string `json:"action"`
		WalletIdList []*markerUserInfo `json:"walletIdList"`
		Reason string `json:"reason"`
	}
	t.Run(common.ACTION_MARKET_USER_LIST, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}

		req := &marketUserListRequest{
			Action:        common.ACTION_MARKET_USER_LIST,
			Uuid: testMobileUserUuid,
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
				for _,w:=range  walletInfo {
					if w.Uuid == testWebSiteUserUuid {
						o:=orm.NewOrm()
						var userMarketInfo models.UserMarketInfo
						err=o.QueryTable("user_market_info").
							Filter("uuid",w.Uuid).RelatedSel("UserInfo").
							One(&userMarketInfo)

						if err!=nil {
							t.Error(err)
							return
						}
						if userMarketInfo.Count != 3 {
							t.Error("unexpected count")
						}
						userIconPath:= util.PathPrefixOfNFT("",common.PATH_KIND_USER_ICON)+userMarketInfo.UserInfo.AvatarFileName+"default.jpg"
						if userIconPath!= w.Thumbnail {
							t.Error("wrong user icon path",w.Thumbnail)
						}
						logs.Info("user icon path from market",userIconPath)
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
			Action:        common.ACTION_FOLLOW_LIST_OPERATION,
			Uuid: testMobileUserUuid,
			FolloweeUuid: testWebSiteUserUuid,
			Operation: common.FOLLOW_LIST_DELETE,
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
				o:=orm.NewOrm()
				err=o.QueryTable("follow_table").
					Filter("followee_uuid",testWebSiteUserUuid).
					Filter("follower_uuid",testMobileUserUuid).
					One(&followInfo)
				if err!=orm.ErrNoRows {
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
		Action string `json:"action"`
		Nickname string `json:"nickname"`
	}
	type duplicateNicknameResponse struct {
		Status int `json:"status"`
		Action string `json:"action"`
		Reason string `json:"reason"`
		Duplicated bool `json:"duplicated"`
		Nickname string `json:"nickname"`
	}
	t.Run(common.ACTION_IS_NICKNAME_DUPLICATED, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		o:=orm.NewOrm()
		userInfo:=models.UserInfo{
			Uuid: testWebSiteUserUuid,
		}
		err = o.Read(&userInfo)
		if err!=nil {
			t.Error(err)
			return
		}
		nickname:= userInfo.Nickname
		req := &duplicateNicknameRequest{
			Action:        common.ACTION_IS_NICKNAME_DUPLICATED,
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
		Action string `json:"action"`
		Uuid string `json:"uuid"`
		Nickname string `json:"nickname"`
	}
	type setNicknameResponse struct {
		Status int `json:"status"`
		Action string `json:"action"`
		Reason string `json:"reason"`
		Nickname string `json:"nickname"`
	}
	testNewUserUuid:= "89043850943860xxxx"
	t.Run(common.ACTION_NFT_SET_NICKNAME, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		o:=orm.NewOrm()
		userInfo:=models.UserInfo{
			Uuid: testWebSiteUserUuid,
		}
		err = o.Read(&userInfo)
		if err!=nil {
			t.Error(err)
			return
		}
		nickname:= userInfo.Nickname
		req := &setNicknameRequest{
			Action:        common.ACTION_NFT_SET_NICKNAME,
			Uuid: testNewUserUuid,
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
	testNewUserUuid2:= "890438509438fdsffds"
	testNickname:= "baobao"
	t.Run(common.ACTION_NFT_SET_NICKNAME, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &setNicknameRequest{
			Action:        common.ACTION_NFT_SET_NICKNAME,
			Uuid: testNewUserUuid2,
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
				o:=orm.NewOrm()
				userInfo:=models.UserInfo{
					Uuid: testNewUserUuid2,
				}
				err = o.Read(&userInfo)
				if err!=nil {
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
		Uuid string `json:"uuid"`
	}
	type isNicknameSetResponse struct {
		Status int `json:"status"`
		Action string `json:"action"`
		Reason string `json:"reason"`
		Set bool `json:"set"`
	}
	t.Run(common.ACTION_IS_NICKNAME_SET, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		o:=orm.NewOrm()
		userInfo:=models.UserInfo{
			Uuid: testWebSiteUserUuid,
		}
		err = o.Read(&userInfo)
		if err!=nil {
			t.Error(err)
			return
		}
		req := &isNicknameSetRequest{
			Action:        common.ACTION_IS_NICKNAME_SET,
			Uuid: testNewUserUuid2,
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
		Action string `json:"action"`
		Uuid string `json:"uuid"`
		WalletId string `json:"walletId"`
	}
	type bindWalletResponse struct {
		Status int `json:"status"`
		Action string `json:"action"`
		Reason string `json:"reason"`
		WalletId string `json:"walletId"`
		Set bool `json:"set"`
	}
	testWallet:= "0x3c62aa7913bc303ee4b9c07df87b556b6770e3fc"
	t.Run(common.ACTION_NFT_BIND_WALLET, func(t *testing.T) {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer c.Close()
		if err != nil {
			t.Error("can not dail to ", u.String())
		}
		req := &bindWalletRequest{
			Action:        common.ACTION_NFT_BIND_WALLET,
			Uuid: testNewUserUuid2,
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
				userMarketInfo:= models.UserMarketInfo{
					Uuid: testNewUserUuid2,
				}
				o:= orm.NewOrm()
				err = o.Read(&userMarketInfo)
				if err!= nil {
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
		Action string `json:"action"`
		SenderUuid string `json:"senderUuid"`
		ReceiverUuid string `json:"receiverUuid"`
		NftLdefIndex string `json:"nftLdefIndex"`
	}
	type nftTransferResponse struct {
		Status int `json:"status"`
		Action string `json:"action"`
		Reason string `json:"reason"`
		SenderUuid string `json:"senderUuid"`
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
			Action:        common.ACTION_NFT_TRANSFER,
			SenderUuid: testWebSiteUserUuid,
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
				// wait for transaction to finish
				<-time.After(1*time.Second)
				return
			default:
				t.Error("wrong action")
				return
			}
		}
	})
}
