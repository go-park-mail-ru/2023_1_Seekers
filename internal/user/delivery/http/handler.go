package http

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
)

type handlers struct {
	useCase user.UseCase
}

func New(uc user.UseCase) user.Handlers {
	return &handlers{
		useCase: uc,
	}
}

// TODO implement
func (h *handlers) CreateProfile(profile models.Profile) error {
	return nil
}
func (h *handlers) GetProfileByID(id uint64) (*models.Profile, error) {
	return nil, nil
}
