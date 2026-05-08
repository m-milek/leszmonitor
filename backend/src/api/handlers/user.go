package handlers

import (
	"fmt"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/services"
)

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload services.UserRegisterPayload
	if !util.DecodeJSONOrRespond(ctx, w, r, &payload) {
		return
	}

	err := services.UserService.RegisterUser(ctx, &payload)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondMessage(ctx, w, http.StatusOK, "")
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload services.LoginPayload
	if !util.DecodeJSONOrRespond(ctx, w, r, &payload) {
		return
	}

	loginResponse, err := services.UserService.Login(ctx, payload)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, loginResponse)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := services.UserService.GetAllUsers(ctx)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, users)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := r.PathValue("username")
	if username == "" {
		util.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("username is required"))
		return
	}

	user, err := services.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, user)
}
