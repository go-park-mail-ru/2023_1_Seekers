package ws

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	pkgSmtp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/smtp"
)

type Subscription struct {
	conn      *connection
	UserEmail string
}

type Hub struct {
	cfg        *config.Config
	rooms      map[string]map[*connection]bool
	broadcast  chan WsItem
	register   chan Subscription
	unregister chan Subscription
}

func NewHub(c *config.Config) *Hub {
	return &Hub{
		broadcast:  make(chan WsItem),
		register:   make(chan Subscription),
		unregister: make(chan Subscription),
		rooms:      make(map[string]map[*connection]bool),
		cfg:        c,
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

func (h *Hub) SendNotifications(message *models.MessageInfo, mailUC mail.UseCaseI) {
	item := WsItem{
		messageInfo: *message,
		userEmail:   message.FromUser.Email,
	}
	domain, err := pkgSmtp.ParseDomain(message.FromUser.Email)
	if err != nil {
		return
	}

	if domain == h.cfg.Mail.PostDomain {
		h.broadcast <- item
	}

	item.messageInfo.Seen = false

	for _, recipient := range item.messageInfo.Recipients {
		domain, err = pkgSmtp.ParseDomain(recipient.Email)
		if err != nil {
			return
		}

		if domain == h.cfg.Mail.PostDomain {
			toEmail, err := mailUC.GetOwnerEmailByFakeEmail(recipient.Email)
			if err != nil {
				return
			}

			item.userEmail = toEmail
			h.broadcast <- item
		}
	}
}
