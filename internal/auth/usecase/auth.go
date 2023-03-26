package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
)

type authUC struct {
	sessionUC auth.SessionUseCaseI
	userUC    _user.UseCaseI
	mailUC    mail.UseCaseI
}

func NewAuthUC(sUC auth.SessionUseCaseI, uc _user.UseCaseI, mUC mail.UseCaseI) auth.UseCaseI {
	return &authUC{
		sessionUC: sUC,
		userUC:    uc,
		mailUC:    mUC,
	}
}

func (u *authUC) SignIn(form models.FormLogin) (*models.AuthResponse, *models.Session, error) {
	email, err := pkg.ValidateLogin(form.Login)
	if err != nil {
		return nil, nil, auth.ErrInvalidLogin
	}
	user, err := u.userUC.GetByEmail(email)
	if err != nil {
		return nil, nil, auth.ErrWrongPw
	}

	if user.Password != form.Password {
		return nil, nil, auth.ErrWrongPw
	}

	session, err := u.sessionUC.CreateSession(user.ID)
	if err != nil {
		return nil, nil, auth.ErrFailedCreateSession
	}

	return &models.AuthResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, session, nil
}

func (u *authUC) SignUp(form models.FormSignUp) (*models.AuthResponse, *models.Session, error) {
	if form.RepeatPw != form.Password {
		return nil, nil, auth.ErrPwDontMatch
	}

	email, err := pkg.ValidateLogin(form.Login)
	if err != nil || len(form.Login) > 30 || len(form.Login) < 3 {
		return nil, nil, auth.ErrInvalidLogin
	}

	user, err := u.userUC.Create(models.User{
		Email:     email,
		Password:  form.Password,
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Avatar:    config.DefaultAvatar,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("cant create user: %w", err)
	}

	err = u.mailUC.SendWelcomeMessage(user.Email)
	if err != nil {
		return nil, nil, auth.ErrInternalHelloMsg
	}

	session, err := u.sessionUC.CreateSession(user.ID)
	if err != nil {
		return nil, nil, auth.ErrFailedCreateSession
	}

	return &models.AuthResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, session, nil
}
