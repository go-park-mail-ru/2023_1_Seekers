package postgres

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"os"
)

type User struct {
	UserID uint64 `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	//HereSince time.Time `json:"hereSince" gorm:"column:here_since"`
	//IsDeleted bool   `json:"isDeleted" gorm:"column:is_deleted"`
	Email     string `validate:"required"`
	Password  []byte `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Avatar    string
	//IsCustomAvatar bool
}

func (*User) TableName() string {
	return os.Getenv(config.DBSchemaNameEnv) + ".users"
}

func (u *User) FromModel(user *models.User) {
	u.Email = user.Email
	u.Password = []byte(user.Password)
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Avatar = user.Avatar
}

type IsCustomAvatar struct {
	IsCustomAvatar bool
}

func (*IsCustomAvatar) TableName() string {
	return os.Getenv(config.DBSchemaNameEnv) + ".users"
}
