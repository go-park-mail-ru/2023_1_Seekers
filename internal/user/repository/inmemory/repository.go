package inmemory

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
)

type profileDB struct {
	profiles []model.Profile
}

func New() user.Repo {
	return &profileDB{
		[]model.Profile{
			{1, "Michail", "Testov", "21.12.2001"},
			{2, "Ivan", "Ivanov", "21.12.2001"},
			{3, "Michail", "Sidorov", "21.12.2001"},
		},
	}
}

func (pDb *profileDB) GetProfileById(id uint64) (*model.Profile, error) {
	for i, p := range pDb.profiles {
		if p.UId == id {
			return &pDb.profiles[i], nil
		}
	}
	return nil, user.ErrUserNotFound
}

func (pDb *profileDB) CreateProfile(profile model.Profile) error {
	_, err := pDb.GetProfileById(profile.UId)
	if err == nil {
		return user.ErrUserExists
	}
	pDb.profiles = append(pDb.profiles, profile)
	return nil
}
