package inmemory

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
)

type usersDB struct {
	users []models.User
}

func New() user.RepoI {
	return &usersDB{
		[]models.User{
			{ID: 0, Email: "support@mailbox.ru", Password: "very_difficult_pw", FirstName: "Michail", LastName: "Testov", Avatar: config.DefaultAvatar},
			{ID: 1, Email: "test@mailbox.ru", Password: "12345", FirstName: "Ivan", LastName: "Ivanov", Avatar: config.DefaultAvatar},
			{ID: 2, Email: "gena@mailbox.ru", Password: "54321", FirstName: "Michail", LastName: "Sidorov", Avatar: config.DefaultAvatar},
			{ID: 3, Email: "max@mailbox.ru", Password: "13795", FirstName: "Michail", LastName: "Testov", Avatar: config.DefaultAvatar},
			{ID: 4, Email: "valera@mailbox.ru", Password: "12345", FirstName: "Michail", LastName: "Testov", Avatar: config.DefaultAvatar},
		},
	}
}

func (uDb *usersDB) GetByID(id uint64) (*models.User, error) {
	for i, u := range uDb.users {
		if u.ID == id {
			return &uDb.users[i], nil
		}
	}
	return nil, user.ErrUserNotFound
}

func (uDb *usersDB) GetByEmail(email string) (*models.User, error) {
	for i, u := range uDb.users {
		if u.Email == email {
			return &uDb.users[i], nil
		}
	}
	return nil, user.ErrUserNotFound
}

func (uDb *usersDB) Create(user models.User) (*models.User, error) {
	_, err := uDb.GetByEmail(user.Email)
	if err == nil {
		return nil, auth.ErrUserExists
	}
	//слой бд отвечает за присваивание id
	// TODO hash pw
	user.ID = uint64(len(uDb.users) + 1)
	uDb.users = append(uDb.users, user)
	return &user, nil
}

func (uDb *usersDB) Delete(ID uint64) error {
	for i, u := range uDb.users {
		if u.ID == ID {
			uDb.users = append(uDb.users[:i], uDb.users[i+1:]...)
			return nil
		}
	}
	return user.ErrUserNotFound
}

func (uDb *usersDB) SetAvatar(ID uint64, avatar string) error {
	for _, u := range uDb.users {
		if u.ID == ID {
			u.Avatar = avatar
			//на слайсах криво обновляется скоро на pg переделаем
			uDb.Delete(u.ID)
			uDb.users = append(uDb.users, u)
			return nil
		}
	}
	return user.ErrUserNotFound
}

func (uDb *usersDB) EditInfo(ID uint64, info models.UserInfo) error {
	for _, u := range uDb.users {
		if u.ID == ID {
			u.FirstName = info.FirstName
			u.LastName = info.LastName
			//на слайсах криво обновляется скоро на pg переделаем
			uDb.Delete(u.ID)
			uDb.users = append(uDb.users, u)
			return nil
		}
	}
	return user.ErrUserNotFound
}

func (uDb *usersDB) EditPw(ID uint64, newPW string) error {
	for _, u := range uDb.users {
		if u.ID == ID {
			u.Password = newPW
			//на слайсах криво обновляется скоро на pg переделаем
			uDb.Delete(u.ID)
			uDb.users = append(uDb.users, u)
			return nil
		}
	}
	return user.ErrUserNotFound
}
