package config

type Config struct {
	// SSO Server configuration
	SSOUrl       string `json:"sso_url" validate:"required"`
	ClientID     string `json:"client_id" validate:"required"`
	CallbackURL  string `json:"callback_url" validate:"required"`
	SignOutURL   string `json:"sign_out_url" validate:"required"`
	PublicKeyAPI string `json:"public_key_api" validate:"required"`

	// Application URLs
	RootURL     string `json:"root_url" validate:"required"`
	SignInURL   string `json:"sign_in_url" validate:"required"`
	FrontendURL string `json:"frontend_url" validate:"required"`

	// Redis configuration for session storage
	RedisURI      string `json:"redis_uri" validate:"required"`
	SessionKey    string `json:"session_key" validate:"required"`
	SessionName   string `json:"session_name" validate:"required"`
	IsRedisSecure bool   `json:"is_redis_secure"`

	// Session configuration
	SessionMaxAge       int  `json:"session_max_age" validate:"required,min=300"` // minimum 5 minutes
	EnableSlidingWindow bool `json:"enable_sliding_window"`

	// Sliding window configuration (required if EnableSlidingWindow is true)
	SessionExtensionDuration  int `json:"session_extension_duration" validate:"required_if=EnableSlidingWindow true,min=300"`
	SessionExtensionThreshold int `json:"session_extension_threshold" validate:"required_if=EnableSlidingWindow true,min=60"`
}

func DefaultConfig() *Config {
	return &Config{
		SSOUrl:       "https://sso.example.com",
		ClientID:     "your-client-id",
		CallbackURL:  "http://localhost:8080/auth/callback",
		SignOutURL:   "http://localhost:8080/auth/signout",
		PublicKeyAPI: "http://localhost:8080/auth/public-key",
		RootURL:      "http://localhost:8080",
		SignInURL:    "http://localhost:8080/auth/signin",
		FrontendURL:  "http://localhost:3000",

		RedisURI:      "redis://:123456@localhost:6379/0",
		SessionKey:    "your-session-key",
		SessionName:   "myapp_session",
		IsRedisSecure: true,

		SessionMaxAge:             3600,
		EnableSlidingWindow:       true,
		SessionExtensionDuration:  1800,
		SessionExtensionThreshold: 1200,
	}
}
