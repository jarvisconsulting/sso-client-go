package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/jarvisconsulting/sso-client-go/pkg/config"
)

type SessionMiddleware struct {
	store  sessions.Store
	config *config.Config
}

func NewSessionMiddleware(store sessions.Store, config *config.Config) *SessionMiddleware {
	return &SessionMiddleware{
		store:  store,
		config: config,
	}
}

func (m *SessionMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := m.store.Get(c.Request, m.config.SessionName)

		if session.IsNew {
			// Set initial expiry time for new sessions
			expiry := time.Now().Add(time.Duration(m.config.SessionMaxAge) * time.Second)
			session.Values["expiry_time"] = expiry.Unix()
			if err := session.Save(c.Request, c.Writer); err != nil {
				c.Error(err)
			}
		} else if m.config.EnableSlidingWindow {
			// Handle session extension for existing sessions
			expiryTime := session.Values["expiry_time"]
			if expiryTime != nil {
				expiry, ok := expiryTime.(int64)
				if ok {
					timeUntilExpiry := time.Until(time.Unix(expiry, 0))
					thresholdDuration := time.Duration(m.config.SessionExtensionThreshold) * time.Second

					// If session is about to expire within the threshold
					if timeUntilExpiry <= thresholdDuration {
						// Extend the session
						newExpiry := time.Now().Add(time.Duration(m.config.SessionExtensionDuration) * time.Second)
						session.Values["expiry_time"] = newExpiry.Unix()
						if err := session.Save(c.Request, c.Writer); err != nil {
							c.Error(err)
						}
					}
				}
			}
		}

		c.Next()
	}
}
