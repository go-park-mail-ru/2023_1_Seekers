package postgres

import (
	"testing"
)

//
//func TestRepositoryCreateUser(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatal(err)
//	}
//	if db == nil {
//		t.Fatal("mock db is null")
//	}
//
//	if mock == nil {
//		t.Fatal("sqlmock is null")
//	}
//	defer db.Close()
//
//	dial := postgres.New(postgres.Config{
//		DSN:                  "sqlmock_test_db",
//		DriverName:           "postgres",
//		Conn:                 db,
//		PreferSimpleProtocol: true,
//	})
//	gdb, err := gorm.Open(dial, &gorm.Config{})
//	if err != nil {
//		t.Fatal(err)
//	}
//	if gdb == nil {
//		t.Fatal("gorm db is nil")
//	}
//
//	gdb.Logger.LogMode(logger.Info)
//	user := models.User{
//		Email:     "test@test.com",
//		Password:  "123456",
//		FirstName: "test",
//		LastName:  "test",
//		Avatar:    "default_avatar",
//	}
//	var mockUser User
//	mockUser.FromModel(&user)
//
//	mock.ExpectBegin()
//
//	mock.ExpectQuery(regexp.QuoteMeta(
//		`SELECT * FROM "mail"."users" WHERE email = $1 LIMIT 1`)).WithArgs(mockUser.Email)
//	mock.ExpectQuery(regexp.QuoteMeta(
//		`INSERT INTO "mail"."users" ("email","password","first_name","last_name","avatar") VALUES ($1,$2,$3,$4,$5) RETURNING "user_id"`)).
//		WithArgs(mockUser.Email, mockUser.Password, mockUser.FirstName, mockUser.LastName, mockUser.Avatar).
//		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
//
//	mock.ExpectCommit()
//
//	pgRepo := New(gdb)
//
//	usr, err := pgRepo.Create(&user)
//	if err != nil {
//		t.Error(err)
//	}
//
//	assert.Equal(t, usr, &mockUser)
//	err = mock.ExpectationsWereMet()
//	assert.NoError(t, err)
//}
//
//func TestRepositoryEditInfo(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatal(err)
//	}
//	if db == nil {
//		t.Fatal("mock db is null")
//	}
//
//	if mock == nil {
//		t.Fatal("sqlmock is null")
//	}
//	defer db.Close()
//
//	dial := postgres.New(postgres.Config{
//		DSN:                  "sqlmock_test_db",
//		DriverName:           "postgres",
//		Conn:                 db,
//		PreferSimpleProtocol: true,
//	})
//	gdb, err := gorm.Open(dial, &gorm.Config{})
//	if err != nil {
//		t.Fatal(err)
//	}
//	if gdb == nil {
//		t.Fatal("gorm db is nil")
//	}
//
//	gdb.Logger.LogMode(logger.Info)
//	mockInfo := models.UserInfo{
//		UserID:    1,
//		FirstName: "test",
//		LastName:  "test",
//		Email:     "test@test.com",
//	}
//	mock.ExpectBegin()
//
//	mock.ExpectExec(regexp.QuoteMeta(
//		`UPDATE "mail"."users" SET "first_name"=$1,"last_name"=$2 WHERE user_id = $3`)).
//		WithArgs(mockInfo.FirstName, mockInfo.LastName, mockInfo.UserID).
//		WillReturnResult(sqlmock.NewResult(int64(mockInfo.UserID), 1))
//	mock.ExpectCommit()
//
//	pgRepo := New(gdb)
//
//	err = pgRepo.EditInfo(mockInfo.UserID, mockInfo)
//	if err != nil {
//		t.Error(err)
//	}
//}

func TestRepositoryDelete(t *testing.T) {
	// TODO fix err
}

func TestRepositoryGetByID(t *testing.T) {
	// не завелось
	//db, mock, err := sqlmock.New()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if db == nil {
	//	t.Fatal("mock db is null")
	//}
	//
	//if mock == nil {
	//	t.Fatal("sqlmock is null")
	//}
	//defer db.Close()
	//
	//dial := postgres.New(postgres.Config{
	//	DSN:                  "sqlmock_test_db",
	//	DriverName:           "postgres",
	//	Conn:                 db,
	//	PreferSimpleProtocol: true,
	//})
	//gdb, err := gorm.Open(dial, &gorm.Config{})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if gdb == nil {
	//	t.Fatal("gorm db is nil")
	//}
	//
	//gdb.Logger.LogMode(logger.Info)
	//var mockID uint64 = 2
	//mockUser := models.User{
	//	UserID:    mockID,
	//	Email:     "test@test.com",
	//	Password:  "123456",
	//	FirstName: "test",
	//	LastName:  "test",
	//	Avatar:    "default_avatar",
	//}
	//
	//mock.ExpectBegin()
	//mock.MatchExpectationsInOrder(false)
	//mock.ExpectQuery(regexp.QuoteMeta(
	//	`SELECT * FROM "mail"."users" WHERE user_id = $1 LIMIT 1`)).
	//	WithArgs(mockID).
	//	WillReturnRows(sqlmock.NewRows([]string{"user_id", "email", "password", "first_name", "last_name", "avatar"}).
	//		AddRow(mockID, mockUser.Email, mockUser.Password, mockUser.FirstName, mockUser.LastName, mockUser.Avatar))
	//mock.ExpectCommit()
	//
	//pgRepo := New(gdb)
	//
	//usr, err := pgRepo.GetByID(mockID)
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//assert.Equal(t, usr, &mockUser)
	//err = mock.ExpectationsWereMet()
	//assert.NoError(t, err)
}

func TestRepositoryGetByEmail(t *testing.T) {

}

func TestRepositorySetAvatar(t *testing.T) {

}

func TestRepositoryEditPw(t *testing.T) {

}

func TestRepositoryGetInfoByID(t *testing.T) {

}

func TestRepositoryGetInfoByEmail(t *testing.T) {

}
