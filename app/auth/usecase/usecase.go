package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/app/user"
	"strconv"
	"time"
)

type useCase struct {
	authRepo auth.Repo
	userUC   _user.UseCase
}

func New(ar auth.Repo, uc _user.UseCase) auth.UseCase {
	return &useCase{
		authRepo: ar,
		userUC:   uc,
	}
}

func (u *useCase) newCookie(userId int) (*model.Cookie, error) {
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

func (u *useCase) getCookie(uId int) (*model.Cookie, error) {
	cookie, err := u.authRepo.GetCookie(uId)
	if err != nil {
		return nil, fmt.Errorf("cant get cookie: %w", err)
	}

	return cookie, nil
}

func (u *useCase) deleteCookie(session string) error {
	err := u.authRepo.DeleteCookie(session)
	if err != nil {
		return fmt.Errorf("cant delete cookie: %w", err)
	}

	return nil
}

func (u *useCase) SignIn(form model.FormAuth) (*model.User, *model.Cookie, error) {
	user, err := u.userUC.GetByEmail(form.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("cant get user: %w", err)
	}

	if user.Password != form.Password {
		return nil, nil, fmt.Errorf("invalid password")
	}
	cookie, err := u.getCookie(user.Id)
	if err != nil {
		cookie, err = u.newCookie(user.Id)
	}
	if err != nil {
		return nil, nil, err
	}

	return user, cookie, nil
}

func (u *useCase) SignUp(form model.FormReg) (*model.User, *model.Cookie, error) {
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

	cookie, err := u.newCookie(user.Id)
	if err != nil {
		return nil, nil, err
	}

	return user, cookie, nil
}

func (u *useCase) Logout(session string) error {
	err := u.deleteCookie(session)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}
