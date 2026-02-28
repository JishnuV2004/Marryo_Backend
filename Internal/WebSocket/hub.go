package websocket

// import "log"

type Hub struct {
	Clients map[uint]*Client
	Register chan *Client
	Unregister chan *Client
	Broadcast chan MessagePayload
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[uint]*Client),
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast: make(chan MessagePayload),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <- h.Register :
			h.Clients[client.UserID] = client
			// log.Println("Registered:", client.UserID)

		case client := <- h.Unregister :
			delete(h.Clients, client.UserID)
			// log.Println("Unregistered:", client.UserID)

		case msg := <- h.Broadcast :
			// log.Println("Broadcast to:", msg.ToUserID)

			if receiver, ok := h.Clients[msg.ToUserID]; ok {
				receiver.Send <- msg
			} else {
				// log.Println("Receiver not connected:", msg.ToUserID)
			}
		}
	}
}