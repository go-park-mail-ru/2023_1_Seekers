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

func (u *useCase) Create(uId int) (*model.Session, error) {
	if uId <= 0 {
		return nil, fmt.Errorf("cant create session, user id <= 0")
	}

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

func (u *useCase) GetSession(uId int) (*model.Session, error) {
	s, err := u.sessionRepo.GetSession(uId)
	if err != nil {
		return nil, fmt.Errorf("cant get session by user %w", err)
	}

	return s, nil
}
