package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth"
	authRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/repository"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/crypto"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/validation"
	pkgErrors "github.com/pkg/errors"
)

//go:generate mockgen -destination=./mocks/mockauthusecase.go -package=mocks -source=../interface.go

type authUC struct {
	cfg         *config.Config
	userUC      _user.UseCaseI
	sessionRepo authRepo.SessionRepoI
}

func NewAuthUC(c *config.Config, uUC _user.UseCaseI, sr authRepo.SessionRepoI) auth.UseCaseI {
	return &authUC{
		cfg:         c,
		userUC:      uUC,
		sessionRepo: sr,
	}
}

func (u *authUC) SignIn(form *models.FormLogin) (*models.AuthResponse, *models.Session, error) {
	email, err := validation.Login(form.Login, u.cfg.Mail.PostAtDomain)
	if err != nil {
		return nil, nil, pkgErrors.WithMessage(errors.ErrInvalidLogin, err.Error())
	}

	user, err := u.userUC.GetByEmail(email)
	if err != nil || user.IsExternal {
		return nil, nil, errors.ErrWrongPw
	}

	if !crypto.ComparePw2Hash(form.Password, user.Password, u.cfg.Password.PasswordSaltLen) {
		return nil, nil, errors.ErrWrongPw
	}

	session, err := u.CreateSession(user.UserID)
	if err != nil {
		return nil, nil, pkgErrors.Wrap(err, "sign in")
	}

	return &models.AuthResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, session, nil
}

func (u *authUC) SignUp(form *models.FormSignUp) (*models.AuthResponse, *models.Session, error) {
	if form.RepeatPw != form.Password {
		return nil, nil, errors.ErrPwDontMatch
	}

	email, err := validation.Login(form.Login, u.cfg.Mail.PostAtDomain)
	if err != nil {
		return nil, nil, errors.ErrInvalidLogin
	}

	user := &models.User{
		Email:     email,
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Avatar:    u.cfg.UserService.DefaultAvatar,
	}

	user.Password, err = crypto.HashPw(form.Password, u.cfg.Password.PasswordSaltLen)
	if err != nil {
		return nil, nil, pkgErrors.Wrap(err, "sign up")
	}

	user, err = u.userUC.Create(user)
	if err != nil {
		return nil, nil, pkgErrors.Wrap(err, "sign up")
	}

	session, err := u.CreateSession(user.UserID)
	if err != nil {
		//надо ли тут откатить прошлое ?
		return nil, nil, pkgErrors.Wrap(err, "sign up")
	}

	return &models.AuthResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, session, nil
}

func (u *authUC) CreateSession(uID uint64) (*models.Session, error) {
	newSession, err := u.sessionRepo.CreateSession(uID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "create session")
	}

	return newSession, nil
}

func (u *authUC) DeleteSession(sessionID string) error {
	err := u.sessionRepo.DeleteSession(sessionID)
	if err != nil {
		return pkgErrors.Wrap(err, "delete avatar")
	}

	return nil
}

func (u *authUC) GetSession(sessionID string) (*models.Session, error) {
	s, err := u.sessionRepo.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	return s, nil
}
