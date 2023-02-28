package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
)

type useCase struct {
	repo user.Repo
}

func New(r user.Repo) user.UseCase {
	return &useCase{repo: r}
}

func (u *useCase) CreateProfile(profile model.Profile) error {
	return u.repo.CreateProfile(profile)
}
func (u *useCase) GetProfileByID(id uint64) (*model.Profile, error) {
	return u.repo.GetProfileByID(id)
}
