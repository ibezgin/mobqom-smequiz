package game

import (
	"crypto/rand"
	"sync"
)

type GameManager struct {
	rooms        map[string]*Room
	mu           sync.Mutex
	createRoomCh chan *Room
	deleteRoomCh chan *Room
}

func NewGameManager() *GameManager {
	return &GameManager{
		rooms: make(map[string]*Room),
	}
}

func generateRoomID() string {
	return "room_" + rand.Text()[:9]

}

func (gm *GameManager) CreateRoom() *Room {
	newRoom := NewRoom(generateRoomID())
	gm.mu.Lock()
	defer gm.mu.Unlock()
	gm.rooms[newRoom.ID] = newRoom
	return newRoom
}

func (gm *GameManager) GetRoom(roomID string) (*Room, bool) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	room, exists := gm.rooms[roomID]
	return room, exists
}

func (gm *GameManager) DeleteRoom(roomID string) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	delete(gm.rooms, roomID)
}
