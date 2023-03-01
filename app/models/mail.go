package models

//db

type Message struct {
	MessageID    uint64
	UserID       uint64
	CreatingDate string
	Title        string
	Text         string
}

type Recipient struct {
	MessageID uint64
	UserID    uint64
}

type Folder struct {
	FolderID uint64
	Name     string
	UserID   uint64
}

type Box struct {
	FolderID  uint64
	MessageID uint64
}

type State struct {
	UserID    uint64
	MessageID uint64
	Read      bool
	Favorite  bool
	Send      bool
}

type User struct {
	UserID   uint64
	Email    string
	Password string
}

// delivery

//type MessageInfo struct {
//	MessageID    uint64   `json:"message_id"`
//	FromUserID   uint64   `json:"from_user_id"`
//	ToUsersIDs   []uint64 `json:"to_users_ids"`
//	CreatingDate string   `json:"creating_date"`
//	Title        string   `json:"title"`
//	Text         string   `json:"text"`
//	Read         bool     `json:"read"`
//	Favorite     bool     `json:"favorite"`
//}

type IncomingMessage struct {
	MessageID    uint64 `json:"message_id"`
	FromUser     string `json:"from_user"`
	CreatingDate string `json:"creating_date"`
	Title        string `json:"title"`
	Text         string `json:"text"`
	Read         bool   `json:"read"`
	Favorite     bool   `json:"favorite"`
}

type OutgoingMessage struct {
	MessageID    uint64   `json:"message_id"`
	ToUsers      []string `json:"to_users"`
	CreatingDate string   `json:"creating_date"`
	Title        string   `json:"title"`
	Text         string   `json:"text"`
	Read         bool     `json:"read"`
	Favorite     bool     `json:"favorite"`
}

type InboxMessages struct {
	Folders  []Folder          `json:"folders"`
	Messages []IncomingMessage `json:"messages"`
}

type OutboxMessages struct {
	Folders  []Folder          `json:"folders"`
	Messages []OutgoingMessage `json:"messages"`
}
