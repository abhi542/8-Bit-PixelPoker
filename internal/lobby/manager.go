package lobby

import (
	"errors"
	"sync"
)

var (
	ErrLobbyFull        = errors.New("lobby is full")
	ErrLobbyNotFound    = errors.New("lobby not found")
	ErrNotEnoughPlayers = errors.New("not enough players to start")
)

type Manager struct {
	lobbies map[string]*Lobby
	mu      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		lobbies: make(map[string]*Lobby),
	}
}

func (m *Manager) CreateLobby(config LobbyConfig) *Lobby {
	m.mu.Lock()
	defer m.mu.Unlock()

	code := GenerateCode()
	// Ensure Uniqueness
	for _, exists := m.lobbies[code]; exists; _, exists = m.lobbies[code] {
		code = GenerateCode()
	}

	lobby := NewLobby(code, config)
	m.lobbies[code] = lobby
	return lobby
}

func (m *Manager) GetLobby(code string) (*Lobby, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	l, ok := m.lobbies[code]
	return l, ok
}

func (m *Manager) RemoveLobby(code string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.lobbies, code)
}
