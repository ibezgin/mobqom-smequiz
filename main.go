package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	WSPort = ":3223"
)

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

type Server struct {
	clients       map[string]*Client
	mu            *sync.RWMutex
	joinServerCh  chan *Client
	leaveServerCh chan *Client
}

func NewServer() *Server {
	return &Server{
		clients:       map[string]*Client{},
		mu:            new(sync.RWMutex),
		joinServerCh:  make(chan *Client, 64),
		leaveServerCh: make(chan *Client, 64),
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
	// fmt.Println("clients count:", len(s.clients))
}
func (s *Server) joinServer(c *Client) {
	s.clients[c.ID] = c
	fmt.Printf("client joined the server, cId = %s\n", c.ID)
}
func (s *Server) leaveSever(c *Client) {
	delete(s.clients, c.ID)
	fmt.Printf("client left the server, cId = %s\n", c.ID)
}

func (s *Server) AcceptLoop() {
	for {
		select {
		case c := <-s.joinServerCh:
			s.joinServer(c)
		case c := <-s.leaveServerCh:
			s.leaveSever(c)
		}
	}
}
func createWSServer() {
	s := NewServer()
	go s.AcceptLoop()
	go fmt.Printf("Starting ws server on port %s\n", WSPort)
	http.HandleFunc("/", s.handleWs)
	log.Fatal(http.ListenAndServe(WSPort, nil))
}
func main() {
	createWSServer()
}
