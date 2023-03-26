package models

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"os"
	"time"
)

type User struct {
	ID        uint64    `json:"id" gorm:"column:user_id"`
	HereSince time.Time `json:"hereSince" gorm:"column:here_since"`
	IsDeleted bool      `json:"isDeleted" gorm:"column:is_deleted"`
	Email     string    `json:"email" validate:"required" gorm:"column:email"`
	Password  string    `json:"password" validate:"required" gorm:"column:password"`
	FirstName string    `json:"firstName" validate:"required" gorm:"column:first_name"`
	LastName  string    `json:"lastName" validate:"required" gorm:"column:last_name"`
	Avatar    string    `json:"avatar,omitempty" gorm:"column:avatar"`
}

type FormSignUp struct {
	Login     string `json:"login" validate:"required"`
	Password  string `json:"password" validate:"required"`
	RepeatPw  string `json:"repeatPw" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type FormLogin struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	Remember bool   `json:"remember" validate:"required"`
}

type UserInfo struct {
	UserID    uint64 `json:"-" gorm:"column:user_id"`
	FirstName string `json:"firstName" validate:"required" gorm:"column:first_name"`
	LastName  string `json:"lastName" validate:"required" gorm:"column:last_name"`
	Email     string `json:"email" validate:"required" gorm:"column:email"`
}

type EditUserInfoResponse struct {
	Email string `json:"email" validate:"required"`
}

type EditPasswordRequest struct {
	Password string `json:"password" validate:"required"`
}

func (User) TableName() string {
	return os.Getenv(config.DBSchemaNameEnv) + ".users"
}
func (UserInfo) TableName() string {
	return os.Getenv(config.DBSchemaNameEnv) + ".users"
}
