package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	clients       map[string]*Client
	mu            *sync.RWMutex
	joinServerCh  chan *Client
	leaveServerCh chan *Client
	broadcastCh   chan *ReqMsg
}

func NewServer() *Server {
	return &Server{
		clients:       map[string]*Client{},
		mu:            new(sync.RWMutex),
		joinServerCh:  make(chan *Client, 64),
		leaveServerCh: make(chan *Client, 64),
		broadcastCh:   make(chan *ReqMsg, 64),
	}

}

func (s *Server) handleWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error on http conn upgrade %v\n", err)
		return
	}
	client := NewClient(conn)
	s.joinServerCh <- client
	go client.readMsgLoop(s)
}
func (s *Server) joinServer(c *Client) {
	s.clients[c.ID] = c
	fmt.Printf("client joined the server, cId = %s\n", c.ID)
}
func (s *Server) leaveSever(c *Client) {
	delete(s.clients, c.ID)
	fmt.Printf("client left the server, cId = %s\n", c.ID)
}
func (s *Server) broadcast(msg *ReqMsg) {
	cls := []*Client{}
	s.mu.RLock()
	for _, c := range s.clients {
		if c.ID != msg.Client.ID {
			cls = append(cls, c)
		}
	}
	s.mu.RUnlock()
	resp := NewResMsg(msg)
	for _, c := range cls {
		err := c.conn.WriteJSON(resp)
		if err != nil {
			fmt.Printf("error sending message to cID = %s", c.ID)
			continue
		}
	}
	fmt.Println("broadcast was sent")
}

func (s *Server) AcceptLoop() {
	for {
		select {
		case c := <-s.joinServerCh:
			s.joinServer(c)
		case c := <-s.leaveServerCh:
			s.leaveSever(c)
		case msg := <-s.broadcastCh:
			go s.broadcast(msg)
		}
	}
}
