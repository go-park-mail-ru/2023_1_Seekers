package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-playground/validator/v10"
	"net/mail"
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

func (u *useCase) Create(user models.User) (*models.User, error) {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	if len(user.Password) < config.PasswordMinLen {
		return nil, _user.ErrTooShortPw
	}
	if _, ok := validMailAddress(user.Email); !ok {
		return nil, _user.ErrInvalidEmail
	}
	return u.userRepo.Create(user)
}

func (u *useCase) Delete(user models.User) error {
	return u.userRepo.Delete(user)
}

func (u *useCase) GetByID(ID uint64) (*models.User, error) {
	return u.userRepo.GetByID(ID)
}

func (u *useCase) GetByEmail(email string) (*models.User, error) {
	return u.userRepo.GetByEmail(email)
}

func (u *useCase) EditInfo(user models.User) (*models.User, error) {
	//TODO
	return nil, nil
}
func (u *useCase) EditPw(ID uint64, newPW string) (*models.User, error) {
	//TODO
	return nil, nil
}

func (u *useCase) EditAvatar(ID uint64, newAvatar *models.Image) error {
	user, err := u.GetByID(ID)
	if err != nil {
		return fmt.Errorf("failed get user : %w", err)
	}
	f := models.S3File{
		Bucket: config.S3AvatarBucket,
		Name:   user.Email,
		Data:   newAvatar.Data,
	}
	err = u.fileUC.Upload(&f)
	if err = u.fileUC.Upload(&f); err != nil {
		return err
	}
	if err = u.userRepo.SetAvatar(ID, f.Name); err != nil {
		return err
	}
	return nil
}

func (u *useCase) GetAvatar(email string) (*models.Image, error) {
	user, err := u.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed get user : %w", err)
	}
	f, err := u.fileUC.Get(config.S3AvatarBucket, user.Avatar)
	if err != nil {
		f, err = u.fileUC.Get(config.S3AvatarBucket, config.DefaultAvatar)
		if err != nil {
			return nil, err
		}
	}

	return &models.Image{
		Name: f.Name,
		Data: f.Data,
	}, nil
}
