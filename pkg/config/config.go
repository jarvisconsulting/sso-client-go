package config

type Config struct {
	// Core SSO URLs
	CallbackURL string `json:"callback_url" validate:"required"`
	SignInURL   string `json:"sign_in_url" validate:"required"`
	RootURL     string `json:"root_url" validate:"required"`

	// Redis configuration for session storage
	RedisURI      string `json:"redis_uri" validate:"required"`
	SessionKey    string `json:"session_key" validate:"required"`
	SessionName   string `json:"session_name" validate:"required"`
	IsRedisSecure bool   `json:"is_redis_secure"`

	// Session configuration
	SessionMaxAge int `json:"session_max_age" validate:"required,min=300"` // minimum 5 minutes

	// Optional: Sliding window configuration
	EnableSlidingWindow       bool `json:"enable_sliding_window"`
	SessionExtensionDuration  int  `json:"session_extension_duration,omitempty" validate:"required_if=EnableSlidingWindow true,min=300"`
	SessionExtensionThreshold int  `json:"session_extension_threshold,omitempty" validate:"required_if=EnableSlidingWindow true,min=60"`
}

func DefaultConfig() *Config {
	return &Config{
		CallbackURL: "http://localhost:8080/auth/callback",
		RootURL:     "http://localhost:8080",
		SignInURL:   "http://localhost:8080/auth/signin",

		RedisURI:      "redis://:123456@localhost:6379/0",
		SessionKey:    "your-session-key",
		SessionName:   "myapp_session",
		IsRedisSecure: true,

		SessionMaxAge: 3600, // 1 hour

		// Optional sliding window configuration
		EnableSlidingWindow:       false,
		SessionExtensionDuration:  1800, // 30 minutes
		SessionExtensionThreshold: 1200, // 20 minutes
	}
}
