package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	_session "github.com/go-park-mail-ru/2023_1_Seekers/app/session"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/app/user"
)

type useCase struct {
	sessionUC _session.UseCase
	userUC    _user.UseCase
}

func New(sc _session.UseCase, uc _user.UseCase) auth.UseCase {
	return &useCase{
		sessionUC: sc,
		userUC:    uc,
	}
}

func (u *useCase) SignIn(form model.FormLogin) (*model.User, *model.Session, error) {
	user, err := u.userUC.GetByEmail(form.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("cant get user: %w", err)
	}

	if user.Password != form.Password {
		return nil, nil, fmt.Errorf("invalid password")
	}
	session, err := u.sessionUC.GetSessionById(user.Id)
	if err != nil {
		session, err = u.sessionUC.Create(user.Id)
		if err != nil {
			return nil, nil, err
		}
	}

	return user, session, nil
}

func (u *useCase) SignUp(form model.FormSignUp) (*model.User, *model.Session, error) {
	if form.RepeatPw != form.Password {
		return nil, nil, fmt.Errorf("passwords dont match")
	}

	user, err := u.userUC.Create(model.User{
		Email:    form.Email,
		Password: form.Password,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("cant create user: %w", err)
	}
	// TODO store profile
	profile := model.Profile{
		UId:       user.Id,
		FirstName: form.FirstName,
		LastName:  form.LastName,
		BirthDate: form.BirthDate,
	}
	fmt.Println(profile)

	session, err := u.sessionUC.Create(user.Id)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (u *useCase) Logout(sessionId string) error {
	err := u.sessionUC.Delete(sessionId)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}

func (u *useCase) Auth(sessionId string) error {
	_, err := u.sessionUC.GetSession(sessionId)
	if err != nil {
		return fmt.Errorf("failed to auth: %w", err)
	}
	return nil
}
