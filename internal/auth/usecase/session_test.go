package usecase

import (
	mockSessionRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/repository/mocks"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/golang/mock/gomock"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUseCase_CreateSession(t *testing.T) {
	var fakeSession *models.Session
	generateFakeData(&fakeSession)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRepo := mockSessionRepo.NewMockSessionRepoI(ctrl)
	sUC := NewSessionUC(sessionRepo)

	sessionRepo.EXPECT().CreateSession(fakeSession.UID).Return(fakeSession, nil)
	response, err := sUC.CreateSession(fakeSession.UID)

	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeSession, response)
	}
}

func TestUseCase_DeleteSession(t *testing.T) {
	sessionID := "adjfkfldlkld"

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRepo := mockSessionRepo.NewMockSessionRepoI(ctrl)
	sUC := NewSessionUC(sessionRepo)

	sessionRepo.EXPECT().DeleteSession(sessionID).Return(nil)
	err := sUC.DeleteSession(sessionID)

	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	}
}

func TestUseCase_GetSession(t *testing.T) {
	var fakeSession *models.Session
	generateFakeData(&fakeSession)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRepo := mockSessionRepo.NewMockSessionRepoI(ctrl)
	sUC := NewSessionUC(sessionRepo)

	sessionRepo.EXPECT().GetSession(fakeSession.SessionID).Return(fakeSession, nil)
	response, err := sUC.GetSession(fakeSession.SessionID)

	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", err, causeErr)
	} else {
		require.Equal(t, fakeSession, response)
	}
}
