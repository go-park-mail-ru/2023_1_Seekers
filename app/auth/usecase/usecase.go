package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/user"
	"strconv"
	"time"
)

type useCase struct {
	authRepo auth.Repo
	userRepo user.Repo
}

func New(ar auth.Repo, ur user.Repo) auth.UseCase {
	return &useCase{
		authRepo: ar,
		userRepo: ur,
	}
}

func (u *useCase) NewCookie(userId int) (*model.Cookie, error) {
	// TODO hash
	cookie := model.Cookie{
		UId:     userId,
		Session: "session_cookie_" + strconv.Itoa(userId),
		Expire:  time.Now().AddDate(0, 0, 1)}

	err := u.authRepo.NewCookie(cookie)
	if err != nil {
		return nil, fmt.Errorf("cant create cookie: %w", err)
	}

	return &cookie, nil
}

func (u *useCase) GetCookie(value string) (*model.Cookie, error) {
	cookie, err := u.authRepo.GetCookie(value)
	if err != nil {
		return nil, fmt.Errorf("cant get cookie: %w", err)
	}

	return cookie, nil
}

func (u *useCase) DeleteCookie(value string) error {
	err := u.authRepo.DeleteCookie(value)
	if err != nil {
		return fmt.Errorf("cant delete cookie: %w", err)
	}

	return nil
}

func (u *useCase) SignIn(form model.FormAuth) (*model.User, *model.Cookie, error) {
	user, err := u.userRepo.GetByEmail(form.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("cant get user: %w", err)
	}

	if user.Password != form.Password {
		return nil, nil, fmt.Errorf("invalid password")
	}

	cookie, err := u.NewCookie(user.Id)
	if err != nil {
		return nil, nil, err
	}

	return user, cookie, nil
}

func (u *useCase) SignUp(form model.FormReg) (*model.User, *model.Cookie, error) {
	if form.RepeatPw != form.Password {
		return nil, nil, fmt.Errorf("passwords dont match")
	}

	user, err := u.userRepo.Create(model.User{
		Email:    form.Email,
		Password: form.Password,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("cant create user: %w", err)
	}

	cookie, err := u.NewCookie(user.Id)
	if err != nil {
		return nil, nil, err
	}

	return user, cookie, nil
}
