package mobile

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"log"
	"net/http"
)

type WebSocketHandler struct {
	beego.Controller
	M *Manager
}

func (this *WebSocketHandler) Get() {
	upgrader:=websocket.Upgrader{}
	upgrader.CheckOrigin = func(rq *http.Request) bool { return true }
	conn,err:= upgrader.Upgrade(this.Ctx.ResponseWriter,this.Ctx.Request, nil)
	if err!=nil {
		logs.Error(err.Error())
		this.Data["json"] = &ConnectErrorResponse{
			Reason: err.Error(),
		}
		this.ServeJSON()
		return
	}
	defer conn.Close()

	m:= this.M
	c := client.NewClient()
	c.Conn = conn

	//m.connectInHandler(c)

	go c.Sender()
	for {
		_, data, err := conn.ReadMessage()
		var kvs map[string]interface{}
		if err != nil {
			log.Println(err.Error())
			break;
		}
		json.Unmarshal(data, &kvs)
		actionI, ok := kvs["action"]
		if !ok {
			logs.Error("action not exist")
			continue
		}
		eventI, ok := kvs["event"]
		if !ok {
			logs.Error("event not exist")
			continue
		}
		actIdI, ok := kvs["actId"]
		if !ok {
			logs.Error("actId not exist")
			continue
		}

		action := actionI.(string)
		event := eventI.(string)
		actId :=actIdI.(string)
		bq:= &RQBaseInfo{
			Event:event,
			ActId: actId,
			Action:action,
		}
		if h, ok := m.handlers[action]; ok {
			go h(c, bq, data)
		} else {
			log.Println("unknown message")
		}
	}
}


