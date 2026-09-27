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

const messageParameterKeyIsRequired = "Parameter key is required"

type SetParameterPayload struct {
	Value any `json:"value"`
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

func (c *GlobalParameterAPIController) SetParameterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := r.PathValue("key")
	if key == "" {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, messageParameterKeyIsRequired)
		return
	}

	payload := SetParameterPayload{}
	if ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	if _, ok := auth.ExtractUserOrRespond(ctx, w, r); !ok {
		return
	}

	param, err := c.service.SetParameter(ctx, GlobalParameterKey(key), payload.Value)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, param)
}
