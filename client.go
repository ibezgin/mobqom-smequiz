package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type MsgType string

const (
	MsgType_Broadcast MsgType = "broadcast"
)

type ReqMsg struct {
	MsgType MsgType
	Client  *Client
	Data    string
}
type ResMsg struct {
	MsgType  MsgType
	Data     string
	SenderId string
}

func NewResMsg(msg *ReqMsg) *ResMsg {
	return &ResMsg{
		MsgType:  msg.MsgType,
		Data:     msg.Data,
		SenderId: msg.Client.ID,
	}
}

type Client struct {
	ID   string
	mu   *sync.RWMutex
	conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	ID := rand.Text()
	return &Client{
		ID:   ID,
		mu:   new(sync.RWMutex),
		conn: conn,
	}
}

func (c *Client) readMsgLoop(srv *Server) {
	defer func() {
		c.conn.Close()
		srv.leaveServerCh <- c
	}()
	for {
		_, b, err := c.conn.ReadMessage()
		if err != nil {

			return
		}
		msg := new(ReqMsg)
		err = json.Unmarshal(b, msg)
		if err != nil {
			fmt.Printf("error unmarshall the msg %v\n", err)
			continue
		}
		msg.Client = c
		srv.broadcastCh <- msg
		// fmt.Println((string)p)
	}
}
