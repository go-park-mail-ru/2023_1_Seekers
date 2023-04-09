package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/image"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/validation"
	pkgErrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=./mocks_auth/mockusecase.go -package=mocks github.com/go-park-mail-ru/2023_1_Seekers/internal/auth UseCaseI

type authUC struct {
	sessionUC auth.SessionUseCaseI
	mailUC    mail.UseCaseI
	userUC    _user.UseCaseI
}

func NewAuthUC(sUC auth.SessionUseCaseI, mUC mail.UseCaseI, uUC _user.UseCaseI) auth.UseCaseI {
	return &authUC{
		sessionUC: sUC,
		mailUC:    mUC,
		userUC:    uUC,
	}
}

func (u *authUC) SignIn(form models.FormLogin) (*models.AuthResponse, *models.Session, error) {
	email, err := validation.Login(form.Login)
	if err != nil {
		return nil, nil, errors.ErrInvalidLogin
	}
	user, err := u.userUC.GetByEmail(email)
	if err != nil {
		return nil, nil, errors.ErrWrongPw
	}

	if !pkg.ComparePw2Hash(form.Password, user.Password) {
		return nil, nil, errors.ErrWrongPw
	}

	session, err := u.sessionUC.CreateSession(user.UserID)
	if err != nil {
		return nil, nil, pkgErrors.Wrap(err, "sign in")
	}

	return &models.AuthResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, session, nil
}

func (u *authUC) SignUp(form models.FormSignUp) (*models.AuthResponse, *models.Session, error) {
	if form.RepeatPw != form.Password {
		return nil, nil, errors.ErrPwDontMatch
	}

	email, err := validation.Login(form.Login)
	if err != nil || len(form.Login) > 30 || len(form.Login) < 3 {
		return nil, nil, errors.ErrInvalidLogin
	}

	user := &models.User{
		Email:     email,
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Avatar:    config.DefaultAvatar,
	}

	user.Password, err = pkg.HashPw(form.Password)
	if err != nil {
		return nil, nil, pkgErrors.Wrap(err, "sign up")
	}

	user, err = u.userUC.Create(user)
	if err != nil {
		return nil, nil, pkgErrors.Wrap(err, "sign up")
	}

	col := image.GetRandColor()
	label := pkg.GetFirstUtf(user.FirstName)
	img, err := image.GenImage(col, label)
	err = u.userUC.EditAvatar(user.UserID, &models.Image{Data: img}, false)
	if err != nil {
		log.Warn(err, "edit avatar")
	}

	_, err = u.mailUC.CreateDefaultFolders(user.UserID)
	if err != nil {
		return nil, nil, pkgErrors.Wrap(err, "sign up")
	}

	err = u.mailUC.SendWelcomeMessage(user.Email)
	if err != nil {
		return nil, nil, pkgErrors.Wrap(err, "sign up")
	}

	session, err := u.sessionUC.CreateSession(user.UserID)
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
