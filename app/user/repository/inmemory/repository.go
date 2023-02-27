package inmemory

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/user"
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

func (pDb *profileDB) GetProfileById(id int) (*model.Profile, error) {
	for i, p := range pDb.profiles {
		if p.UId == id {
			return &pDb.profiles[i], nil
		}
	}
	return nil, fmt.Errorf("no user with id %v", id)
}

func (pDb *profileDB) CreateProfile(profile model.Profile) error {
	_, err := pDb.GetProfileById(profile.UId)
	if err == nil {
		return fmt.Errorf("such user exists")
	}
	pDb.profiles = append(pDb.profiles, profile)
	return nil
}
