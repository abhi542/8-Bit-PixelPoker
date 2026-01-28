package network

import "encoding/json"

// MessageType constants
const (
	MsgCreateLobby = "CREATE_LOBBY"
	MsgJoinLobby   = "JOIN_LOBBY"
	MsgStartGame   = "START_GAME"
	MsgJoinTable   = "JOIN_TABLE"
	MsgAction      = "ACTION"
	MsgReconnect   = "RECONNECT"
	MsgState       = "STATE"
	MsgError       = "ERROR"
)

// IncomingMessage is the generic envelope for client->server
type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// OutgoingMessage is the generic envelope for server->client
type OutgoingMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	Seq     int64       `json:"seq,omitempty"`
}

// Payload Structs

type JoinTablePayload struct {
	TableID string `json:"table_id"`
	SeatIdx int    `json:"seat_idx"`
}

type CreateLobbyPayload struct {
	MaxPlayers int `json:"max_players"`
}

type JoinLobbyPayload struct {
	Code     string `json:"code"`
	Username string `json:"username"`
}

type ReconnectPayload struct {
	LobbyCode string `json:"lobby_code"`
	PlayerID  string `json:"player_id"`
}

type ActionPayload struct {
	ActionType string `json:"action_type"` // CHECK, BET, etc.
	Amount     int64  `json:"amount"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
