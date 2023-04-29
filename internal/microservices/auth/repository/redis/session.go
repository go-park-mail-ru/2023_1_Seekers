package redis

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/repository"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/rand"
	pkgErrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
)

type sessionsDB struct {
	cfg           *config.Config
	redisSessions *redis.Client
}

func NewSessionRepo(c *config.Config, redisClient *redis.Client) repository.SessionRepoI {
	err := redisClient.Set(context.Background(), "randgeneratedcookie12334524524523542", 1, c.Sessions.CookieTTL).Err()
	if err != nil {
		log.Fatal("failed init session repo")
	}

	return &sessionsDB{
		cfg:           c,
		redisSessions: redisClient,
	}
}

func (sDb *sessionsDB) CreateSession(uID uint64) (*models.Session, error) {
	value, err := rand.String(sDb.cfg.Sessions.CookieLen)
	if err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInternal, "cant generate cookie")
	}
	err = sDb.redisSessions.Set(context.Background(), value, uID, sDb.cfg.Sessions.CookieTTL).Err()
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
