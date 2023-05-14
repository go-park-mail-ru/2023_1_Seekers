package ws

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsItem struct {
	messageInfo models.MessageInfo
	userEmail   string
}

type connection struct {
	ws   *websocket.Conn
	send chan WsItem
}

func (c *connection) write(msg models.MessageInfo) error {
	if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return c.ws.WriteJSON(msg)
	}

	return c.ws.WriteJSON(msg)
}

func (c *connection) writeType(mt int) error {
	if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return c.ws.WriteMessage(mt, []byte{})
	}
	return c.ws.WriteMessage(mt, []byte{})
}

func (s Subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case item, ok := <-c.send:
			if !ok {
				if err := c.writeType(websocket.CloseMessage); err != nil {
					return
				}
				return
			}

			if err := c.write(item.messageInfo); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.writeType(websocket.PingMessage); err != nil {
				return
			}
		}
	}
}

func ServeWs(w http.ResponseWriter, r *http.Request, userEmail string, hub *Hub) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//log.Println(err.Error())
		return
	}

	conn := &connection{
		send: make(chan WsItem, 256),
		ws:   ws,
	}
	sub := Subscription{
		conn:      conn,
		UserEmail: userEmail,
	}
	hub.register <- sub

	go sub.writePump()
}
