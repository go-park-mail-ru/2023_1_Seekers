package models

import (
	"github.com/microcosm-cc/bluemonday"
	"html"
)

//go:generate easyjson -all

type Folder struct {
	FolderID       uint64 `json:"folder_id" gorm:"primaryKey"`
	UserID         uint64 `json:"-"`
	LocalName      string `json:"folder_slug"`
	Name           string `json:"name"`
	MessagesUnseen int    `json:"messages_unseen"`
	MessagesCount  int    `json:"messages_count"`
}

type MessageInfo struct {
	MessageID        uint64           `json:"message_id"`
	FromUser         UserInfo         `json:"from_user_id" gorm:"embedded;embeddedPrefix:from_"`
	Recipients       []UserInfo       `json:"recipients" gorm:"-"`
	Attachments      []AttachmentInfo `json:"attachments" gorm:"-"`
	AttachmentsSize  string           `json:"attachmentsSize" gorm:"-"`
	Title            string           `json:"title"`
	CreatedAt        string           `json:"created_at"`
	Text             string           `json:"text"`
	ReplyToMessageID *uint64          `json:"-" gorm:"null"`
	ReplyTo          *MessageInfo     `json:"reply_to" gorm:"-"`
	Seen             bool             `json:"seen"`
	Favorite         bool             `json:"favorite"`
	Deleted          bool             `json:"-"`
	IsDraft          bool             `json:"is_draft"`
}

type Recipients struct {
	Users []UserInfo `json:"users"`
}

type User2Folder struct {
	UserID   uint64
	FolderID uint64
}

type Attachment struct {
	FileName string `json:"fileName"`
	FileData string `json:"fileData"`
}

type AttachmentInfo struct {
	AttachID  uint64 `json:"attachID"`
	FileName  string `json:"fileName"`
	FileData  []byte `json:"-"`
	S3FName   string `json:"-"`
	Type      string `json:"type"`
	SizeStr   string `json:"sizeStr"`
	SizeCount int64  `json:"sizeCount"`
}

type FormMessage struct {
	FromUser         string       `json:"from_user" validate:"required"`
	Recipients       []string     `json:"recipients" validate:"required"`
	Title            string       `json:"title"`
	Text             string       `json:"text"`
	ReplyToMessageID *uint64      `json:"reply_to"`
	Attachments      []Attachment `json:"attachments"`
}

type FormSearchMessages struct {
	FromUser string `json:"fromUser"`
	ToUser   string `json:"toUser"`
	Folder   string `json:"folder"`
	Filter   string `json:"filter"`
}

func (form *FormMessage) Sanitize() {
	form.Title = html.EscapeString(form.Title)
	sanitizer := bluemonday.UGCPolicy()
	form.Text = sanitizer.Sanitize(form.Text)
	for i, s := range form.Recipients {
		form.Recipients[i] = html.EscapeString(s)
	}
}

type FormFolder struct {
	Name string `json:"name"`
}

func (form *FormFolder) Sanitize() {
	form.Name = html.EscapeString(form.Name)
}

type FolderResponse struct {
	Folder   Folder        `json:"folder"`
	Messages []MessageInfo `json:"messages"`
}

type MessagesResponse struct {
	Messages []MessageInfo `json:"messages"`
}

type FoldersResponse struct {
	Folders []Folder `json:"folders"`
	Count   int      `json:"count"`
}

type MessageResponse struct {
	Message MessageInfo `json:"message"`
}
