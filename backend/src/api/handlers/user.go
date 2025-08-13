package handlers

import (
	"encoding/json"
	"fmt"
	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/services"
	"github.com/m-milek/leszmonitor/logger"
	"net/http"
)

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var payload services.UserRegisterPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := services.UserService.RegisterUser(r.Context(), &payload)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "")
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload services.LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		msg := "Failed to decode login request payload"
		util.RespondMessage(w, http.StatusBadRequest, msg)
		return
	}

	loginResponse, err := services.UserService.Login(r.Context(), payload)

	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, loginResponse)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := services.UserService.GetAllUsers(r.Context())

	if err != nil {
		util.RespondError(w, err.Code, err.Err)
	}

	util.RespondJSON(w, http.StatusOK, users)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		logger.Api.Trace().Msg("Username is required")
		util.RespondError(w, http.StatusBadRequest, fmt.Errorf("username is required"))
		return
	}

	user, err := services.UserService.GetUserByUsername(r.Context(), username)

	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, user)
}
