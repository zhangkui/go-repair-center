package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID      int64    `json:"user_id"`
	Username    string   `json:"username"`
	Permissions []string `json:"permissions,omitempty"`
	TokenType   string   `json:"token_type"`
	jwtv5.RegisteredClaims
}

type Manager struct {
	secret      []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func NewManager(secret string, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (m *Manager) Issue(userID int64, username string, permissions []string, tokenType string) (string, time.Time, error) {
	ttl := m.accessTTL
	if tokenType == "refresh" {
		ttl = m.refreshTTL
	}
	now := time.Now()
	expires := now.Add(ttl)
	claims := Claims{
		UserID:      userID,
		Username:    username,
		Permissions: permissions,
		TokenType:   tokenType,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   username,
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(expires),
		},
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	return signed, expires, err
}

func (m *Manager) Parse(raw string, expectedType string) (*Claims, error) {
	token, err := jwtv5.ParseWithClaims(raw, &Claims{}, func(token *jwtv5.Token) (any, error) {
		if token.Method != jwtv5.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.TokenType != expectedType {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}