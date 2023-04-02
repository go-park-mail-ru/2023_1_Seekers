package models

import "html"

type Folder struct {
	FolderID       uint64 `json:"folder_id" gorm:"primaryKey"`
	UserID         uint64 `json:"-"`
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

func (form *FormMessage) Sanitize() {
	form.Title = html.EscapeString(form.Title)
	form.Text = html.EscapeString(form.Text)
	for i, s := range form.Recipients {
		form.Recipients[i] = html.EscapeString(s)
	}
}

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
