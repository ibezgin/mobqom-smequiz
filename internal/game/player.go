package game

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID     string
	Conn   *websocket.Conn
	Screen Screen
	score  uint
}

func NewPlayer(id string, conn *websocket.Conn) *Player {
	return &Player{
		ID:     id,
		Conn:   conn,
		Screen: Screen_WaitPlayers,
		score:  0,
	}
}

func (p *Player) ReadMessage(gm *GameManager) error {
	for {
		_, dt, err := p.Conn.ReadMessage()
		if err != nil {
			return err
		}
		msg := new(ReqMsg)
		json.Unmarshal(dt, msg)
		msg.Player = p

		switch msg.MsgType {
		case MsgType_JoinRoom:
			p.joinRoom(gm, msg)
		case MsgType_LeaveRoom:
			p.leaveRoom(gm, msg)
		}
	}
}
func (p *Player) joinRoom(gm *GameManager, msg *ReqMsg) {
	room, exist := gm.GetRoom(msg.RoomID)
	if !exist {
		fmt.Printf("Room %s does not exist\n", msg.RoomID)
		return
	}
	room.AddPlayer(p)
	fmt.Printf("Player %s joined room %s\n", p.ID, msg.RoomID)
}

func (p *Player) leaveRoom(gm *GameManager, msg *ReqMsg) {
	room, exist := gm.GetRoom(msg.RoomID)
	if !exist {
		fmt.Printf("Room %s does not exist\n", msg.RoomID)
		return
	}

	room.RemovePlayer(p)
	fmt.Printf("Player %s left room %s\n", p.ID, msg.RoomID)

	if playersCount := len(room.players); playersCount == 0 {
		gm.DeleteRoom(room.ID)
		fmt.Printf("Room %s deleted as it became empty\n", msg.RoomID)
	}
}
