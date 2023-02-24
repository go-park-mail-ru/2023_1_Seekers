package user

import "github.com/go-park-mail-ru/2023_1_Seekers/app/model"

type UseCase interface {
	Create(user model.User) error
	Delete(user model.User) error
	GetById(id int) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}

type Repo interface {
	Create(user model.User) error
	Delete(user model.User) error
	GetById(id int) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}
