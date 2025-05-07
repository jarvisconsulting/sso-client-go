# SSO Go Client Library

A robust Go client library for Single Sign-On (SSO) integration with sliding window session management, Redis-based session storage, and database failover support.

## Features

- 🔐 Secure SSO integration with JWT validation
- 🔄 Sliding window session management
- 📦 Redis-based session storage
- 💾 Database failover support (Primary/Secondary)
- 🚀 Easy integration with Gin web framework



## Installation

```bash
go get github.com/yourusername/sso-go-client
```

## Required Dependencies

```go
require (
    github.com/gin-gonic/gin v1.9.x
    github.com/golang-jwt/jwt/v5 v5.x.x
    github.com/gomodule/redigo v1.8.x
    github.com/gorilla/sessions v1.2.x
    gorm.io/gorm v1.25.x
    gorm.io/driver/postgres v1.5.x
)
```

## Basic Usage

1. Initialize the client:

```go
import (
    ssoclient "github.com/yourusername/sso-go-client"
    "github.com/yourusername/sso-go-client/pkg/config"
)

func main() {
    // Configure the client
    cfg := &config.Config{
        // SSO Server Configuration
        SSOUrl:       "https://sso.example.com",
        ClientID:     "your-client-id",
        CallbackURL:  "http://localhost:8080/auth/callback",
        SignOutURL:   "http://localhost:8080/auth/signout",
        PublicKeyAPI: "https://sso.example.com/api/public-key",

        // Application URLs
        RootURL:     "http://localhost:8080",
        SignInURL:   "http://localhost:8080/auth/signin",
        FrontendURL: "http://localhost:3000",

        // Redis Configuration
        RedisURI:      "redis://:password@localhost:6379/0",
        SessionKey:    "your-session-key",
        SessionName:   "myapp_session",
        IsRedisSecure: true,

        // Session Configuration
        SessionMaxAge:             3600, // 1 hour
        EnableSlidingWindow:       true,
        SessionExtensionDuration:  1800, // 30 minutes
        SessionExtensionThreshold: 1200, // 20 minutes
    }

    // Initialize databases
    primaryDB, err := gorm.Open(postgres.Open("postgres://user:pass@primary:5432/db"))
    if err != nil {
        log.Fatal(err)
    }

    secondaryDB, err := gorm.Open(postgres.Open("postgres://user:pass@secondary:5432/db"))
    if err != nil {
        log.Printf("Warning: Secondary DB not available: %v", err)
    }

    // Create client
    client, err := ssoclient.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Initialize with repositories
    client = client.WithRepository(primaryDB, secondaryDB)

    // Get middleware and handlers
    middleware := client.GetMiddleware()
    handlers := client.GetHandlers()
}
```

2. Set up routes with Gin:

```go
func setupRouter(middleware *ssoclient.Middleware, handlers *ssoclient.Handlers) *gin.Engine {
    router := gin.Default()

    // Apply session middleware globally
    router.Use(middleware.Session)

    // Public routes
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })

    // Auth routes
    auth := router.Group("/auth")
    {
        auth.GET("/signin", handlers.SignIn)
        auth.GET("/callback", handlers.Callback)
        auth.POST("/signout", handlers.SignOut)
    }

    // Protected routes
    api := router.Group("/api")
    api.Use(middleware.RequireAuth)
    api.Use(middleware.SetUserID)
    {
        api.GET("/user", handlers.User)
    }

    return router
}
```

## Database Schema

The library requires the following tables:

```sql
-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- SSH Keys table
CREATE TABLE ssh_keys (
    id SERIAL PRIMARY KEY,
    private_key TEXT NOT NULL,
    public_key TEXT NOT NULL,
    status VARCHAR(20) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User Access Tokens table
CREATE TABLE user_access_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    jti VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Session Management

The library implements a sliding window session mechanism:

1. Initial session duration is set by `SessionMaxAge`
2. When a request is made within `SessionExtensionThreshold` of expiry:
   - Session is extended by `SessionExtensionDuration`
   - New expiry time is saved to Redis

Example timeline:
- T+0: Session created (expires in 1 hour)
- T+40min: Request made (within 20min threshold)
- T+40min: Session extended by 30min (new expiry 1h10min from now)

## Database Failover

The library implements automatic failover between primary and secondary databases:

1. All queries first attempt the primary database
2. On primary database failure, automatically retry with secondary
3. Both databases must have the same schema and be in sync
4. Writes are attempted on both databases when available

## Error Handling

The library provides detailed error types for different failure scenarios:

```go
// Check specific error types
if errors.Is(err, ssoclient.ErrInvalidToken) {
    // Handle invalid token
}

if errors.Is(err, ssoclient.ErrDatabaseFailover) {
    // Handle database failover
}
```

## Monitoring

The library exposes metrics for monitoring:

- Session creation/deletion rates
- Database failover events
- Authentication failures

## Security Considerations

1. Always use HTTPS in production
2. Secure Redis with authentication and encryption
3. Regularly rotate SSH keys
4. Monitor for suspicious authentication patterns
5. Keep all dependencies updated

## Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 