package persistence

import "poker-server/internal/lobby"

type Store interface {
	SaveLobby(l *lobby.Lobby) error
	GetLobby(code string) (*lobby.Lobby, error)
	// RemoveLobby? For now we keep history or specific cleanup later.
}
