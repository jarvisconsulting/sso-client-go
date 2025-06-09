package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jarvisconsulting/sso-client-go/pkg/config"
	"github.com/jarvisconsulting/sso-client-go/pkg/models"
	"github.com/jarvisconsulting/sso-client-go/pkg/store"
)

const (
	SessionUserIDKey = "session_user_id" // Key used to store the user ID in the session
	SessionIsMobileKey = "is_mobile"
)

type UserRepository interface {
	FindByID(id uint) (*models.User, error)
	FindByJTI(jti string) (uint, error)
	GetLastSshKey() (*models.SshKey, error)
}

type AuthService struct {
	userRepo     UserRepository
	config       *config.Config
	sessionStore store.SessionStore
}

func NewAuthService(userRepo UserRepository, cfg *config.Config, sessionStore store.SessionStore) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		config:       cfg,
		sessionStore: sessionStore,
	}
}

func (s *AuthService) IsUserSignedIn(r *http.Request) bool {
	session, err := s.sessionStore.GetStore().Get(r, s.config.SessionName)
	if err != nil {
		return false
	}

	_, ok := session.Values[SessionUserIDKey]
	return ok
}

func (s *AuthService) SignInUser(w http.ResponseWriter, r *http.Request, userID uint, isMobile bool) error {
	session, err := s.sessionStore.GetStore().Get(r, s.config.SessionName)
	if err != nil {
		return err
	}

	session.Values[SessionUserIDKey] = userID
	session.Values[SessionIsMobileKey] = isMobile
	return session.Save(r, w)
}

func (s *AuthService) SignOutUser(w http.ResponseWriter, r *http.Request) error {
	session, err := s.sessionStore.GetStore().Get(r, s.config.SessionName)
	if err != nil {
		return err
	}

	delete(session.Values, SessionUserIDKey)
	return session.Save(r, w)
}

// func (s *AuthService) HandleCallback(params map[string]string) (uint, error) {
// 	idToken, ok := params["id_token"]
// 	if !ok || idToken == "" {
// 		return 0, errors.New("id_token not provided")
// 	}

// 	sshKey, err := s.userRepo.GetLastSshKey()
// 	if err != nil {
// 		return 0, fmt.Errorf("error getting SSH key: %w", err)
// 	}

// 	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(sshKey.PrivateRsaKey))
// 	if err != nil {
// 		return 0, fmt.Errorf("error parsing private key: %w", err)
// 	}

// 	token, err := jwt.Parse(idToken, func(token *jwt.Token) (any, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return &privateKey.PublicKey, nil
// 	})

// 	var jti string
// 	var extractErr error

// 	if err != nil {
// 		parts := strings.Split(idToken, ".")
// 		if len(parts) >= 2 {
// 			payload := parts[1]
// 			if len(payload)%4 != 0 {
// 				payload += strings.Repeat("=", 4-len(payload)%4)
// 			}

// 			decoded, err := base64.URLEncoding.DecodeString(payload)
// 			if err == nil {
// 				payloadStr := string(decoded)

// 				if len(payloadStr) >= 2 && payloadStr[0] == '"' && payloadStr[len(payloadStr)-1] == '"' {
// 					jti = payloadStr[1 : len(payloadStr)-1]
// 				}
// 			}
// 		}

// 		if jti == "" {
// 			extractErr = fmt.Errorf("couldn't extract JTI from malformed token: %w", err)
// 		}
// 	} else if !token.Valid {
// 		extractErr = errors.New("invalid token")
// 	} else {
// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if !ok {
// 			extractErr = errors.New("invalid claims format")
// 		} else {
// 			jtiClaim, ok := claims["jti"]
// 			if !ok {
// 				extractErr = errors.New("jti not found in token claims")
// 			} else {
// 				jti, ok = jtiClaim.(string)
// 				if !ok {
// 					extractErr = errors.New("jti is not a string")
// 				}
// 			}
// 		}
// 	}

// 	if extractErr != nil && jti == "" {
// 		return 0, extractErr
// 	}

// 	userID, err := s.userRepo.FindByJTI(jti)
// 	if err != nil {
// 		return 0, fmt.Errorf("error finding user by JTI: %w", err)
// 	}

// 	return userID, nil
// }

func (s *AuthService) HandleCallback(params map[string]string) (uint, error) {
	idToken, ok := params["id_token"]
	if !ok || idToken == "" {
		return 0, errors.New("id_token not provided")
	}

	sshKey, err := s.userRepo.GetLastSshKey()
	if err != nil {
		return 0, fmt.Errorf("error getting SSH key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(sshKey.PrivateRsaKey))
	if err != nil {
		return 0, fmt.Errorf("error parsing private key: %w", err)
	}

	// Parse and validate the token
	token, err := jwt.Parse(idToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return &privateKey.PublicKey, nil
	})

	// If token is invalid, return immediately
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	fmt.Printf("Token claims: %+v\n", token.Claims)

	var jti string

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims format")
	}

	if jtiClaim, exists := claims["jti"]; exists {
		if jtiStr, ok := jtiClaim.(string); ok {
			jti = jtiStr
			fmt.Printf("Found JTI in jti claim: %s\n", jti)
		}
	} else if len(claims) == 1 {
		for k, v := range claims {
			fmt.Printf("Checking claim key: %s, value: %v\n", k, v)
			if strVal, ok := v.(string); ok {
				jti = strVal
				fmt.Printf("Found JTI as direct string value: %s\n", jti)
				break
			}
		}
	}

	if jti == "" {
		return 0, errors.New("could not extract JTI from token - token must contain either a jti claim or a single string value")
	}

	userID, err := s.userRepo.FindByJTI(jti)
	if err != nil {
		return 0, fmt.Errorf("error finding user by JTI: %w", err)
	}

	return userID, nil
}

func (s *AuthService) GetUserIDFromSession(r *http.Request) (uint, error) {
	session, err := s.sessionStore.GetStore().Get(r, s.config.SessionName)
	if err != nil {
		return 0, err
	}

	userID, ok := session.Values[SessionUserIDKey].(uint)
	if !ok {
		return 0, errors.New("user ID not found in session")
	}

	return userID, nil
}

func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *AuthService) IsUserMobile(r *http.Request) (bool, error) {
	session, err := s.sessionStore.GetStore().Get(r, s.config.SessionName)
	if err != nil {
		return false, err
	}

	isMobile, ok := session.Values[SessionIsMobileKey].(bool)
	if !ok {
		return false, nil
	}

	return isMobile, nil
}
