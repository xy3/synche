package auth

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Service wraps the signing key and the issuer
type Service struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type AccessTokenClaims struct {
	Email     string
	TokenType string
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	Email     string
	TokenType string
	CustomKey string
	jwt.StandardClaims
}

// GenerateAccessToken generates a jwt token
func (auth *Service) GenerateAccessToken(email string) (signedToken string, err error) {
	claims := &AccessTokenClaims{
		Email:     email,
		TokenType: "access",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(auth.ExpirationHours)).Unix(),
			Issuer:    auth.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(auth.SecretKey))
}

// GenerateRefreshToken generate a new refresh token for the given user
func (auth *Service) GenerateRefreshToken(email, tokenHash string) (string, error) {
	customKey := auth.GenerateCustomKey(email, tokenHash)
	tokenType := "refresh"

	claims := RefreshTokenClaims{
		Email:     email,
		CustomKey: customKey,
		TokenType: tokenType,
		StandardClaims: jwt.StandardClaims{
			Issuer: auth.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(auth.SecretKey))
}

func (auth *Service) GenerateCustomKey(email, tokenHash string) string {
	hash := md5.Sum([]byte(email + tokenHash + auth.SecretKey))
	return hex.EncodeToString(hash[:])
}

// ValidateAccessToken validates the jwt token
func (auth *Service) ValidateAccessToken(signedToken string) (claims *AccessTokenClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&AccessTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(auth.SecretKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid || claims.Email == "" || claims.TokenType != "access" {
		return nil, errors.New("invalid token, authentication failed")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("JWT has expired")
	}
	return
}

// ValidateRefreshToken validates the jwt token
func (auth *Service) ValidateRefreshToken(signedToken string) (claims *RefreshTokenClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&RefreshTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(auth.SecretKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid || claims.Email == "" || claims.TokenType != "refresh" {
		return nil, errors.New("invalid token, authentication failed")
	}
	return
}
