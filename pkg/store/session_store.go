package store

import (
	redistore "github.com/boj/redistore"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"
	"github.com/jarvisconsulting/sso-client-go/pkg/config"
)

type SessionStore interface {
	GetStore() sessions.Store
	Close() error
}

type RedisSessionStore struct {
	store  *redistore.RediStore
	config *config.Config
}

func NewRedisSessionStore(redisURI, sessionKey string, isRedisSecure bool, sessionMaxAge int) (SessionStore, error) {
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisURI)
		},
	}

	store, err := redistore.NewRediStoreWithPool(pool, []byte(sessionKey))
	if err != nil {
		return nil, err
	}

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionMaxAge, // 1 hour
		HttpOnly: true,
		Secure:   isRedisSecure,
	}

	return &RedisSessionStore{
		store: store,
	}, nil
}

func (s *RedisSessionStore) GetStore() sessions.Store {
	return s.store
}

func (s *RedisSessionStore) Close() error {
	return s.store.Close()
}
