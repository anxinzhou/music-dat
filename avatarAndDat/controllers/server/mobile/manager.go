package mobile

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
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

func (m *Manager) errorHandler(c *client.Client, bq *RQBaseInfo, err error) {
	bq.Event = "failed"
	res := &ErrorResponse{
		RQBaseInfo: *bq,
		Reason:     err.Error(),
	}
	resWrapper, err := json.Marshal(res)
	if err != nil {
		panic(err)
		return
	}
	c.Send(resWrapper)
}

func (m *Manager) wrapperAndSend(c *client.Client, bq *RQBaseInfo, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	c.Send(data)
}