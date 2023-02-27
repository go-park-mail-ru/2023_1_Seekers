package inmemory

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/user"
)

type usersDB struct {
	users []model.User
}

func New() user.Repo {
	return &usersDB{}
}

func (u *usersDB) CreateProfile(profile model.Profile) error {
	return nil
}
func (u *usersDB) GetProfileById(id int) (*model.Profile, error) {
	return nil, nil
}
