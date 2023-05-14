package postgres

import (
	"database/sql"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

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

func createConfig() *config.Config {
	cfg := new(config.Config)
	cfg.DB.DBSchemaName = "mail"

	return cfg
}

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestRepository_SelectFolderByUserNFolderSlug(t *testing.T) {
	cfg := createConfig()

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

	mailRepo := New(cfg, gormDB)
	response, err := mailRepo.SelectFolderByUserNFolderSlug(userID, folderSlug)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeFolder, response)
	}
}

func TestRepository_SelectFolderByUserNFolderName(t *testing.T) {
	cfg := createConfig()

	var fakeFolder *models.Folder
	generateFakeData(&fakeFolder)
	userID := uint64(1)
	folderName := "my"

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"folder_id", "user_id", "local_name", "name", "messages_unseen", "messages_count"}).
		AddRow(fakeFolder.FolderID, fakeFolder.UserID, fakeFolder.LocalName, fakeFolder.Name,
			fakeFolder.MessagesUnseen, fakeFolder.MessagesCount)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "folders" WHERE user_id = $1 AND name = $2`)).
		WithArgs(userID, folderName).WillReturnRows(rows)

	mailRepo := New(cfg, gormDB)
	response, err := mailRepo.SelectFolderByUserNFolderName(userID, folderName)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeFolder, response)
	}
}

func TestRepository_SelectFoldersByUser(t *testing.T) {
	cfg := createConfig()

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

	mailRepo := New(cfg, gormDB)
	response, err := mailRepo.SelectFoldersByUser(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeFolders, response)
	}
}

func TestRepository_SelectFolderByUserNMessage(t *testing.T) {
	cfg := createConfig()

	var fakeFolder *models.Folder
	generateFakeData(&fakeFolder)
	userID := uint64(1)
	messageID := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"folder_id", "user_id", "local_name", "name", "messages_unseen", "messages_count"}).
		AddRow(fakeFolder.FolderID, fakeFolder.UserID, fakeFolder.LocalName, fakeFolder.Name, fakeFolder.MessagesUnseen, fakeFolder.MessagesCount)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT folders.* FROM "boxes" JOIN mail.folders using(folder_id) WHERE boxes.user_id = $1 AND message_id = $2 ORDER BY "boxes"."user_id" LIMIT 1`)).
		WithArgs(userID, messageID).WillReturnRows(rows)

	mailRepo := New(cfg, gormDB)
	response, err := mailRepo.SelectFolderByUserNMessage(userID, messageID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeFolder, response)
	}
}

