package globalparameters

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/httpx"
)

type GlobalParameterAPIController struct {
	service IGlobalParameterService
}

func NewGlobalParameterAPIController(service IGlobalParameterService) GlobalParameterAPIController {
	return GlobalParameterAPIController{
		service: service,
	}
}

// @Summary	List global parameters
// @Tags		global-parameters
// @Produce	json
// @Success	200	{array}		GlobalParameter
// @Failure	401	{object}	httpx.ErrorResponse
// @Failure	403	{object}	httpx.ErrorResponse
// @Failure	500	{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/global-parameters [get]
func (c *GlobalParameterAPIController) GetAllParametersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params, err := c.service.GetAllParameters(ctx)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, params)
}

// @Summary	Set global parameters
// @Tags		global-parameters
// @Accept		json
// @Produce	json
// @Param		body	body		map[string]interface{}	true	"Parameter key to new value"
// @Success	200		{array}		GlobalParameter
// @Failure	400		{object}	httpx.ErrorResponse
// @Failure	401		{object}	httpx.ErrorResponse
// @Failure	403		{object}	httpx.ErrorResponse
// @Failure	409		{object}	httpx.ErrorResponse
// @Failure	500		{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/global-parameters [patch]
func (c *GlobalParameterAPIController) SetParametersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload := map[GlobalParameterKey]any{}
	if ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	params, err := c.service.SetParameters(ctx, payload)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, params)
}
