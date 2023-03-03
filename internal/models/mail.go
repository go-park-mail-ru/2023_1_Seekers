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
	FolderID uint64 `json:"folder_id"`
	Name     string `json:"name"`
	UserID   uint64 `json:"user_id"`
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
