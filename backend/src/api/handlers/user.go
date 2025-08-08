package handlers

import (
	"encoding/json"
	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/common"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logger"
	"net/http"
)

type UserRegisterPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var payload UserRegisterPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	logger.Api.Debug().
		Str("username", payload.Username).
		Str("email", payload.Email).
		Msg("RawUser registration")

	hashedPassword, err := hashPassword(payload.Password)
	if err != nil {
		msg := "Failed to hash password"
		logger.Api.Error().Err(err).Str("username", payload.Username).Msg(msg)
		util.RespondMessage(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	user := common.NewUser(payload.Username, hashedPassword, payload.Email)

	_, err = db.AddUser(user)

	if err != nil {
		msg := "Failed to register user"
		logger.Api.Error().Err(err).Str("username", payload.Username).Msg(msg)
		util.RespondMessage(w, http.StatusInternalServerError, msg)
		return
	}

	msg := "User registered successfully"
	logger.Api.Info().
		Str("username", payload.Username).
		Msg(msg)

	util.RespondMessage(w, http.StatusOK, msg)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetAllUsers()
	if err != nil {
		logger.Api.Error().Err(err).Msg("Failed to retrieve users")
		util.RespondMessage(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	util.RespondJSON(w, http.StatusOK, users)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		logger.Api.Trace().Msg("Username is required")
		util.RespondMessage(w, http.StatusBadRequest, "Username is required")
		return
	}

	user, err := db.GetUser(username)
	if err != nil {
		if err == db.ErrNotFound {
			msg := "User not found"
			logger.Api.Warn().Str("username", username).Msg(msg)
			util.RespondMessage(w, http.StatusNotFound, msg)
			return
		}
		logger.Api.Error().Err(err).Str("username", username).Msg("Failed to retrieve user")
		util.RespondMessage(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	util.RespondJSON(w, http.StatusOK, user)
}
