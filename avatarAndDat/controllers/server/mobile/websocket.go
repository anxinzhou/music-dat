package mobile

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/chainHelper"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/transactionQueue"
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
	chainHandler *chainHelper.ChainHandler
	TransactionQueue *transactionQueue.TransactionQueue
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

func (m* Manager) SetTransactionQueue(transactionQueue *transactionQueue.TransactionQueue) {
	m.TransactionQueue = transactionQueue
}

func (m *Manager) RegisterHandler(action string, h handler) {
	m.handlers[action] = h
}

func (m *Manager) UnregisterHandler(action string, h handler) {
	delete(m.handlers, action)
}

func (m *Manager) SetChainHandler(chainHandler *chainHelper.ChainHandler) {
	m.chainHandler = chainHandler
}

func (m* Manager) ChainHandler() *chainHelper.ChainHandler{
	return m.chainHandler
}

func (m *Manager) Init() {
	m.RegisterHandler(common.ACTION_MP_LIST,m.GetMPListHandler)
	m.RegisterHandler(common.ACTION_ITEM_DETAILS,m.ItemDetailsHandler)
	m.RegisterHandler(common.ACTION_NFT_PUCHASE_CONFIRM,m.PurchaseConfirmHandler)
	m.RegisterHandler(common.ACTION_TOKENBUY_PAID,m.TokenBuyPaidHandler)
	m.RegisterHandler(common.ACTION_NFT_DISPLAY,m.NFTDisplayHandler)
	m.RegisterHandler(common.ACTION_MARKET_USER_LIST,m.MarketUserListHandler)
	m.RegisterHandler(common.ACTION_NFT_PURCHASE_HISTORY,m.NFTPurchaseHistoryHandler)
	m.RegisterHandler(common.ACTION_NFT_SHOPPING_CART_CHANGE, m.ShoppingCartChangeHandler)
	m.RegisterHandler(common.ACTION_NFT_SHOPPING_CART_LIST,m.ShoppingCartListHandler)
	m.RegisterHandler(common.ACTION_NFT_TRANSFER,m.NFTTransferHandler)
	m.RegisterHandler(common.ACTION_NFT_BIND_WALLET,m.BindWalletHandler)
	m.RegisterHandler(common.ACTION_NFT_SET_NICKNAME,m.SetNicknameHandler)
	m.RegisterHandler(common.ACTION_IS_NICKNAME_DUPLICATED,m.IsNicknameDuplicatedHandler)
	m.RegisterHandler(common.ACTION_FOLLOW_LIST,m.FollowListHandler)
	m.RegisterHandler(common.ACTION_FOLLOW_LIST_OPERATION,m.FollowListOperationHandler)
	m.RegisterHandler(common.ACTION_IS_NICKNAME_SET,m.IsNicknameSetHandler)
	m.RegisterHandler(common.ACTION_USER_MARKET_INFO,m.UserMarketInfoHandler)
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