func TestRepository_SelectFolderMessagesByUserNFolderID(t *testing.T) {
	cfg := createConfig()

	var fakeMessages []models.MessageInfo
	generateFakeData(&fakeMessages)
	for i := range fakeMessages {
		fakeMessages[i].IsDraft = false
		fakeMessages[i].Recipients = nil
		fakeMessages[i].FromUser.Email = ""
		fakeMessages[i].FromUser.FirstName = ""
		fakeMessages[i].FromUser.LastName = ""
		fakeMessages[i].ReplyTo = nil
		fakeMessages[i].Attachments = nil
		fakeMessages[i].AttachmentsSize = ""
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

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "boxes" JOIN mail.messages using(message_id) WHERE user_id = $1
AND folder_id = $2 AND is_draft = $3 ORDER BY created_at DESC`)).
		WithArgs(userID, folderID, false).WillReturnRows(rows)

	mailRepo := New(cfg, gormDB)
	response, err := mailRepo.SelectFolderMessagesByUserNFolderID(userID, folderID, false)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeMessages, response)
	}
}

func TestRepository_DeleteFolder(t *testing.T) {
	cfg := createConfig()

	folderID := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "folders" WHERE folder_id = $1`)).
		WithArgs(folderID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mailRepo := New(cfg, gormDB)
	err = mailRepo.DeleteFolder(folderID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_DeleteMessageForUser(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	messageID := uint64(1)
	folderID := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "boxes" WHERE user_id = $1 AND message_id = $2 AND folder_id = $3`)).
		WithArgs(userID, messageID, folderID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mailRepo := New(cfg, gormDB)
	err = mailRepo.DeleteBox(userID, messageID, folderID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_DeleteMessageFromMessages(t *testing.T) {
	cfg := createConfig()

	messageID := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "messages" WHERE message_id = $1`)).
		WithArgs(messageID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mailRepo := New(cfg, gormDB)
	err = mailRepo.DeleteMessageFromMessages(messageID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_UpdateFolder(t *testing.T) {
	cfg := createConfig()

	var fakeFolder models.Folder
	generateFakeData(&fakeFolder)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "folders" SET "folder_id"=$1,"user_id"=$2,"local_name"=$3,"name"=$4,"messages_unseen"=$5,"messages_count"=$6 
WHERE "folder_id" = $7`)).WithArgs(fakeFolder.FolderID, fakeFolder.UserID, fakeFolder.LocalName, fakeFolder.Name, fakeFolder.MessagesUnseen,
		fakeFolder.MessagesCount, fakeFolder.FolderID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	mailRepo := New(cfg, gormDB)
	err = mailRepo.UpdateFolder(fakeFolder)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_SelectRecipientsByMessage(t *testing.T) {
	cfg := createConfig()

	var fakeRecipients []uint64
	generateFakeData(&fakeRecipients)

	messageID := uint64(1)
	fromUserID := uint64(1)

	for i, other := range fakeRecipients {
		if other == fromUserID {
			fakeRecipients = append(fakeRecipients[:i], fakeRecipients[i+1:]...)
		}
	}

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"user_id"})
	for _, recipient := range fakeRecipients {
		rows.AddRow(recipient)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "user_id" FROM "boxes" WHERE message_id = $1`)).
		WithArgs(messageID).WillReturnRows(rows)

	mailRepo := New(cfg, gormDB)
	response, err := mailRepo.SelectRecipientsByMessage(messageID, fromUserID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeRecipients, response)
	}
}

func TestRepository_SelectMessageByUserNMessage(t *testing.T) {
	cfg := createConfig()

	var fakeMessage *models.MessageInfo
	generateFakeData(&fakeMessage)
	fakeMessage.IsDraft = false
	fakeMessage.Recipients = nil
	fakeMessage.FromUser.Email = ""
	fakeMessage.FromUser.FirstName = ""
	fakeMessage.FromUser.LastName = ""
	fakeMessage.ReplyTo = nil
	fakeMessage.Attachments = nil
	fakeMessage.AttachmentsSize = ""

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

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "boxes" JOIN mail.messages using(message_id) WHERE user_id = $1
AND message_id = $2`)).
		WithArgs(userID, messageID).WillReturnRows(rows)

	mailRepo := New(cfg, gormDB)
	response, err := mailRepo.SelectMessageByUserNMessage(userID, messageID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeMessage, response)
	}
}

func TestRepository_InsertMessage(t *testing.T) {
	cfg := createConfig()

	messageID := uint64(1)

	var fakeMessage *models.MessageInfo
	generateFakeData(&fakeMessage)
	fakeMessage.MessageID = messageID
	fakeMessage.Recipients = fakeMessage.Recipients[:1]
	fakeMessage.Seen = true
	fakeMessage.Attachments = nil
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
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "messages" ("from_user_id","title","text","created_at","reply_to_message_id")
	VALUES ($1,$2,$3,$4,$5)`)).WithArgs(fakeMessage.FromUser.UserID, fakeMessage.Title, fakeMessage.Text,
		fakeMessage.CreatedAt, fakeMessage.ReplyToMessageID).WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "boxes" ("user_id","message_id","folder_id","seen","favorite","deleted","is_draft")
