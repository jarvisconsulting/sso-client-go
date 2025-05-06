package application

import (
	"log"
	"net/http"
	"sso-go-client/config"
	"sso-go-client/database"

	"github.com/boj/redistore"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

type Application struct {
	PrimaryDB   *gorm.DB
	SecondaryDB *gorm.DB
	Config      *config.Config
	Store       sessions.Store
}

func NewApplication(cfg config.Config) *Application {
	PrimaryDB, err := database.DBInit(cfg.PrimaryDBHost, cfg.PrimaryDBUser, cfg.PrimaryDBPassword, cfg.PrimaryDBName, cfg.PrimaryDBPort, cfg.PrimaryDBSslMode)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Printf("Connected to primary database: %s", cfg.PrimaryDBName)

	SecondaryDB, err := database.DBInit(cfg.SecondaryDBHost, cfg.SecondaryDBUser, cfg.SecondaryDBPassword, cfg.SecondaryDBName, cfg.SecondaryDBPort, cfg.SecondaryDBSslMode)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Printf("Connected to secondary database: %s", cfg.SecondaryDBName)
	sessionKey := cfg.SessionKey

	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(cfg.RedisURI)
		},
	}
	store, err := redistore.NewRediStoreWithPool(pool, []byte(sessionKey))
	if err != nil {
		log.Fatal("Failed to connect to Redis for session store:", err)
	}
	log.Printf("Connected to Redis for session store: %s", cfg.RedisURI)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,                 // Prevent client-side JavaScript from accessing the cookie
		Secure:   true,                 // Only send the cookie over HTTPS
		SameSite: http.SameSiteLaxMode, // Prevent the cookie from being sent along with requests to other sites
	}

	return &Application{
		PrimaryDB:   PrimaryDB,
		SecondaryDB: SecondaryDB,
		Config:      &cfg,
		Store:       store,
	}
}
