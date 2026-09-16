package controllers

import (
	"fmt"
	"net/http"

	"github.com/m-milek/leszmonitor/platform/httpx"
	"github.com/m-milek/leszmonitor/services"
)

type UserAPIController struct {
	service services.IUserService
}

func NewUserAPIController(service services.IUserService) UserAPIController {
	return UserAPIController{
		service: service,
	}
}

func (c *UserAPIController) UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload services.UserRegisterPayload
	if !httpx.DecodeJSONOrRespond(ctx, w, r, &payload) {
		return
	}

	err := c.service.RegisterUser(ctx, &payload)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondMessage(ctx, w, http.StatusOK, "")
}

func (c *UserAPIController) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload services.LoginPayload
	if !httpx.DecodeJSONOrRespond(ctx, w, r, &payload) {
		return
	}

	loginResponse, err := c.service.Login(ctx, payload)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, loginResponse)
}

func (c *UserAPIController) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := c.service.GetAllUsers(ctx)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, users)
}

func (c *UserAPIController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := r.PathValue("username")
	if username == "" {
		httpx.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("username is required"))
		return
	}

	user, err := c.service.GetUserByUsername(ctx, username)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, user)
}

func (c *UserAPIController) SetUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := r.PathValue("username")
	if username == "" {
		httpx.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("username is required"))
		return
	}

	var payload services.SetUserRolePayload
	if !httpx.DecodeJSONOrRespond(ctx, w, r, &payload) {
		return
	}

	user, err := c.service.SetUserRole(ctx, username, payload)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, user)
}
