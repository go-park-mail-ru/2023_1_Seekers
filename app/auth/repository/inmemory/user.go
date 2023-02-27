package inmemory

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
)

type usersDB struct {
	users []model.User
}

func New() auth.Repo {
	return &usersDB{
		[]model.User{
			{1, "test@example.com", "1234"},
			{2, "gena@example.com", "4321"},
			{3, "max@example.com", "1379"},
		},
	}
}

func (uDb *usersDB) GetById(id uint64) (*model.User, error) {
	for i, u := range uDb.users {
		if u.Id == id {
			return &uDb.users[i], nil
		}
	}
	return nil, auth.ErrUserNotFound
}

func (uDb *usersDB) GetByEmail(email string) (*model.User, error) {
	for i, u := range uDb.users {
		if u.Email == email {
			return &uDb.users[i], nil
		}
	}
	return nil, auth.ErrUserNotFound
}

func (uDb *usersDB) Create(user model.User) (*model.User, error) {
	_, err := uDb.GetById(user.Id)
	if err == nil {
		return nil, auth.ErrUserExists
	}
	_, err = uDb.GetByEmail(user.Email)
	if err == nil {
		return nil, auth.ErrUserExists
	}
	//слой бд отвечает за присваивание id
	// TODO hash pw
	user.Id = uint64(len(uDb.users) + 1)
	uDb.users = append(uDb.users, user)
	return &user, nil
}

func (uDb *usersDB) Delete(user model.User) error {
	for i, u := range uDb.users {
		if u.Id == user.Id {
			uDb.users = append(uDb.users[:i], uDb.users[i+1:]...)
			return nil
		}
	}
	return auth.ErrUserNotFound
}
