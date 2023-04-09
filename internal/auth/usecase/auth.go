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

type authUC struct {
	sessionUC auth.SessionUseCaseI
	userRepo  _user.RepoI
	mailUC    mail.UseCaseI
	userUC    _user.UseCaseI
}

func NewAuthUC(sUC auth.SessionUseCaseI, uRepo _user.RepoI, mUC mail.UseCaseI, uUC _user.UseCaseI) auth.UseCaseI {
	return &authUC{
		sessionUC: sUC,
		userRepo:  uRepo,
		mailUC:    mUC,
		userUC:    uUC,
	}
}

func (u *authUC) SignIn(form models.FormLogin) (*models.AuthResponse, *models.Session, error) {
	email, err := validation.Login(form.Login)
	if err != nil {
		return nil, nil, errors.ErrInvalidLogin
	}
	user, err := u.userRepo.GetByEmail(email)
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
		pkgErrors.Wrap(err, "sign up")
	}

	user, err = u.userRepo.Create(user)
	if err != nil {
		return nil, nil, pkgErrors.Wrap(err, "sign up")
	}

	col := image.GetRandColor()
	label := user.FirstName[0:1]
	img, err := image.GenImage(col, label)
	err = u.userUC.EditAvatar(user.UserID, &models.Image{Data: img})
	if err != nil {
		log.Info(err, "edit ava")
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

func (u *authUC) EditPw(ID uint64, form models.EditPasswordRequest) error {
	if form.RepeatPw != form.Password {
		return errors.ErrPwDontMatch
	}

	user, err := u.userRepo.GetByID(ID)
	if !pkg.ComparePw2Hash(form.PasswordOld, user.Password) {
		return errors.ErrWrongPw
	}

	if err := validation.Password(form.Password); err != nil {
		return pkgErrors.Wrap(err, "create")
	}
	hashPw, err := pkg.HashPw(form.Password)
	if err != nil {
		return pkgErrors.Wrap(err, "edit password")
	}
	err = u.userRepo.EditPw(ID, hashPw)
	if err != nil {
		return pkgErrors.Wrap(err, "edit password")
	}
	return nil
}
