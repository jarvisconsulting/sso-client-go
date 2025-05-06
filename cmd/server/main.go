package main

import (
	"log"
	"net/http"
	"sso-go-client/cmd/application"
	"sso-go-client/config"
	"sso-go-client/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	app := application.NewApplication(cfg)

	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	sessionMiddleware := middleware.NewSessionMiddleware(app.Store)
	router.Use(sessionMiddleware.Handler)
	app.RegisterRoutes(router)

	log.Printf("Trying very hard to start the server on PORT:  %s...\n", cfg.SSOPort)
	if err := http.ListenAndServe(":"+cfg.SSOPort, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
