package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
)

type TokenClaims struct {
	TokenId  string     `json:"token_id"`
	UserRole model.Role `json:"role"`
	jwt.RegisteredClaims
}

type TokenDTO struct {
	UserID   string
	UserRole model.Role
	Expires  time.Duration
}

var ErrInvalidTokenPayload = errors.New("Invalid token payload")

type TokenManager interface {
	GenerateToken(dto TokenDTO) (string, string, error)
	ValidateToken(token string) (*TokenClaims, error)
}

type tokenManager struct {
	secretKey string
}

func (t *tokenManager) GenerateToken(dto TokenDTO) (string, string, error) {
	tokenId := uuid.New().String()

	claims := &TokenClaims{
		TokenId:  tokenId,
		UserRole: dto.UserRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(dto.Expires)),
			Subject:   dto.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(t.secretKey))

	return signedToken, tokenId, err
}

func (t *tokenManager) ValidateToken(token string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(tkn *jwt.Token) (any, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tkn.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, ErrInvalidTokenPayload
	}

	return claims, nil
}

func NewTokenManager(secretKey string) TokenManager {
	return &tokenManager{
		secretKey: secretKey,
	}
}
