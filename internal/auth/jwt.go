package auth

import (
	"errors"
	"fmt"
	appcfg "gpu-scheduler-platform/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	issuer string
	secret []byte
	expire time.Duration
}

type JWTClaims struct {
	Name        string   `json:"name,omitempty"`
	TenantID    string   `json:"tenant_id,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	jwt.RegisteredClaims
}

func NewJWTManager(cfg appcfg.JWTConfig) *JWTManager {
	expire := cfg.Expire
	if expire <= 0 {
		expire = 24 * time.Hour
	}
	return &JWTManager{
		issuer: cfg.Issuer,
		secret: []byte(cfg.Secret),
		expire: expire,
	}
}

func (m *JWTManager) ParseToken(tokenStr string) (*Subject, error) {
	if m == nil || len(m.secret) == 0 {
		return nil, errors.New("jwt manager is not configured")
	}
	if tokenStr == "" {
		return nil, errors.New("empty token")
	}

	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	}, jwt.WithIssuer(m.issuer))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &Subject{
		SubjectID:   claims.Subject,
		Name:        claims.Name,
		TenantID:    claims.TenantID,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
		TokenID:     claims.ID,
		Issuer:      claims.Issuer,
	}, nil
}

func (m *JWTManager) IssueToken(sub *Subject) (string, error) {
	if m == nil || len(m.secret) == 0 {
		return "", errors.New("jwt manager is not configured")
	}
	if sub == nil || sub.SubjectID == "" {
		return "", errors.New("invalid subject")
	}

	now := time.Now().UTC()
	claims := JWTClaims{
		Name:        sub.Name,
		TenantID:    sub.TenantID,
		Roles:       sub.Roles,
		Permissions: sub.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   sub.SubjectID,
			ID:        sub.TokenID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}
