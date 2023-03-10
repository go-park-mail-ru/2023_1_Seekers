package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-playground/validator/v10"
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
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		fmt.Println(5)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	if len(user.Password) < config.PasswordMinLen {
		fmt.Println(6)
		return nil, _user.ErrTooShortPw
	}
	if _, ok := validMailAddress(user.Email); !ok {
		fmt.Println(7)
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
	validate := validator.New()
	err := validate.Struct(profile)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}
	return u.userRepo.CreateProfile(profile)
}
func (u *useCase) GetProfileByID(id uint64) (*models.Profile, error) {
	return u.userRepo.GetProfileByID(id)
}
