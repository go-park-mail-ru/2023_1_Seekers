package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/user"
)

type useCase struct {
	repo user.Repo
}

func New(r user.Repo) user.UseCase {
	return &useCase{repo: r}
}

func (u *useCase) Create(user model.User) (*model.User, error) {
	if user.Email == "" {
		return nil, fmt.Errorf("cant create user, email is empty")
	}
	if user.Password == "" {
		return nil, fmt.Errorf("cant create user, email is empty")
	}
	return u.repo.Create(user)
}

func (u *useCase) Delete(user model.User) error {
	return u.repo.Delete(user)
}

func (u *useCase) GetById(id int) (*model.User, error) {
	return u.repo.GetById(id)
}

func (u *useCase) GetByEmail(email string) (*model.User, error) {
	return u.repo.GetByEmail(email)
}
