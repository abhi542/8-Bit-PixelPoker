package lobby

import (
	"poker-server/internal/poker"
	"sync"
)

type LobbyState int

const (
	StateWaiting LobbyState = iota
	StateReady
	StateInGame
)

type Lobby struct {
	ID        string
	Code      string
	Config    LobbyConfig
	State     LobbyState
	Players   map[string]string // ClientID -> Username
	GameTable *poker.Table

	mu sync.RWMutex
}

func (l *Lobby) Mu() *sync.RWMutex {
	return &l.mu
}

type LobbyConfig struct {
	MaxPlayers int
	SmallBlind int64
	BigBlind   int64
}

func NewLobby(code string, config LobbyConfig) *Lobby {
	return &Lobby{
		ID:      code, // Use code as ID for simplicity in MVP
		Code:    code,
		Config:  config,
		State:   StateWaiting,
		Players: make(map[string]string),
	}
}

func (l *Lobby) Join(clientID, username string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.Players) >= l.Config.MaxPlayers {
		return ErrLobbyFull
	}

	l.Players[clientID] = username

	// Auto-ready check?
	if len(l.Players) >= 2 {
		l.State = StateReady
	}

	return nil
}

func (l *Lobby) HasPlayer(playerID string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, ok := l.Players[playerID]
	return ok
}

func (l *Lobby) StartGame() (*poker.Table, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.Players) < 2 {
		return nil, ErrNotEnoughPlayers
	}

	// Create Table
	table := poker.NewTable(l.ID, poker.TableConfig{
		SmallBlind: l.Config.SmallBlind,
		BigBlind:   l.Config.BigBlind,
	})

	// Seat players randomly or in order
	seatIdx := 0
	for cid, name := range l.Players {
		// BuyIn defaults to 100BB for now
		table.JoinSeat(seatIdx, cid, name, l.Config.BigBlind*100)
		seatIdx++
	}

	l.GameTable = table
	l.State = StateInGame

	table.StartHand()

	return table, nil
}
