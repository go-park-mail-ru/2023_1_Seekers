package models

//db

type Folder struct {
	FolderID       uint64 `json:"folder_id"`
	UserID         uint64 `json:"user_id"`
	LocalName      string `json:"folder_slug"`
	Name           string `json:"name"`
	MessagesUnseen int    `json:"messages_unseen"`
	MessagesCount  int    `json:"messages_count"`
}

type MessageInfo struct {
	MessageID        uint64       `json:"message_id"`
	FromUser         UserInfo     `json:"from_user_id" gorm:"embedded;embeddedPrefix:from_"`
	Recipients       []UserInfo   `json:"recipients" gorm:"-"`
	Title            string       `json:"title"`
	CreatedAt        string       `json:"created_at"`
	Text             string       `json:"text"`
	ReplyToMessageID *uint64      `json:"-" gorm:"null"`
	ReplyTo          *MessageInfo `json:"reply_to" gorm:"-"`
	Seen             bool         `json:"seen"`
	Favorite         bool         `json:"favorite"`
	Deleted          bool         `json:"deleted"`
}

type FormMessage struct {
	Recipients       []string `json:"recipients" validate:"required"`
	Title            string   `json:"title" validate:"required"`
	Text             string   `json:"text" validate:"required"`
	ReplyToMessageID *uint64  `json:"reply_to"`
}

//type IncomingMessage struct {
//	MessageID    uint64 `json:"message_id"`
//	FromUser     string `json:"from_user"`
//	CreatingDate string `json:"creating_date"`
//	Title        string `json:"title"`
//	Text         string `json:"text"`
//	Read         bool   `json:"read"`
//	Favorite     bool   `json:"favorite"`
//}
//
//type OutgoingMessage struct {
//	MessageID    uint64   `json:"message_id"`
//	ToUsers      []string `json:"to_users"`
//	CreatingDate string   `json:"creating_date"`
//	Title        string   `json:"title"`
//	Text         string   `json:"text"`
//	Read         bool     `json:"read"`
//	Favorite     bool     `json:"favorite"`
//}

type FolderResponse struct {
	Folder   Folder        `json:"folder"`
	Messages []MessageInfo `json:"messages"`
}

type FoldersResponse struct {
	Folders []Folder `json:"folders"`
	Count   int      `json:"count"`
}

type MessageResponse struct {
	Message MessageInfo `json:"message"`
}
