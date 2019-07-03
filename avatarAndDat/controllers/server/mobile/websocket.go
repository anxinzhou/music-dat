package mobile

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"log"
	"net/http"
)

const CLIENT_BUFFER = 1024

type handler func(user *client.Client, action string, data []byte)

type Manager struct {
	loginClients map[string]*client.Client
	clients      map[*client.Client]bool
	broadcast    chan []byte
	register     chan *client.Client
	unregister   chan *client.Client
	handlers     map[string]handler
	chainHandler *common.ChainHandler
}

func NewManager() *Manager {
	m := &Manager{
		loginClients: make(map[string]*client.Client),
		clients:      make(map[*client.Client]bool),
		register:     make(chan *client.Client, CLIENT_BUFFER),
		unregister:   make(chan *client.Client, CLIENT_BUFFER),
		broadcast:    make(chan []byte, CLIENT_BUFFER),
		handlers:     make(map[string]handler),
	}
	return m
}

func (m *Manager) RegisterHandler(action string, h handler) {
	m.handlers[action] = h
}

func (m *Manager) UnregisterHandler(action string, h handler) {
	delete(m.handlers, action)
}

func (m *Manager) SetChainHandler(chainHandler *common.ChainHandler) {
	m.chainHandler = chainHandler
}

func (m* Manager) ChainHandler() *common.ChainHandler{
	return m.chainHandler
}

func (m *Manager) Init() {
	m.RegisterHandler("mp_list",m.GetMPListHandler)
	m.RegisterHandler("item_details",m.ItemDetailsHandler)
	m.RegisterHandler("NFT_purchase_confirm",m.PurchaseConfirmHandler)
	m.RegisterHandler("tokenbuy_paid",m.TokenBuyPaidHandler)
	m.RegisterHandler("NFT_display",m.NFTDisplayHandler)
	m.RegisterHandler("market_user_list",m.MarketUserListHandler)
	m.RegisterHandler("nft_purchase_history",m.NFTPurchaseHistoryHandler)
	m.RegisterHandler("nft_shopping_cart_change", m.ShoppingCartChangeHandler)
	m.RegisterHandler("nft_shopping_cart_list",m.ShoppingCartListHandler)
	m.RegisterHandler("nft_transfer",m.NFTTransferHandler)
	m.RegisterHandler("bind_wallet",m.BindWalletHandler)
	m.RegisterHandler("set_nickname",m.SetNicknameHandler)
	m.RegisterHandler("is_nickname_duplicated",m.IsNicknameDuplicatedHandler)
	m.RegisterHandler("follow_list",m.FollowListHandler)
	m.RegisterHandler("follow_list_operation",m.FollowListOperationHandler)
	m.RegisterHandler("is_nickname_set",m.IsNicknameSetHandler)
}

func (m *Manager) errorHandler(c *client.Client, action string, err error) {
	type response struct {
		Status int `json:"status"`
		Action string `json:"action"`
		Reason string `json:"reason"`
	}
	res:= response{
		Status: common.RESPONSE_STATUS_FAIL,
		Action: action,
		Reason: err.Error(),
	}

	resWrapper, err := json.Marshal(&res)
	if err != nil {
		panic(err)
	}
	c.Send(resWrapper)
}

func (m *Manager) wrapperAndSend(c *client.Client, action string, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	c.Send(data)
}

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
		err:=errors.New("cannot cannect")
		this.Data["json"]= err
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

		action := actionI.(string)
		if h, ok := m.handlers[action]; ok {
			go h(c, action, data)
		} else {
			log.Println("unknown message")
		}
	}
}


