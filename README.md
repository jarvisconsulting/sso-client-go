# SSO Client Example

This example demonstrates how to use the SSO client package with sliding window session management.

## Prerequisites

- Go 1.16 or later
- PostgreSQL database
- Redis server
- SSO server (configured and running)

## Dependencies

```bash
go get -u github.com/gin-gonic/gin
go get -u github.com/lib/pq
```

## Configuration

1. Update the database connection string in `main.go`:
```go
db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/myapp?sslmode=disable")
```

2. Configure your SSO settings in the `config.Config` struct:
```go
cfg := &config.Config{
    SSOUrl:       "https://sso.example.com",
    ClientID:     "your-client-id",
    // ... other settings
}
```

3. Configure Redis for session storage:
```go
RedisURI:     "redis://:password@localhost:6379/0",
SessionKey:   "your-session-key",
```

## Session Management

The example implements sliding window session management with the following configuration:

- Base session duration: 1 hour
- Extension duration: 30 minutes
- Extension threshold: 20 minutes

When a user makes a request and their session is within 20 minutes of expiring, it will automatically be extended by 30 minutes.

## Running the Example

1. Start your Redis server
2. Start your PostgreSQL server
3. Run the example:
```bash
go run main.go
```

The server will start on port 8080 with the following endpoints:

- `GET /health` - Public health check endpoint
- `GET /auth/signin` - Initiates SSO login
- `GET /auth/callback` - SSO callback endpoint
- `POST /auth/signout` - Handles sign out
- `GET /api/user` - Protected endpoint, returns user details
- `GET /api/profile` - Protected endpoint example

## Session Flow

1. User accesses a protected endpoint
2. If not authenticated, redirected to SSO login
3. After successful login, user gets a session
4. Session automatically extends if user remains active
5. Session expires after inactivity period

## Frontend Integration

The frontend should handle:

1. Redirecting to `/auth/signin` for login
2. Handling 401 responses by redirecting to login
3. Making authenticated requests to protected endpoints

Example frontend error handling:
```javascript
// Axios interceptor example
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      window.location.href = '/auth/signin';
    }
    return Promise.reject(error);
  }
);
``` 