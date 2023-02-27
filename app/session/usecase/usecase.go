package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/model"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/session"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/utils"
	"github.com/go-park-mail-ru/2023_1_Seekers/config"
)

type useCase struct {
	sessionRepo session.Repo
}

func New(sr session.Repo) session.UseCase {
	return &useCase{
		sessionRepo: sr,
	}
}

func (u *useCase) Create(uId uint64) (*model.Session, error) {
	value := utils.String(config.CookieLen)
	newSession := model.Session{
		UId:       uId,
		SessionId: value,
	}

	err := u.sessionRepo.Create(newSession)
	if err != nil {
		return nil, fmt.Errorf("cant create session %w", err)
	}

	return &newSession, nil
}

func (u *useCase) Delete(sessionId string) error {
	err := u.sessionRepo.Delete(sessionId)
	if err != nil {
		return fmt.Errorf("cant delete session %w", err)
	}

	return nil
}

func (u *useCase) DeleteByUId(uId uint64) error {
	err := u.sessionRepo.DeleteByUId(uId)
	if err != nil {
		return fmt.Errorf("cant delete session by id %w", err)
	}

	return nil
}

func (u *useCase) GetSession(sessionId string) (*model.Session, error) {
	s, err := u.sessionRepo.GetSession(sessionId)
	if err != nil {
		return nil, fmt.Errorf("cant get session: %w", err)
	}

	return s, nil
}

func (u *useCase) GetSessionByUId(uId uint64) (*model.Session, error) {
	s, err := u.sessionRepo.GetSessionByUId(uId)
	if err != nil {
		return nil, fmt.Errorf("cant get session by user %w", err)
	}

	return s, nil
}
