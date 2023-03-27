package redis

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type sessionsDB struct {
	redisSessions *redis.Client
}

func NewSessionRepo(redisClient *redis.Client) auth.SessionRepoI {
	redisClient.Set(context.Background(), "randgeneratedcookie12334524524523542", 1, config.CookieTTL).Err()
	return &sessionsDB{
		redisSessions: redisClient,
	}
}

func (sDb *sessionsDB) CreateSession(uID uint64) (*models.Session, error) {
	value, err := pkg.String(config.CookieLen)
	if err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, "cant generate cookie")
	}
	err = sDb.redisSessions.Set(context.Background(), value, uID, config.CookieTTL).Err()
	if err != nil {
		return nil, pkgErrors.WithMessagef(errors.ErrInternal, "cant set cookie : %v", err.Error())
	}

	return &models.Session{
		UID:       uID,
		SessionID: value,
	}, nil
}

func (sDb *sessionsDB) DeleteSession(sessionID string) error {
	err := sDb.redisSessions.Del(context.Background(), sessionID).Err()
	if err != nil {
		return pkgErrors.WithMessagef(errors.ErrFailedDeleteSession, "delete cookie %v", err.Error())
	}

	return nil
}

func (sDb *sessionsDB) GetSession(sessionID string) (*models.Session, error) {
	uIDstr, err := sDb.redisSessions.Get(context.Background(), sessionID).Result()
	if err != nil {
		return nil, pkgErrors.WithMessagef(errors.ErrFailedGetSession, "get cookie %v", err.Error())
	}
	uID, err := strconv.Atoi(uIDstr)
	if err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return &models.Session{
		UID:       uint64(uID),
		SessionID: sessionID,
	}, nil
}
