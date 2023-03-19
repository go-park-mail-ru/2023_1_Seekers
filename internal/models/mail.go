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

type UserInfo struct {
	UserID    uint64 `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type MessageInfo struct {
	MessageID  uint64     `json:"message_id"`
	FromUser   UserInfo   `json:"from_user_id" gorm:"embedded;embeddedPrefix:;embeddedPrefix:from_"`
	Recipients []UserInfo `json:"recipients" gorm:"-"`
	Title      string     `json:"title"`
	CreatedAt  string     `json:"created_at"`
	Text       string     `json:"text"`
	Seen       bool       `json:"seen"`
	Favorite   bool       `json:"favorite"`
	Deleted    bool       `json:"deleted"`
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
	Folder         Folder        `json:"folder"`
	Messages       []MessageInfo `json:"messages"`
	MessagesUnseen int           `json:"messages_unseen"`
	MessagesCount  int           `json:"messages_count"`
}

type FoldersResponse struct {
	Folders []Folder `json:"folders"`
}
