package postgres

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestRepositoryCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	if db == nil {
		t.Fatal("mock db is null")
	}

	if mock == nil {
		t.Fatal("sqlmock is null")
	}
	defer db.Close()

	dial := postgres.New(postgres.Config{
		DSN:                  "sqlmock_test_db",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gdb, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if gdb == nil {
		t.Fatal("gorm db is nil")
	}

	gdb.Logger.LogMode(logger.Info)
	mockUser := models.User{
		Email:     "test@test.com",
		Password:  "123456",
		FirstName: "test",
		LastName:  "test",
		Avatar:    "default_avatar",
	}
	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "mail"."users" ("email","password","first_name","last_name","avatar") VALUES ($1,$2,$3,$4,$5) RETURNING "user_id"`)).
		WithArgs(mockUser.Email, mockUser.Password, mockUser.FirstName, mockUser.LastName, mockUser.Avatar).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

	mock.ExpectCommit()

	pgRepo := New(gdb)

	usr, err := pgRepo.Create(&mockUser)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, usr, &mockUser)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
