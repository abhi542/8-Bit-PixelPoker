package network

import (
	"encoding/json"
	"log"
	"sync"

	"poker-server/internal/lobby"
	"poker-server/internal/poker"

	"github.com/google/uuid"
)

// Server handles all network logic and game management.
type Server struct {
	Hub          *Hub
	LobbyManager *lobby.Manager

	// Map Client -> LobbyCode
	clientLobby map[*Client]string
	mu          sync.RWMutex
}

func NewServer(hub *Hub) *Server {
	return &Server{
		Hub:          hub,
		LobbyManager: lobby.NewManager(),
		clientLobby:  make(map[*Client]string),
	}
}

// HandleMessage processes protocols for Lobby and Game interactions.
func (s *Server) HandleMessage(client *Client, msg IncomingMessage) {
	switch msg.Type {
	case MsgCreateLobby:
		s.handleCreateLobby(client, msg.Payload)
	case MsgJoinLobby:
		s.handleJoinLobby(client, msg.Payload)
	case MsgStartGame:
		s.handleStartGame(client)
	case MsgAction:
		s.handleGameAction(client, msg.Payload)
	case MsgReconnect:
		s.handleReconnect(client, msg.Payload)
	default:
		s.sendError(client, "UNKNOWN_TYPE", "Unknown message type")
	}
}

func (s *Server) HandleDisconnect(client *Client) {
	log.Printf("Client disconnected: %v (User: %s)", client, client.UserID)

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clientLobby, client)

	// Logic to notify lobby of disconnection could go here (e.g. mark seat as Disconnected)
	// For MVP, we just leave them hanging in the game state.
}

func (s *Server) handleReconnect(client *Client, payload json.RawMessage) {
	var req ReconnectPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		s.sendError(client, "INVALID_PAYLOAD", "Bad Reconnect Payload")
		return
	}

	l, ok := s.LobbyManager.GetLobby(req.LobbyCode)
	if !ok {
		s.sendError(client, "LOBBY_NOT_FOUND", "Lobby does not exist")
		return
	}

	// Verify player exists
	if !l.HasPlayer(req.PlayerID) {
		s.sendError(client, "PLAYER_NOT_FOUND", "Player not in lobby")
		return
	}

	// Bind Client
	client.UserID = req.PlayerID

	s.mu.Lock()
	s.clientLobby[client] = l.Code
	s.mu.Unlock()

	// Send Success
	s.sendJSON(client, MsgJoinLobby, map[string]string{
		"code":      l.Code,
		"player_id": req.PlayerID,
		"status":    "RECONNECTED",
	})

	// Send Game State
	if l.GameTable != nil {
		s.sendJSON(client, MsgState, l.GameTable)
	}
}

func (s *Server) handleCreateLobby(client *Client, payload json.RawMessage) {
	var req CreateLobbyPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		s.sendError(client, "INVALID_PAYLOAD", "Bad Create Payload")
		return
	}

	cfg := lobby.LobbyConfig{
		MaxPlayers: req.MaxPlayers,
		SmallBlind: 10,
		BigBlind:   20,
	}
	if cfg.MaxPlayers < 2 {
		cfg.MaxPlayers = 2
	}
	if cfg.MaxPlayers > 9 {
		cfg.MaxPlayers = 9
	}

	l := s.LobbyManager.CreateLobby(cfg)

	s.sendJSON(client, MsgCreateLobby, map[string]string{
		"code": l.Code,
		"id":   l.ID,
	})
}

func (s *Server) handleJoinLobby(client *Client, payload json.RawMessage) {
	var req JoinLobbyPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		s.sendError(client, "INVALID_PAYLOAD", "Bad Join Payload")
		return
	}

	l, ok := s.LobbyManager.GetLobby(req.Code)
	if !ok {
		s.sendError(client, "LOBBY_NOT_FOUND", "Lobby does not exist")
		return
	}

	// Generate persistent ID if new
	pid := uuid.New().String()
	client.UserID = pid

	if err := l.Join(client.UserID, req.Username); err != nil {
		s.sendError(client, "JOIN_FAILED", err.Error())
		return
	}

	s.mu.Lock()
	s.clientLobby[client] = l.Code
	s.mu.Unlock()

	s.sendJSON(client, MsgJoinLobby, map[string]string{
		"code":      l.Code,
		"player_id": pid,
		"status":    "JOINED",
	})
}

func (s *Server) handleStartGame(client *Client) {
	s.mu.RLock()
	code, ok := s.clientLobby[client]
	s.mu.RUnlock()

	if !ok {
		s.sendError(client, "NO_LOBBY", "Not in a lobby")
		return
	}

	l, _ := s.LobbyManager.GetLobby(code)
	if l.State == lobby.StateInGame {
		s.sendError(client, "GAME_RUNNING", "Game already running")
		return
	}

	table, err := l.StartGame()
	if err != nil {
		s.sendError(client, "START_FAILED", err.Error())
		return
	}

	s.broadcastToLobby(l.Code, MsgState, table)
}

func (s *Server) handleGameAction(client *Client, payload json.RawMessage) {
	s.mu.RLock()
	code, ok := s.clientLobby[client]
	s.mu.RUnlock()

	if !ok {
		s.sendError(client, "NO_LOBBY", "Not in a lobby")
		return
	}

	l, _ := s.LobbyManager.GetLobby(code)
	if l.GameTable == nil {
		s.sendError(client, "NO_GAME", "Game not started")
		return
	}

	var actReq ActionPayload
	if err := json.Unmarshal(payload, &actReq); err != nil {
		return
	}

	seatIdx := -1
	for _, seat := range l.GameTable.Seats {
		if seat != nil && seat.PlayerID == client.UserID {
			seatIdx = seat.Index
			break
		}
	}

	if seatIdx == -1 {
		s.sendError(client, "NOT_SEATED", "You are not valid player")
		return
	}

	actType := parseActionType(actReq.ActionType)
	act := poker.Action{Type: actType, Amount: actReq.Amount}

	if err := l.GameTable.HandleAction(seatIdx, act); err != nil {
		s.sendError(client, "ACTION_FAILED", err.Error())
		return
	}

	s.broadcastToLobby(l.Code, MsgState, l.GameTable)
}

// broadcastToLobby sends a message to all players in a lobby.
func (s *Server) broadcastToLobby(lobbyCode string, msgType string, payload interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for client, code := range s.clientLobby {
		if code == lobbyCode {
			// Generate Payload
			var finalPayload interface{} = payload

			// If payload is Table, Sanitize
			if table, ok := payload.(*poker.Table); ok {
				finalPayload = table.GetViewFor(client.UserID)
			}

			msg := OutgoingMessage{
				Type:    msgType,
				Payload: finalPayload,
			}

			// Non-blocking send
			select {
			case client.send <- msg:
			default:
				// Dropped
			}
		}
	}
}

func (s *Server) sendJSON(c *Client, msgType string, payload interface{}) {
	c.send <- OutgoingMessage{Type: msgType, Payload: payload}
}

func (s *Server) sendError(c *Client, code, msg string) {
	c.send <- OutgoingMessage{
		Type:    MsgError,
		Payload: ErrorPayload{Code: code, Message: msg},
	}
}

func parseActionType(s string) poker.ActionType {
	switch s {
	case "CHECK":
		return poker.ActionCheck
	case "BET":
		return poker.ActionBet
	case "CALL":
		return poker.ActionCall
	case "RAISE":
		return poker.ActionRaise
	case "FOLD":
		return poker.ActionFold
	case "ALL_IN":
		return poker.ActionAllIn
	default:
		return poker.ActionFold
	}
}
