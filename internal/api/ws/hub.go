package ws

type Subscription struct {
	conn      *connection
	UserEmail string
}

type Hub struct {
	rooms      map[string]map[*connection]bool
	broadcast  chan WsItem
	register   chan Subscription
	unregister chan Subscription
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan WsItem),
		register:   make(chan Subscription),
		unregister: make(chan Subscription),
		rooms:      make(map[string]map[*connection]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.UserEmail]
			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[s.UserEmail] = connections
			}

			h.rooms[s.UserEmail][s.conn] = true
		case s := <-h.unregister:
			connections := h.rooms[s.UserEmail]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)

					if len(connections) == 0 {
						delete(h.rooms, s.UserEmail)
					}
				}
			}
		case m := <-h.broadcast:
			connections := h.rooms[m.userEmail]

			for c := range connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(connections, c)

					if len(connections) == 0 {
						delete(h.rooms, m.userEmail)
					}
				}
			}
		}
	}
}
