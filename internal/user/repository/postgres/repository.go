package postgres

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=mocks/repository.go -package=mocks github.com/go-park-mail-ru/2023_1_Seekers/internal/user RepoI

type userDB struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.RepoI {
	return &userDB{
		db: db,
	}
}

func (uDB *userDB) Create(user *models.User) (*models.User, error) {
	_, err := uDB.GetByEmail(user.Email)
	if err == nil {
		return nil, errors.ErrUserExists
	}
	dbUser := User{}
	dbUser.FromModel(user)
	tx := uDB.db.Create(&dbUser)
	if err := tx.Error; err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	user.UserID = dbUser.UserID
	return user, nil
}

func (uDB *userDB) EditInfo(ID uint64, info models.UserInfo) error {
	tx := uDB.db.Omit("user_id", "email", "password").Where("user_id = ?", ID).Updates(&info)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.WithMessage(errors.ErrUserNotFound, err.Error())
		}
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}
	return nil
}

func (uDB *userDB) Delete(ID uint64) error {
	tx := uDB.db.Where("user_id = ?", ID).Delete(models.User{})
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.WithMessage(errors.ErrUserNotFound, err.Error())
		}
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (uDB *userDB) GetByID(ID uint64) (*models.User, error) {
	usr := models.User{}

	tx := uDB.db.Where("user_id = ?", ID).Take(&usr)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.WithMessage(errors.ErrUserNotFound, err.Error())
		}
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return &usr, nil
}

func (uDB *userDB) GetByEmail(email string) (*models.User, error) {
	usr := models.User{}

	tx := uDB.db.Where("email = ?", email).Take(&usr)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.WithMessage(errors.ErrUserNotFound, err.Error())
		}
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return &usr, nil
}

func (uDB *userDB) SetAvatar(ID uint64, avatar string) error {
	tx := uDB.db.Model(&models.User{}).Omit("user_id", "email", "password").Where("user_id = ?", ID).Update("avatar", avatar)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.WithMessage(errors.ErrUserNotFound, err.Error())
		}
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (uDB *userDB) EditPw(ID uint64, newPW string) error {
	tx := uDB.db.Model(&models.User{}).Omit("user_id", "email").Where("user_id = ?", ID).Update("password", []byte(newPW))
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.WithMessage(errors.ErrUserNotFound, err.Error())
		}

		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (uDB *userDB) GetInfoByID(ID uint64) (*models.UserInfo, error) {
	userInfo := models.UserInfo{}
	tx := uDB.db.Where("user_id = ?", ID).Take(&userInfo)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.WithMessage(errors.ErrUserNotFound, err.Error())
		}

		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return &userInfo, nil
}

func (uDB *userDB) GetInfoByEmail(email string) (*models.UserInfo, error) {
	userInfo := models.UserInfo{}
	tx := uDB.db.Where("email = ?", email).Take(&userInfo)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.WithMessage(errors.ErrUserNotFound, err.Error())
		}

		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return &userInfo, nil
}
