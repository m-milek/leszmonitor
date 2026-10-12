package users

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/httpx"
)

type UserAPIController struct {
	service IUserService
}

func NewUserAPIController(service IUserService) UserAPIController {
	return UserAPIController{
		service: service,
	}
}

// @Summary	Register a user
// @Tags		auth
// @Accept		json
// @Produce	json
// @Param		body	body	UserRegisterPayload	true	"Registration data"
// @Success	204
// @Failure	400	{object}	httpx.ErrorResponse
// @Failure	401	{object}	httpx.ErrorResponse
// @Failure	500	{object}	httpx.ErrorResponse
// @Router		/auth/register [post]
func (c *UserAPIController) UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload UserRegisterPayload
	if !httpx.DecodeJSONOrRespond(ctx, w, r, &payload) {
		return
	}

	err := c.service.RegisterUser(ctx, &payload)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary	Log in
// @Tags		auth
// @Accept		json
// @Produce	json
// @Param		body	body		LoginPayload	true	"Credentials"
// @Success	200		{object}	LoginResponse
// @Failure	400		{object}	httpx.ErrorResponse
// @Failure	401		{object}	httpx.ErrorResponse
// @Failure	500		{object}	httpx.ErrorResponse
// @Router		/auth/login [post]
func (c *UserAPIController) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload LoginPayload
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

// @Summary	List users
// @Tags		users
// @Produce	json
// @Success	200	{array}		User
// @Failure	401	{object}	httpx.ErrorResponse
// @Failure	500	{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/users [get]
func (c *UserAPIController) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := c.service.GetAllUsers(ctx)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, users)
}

// @Summary	Get a user
// @Tags		users
// @Produce	json
// @Param		username	path		string	true	"Username"
// @Success	200			{object}	User
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	404			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/users/{username} [get]
func (c *UserAPIController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := r.PathValue("username")

	user, err := c.service.GetUserByUsername(ctx, username)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, user)
}

// @Summary	Set a user's role
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		username	path		string				true	"Username"
// @Param		body		body		SetUserRolePayload	true	"New role"
// @Success	200			{object}	User
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	404			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/users/{username}/role [patch]
func (c *UserAPIController) SetUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := r.PathValue("username")

	var payload SetUserRolePayload
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
