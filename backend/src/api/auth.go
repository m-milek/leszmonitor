package api

import (
	"fmt"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/m-milek/leszmonitor/env"
	"github.com/m-milek/leszmonitor/logger"
	"net/http"
	"os"
)

type JwtClaims struct {
	jwt2.MapClaims
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

func JwtFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		return "", nil
	}

	return authHeader[len(prefix):], nil
}

func DecodeJwtClaims(jwtString string) (JwtClaims, error) {
	claims := JwtClaims{}
	token, err := jwt2.ParseWithClaims(jwtString, &claims, func(token *jwt2.Token) (interface{}, error) {
		return []byte(os.Getenv(env.JwtSecret)), nil
	})

	if err != nil {
		logger.Api.Error().Err(err).Msg("Failed to parse JWT token")
		return JwtClaims{}, err
	}

	if !token.Valid {
		logger.Api.Warn().Msg("Invalid JWT token")
		return JwtClaims{}, fmt.Errorf("invalid JWT token")
	}

	return claims, nil
}
