package handlers

import (
	"encoding/json"
	"net/http"

	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/env"
	"github.com/m-milek/leszmonitor/logger"
	"golang.org/x/crypto/bcrypt"

	"os"
	"strconv"
	"time"
)

type LoginPayload struct {
	jwt2.MapClaims
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		msg := "Failed to decode login request payload"
		logger.Api.Error().Err(err).Msg(msg)
		util.RespondMessage(w, http.StatusBadRequest, msg)
		return
	}

	logger.Api.Debug().
		Str("username", payload.Username).
		Msg("User login attempt")

	user, err := db.GetRawUser(payload.Username)

	if err != nil {
		msg := "User not found"
		logger.Api.Error().Err(err).Str("username", payload.Username).Msg(msg)
		util.RespondMessage(w, http.StatusInternalServerError, msg)
		return
	}

	matches := checkPasswordHash(payload.Password, user.PasswordHash)

	if !matches {
		msg := "Invalid password"
		logger.Api.Error().Str("username", payload.Username).Msg(msg)
		util.RespondMessage(w, http.StatusUnauthorized, msg)
		return
	}

	expiryHours, err := strconv.Atoi(os.Getenv(env.JwtExpiryHours))
	if err != nil {
		logger.Api.Error().Err(err).Msg("Invalid JwtExpiryHours value")
		return
	}
	validFor := time.Duration(expiryHours) * time.Hour
	expiryDate := time.Now().Add(validFor)

	logger.Api.Debug().Msg("Creating JWT token valid for " + validFor.String())

	jwt := jwt2.NewWithClaims(
		jwt2.SigningMethodHS256,
		jwt2.MapClaims{
			"username": payload.Username,
			"exp":      jwt2.NewNumericDate(expiryDate),
		},
	)
	token, err := jwt.SignedString([]byte(os.Getenv(env.JwtSecret)))

	if err != nil {
		msg := "Failed to log in"
		logger.Api.Error().Err(err).Msg(msg)
		util.RespondMessage(w, http.StatusInternalServerError, msg)
		return
	}

	msg := "User logged in successfully"
	logger.Api.Info().
		Str("username", payload.Username).
		Msg(msg)

	util.RespondJSON(w, http.StatusOK, map[string]string{"jwt": token})
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
