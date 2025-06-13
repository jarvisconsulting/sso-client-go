package middleware

import (
	"net/http"
	"strconv"

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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		_, ok := session.Values[auth.SessionUserIDKey]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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

func (m *AuthMiddleware) SetIsMobile() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := m.sessionStore.GetStore().Get(c.Request, m.sessionName)
		if err == nil {
			if isMobile, ok := session.Values[auth.SessionIsMobileKey].(bool); ok {
				c.Set("is_mobile", isMobile)
				c.Header("Is-Mobile", strconv.FormatBool(isMobile))
			} else {
				c.Set("is_mobile", false)
				c.Header("Is-Mobile", "false")
			}
		}
		c.Next()
	}
}
