package postgres

import (
	"database/sql"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func createConfig() *config.Config {
	cfg := new(config.Config)
	cfg.DB.DBSchemaName = "mail"

	return cfg
}

func mockDB() (*sql.DB, *gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("mocking database error: %s", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("opening gorm error: %s", err)
	}

	return db, gormDB, mock, nil
}

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestRepository_CreateAlreadyExists(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	generateFakeData(&fakeUser)
	fakeUser.UserID = 10
	fakeUser.IsExternal = false
	fakeUser.Email = "mock_valid@mailbx.ru"

	fakeDBUser := User{}
	fakeDBUser.FromModel(fakeUser)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"user_id", "email", "password", "first_name", "last_name", "avatar", "is_external"}).
		AddRow(fakeUser.UserID, fakeUser.Email, fakeUser.Password, fakeUser.FirstName, fakeUser.LastName, fakeUser.Avatar, fakeUser.IsExternal)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mail"."users" WHERE email = $1 LIMIT 1`)).WithArgs(fakeUser.Email).
		WillReturnRows(rows)

	userRepo := New(cfg, gormDB)
	_, err = userRepo.Create(fakeUser)
	causeErr := pkgErr.Cause(err)

	require.Equal(t, errors.ErrUserExists, causeErr)
}

func TestRepository_EditInfo(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.UserInfo
	generateFakeData(&fakeUser)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mail"."users" SET "first_name"=$1,"last_name"=$2 WHERE user_id = $3`)).
		WithArgs(fakeUser.FirstName, fakeUser.LastName, fakeUser.UserID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	userRepo := New(cfg, gormDB)
	err = userRepo.EditInfo(fakeUser.UserID, fakeUser)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_Delete(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "mail"."users" WHERE user_id = $1`)).
		WithArgs(userID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	userRepo := New(cfg, gormDB)
	err = userRepo.Delete(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_SetAvatar(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	avatar := "avatar.png"

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mail"."users" SET "avatar"=$1 WHERE user_id = $2`)).
		WithArgs(avatar, userID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	userRepo := New(cfg, gormDB)
	err = userRepo.SetAvatar(userID, avatar)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_EditPw(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	newPw := "too_long_password!"
	bytesPw := []byte(newPw)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mail"."users" SET "password"=$1 WHERE user_id = $2`)).
		WithArgs(bytesPw, userID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	userRepo := New(cfg, gormDB)
	err = userRepo.EditPw(userID, newPw)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_GetInfoByID(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.UserInfo
	generateFakeData(&fakeUser)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"user_id", "email", "first_name", "last_name"}).
		AddRow(fakeUser.UserID, fakeUser.Email, fakeUser.FirstName, fakeUser.LastName)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mail"."users" WHERE user_id = $1 LIMIT 1`)).
		WithArgs(fakeUser.UserID).WillReturnRows(rows)

	userRepo := New(cfg, gormDB)
	response, err := userRepo.GetInfoByID(fakeUser.UserID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUser, response)
	}
}

func TestRepository_GetInfoByEmail(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.UserInfo
	generateFakeData(&fakeUser)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"user_id", "email", "first_name", "last_name"}).
		AddRow(fakeUser.UserID, fakeUser.Email, fakeUser.FirstName, fakeUser.LastName)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mail"."users" WHERE email = $1 LIMIT 1`)).
		WithArgs(fakeUser.Email).WillReturnRows(rows)

	userRepo := New(cfg, gormDB)
	response, err := userRepo.GetInfoByEmail(fakeUser.Email)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUser, response)
	}
}

func TestRepository_IsCustomAvatar(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"is_custom_avatar"}).
		AddRow(true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users"."is_custom_avatar" FROM "users" WHERE user_id = $1 LIMIT 1`)).
		WithArgs(userID).WillReturnRows(rows)

	userRepo := New(cfg, gormDB)
	response, err := userRepo.IsCustomAvatar(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, true, response)
	}
}

func TestRepository_SetCustomAvatar(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mail"."users" SET "is_custom_avatar"=$1 WHERE user_id = $2`)).
		WithArgs(true, userID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	userRepo := New(cfg, gormDB)
	err = userRepo.SetCustomAvatar(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}
