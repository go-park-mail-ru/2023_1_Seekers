package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
)

type useCase struct {
	repo user.Repo
}

func New(r user.Repo) user.UseCase {
	return &useCase{repo: r}
}

func (u *useCase) CreateProfile(profile models.Profile) error {
	return u.repo.CreateProfile(profile)
}
func (u *useCase) GetProfileByID(id uint64) (*models.Profile, error) {
	return u.repo.GetProfileByID(id)
}
