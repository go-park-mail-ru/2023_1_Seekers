package http

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/user"
)

type handlers struct {
	useCase user.UseCase
}

func New(uc user.UseCase) user.Handlers {
	return &handlers{
		useCase: uc,
	}
}

func (h *handlers) CreateProfile(profile model.Profile) error {
	return nil
}
func (h *handlers) GetProfileById(id int) (*model.Profile, error) {
	return nil, nil
}
