package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

const (
	SEND_BUFFER_SIZE = 32
)

const (
	STATUS_LOGIN = iota
	STATUS_LOGOUT
)

type Client struct {
	Conn       *websocket.Conn
	Address    common.Address
	sendBuffer chan []byte
	Status     uint64
	close      chan int
	mutex      *sync.Mutex // used to protect modify client state
}

func NewClient() *Client {
	return &Client{
		sendBuffer: make(chan []byte, SEND_BUFFER_SIZE),
		close:      make(chan int, 1),
		mutex:      &sync.Mutex{},
	}
}

func (c *Client) Close() {
	c.close <-1
}

func (c *Client) Login(address common.Address) {
	c.Address = address
	c.Status = STATUS_LOGIN
}

func (c *Client) Logout() {
	c.Status = STATUS_LOGOUT
}

func (c *Client) IsLogin() bool {
	return c.Status == STATUS_LOGIN
}

func (c *Client) IsLogout() bool {
	return c.Status == STATUS_LOGOUT
}

func (c *Client) Send(data []byte) {
	c.sendBuffer <- data
}

func (c *Client) Sender() {
	for {
		select {
		case data := <-c.sendBuffer:
			err := c.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		case <-c.close:
			return
		}
	}
}

