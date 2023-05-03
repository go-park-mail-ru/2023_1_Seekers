package redis

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-redis/redismock/v9"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createConfig() *config.Config {
	cfg := new(config.Config)
	ttl, _ := time.ParseDuration("2400h0m0s")
	cfg.Sessions.CookieTTL = ttl
	cfg.Sessions.CookieLen = 0

	return cfg
}

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestDelivery_CreateSession(t *testing.T) {
	cfg := createConfig()
	var fakeSession *models.Session
	generateFakeData(&fakeSession)
	fakeSession.SessionID = ""

	db, mock := redismock.NewClientMock()
	mock.Regexp().ExpectSet(`.*`, 1, cfg.Sessions.CookieTTL).SetVal("admin")
	mock.Regexp().ExpectSet(`.*`, fakeSession.UID, cfg.Sessions.CookieTTL).SetVal(fakeSession.SessionID)

	sessionRepo := NewSessionRepo(cfg, db)
	response, err := sessionRepo.CreateSession(fakeSession.UID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, response, fakeSession)
	}
}

func TestDelivery_DeleteSession(t *testing.T) {
	cfg := createConfig()

	var fakeSession *models.Session
	generateFakeData(&fakeSession)

	db, mock := redismock.NewClientMock()
	mock.Regexp().ExpectSet(`.*`, 1, cfg.Sessions.CookieTTL).SetVal("admin")
	mock.ExpectDel(fakeSession.SessionID).SetVal(int64(fakeSession.UID))

	sessionRepo := NewSessionRepo(cfg, db)
	err := sessionRepo.DeleteSession(fakeSession.SessionID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestDelivery_GetSession(t *testing.T) {
	cfg := createConfig()

	var fakeSession *models.Session
	generateFakeData(&fakeSession)

	db, mock := redismock.NewClientMock()
	mock.Regexp().ExpectSet(`.*`, 1, cfg.Sessions.CookieTTL).SetVal("admin")
	mock.ExpectGet(fakeSession.SessionID).SetVal(fmt.Sprintf("%d", fakeSession.UID))

	sessionRepo := NewSessionRepo(cfg, db)
	response, err := sessionRepo.GetSession(fakeSession.SessionID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, response, fakeSession)
	}
}
