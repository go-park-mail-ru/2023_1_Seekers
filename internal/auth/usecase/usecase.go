package usecase

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
)

type useCase struct {
	authRepo auth.RepoI
	userUC   _user.UseCaseI
}

func New(ar auth.RepoI, uc _user.UseCaseI) auth.UseCaseI {
	return &useCase{
		authRepo: ar,
		userUC:   uc,
	}
}

func (u *useCase) SignIn(form models.FormLogin) (*models.User, error) {
	user, err := u.userUC.GetUserByEmail(form.Email)
	if err != nil {
		return nil, fmt.Errorf("cant get user: %w", err)
	}

	if user.Password != form.Password {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (u *useCase) SignUp(form models.FormSignUp) (*models.User, error) {
	if form.RepeatPw != form.Password {
		return nil, auth.ErrPwDontMatch
	}

	user, err := u.userUC.CreateUser(models.User{
		Email:    form.Email,
		Password: form.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("cant create user: %w", err)
	}

	return user, nil
}

func (u *useCase) CreateSession(uID uint64) (*models.Session, error) {
	value, err := pkg.String(config.CookieLen)
	if err != nil {
		return nil, fmt.Errorf("cant create session: %w", err)
	}
	newSession := models.Session{
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

func (u *useCase) GetSession(sessionID string) (*models.Session, error) {
	s, err := u.authRepo.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("cant get session: %w", err)
	}

	return s, nil
}

func (u *useCase) GetSessionByUID(uID uint64) (*models.Session, error) {
	s, err := u.authRepo.GetSessionByUID(uID)
	if err != nil {
		return nil, fmt.Errorf("cant get session by user %w", err)
	}

	return s, nil
}
