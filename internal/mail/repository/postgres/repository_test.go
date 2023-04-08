package postgres

import (
	"database/sql"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"regexp"
	"testing"
)

func mockDB() (*sql.DB, *gorm.DB, sqlmock.Sqlmock, error) {
	os.Setenv(config.DBSchemaNameEnv, "mail")

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

func TestRepository_SelectFolderByUserNFolder(t *testing.T) {
	var fakeFolder *models.Folder
	generateFakeData(&fakeFolder)
	userID := uint64(1)
	folderSlug := "inbox"

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"folder_id", "user_id", "local_name", "name", "messages_unseen", "messages_count"}).
		AddRow(fakeFolder.FolderID, fakeFolder.UserID, fakeFolder.LocalName, fakeFolder.Name,
			fakeFolder.MessagesUnseen, fakeFolder.MessagesCount)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "folders" WHERE user_id = $1 AND local_name = $2`)).
		WithArgs(userID, folderSlug).WillReturnRows(rows)

	mailRepo := New(gormDB)
	response, err := mailRepo.SelectFolderByUserNFolder(userID, folderSlug)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeFolder, response)
	}
}

func TestRepository_SelectFoldersByUser(t *testing.T) {
	var fakeFolders []models.Folder
	generateFakeData(&fakeFolders)
	userID := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"folder_id", "user_id", "local_name", "name", "messages_unseen", "messages_count"})

	for _, folder := range fakeFolders {
		rows.AddRow(folder.FolderID, folder.UserID, folder.LocalName, folder.Name,
			folder.MessagesUnseen, folder.MessagesCount)
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "folders" WHERE user_id = $1`)).WithArgs(userID).WillReturnRows(rows)

	mailRepo := New(gormDB)
	response, err := mailRepo.SelectFoldersByUser(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeFolders, response)
	}
}

func TestRepository_SelectFolderMessagesByUserNFolder(t *testing.T) {
	var fakeMessages []models.MessageInfo
	generateFakeData(&fakeMessages)
	for i, _ := range fakeMessages {
		fakeMessages[i].Recipients = nil
		fakeMessages[i].FromUser.Email = ""
		fakeMessages[i].FromUser.FirstName = ""
		fakeMessages[i].FromUser.LastName = ""
		fakeMessages[i].ReplyTo = nil
	}

	userID := uint64(1)
	folderID := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"message_id", "from_user_id", "title", "created_at", "text",
		"reply_to_message_id", "seen", "favorite", "deleted"})
	for _, message := range fakeMessages {
		rows.AddRow(message.MessageID, message.FromUser.UserID, message.Title, message.CreatedAt,
			message.Text, message.ReplyToMessageID, message.Seen, message.Favorite, message.Deleted)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mail"."boxes" JOIN mail.messages using(message_id) WHERE user_id = $1 AND folder_id = $2`)).
		WithArgs(userID, folderID).WillReturnRows(rows)

	mailRepo := New(gormDB)
	response, err := mailRepo.SelectFolderMessagesByUserNFolder(userID, folderID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeMessages, response)
	}
}

func TestRepository_SelectRecipientsByMessage(t *testing.T) {
	var fakeRecipients []uint64
	generateFakeData(&fakeRecipients)

	messageID := uint64(1)
	fromUserID := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"user_id"})
	for _, recipient := range fakeRecipients {
		rows.AddRow(recipient)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "user_id" FROM "mail"."boxes" WHERE message_id = $1 AND user_id != $2`)).
		WithArgs(messageID, fromUserID).WillReturnRows(rows)

	mailRepo := New(gormDB)
	response, err := mailRepo.SelectRecipientsByMessage(messageID, fromUserID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeRecipients, response)
	}
}

func TestRepository_SelectMessageByUserNMessage(t *testing.T) {
	var fakeMessage *models.MessageInfo
	generateFakeData(&fakeMessage)
	fakeMessage.Recipients = nil
	fakeMessage.FromUser.Email = ""
	fakeMessage.FromUser.FirstName = ""
	fakeMessage.FromUser.LastName = ""
	fakeMessage.ReplyTo = nil

	userID := uint64(1)
	messageID := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"message_id", "from_user_id", "title", "created_at", "text",
		"reply_to_message_id", "seen", "favorite", "deleted"})
	rows.AddRow(fakeMessage.MessageID, fakeMessage.FromUser.UserID, fakeMessage.Title, fakeMessage.CreatedAt,
		fakeMessage.Text, fakeMessage.ReplyToMessageID, fakeMessage.Seen, fakeMessage.Favorite, fakeMessage.Deleted)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mail"."boxes" JOIN mail.messages using(message_id) WHERE user_id = $1 AND message_id = $2`)).
		WithArgs(userID, messageID).WillReturnRows(rows)

	mailRepo := New(gormDB)
	response, err := mailRepo.SelectMessageByUserNMessage(userID, messageID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeMessage, response)
	}
}

