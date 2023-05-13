package ws

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
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

func (s Subscription) readPump(hub *Hub, mailUC mail.UseCaseI) {
	c := s.conn
	defer func() {
		hub.unregister <- s
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	if err := c.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}
	c.ws.SetPongHandler(func(string) error {
		if err := c.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			return nil
		}
		return nil
	})
	for {
		var form models.FormMessage

		err := c.ws.ReadJSON(&form)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				//log.Printf("error: %v", err)
			}

			break
		}

		message, err := mailUC.SendMessage(form)
		if err != nil {
			return
		}

		item := WsItem{
			messageInfo: *message,
			userEmail:   message.FromUser.Email,
		}

		hub.broadcast <- item
		item.messageInfo.Seen = false

		for _, recipient := range item.messageInfo.Recipients {
			item.userEmail = recipient.Email
			hub.broadcast <- item
		}
	}
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

func (s Subscription) writePump(hub *Hub) {
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

func ServeWs(w http.ResponseWriter, r *http.Request, userEmail string, hub *Hub, mailUC mail.UseCaseI) {
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

	go sub.writePump(hub)
	go sub.readPump(hub, mailUC)
}
