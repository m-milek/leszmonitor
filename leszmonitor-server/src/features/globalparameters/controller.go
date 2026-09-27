package globalparameters

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
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

func (c *GlobalParameterAPIController) GetAllParametersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params, err := c.service.GetAllParameters(ctx)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, params)
}

func (c *GlobalParameterAPIController) SetParametersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload := map[GlobalParameterKey]any{}
	if ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	if _, ok := auth.ExtractUserOrRespond(ctx, w, r); !ok {
		return
	}

	params, err := c.service.SetParameters(ctx, payload)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, params)
}
