package inmemory

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
)

type usersDB struct {
	users []models.User
}

func New() user.RepoI {
	return &usersDB{
		[]models.User{
			{0, "support@mailbox.ru", "very_difficult_pw", "Michail", "Testov", config.DefaultAvatar},
			{1, "test@mailbox.ru", "12345", "Ivan", "Ivanov", config.DefaultAvatar},
			{2, "gena@mailbox.ru", "54321", "Michail", "Sidorov", config.DefaultAvatar},
			{3, "max@mailbox.ru", "13795", "Michail", "Testov", config.DefaultAvatar},
			{4, "valera@mailbox.ru", "12345", "Michail", "Testov", config.DefaultAvatar},
		},
	}
}

func (uDb *usersDB) GetUserByID(id uint64) (*models.User, error) {
	for i, u := range uDb.users {
		if u.ID == id {
			return &uDb.users[i], nil
		}
	}
	return nil, auth.ErrUserNotFound
}

func (uDb *usersDB) GetUserByEmail(email string) (*models.User, error) {
	for i, u := range uDb.users {
		if u.Email == email {
			return &uDb.users[i], nil
		}
	}
	return nil, auth.ErrUserNotFound
}

func (uDb *usersDB) CreateUser(user models.User) (*models.User, error) {
	_, err := uDb.GetUserByEmail(user.Email)
	if err == nil {
		return nil, auth.ErrUserExists
	}
	//слой бд отвечает за присваивание id
	// TODO hash pw
	user.ID = uint64(len(uDb.users) + 1)
	uDb.users = append(uDb.users, user)
	return &user, nil
}

func (uDb *usersDB) DeleteUser(user models.User) error {
	for i, u := range uDb.users {
		if u.ID == user.ID {
			uDb.users = append(uDb.users[:i], uDb.users[i+1:]...)
			return nil
		}
	}
	return auth.ErrUserNotFound
}
