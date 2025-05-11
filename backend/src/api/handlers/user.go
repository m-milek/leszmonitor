package handlers

import (
	"encoding/json"
	"github.com/m-milek/leszmonitor/api/util"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/model"
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
		Msg("User registration")

	user, err := model.NewUser(payload.Username, payload.Password, payload.Email)
	if err != nil {
		msg := "Failed to create user"
		logger.Api.Error().Err(err).Str("username", payload.Username).Msg(msg)
		util.RespondMessage(w, http.StatusInternalServerError, msg)
		return
	}

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
