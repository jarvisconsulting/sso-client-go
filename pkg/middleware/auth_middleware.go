package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jarvisconsulting/sso-client-go/pkg/auth"
	"github.com/jarvisconsulting/sso-client-go/pkg/store"
)

type AuthMiddleware struct {
	sessionStore store.SessionStore
	sessionName  string
	signInURL    string
}

func NewAuthMiddleware(sessionStore store.SessionStore, sessionName, signInURL string) *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore: sessionStore,
		sessionName:  sessionName,
		signInURL:    signInURL,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := m.sessionStore.GetStore().Get(c.Request, m.sessionName)
		if err != nil {
			c.Redirect(http.StatusFound, m.signInURL)
			c.Abort()
			return
		}

		_, ok := session.Values[auth.SessionUserIDKey]
		if !ok {
			c.Redirect(http.StatusFound, m.signInURL)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) SetUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := m.sessionStore.GetStore().Get(c.Request, m.sessionName)
		if err == nil {
			if userID, ok := session.Values[auth.SessionUserIDKey].(uint); ok {
				c.Set("user_id", userID)
			}
		}
		c.Next()
	}
}
