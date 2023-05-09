package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/repository"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/crypto"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/image"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/validation"
	"github.com/go-playground/validator/v10"
	pkgErrors "github.com/pkg/errors"
	"net/mail"
	"path/filepath"
)

//go:generate mockgen -destination=./mocks/mockusecase.go -source=../interface.go -package=mocks

type useCase struct {
	cfg      *config.Config
	userRepo _userRepo.UserRepoI
	fileUC   file_storage.UseCaseI
}

func New(c *config.Config, r _userRepo.UserRepoI, fUC file_storage.UseCaseI) _user.UseCaseI {
	image.Init(c.UserService.AvatarTTFPath)
	return &useCase{cfg: c, userRepo: r, fileUC: fUC}
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
	if len(user.Password) < u.cfg.Password.PasswordMinLen {
		return nil, errors.ErrTooShortPw
	}
	if _, ok := validMailAddress(user.Email); !ok {
		return nil, errors.ErrInvalidEmail
	}

	user, err = u.userRepo.Create(user)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "Create user")
	}

	col := image.GetRandColor()
	label := common.GetFirstUtf(user.FirstName)
	img, err := image.GenImage(col, label, u.cfg.UserService.UserDefaultAvatarSize, u.cfg.UserService.UserDefaultAvatarSize)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "Create user - generate avatar")
	}

	err = u.EditAvatar(user.UserID, &models.Image{Data: img}, false)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "Create user - edit avatar")
	}

	return user, nil
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

func (u *useCase) GetInfoByEmail(email string) (*models.UserInfo, error) {
	user, err := u.GetByEmail(email)
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

func (u *useCase) EditInfo(ID uint64, info *models.UserInfo) (*models.UserInfo, error) {
	user, err := u.userRepo.GetByID(ID)
	if err != nil || user.UserID != ID {
		return nil, pkgErrors.Wrap(err, "get user by id")
	}

	err = u.userRepo.EditInfo(ID, info)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "edit info")
	}
	ok, err := u.userRepo.IsCustomAvatar(ID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "edit info - check is custom avatar")
	}

	if !ok {
		avatar, err := u.GetAvatar(user.Email)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "edit info - get avatar")
		}

		label := common.GetFirstUtf(info.FirstName)
		updAvatar, err := image.UpdateImgText(avatar.Data, label, u.cfg.UserService.UserDefaultAvatarSize, u.cfg.UserService.UserDefaultAvatarSize)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "edit info - update avatar")
		}

		err = u.EditAvatar(user.UserID, &models.Image{Data: updAvatar}, false)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "edit info - update avatar")
		}
	}
	return info, nil
}

func (u *useCase) EditAvatar(ID uint64, newAvatar *models.Image, isCustom bool) error {
	user, err := u.GetByID(ID)
	if err != nil {
		return pkgErrors.Wrap(err, "get user by id")
	}
	f := models.S3File{
		Bucket: u.cfg.S3.S3AvatarBucket,
		Name:   user.Email + filepath.Ext(newAvatar.Name),
		Data:   newAvatar.Data,
	}

	if err = u.fileUC.Upload(&f); err != nil {
		return pkgErrors.Wrap(err, "edit avatar")
	}
	if err = u.userRepo.SetAvatar(ID, f.Name); err != nil {
		return pkgErrors.Wrap(err, "edit avatar - set")
	}
	if isCustom {
		if err = u.userRepo.SetCustomAvatar(ID); err != nil {
			return pkgErrors.Wrap(err, "edit avatar - set custom")
		}
	}
	return nil
}

func (u *useCase) GetAvatar(email string) (*models.Image, error) {
	user, err := u.GetByEmail(email)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "get user by email")
	}
	f, err := u.fileUC.Get(u.cfg.S3.S3AvatarBucket, user.Avatar)
	if err != nil {
		//надо ли отправлять дефолтный если что-то пошло не так
		f, err = u.fileUC.Get(u.cfg.S3.S3AvatarBucket, u.cfg.UserService.DefaultAvatar)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "get get avatar")
		}
	}

	return &models.Image{
		Name: f.Name,
		Data: f.Data,
	}, nil
}

func (u *useCase) EditPw(ID uint64, form *models.EditPasswordRequest) error {
	if form.RepeatPw != form.Password {
		return errors.ErrPwDontMatch
	}

	user, err := u.userRepo.GetByID(ID)
	if err != nil {
		return pkgErrors.Wrap(err, "edit password")
	}

	if !crypto.ComparePw2Hash(form.PasswordOld, user.Password, u.cfg.Password.PasswordSaltLen) {
		return errors.ErrWrongPw
	}

	if err := validation.Password(form.Password, u.cfg.Password.PasswordMinLen); err != nil {
		return pkgErrors.Wrap(err, "edit password")
	}
	hashPw, err := crypto.HashPw(form.Password, u.cfg.Password.PasswordSaltLen)
	if err != nil {
		return pkgErrors.Wrap(err, "edit password")
	}
	err = u.userRepo.EditPw(ID, hashPw)
	if err != nil {
		return pkgErrors.Wrap(err, "edit password")
	}
	return nil
}
