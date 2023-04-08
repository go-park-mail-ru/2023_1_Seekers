package redis

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-redis/redismock/v9"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	//"github.com/alicebob/miniredis/v2"
	_ "github.com/go-redis/redismock/v9"

	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"testing"
)

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestDelivery_CreateSession(t *testing.T) {
	config.CookieLen = 0
	var fakeSession *models.Session
	generateFakeData(&fakeSession)
	fakeSession.SessionID = ""

	db, mock := redismock.NewClientMock()
	mock.Regexp().ExpectSet(`.*`, fakeSession.UID, config.CookieTTL).SetVal(fakeSession.SessionID)

	sessionRepo := NewSessionRepo(db)
	response, err := sessionRepo.CreateSession(fakeSession.UID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, response, fakeSession)
	}
}

func TestDelivery_DeleteSession(t *testing.T) {
	var fakeSession *models.Session
	generateFakeData(&fakeSession)

	db, mock := redismock.NewClientMock()
	mock.ExpectDel(fakeSession.SessionID).SetVal(int64(fakeSession.UID))

	sessionRepo := NewSessionRepo(db)
	err := sessionRepo.DeleteSession(fakeSession.SessionID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestDelivery_GetSession(t *testing.T) {
	var fakeSession *models.Session
	generateFakeData(&fakeSession)

	db, mock := redismock.NewClientMock()
	mock.ExpectGet(fakeSession.SessionID).SetVal(fmt.Sprintf("%d", fakeSession.UID))

	sessionRepo := NewSessionRepo(db)
	response, err := sessionRepo.GetSession(fakeSession.SessionID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, response, fakeSession)
	}
}
