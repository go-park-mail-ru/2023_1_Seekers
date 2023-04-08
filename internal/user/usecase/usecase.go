package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-playground/validator/v10"
	pkgErrors "github.com/pkg/errors"
	"net/mail"
	"path/filepath"
)

type useCase struct {
	userRepo _user.RepoI
	fileUC   file_storage.UseCaseI
}

func New(r _user.RepoI, fUC file_storage.UseCaseI) _user.UseCaseI {
	return &useCase{userRepo: r, fileUC: fUC}
}

func validMailAddress(email string) (string, bool) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", false
	}
	return addr.Address, true
}

func (u *useCase) Create(user *models.User) (*models.User, error) {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return nil, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()) //fmt.Errorf("failed to create user: %w", err)
	}
	if len(user.Password) < config.PasswordMinLen {
		return nil, errors.ErrTooShortPw
	}
	if _, ok := validMailAddress(user.Email); !ok {
		return nil, errors.ErrInvalidEmail
	}
	_, err = u.GetByEmail(user.Email)
	if err == nil {
		return nil, errors.ErrUserExists
	}
	return u.userRepo.Create(user)
}

func (u *useCase) GetInfo(ID uint64) (*models.UserInfo, error) {
	user, err := u.GetByID(ID)
	if err != nil {
		return nil, fmt.Errorf("failed get user : %w", err)
	}
	return &models.UserInfo{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}, nil
}

func (u *useCase) Delete(ID uint64) error {
	err := u.userRepo.Delete(ID)
	if err != nil {
		return pkgErrors.Wrap(err, "delete user")
	}
	return nil
}

func (u *useCase) GetByID(ID uint64) (*models.User, error) {
	user, err := u.userRepo.GetByID(ID)
	if err != nil {
		return user, pkgErrors.Wrap(err, "get user by id")
	}
	return user, nil
}

func (u *useCase) GetByEmail(email string) (*models.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return user, pkgErrors.Wrap(err, "get user by email")
	}
	return user, nil
}

func (u *useCase) EditInfo(ID uint64, info models.UserInfo) (*models.UserInfo, error) {
	user, err := u.userRepo.GetByID(ID)
	if err != nil || user.UserID != ID {
		return nil, pkgErrors.Wrap(err, "get user by id")
	}

	err = u.userRepo.EditInfo(ID, info)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "edit info")
	}
	return &info, nil
}

func (u *useCase) EditAvatar(ID uint64, newAvatar *models.Image) error {
	user, err := u.GetByID(ID)
	if err != nil {
		return pkgErrors.Wrap(err, "get user by id")
	}
	f := models.S3File{
		Bucket: config.S3AvatarBucket,
		Name:   user.Email + filepath.Ext(newAvatar.Name),
		Data:   newAvatar.Data,
	}

	if err = u.fileUC.Upload(&f); err != nil {
		return pkgErrors.Wrap(err, "edit avatar")
	}
	if err = u.userRepo.SetAvatar(ID, f.Name); err != nil {
		return pkgErrors.Wrap(err, "set avatar")
	}
	return nil
}

func (u *useCase) GetAvatar(email string) (*models.Image, error) {
	user, err := u.GetByEmail(email)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get user by email")
	}
	f, err := u.fileUC.Get(config.S3AvatarBucket, user.Avatar)
	if err != nil {
		//надо ли отправлять дефолтный если что-то пошло не так
		f, err = u.fileUC.Get(config.S3AvatarBucket, config.DefaultAvatar)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "get get avatar")
		}
	}

	return &models.Image{
		Name: f.Name,
		Data: f.Data,
	}, nil
}
