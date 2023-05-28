package postgres

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
)

type User struct {
	UserID uint64 `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	//HereSince time.Time `json:"hereSince" gorm:"column:here_since"`
	//IsDeleted bool   `json:"isDeleted" gorm:"column:is_deleted"`
	Email      string
	Password   []byte
	FirstName  string
	LastName   string
	Avatar     string
	IsExternal bool
	IsFake     bool
	//IsCustomAvatar bool
}

func (*User) TableName(schemaName string) string {
	return schemaName + ".users"
}

func (u *User) FromModel(user *models.User) {
	u.Email = user.Email
	u.Password = []byte(user.Password)
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Avatar = user.Avatar
	u.IsExternal = user.IsExternal
	u.IsFake = user.IsFake
}

type IsCustomAvatar struct {
	IsCustomAvatar bool
}

func (*IsCustomAvatar) TableName() string {
	return "mail.users"
}
