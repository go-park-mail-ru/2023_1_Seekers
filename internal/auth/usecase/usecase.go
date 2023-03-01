package usecase

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/build/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
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
		return nil, errors.New("user not found")
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

func (u *useCase) CreateSession(uID uint64) (*model.Session, error) {
	value, err := pkg.String(config.CookieLen)
	if err != nil {
		return nil, fmt.Errorf("cant create session: %w", err)
	}
	newSession := model.Session{
		UID:       uID,
		SessionID: value,
	}

	err = u.authRepo.CreateSession(newSession)
	if err != nil {
		return nil, fmt.Errorf("cant create session: %w", err)
	}

	return &newSession, nil
}

func (u *useCase) DeleteSession(sessionID string) error {
	err := u.authRepo.DeleteSession(sessionID)
	if err != nil {
		return fmt.Errorf("cant delete session: %w", err)
	}

	return nil
}

func (u *useCase) DeleteSessionByUID(uID uint64) error {
	err := u.authRepo.DeleteSessionByUID(uID)
	if err != nil {
		return fmt.Errorf("cant delete session by id: %w", err)
	}

	return nil
}

func (u *useCase) GetSession(sessionID string) (*model.Session, error) {
	s, err := u.authRepo.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("cant get session: %w", err)
	}

	return s, nil
}

func (u *useCase) GetSessionByUID(uID uint64) (*model.Session, error) {
	s, err := u.authRepo.GetSessionByUID(uID)
	if err != nil {
		return nil, fmt.Errorf("cant get session by user %w", err)
	}

	return s, nil
}