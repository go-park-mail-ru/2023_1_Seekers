package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/user"
)

type useCase struct {
	repo user.Repo
}

func New(r user.Repo) user.UseCase {
	return &useCase{repo: r}
}

func (u *useCase) CreateProfile(profile model.Profile) error {
	return nil
}
func (u *useCase) GetProfileById(id int) (*model.Profile, error) {
	return nil, nil
}
