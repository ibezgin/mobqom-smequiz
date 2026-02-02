package game

import "sync"

type Room struct {
	ID       string
	players  map[string]*Player
	Stage    Stage
	mu       sync.Mutex
	resMsgCh chan *ResMsg
	reqMsgCh chan *ReqMsg
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		players: map[string]*Player{},
		Stage:   Stage_Waiting,
	}
}

func (r *Room) AddPlayer(p *Player) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.players[p.ID] = p
}

func (r *Room) RemovePlayer(p *Player) {
	delete(r.players, p.ID)
}

func (r *Room) ChangeStage(stage Stage) {
	switch stage {
	case Stage_Waiting, Stage_Active, Stage_Ended:
		r.Stage = stage
	}
}
