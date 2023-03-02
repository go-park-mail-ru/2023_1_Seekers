package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
)

type useCase struct {
	userRepo user.RepoI
}

func New(r user.RepoI) user.UseCaseI {
	return &useCase{userRepo: r}
}

func (u *useCase) CreateUser(user models.User) (*models.User, error) {
	return u.userRepo.CreateUser(user)
}

func (u *useCase) DeleteUser(user models.User) error {
	return u.userRepo.DeleteUser(user)
}

func (u *useCase) GetUserByID(ID uint64) (*models.User, error) {
	return u.userRepo.GetUserByID(ID)
}

func (u *useCase) GetUserByEmail(email string) (*models.User, error) {
	return u.userRepo.GetUserByEmail(email)
}

func (u *useCase) CreateProfile(profile models.Profile) error {
	// TODO some profile validation
	return u.userRepo.CreateProfile(profile)
}
func (u *useCase) GetProfileByID(id uint64) (*models.Profile, error) {
	return u.userRepo.GetProfileByID(id)
}
