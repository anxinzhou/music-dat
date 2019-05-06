package ws

import (
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
)

const CLIENT_BUFFER = 1024

type handler func(user *client.Client, bq *RQBaseInfo, data []byte)

type Manager struct {
	loginClients map[string]*client.Client
	clients      map[*client.Client]bool
	broadcast    chan []byte
	register     chan *client.Client
	unregister   chan *client.Client
	handlers     map[string]handler
	chainHandler *ChainHandler
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

func (m *Manager) SetChainHandler(chainHandler *ChainHandler) {
	m.chainHandler = chainHandler
}

func (m* Manager) ChainHandler() *ChainHandler{
	return m.chainHandler
}

func (m *Manager) Init() {
	m.RegisterHandler("mp_list",m.GetMPList)
	m.RegisterHandler("NFT_purchase_confirm",m.PurchaseConfirmHandler)
	m.RegisterHandler("item_details", m.ItemDetailsHandler)
	m.RegisterHandler("NFT_display",m.NFTDisplayHandler)
	m.RegisterHandler("tokenbuy_paid",m.TokenBuyPaidHandler)
	m.RegisterHandler("market_user_list",m.MarketUserListHandler)
	m.RegisterHandler("user_market_info",m.UserMarketInfoHandler)
}

//func (m *Manager) DisPatchMsg() {
//	for {
//		select {
//		case c := <-m.register:
//			m.clients[c] = true
//			log.Println("a new user connect")
//		case <-m.unregister:
//			//delete(m.clients, c)
//			log.Println("a user unregister disconnect")
//		case message := <-m.broadcast:
//			log.Println("broadcast a message:", string(message))
//			for c, _ := range m.clients {
//				c.Send(message) // an active user may block other user here, fix in the future
//			}
//		}
//	}
//}
