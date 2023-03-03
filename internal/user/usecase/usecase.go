package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"net/mail"
)

type useCase struct {
	userRepo _user.RepoI
}

func New(r _user.RepoI) _user.UseCaseI {
	return &useCase{userRepo: r}
}

func validMailAddress(email string) (string, bool) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", false
	}
	return addr.Address, true
}

func (u *useCase) CreateUser(user models.User) (*models.User, error) {
	if len(user.Password) < config.PasswordMinLen {
		return nil, _user.ErrTooShortPw
	}
	if _, ok := validMailAddress(user.Email); !ok {
		return nil, _user.ErrInvalidEmail
	}
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
	return u.userRepo.CreateProfile(profile)
}
func (u *useCase) GetProfileByID(id uint64) (*models.Profile, error) {
	return u.userRepo.GetProfileByID(id)
}
