package auth

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService *AuthService
	config      *Config
}

type Config struct {
	SignInURL   string
	SignOutURL  string
	CallbackURL string
	RootURL     string
}

func NewHandler(authService *AuthService, config *Config) *Handler {
	return &Handler{
		authService: authService,
		config:      config,
	}
}

func (h *Handler) SignIn(c *gin.Context) {
	if h.authService.IsUserSignedIn(c.Request) {
		c.Redirect(http.StatusFound, h.config.RootURL)
		return
	}

	log.Printf("Redirecting to SignInURL: %s", h.config.SignInURL)
	c.Redirect(http.StatusFound, h.config.SignInURL)
}

func (h *Handler) SignOut(c *gin.Context) {
	log.Printf("Signing out user")
	err := h.authService.SignOutUser(c.Writer, c.Request)
	if err != nil {
		log.Printf("Failed to sign out: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign out"})
		return
	}

	log.Printf("Redirecting to frontend after signout: %s", h.config.SignInURL)
	c.Redirect(http.StatusFound, h.config.SignInURL)
}

func (h *Handler) Callback(c *gin.Context) {
	params := make(map[string]string)

	pyIdToken := c.Query("py_id_token")
	if pyIdToken != "" {
		params["id_token"] = pyIdToken
		log.Printf("Using py_id_token for authentication")
	} else {
		params["id_token"] = c.Query("id_token")
		log.Printf("Using id_token for authentication")
	}

	endpoint := c.Query("endpoint")

	userID, err := h.authService.HandleCallback(params)
	if err != nil {
		log.Printf("Failed to process callback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process callback"})
		return
	}

	err = h.authService.SignInUser(c.Writer, c.Request, userID)
	if err != nil {
		log.Printf("Failed to sign in user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign in user"})
		return
	}

	redirectURL := h.config.RootURL
	if endpoint != "" {
		redirectURL += "?endpoint=" + url.QueryEscape(endpoint)
	}

	c.Redirect(http.StatusFound, redirectURL)
}

func (h *Handler) User(c *gin.Context) {
	if !h.authService.IsUserSignedIn(c.Request) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not signed in"})
		return
	}

	userID, err := h.authService.GetUserIDFromSession(c.Request)
	if err != nil {
		log.Printf("Failed to get user ID from session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from session"})
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}

func (h *Handler) WebhookSignOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "logout"})
}
