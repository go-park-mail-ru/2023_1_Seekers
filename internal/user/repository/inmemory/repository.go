package inmemory

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
)

type usersDB struct {
	users    []models.User
	profiles []models.Profile
}

func New() user.RepoI {
	return &usersDB{
		[]models.User{
			{1, "test@example.com", "12345"},
			{2, "gena@example.com", "54321"},
			{3, "max@example.com", "13795"},
		},
		[]models.Profile{
			{1, "Michail", "Testov", "21.12.2001"},
			{2, "Ivan", "Ivanov", "21.12.2001"},
			{3, "Michail", "Sidorov", "21.12.2001"},
		},
	}
}

func (uDb *usersDB) GetUserByID(id uint64) (*models.User, error) {
	for i, u := range uDb.users {
		if u.ID == id {
			return &uDb.users[i], nil
		}
	}
	return nil, auth.ErrUserNotFound
}

func (uDb *usersDB) GetUserByEmail(email string) (*models.User, error) {
	for i, u := range uDb.users {
		if u.Email == email {
			return &uDb.users[i], nil
		}
	}
	return nil, auth.ErrUserNotFound
}

func (uDb *usersDB) CreateUser(user models.User) (*models.User, error) {
	_, err := uDb.GetUserByEmail(user.Email)
	if err == nil {
		return nil, auth.ErrUserExists
	}
	//слой бд отвечает за присваивание id
	// TODO hash pw
	user.ID = uint64(len(uDb.users) + 1)
	uDb.users = append(uDb.users, user)
	return &user, nil
}

func (uDb *usersDB) DeleteUser(user models.User) error {
	for i, u := range uDb.users {
		if u.ID == user.ID {
			uDb.users = append(uDb.users[:i], uDb.users[i+1:]...)
			return nil
		}
	}
	return auth.ErrUserNotFound
}

func (uDb *usersDB) GetProfileByID(id uint64) (*models.Profile, error) {
	for i, p := range uDb.profiles {
		if p.UID == id {
			return &uDb.profiles[i], nil
		}
	}
	return nil, user.ErrUserNotFound
}

func (uDb *usersDB) CreateProfile(profile models.Profile) error {
	_, err := uDb.GetProfileByID(profile.UID)
	if err == nil {
		return user.ErrUserExists
	}
	uDb.profiles = append(uDb.profiles, profile)
	return nil
}
