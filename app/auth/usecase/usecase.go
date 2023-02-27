package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	_session "github.com/go-park-mail-ru/2023_1_Seekers/app/session"
)

type useCase struct {
	sessionUC _session.UseCase
	authRepo  auth.Repo
}

func New(sc _session.UseCase, ar auth.Repo) auth.UseCase {
	return &useCase{
		sessionUC: sc,
		authRepo:  ar,
	}
}

func (u *useCase) SignIn(form model.FormLogin) (*model.User, error) {
	user, err := u.authRepo.GetByEmail(form.Email)
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

func (u *useCase) SignUp(form model.FormSignUp) (*model.User, error) {
	if form.RepeatPw != form.Password {
		return nil, nil, fmt.Errorf("passwords dont match")
	}

	user, err := u.authRepo.Create(model.User{
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
