package inmemory

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/user"
)

type UsersDB struct {
	Users []model.User
}

func New() user.Repo {
	return &UsersDB{
		[]model.User{
			{1, "test@example.com", "1234"},
			{2, "gena@example.com", "4321"},
			{3, "max@example.com", "1379"},
		},
	}
}

func (u *UsersDB) GetById(id int) (*model.User, error) {
	for i, user := range u.Users {
		if user.Id == id {
			return &u.Users[i], nil
		}
	}
	return nil, fmt.Errorf("no user with id %v", id)
}

func (u *UsersDB) GetByEmail(email string) (*model.User, error) {
	for i, user := range u.Users {
		if user.Email == email {
			return &u.Users[i], nil
		}
	}
	return nil, fmt.Errorf("no user with email %v", email)
}

func (u *UsersDB) Create(user model.User) error {
	_, err := u.GetById(user.Id)
	if err == nil {
		return fmt.Errorf("such user exists")
	}
	_, err = u.GetByEmail(user.Email)
	if err == nil {
		return fmt.Errorf("such user exists")
	}
	//слой бд отвечает за присваивание id
	user.Id = len(u.Users) + 1
	u.Users = append(u.Users, user)
	return nil
}

func (u *UsersDB) Delete(user model.User) error {
	for i, usr := range u.Users {
		if usr.Id == user.Id {
			u.Users = append(u.Users[:i], u.Users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("no such user %v", user)
}