func TestRepository_InsertMessage(t *testing.T) {
	messageID := uint64(1)

	var fakeMessage *models.MessageInfo
	generateFakeData(&fakeMessage)
	fakeMessage.MessageID = messageID
	fakeMessage.Recipients = fakeMessage.Recipients[:1]
	fakeMessage.Seen = true
	user2folder := []models.User2Folder{
		{
			UserID:   fakeMessage.FromUser.UserID,
			FolderID: 1,
		},
		{
			UserID:   fakeMessage.Recipients[0].UserID,
			FolderID: 2,
		},
	}

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"message_id"})
	rows.AddRow(fakeMessage.MessageID)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "mail"."messages" ("from_user_id","title","text","created_at","reply_to_message_id")
	VALUES ($1,$2,$3,$4,$5)`)).WithArgs(fakeMessage.FromUser.UserID, fakeMessage.Title, fakeMessage.Text,
		fakeMessage.CreatedAt, fakeMessage.ReplyToMessageID).WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "mail"."boxes" ("user_id","message_id","folder_id","seen","favorite","deleted")
VALUES ($1,$2,$3,$4,$5,$6)`)).WithArgs(user2folder[0].UserID, fakeMessage.MessageID, user2folder[0].FolderID, fakeMessage.Seen,
		fakeMessage.Favorite, fakeMessage.Deleted).WillReturnResult(sqlmock.NewResult(int64(0), 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "mail"."boxes" ("user_id","message_id","folder_id","seen","favorite","deleted")
	VALUES ($1,$2,$3,$4,$5,$6)`)).WithArgs(user2folder[1].UserID, fakeMessage.MessageID, user2folder[1].FolderID, false,
		fakeMessage.Favorite, fakeMessage.Deleted).WillReturnResult(sqlmock.NewResult(int64(0), 1))
	mock.ExpectCommit()

	mailRepo := New(gormDB)
	err = mailRepo.InsertMessage(fakeMessage.FromUser.UserID, fakeMessage, user2folder)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeMessage.MessageID, messageID)
	}
}

func TestRepository_InsertFolder(t *testing.T) {
	folderID := uint64(1)

	var fakeFolder *models.Folder
	generateFakeData(&fakeFolder)
	fakeFolder.FolderID = 0

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"folder_id"})
	rows.AddRow(folderID)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "folders" ("user_id","local_name","name","messages_unseen","messages_count")
	VALUES ($1,$2,$3,$4,$5)`)).WithArgs(fakeFolder.UserID, fakeFolder.LocalName, fakeFolder.Name,
		fakeFolder.MessagesUnseen, fakeFolder.MessagesCount).WillReturnRows(rows)
	mock.ExpectCommit()

	mailRepo := New(gormDB)
	folderID, err = mailRepo.InsertFolder(fakeFolder)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, folderID, folderID)
	}
}

func TestRepository_UpdateMessageState(t *testing.T) {
	userID := uint64(1)
	messageID := uint64(1)
	stateName := "seen"
	stateValue := false

	var fakeFolder *models.Folder
	generateFakeData(&fakeFolder)
	fakeFolder.FolderID = 0

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE "mail"."boxes" SET "%s"=$1 WHERE user_id = $2 AND message_id = $3`,
		stateName))).
		WithArgs(stateValue, userID, messageID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	mailRepo := New(gormDB)
	err = mailRepo.UpdateMessageState(userID, messageID, stateName, stateValue)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	}
}