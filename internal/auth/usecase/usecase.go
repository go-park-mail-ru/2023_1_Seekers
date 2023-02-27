package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/model"
)

type useCase struct {
	authRepo auth.Repo
}

func New(ar auth.Repo) auth.UseCase {
	return &useCase{
		authRepo: ar,
	}
}

func (u *useCase) SignIn(form model.FormLogin) (*model.User, error) {
	user, err := u.authRepo.GetByEmail(form.Email)
	if err != nil {
		return nil, fmt.Errorf("cant get user: %w", err)
	}

	if user.Password != form.Password {
		return nil, auth.ErrInvalidPw
	}

	return user, nil
}

func (u *useCase) SignUp(form model.FormSignUp) (*model.User, error) {
	if form.RepeatPw != form.Password {
		return nil, auth.ErrPwDontMatch
	}

	user, err := u.authRepo.Create(model.User{
		Email:    form.Email,
		Password: form.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("cant create user: %w", err)
	}

	return user, nil
}