VALUES ($1,$2,$3,$4,$5,$6,$7)`)).WithArgs(user2folder[0].UserID, fakeMessage.MessageID, user2folder[0].FolderID, fakeMessage.Seen,
		fakeMessage.Favorite, fakeMessage.Deleted, fakeMessage.IsDraft).WillReturnResult(sqlmock.NewResult(int64(0), 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "boxes" ("user_id","message_id","folder_id","seen","favorite","deleted","is_draft")
	VALUES ($1,$2,$3,$4,$5,$6,$7)`)).WithArgs(user2folder[1].UserID, fakeMessage.MessageID, user2folder[1].FolderID, false,
		fakeMessage.Favorite, fakeMessage.Deleted, fakeMessage.IsDraft).WillReturnResult(sqlmock.NewResult(int64(0), 1))
	mock.ExpectCommit()

	mailRepo := New(cfg, gormDB)
	err = mailRepo.InsertMessage(fakeMessage.FromUser.UserID, fakeMessage, user2folder)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeMessage.MessageID, messageID)
	}
}

func TestRepository_SelectCustomFoldersByUser(t *testing.T) {
	cfg := createConfig()

	var fakeFolders []models.Folder
	generateFakeData(&fakeFolders)
	defaultLocalNames := []string{"inbox", "outbox", "drafts", "trash", "spam"}
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
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "folders" WHERE user_id = $1 AND local_name NOT IN ($2,$3,$4,$5,$6)`)).WithArgs(userID, "inbox", "outbox", "drafts", "trash", "spam").WillReturnRows(rows)

	mailRepo := New(cfg, gormDB)
	response, err := mailRepo.SelectCustomFoldersByUser(userID, defaultLocalNames)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeFolders, response)
	}
}

func TestRepository_UpdateMessage(t *testing.T) {
	cfg := createConfig()
	var fakeMessage *models.MessageInfo
	generateFakeData(&fakeMessage)
	toInsert := []models.User2Folder{{
		1, 1,
	}}
	toDelete := []models.User2Folder{{
		2, 2,
	}}

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "messages" SET "title"=$1,"text"=$2,"created_at"=$3,"reply_to_message_id"=$4 WHERE message_id = $5`)).
		WithArgs(fakeMessage.Title, fakeMessage.Text, fakeMessage.CreatedAt, fakeMessage.ReplyToMessageID, fakeMessage.MessageID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "boxes" ("user_id","message_id","folder_id","seen","favorite","deleted","is_draft")
VALUES ($1,$2,$3,$4,$5,$6,$7)`)).WithArgs(toInsert[0].UserID, fakeMessage.MessageID, toInsert[0].FolderID, fakeMessage.Seen,
		fakeMessage.Favorite, fakeMessage.Deleted, fakeMessage.IsDraft).WillReturnResult(sqlmock.NewResult(int64(0), 1))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "boxes" WHERE user_id = $1 AND folder_id = $2 AND message_id = $3`)).
		WithArgs(toDelete[0].UserID, toDelete[0].FolderID, fakeMessage.MessageID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mailRepo := New(cfg, gormDB)
	err = mailRepo.UpdateMessage(fakeMessage, toInsert, toDelete)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_InsertFolder(t *testing.T) {
	cfg := createConfig()

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

	mailRepo := New(cfg, gormDB)
	folderID, err = mailRepo.InsertFolder(fakeFolder)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, folderID, folderID)
	}
}

func TestRepository_UpdateMessageState(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	messageID := uint64(1)
	folderID := uint64(1)
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
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE "boxes" SET "%s"=$1 WHERE user_id = $2 AND message_id = $3 AND folder_id = $4`,
		stateName))).
		WithArgs(stateValue, userID, messageID, folderID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	mailRepo := New(cfg, gormDB)
	err = mailRepo.UpdateMessageState(userID, messageID, folderID, stateName, stateValue)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_UpdateMessageFolder(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	messageID := uint64(1)
	oldFolderID := uint64(1)
	newFolderID := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "boxes" SET "folder_id"=$1 WHERE user_id = $2 AND message_id = $3 AND folder_id = $4`)).
		WithArgs(newFolderID, userID, messageID, oldFolderID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	mailRepo := New(cfg, gormDB)
	err = mailRepo.UpdateMessageFolder(userID, messageID, oldFolderID, newFolderID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}
