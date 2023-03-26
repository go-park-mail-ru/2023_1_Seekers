package postgres

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type userDB struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.RepoI {
	return &userDB{
		db: db,
	}
}

func (uDB *userDB) Create(user models.User) (*models.User, error) {
	tx := uDB.db.Create(user)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database [users]")
	}
	return &user, nil
}

func (uDB *userDB) EditInfo(ID uint64, info models.UserInfo) error {
	tx := uDB.db.Omit("user_id", "email", "password").Where("user_id = ?", ID).Updates(info)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database [users]")
	}
	return nil
}

func (uDB *userDB) Delete(ID uint64) error {
	tx := uDB.db.Table("mail.users").Where("user_id = ?", ID).Delete(models.User{})
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database [users]")
	}

	return nil
}

func (uDB *userDB) GetByID(ID uint64) (*models.User, error) {
	usr := models.User{}

	tx := uDB.db.Table("mail.users").Where("user_id = ?", ID).Take(&usr)
	if err := tx.Error; err != nil {
		return nil, errors.Wrap(err, "database [users]")
	}

	return &usr, nil
}

func (uDB *userDB) GetByEmail(email string) (*models.User, error) {
	usr := models.User{}

	tx := uDB.db.Where("email = ?", email).Take(&usr)
	if err := tx.Error; err != nil {
		return nil, errors.Wrap(err, "database [users]")
	}

	return &usr, nil
}

func (uDB *userDB) SetAvatar(ID uint64, avatar string) error {
	tx := uDB.db.Omit("user_id", "email", "password").Where("user_id = ?", ID).Update("avatar", avatar)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database [users]")
	}

	return nil
}

func (uDB *userDB) EditPw(ID uint64, newPW string) error {
	tx := uDB.db.Omit("user_id", "email").Where("user_id = ?", ID).Update("password", newPW)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database [users]")
	}

	return nil
}
