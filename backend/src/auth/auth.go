package auth

import (
	"fmt"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/m-milek/leszmonitor/config"
	"github.com/m-milek/leszmonitor/log"
)

// JwtClaims represents the claims stored in a Leszmonitor JWT token.
// It includes standard claims and a custom Username field.
type JwtClaims struct {
	jwt.MapClaims
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

// UserClaims extends standard JWT claims with custom fields.
type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// jwtFromRequest extracts the JWT token from the Authorization header of the HTTP request.
func jwtFromRequest(r *http.Request) (string, error) {
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

func decodeJwtClaims(jwtString string) (JwtClaims, error) {
	claims := JwtClaims{}
	token, err := jwt.ParseWithClaims(jwtString, &claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(config.JwtSecret)), nil
	})

	if err != nil {
		log.Api.Error().Err(err).Msg("Failed to parse JWT token")
		return JwtClaims{}, err
	}

	if !token.Valid {
		log.Api.Warn().Msg("Invalid JWT token")
		return JwtClaims{}, fmt.Errorf("invalid JWT token")
	}

	return claims, nil
}

func ValidateJwt(token string) (*UserClaims, error) {
	jwtSecret := os.Getenv(config.JwtSecret)
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT secret is not configured")
	}

	// Parse and validate the token with custom claims
	claims := &UserClaims{}

	parsedJwt, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Api.Warn().Msgf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Api.Error().Err(err).Msg("Failed to parse JWT token")
		return nil, fmt.Errorf("invalid JWT token: %w", err)
	}

	if !parsedJwt.Valid {
		log.Api.Warn().Msg("Invalid JWT token")
		return nil, fmt.Errorf("invalid JWT token")
	}

	// Verify we have the expected claims type
	userClaims, ok := parsedJwt.Claims.(*UserClaims)
	if !ok {
		log.Api.Warn().Msg("Unexpected JWT claims type")
		return nil, fmt.Errorf("unexpected JWT claims type")
	}

	// Ensure username is present
	if userClaims.Username == "" {
		log.Api.Warn().Msg("JWT token is missing username claim")
		return nil, fmt.Errorf("JWT token is missing username claim")
	}

	return userClaims, nil
}
