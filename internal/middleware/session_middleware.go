package middleware

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

type SessionMiddleware struct {
	store sessions.Store
}

func NewSessionMiddleware(store sessions.Store) *SessionMiddleware {
	return &SessionMiddleware{
		store: store,
	}
}

func (m *SessionMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.store.Get(r, "sso_session")
		if err != nil {
			log.Printf("Session error (creating new session): %v", err)
			session, _ = m.store.New(r, "sso_session")
			session.Options.Path = "/"
			session.Options.HttpOnly = true
			session.Options.Secure = r.TLS != nil
			session.Options.MaxAge = 3600 // 1 hour
			session.Save(r, w)
		}

		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
