package application

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *Application) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/sso-health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
