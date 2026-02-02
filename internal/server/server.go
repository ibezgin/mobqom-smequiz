package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/ibezgin/mobqom-smequiz/internal/game"
)

type NewRoomResponse struct {
	ID string `json:"id"`
}

type Server struct {
	joinServerCh  chan *game.Player
	leaveServerCh chan *game.Player
}

func NewWSServer() *Server {
	return &Server{
		joinServerCh:  make(chan *game.Player),
		leaveServerCh: make(chan *game.Player),
	}
}

func (s *Server) HandleWs(gm *game.GameManager, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("error upgrade websocket connection\n")
	}
	player := game.NewPlayer(r.RemoteAddr, conn)
	s.joinServerCh <- player
	go player.ReadMessage(gm)
}
func (s *Server) HandleCreateRoom(gm *game.GameManager, w http.ResponseWriter, r *http.Request) {
	newRoom := gm.CreateRoom()
	fmt.Printf("Room created with ID: %s\n", newRoom.ID)
	w.Header().Set("Content-Type", "application/json")
	resp := NewRoomResponse{ID: newRoom.ID}
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) AcceptLoop() {
	for {
		select {
		case player := <-s.joinServerCh:
			s.joinServer(player)
		case player := <-s.leaveServerCh:
			s.leaveServer(player)
		}
	}
}

func (s *Server) joinServer(player *game.Player) {
	fmt.Printf("Client %s joined the server\n", player.ID)
}
func (s *Server) leaveServer(player *game.Player) {
	fmt.Printf("Client %s left the server\n", player.ID)
}
