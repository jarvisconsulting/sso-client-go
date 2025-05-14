package ssoclient

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/jarvisconsulting/sso-client-go/pkg/auth"
	"github.com/jarvisconsulting/sso-client-go/pkg/config"
	"github.com/jarvisconsulting/sso-client-go/pkg/middleware"
	"github.com/jarvisconsulting/sso-client-go/pkg/models"
	"github.com/jarvisconsulting/sso-client-go/pkg/store"
)

type Client struct {
	config       *config.Config
	authService  *auth.AuthService
	authHandler  *auth.Handler
	sessionStore store.SessionStore
}

type Handlers struct {
	SignIn         gin.HandlerFunc
	SignOut        gin.HandlerFunc
	Callback       gin.HandlerFunc
	User           gin.HandlerFunc
	WebhookSignOut gin.HandlerFunc
}

type Middleware struct {
	RequireAuth gin.HandlerFunc
	SetUserID   gin.HandlerFunc
	Session     gin.HandlerFunc
}

func New(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	sessionStore, err := store.NewRedisSessionStore(cfg.RedisURI, cfg.SessionKey, cfg.IsRedisSecure, cfg.SessionMaxAge)
	if err != nil {
		return nil, err
	}

	return &Client{
		config:       cfg,
		sessionStore: sessionStore,
	}, nil
}

func (c *Client) WithRepository(primaryDB *gorm.DB, secondaryDB *gorm.DB) *Client {
	userRepo := NewUserRepository(primaryDB, secondaryDB)

	handlerConfig := &auth.Config{
		SignInURL:   c.config.SignInURL,
		CallbackURL: c.config.CallbackURL,
		RootURL:     c.config.RootURL,
	}

	c.authService = auth.NewAuthService(userRepo, c.config, c.sessionStore)
	c.authHandler = auth.NewHandler(c.authService, handlerConfig)

	return c
}

func (c *Client) GetHandlers() *Handlers {
	if c.authHandler == nil {
		log.Fatal("AuthHandler is nil. Make sure to call WithRepository before GetHandlers")
	}

	return &Handlers{
		SignIn:         c.authHandler.SignIn,
		SignOut:        c.authHandler.SignOut,
		Callback:       c.authHandler.Callback,
		User:           c.authHandler.User,
		WebhookSignOut: c.authHandler.WebhookSignOut,
	}
}

func (c *Client) GetMiddleware() *Middleware {
	sessionMiddleware := middleware.NewSessionMiddleware(c.sessionStore.GetStore(), c.config)
	authMiddleware := middleware.NewAuthMiddleware(c.sessionStore, c.config.SessionName, c.config.SignInURL)

	return &Middleware{
		RequireAuth: authMiddleware.RequireAuth(),
		SetUserID:   authMiddleware.SetUserID(),
		Session:     sessionMiddleware.Handler(),
	}
}

func (c *Client) Close() error {
	if c.sessionStore != nil {
		return c.sessionStore.Close()
	}
	return nil
}

func (c *Client) IsUserSignedIn(r *http.Request) bool {
	return c.authService.IsUserSignedIn(r)
}

func (c *Client) GetUserIDFromSession(r *http.Request) (uint, error) {
	return c.authService.GetUserIDFromSession(r)
}

func (c *Client) GetUserByID(id uint) (*models.User, error) {
	user, err := c.authService.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
