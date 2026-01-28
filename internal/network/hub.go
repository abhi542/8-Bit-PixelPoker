package network

// ClientMessage wraps the raw message with the sender source
type ClientMessage struct {
	Client *Client
	Msg    IncomingMessage
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	incoming chan ClientMessage

	// Outbound messages to all clients.
	broadcast chan OutgoingMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Table Manager (Handling Game Logic)
	tableManager TableManager
}

type TableManager interface {
	HandleMessage(client *Client, msg IncomingMessage)
	HandleDisconnect(client *Client)
}

func NewHub(tm TableManager) *Hub {
	return &Hub{
		incoming:     make(chan ClientMessage),
		broadcast:    make(chan OutgoingMessage),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		clients:      make(map[*Client]bool),
		tableManager: tm,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				// Notify manager of disconnect
				if h.tableManager != nil {
					h.tableManager.HandleDisconnect(client)
				}
			}
		case msg := <-h.incoming:
			// Route to Game Engine
			if h.tableManager != nil {
				h.tableManager.HandleMessage(msg.Client, msg.Msg)
			}
		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
