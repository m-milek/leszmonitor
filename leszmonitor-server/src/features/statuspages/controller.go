package statuspages

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/httpx"
)

type StatusPageAPIController struct {
	service IStatusPageService
}

func NewStatusPageAPIController(service IStatusPageService) StatusPageAPIController {
	return StatusPageAPIController{
		service: service,
	}
}

type StatusPagePayload struct {
	Slug               string      `json:"slug"`
	Title              string      `json:"title"`
	Description        string      `json:"description"`
	GroupMonitorsByTag bool        `json:"groupMonitorsByTag"`
	TagIDs             []uuid.UUID `json:"tagIds"`
}

// @Summary	Create a status page
// @Tags		status-pages
// @Accept		json
// @Produce	json
// @Param		body	body		StatusPagePayload	true	"Status page fields"
// @Success	201		{object}	StatusPage
// @Failure	400		{object}	httpx.ErrorResponse
// @Failure	401		{object}	httpx.ErrorResponse
// @Failure	403		{object}	httpx.ErrorResponse
// @Failure	409		{object}	httpx.ErrorResponse
// @Failure	500		{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/status-pages [post]
func (c *StatusPageAPIController) CreateStatusPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload := StatusPagePayload{}
	if ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	page, serviceErr := c.service.CreateStatusPage(ctx, StatusPage{
		Slug:               payload.Slug,
		Title:              payload.Title,
		Description:        payload.Description,
		GroupMonitorsByTag: payload.GroupMonitorsByTag,
		TagIDs:             payload.TagIDs,
	})
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusCreated, page)
}

// @Summary	List status pages
// @Tags		status-pages
// @Produce	json
// @Success	200	{array}		StatusPage
// @Failure	401	{object}	httpx.ErrorResponse
// @Failure	403	{object}	httpx.ErrorResponse
// @Failure	500	{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/status-pages [get]
func (c *StatusPageAPIController) GetAllStatusPagesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pages, err := c.service.GetAllStatusPages(ctx)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, pages)
}

// @Summary	Get a status page
// @Tags		status-pages
// @Produce	json
// @Param		statusPageId	path		string	true	"Status page ID"
// @Success	200				{object}	StatusPage
// @Failure	400				{object}	httpx.ErrorResponse
// @Failure	401				{object}	httpx.ErrorResponse
// @Failure	403				{object}	httpx.ErrorResponse
// @Failure	404				{object}	httpx.ErrorResponse
// @Failure	500				{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/status-pages/{statusPageId} [get]
func (c *StatusPageAPIController) GetStatusPageByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := httpx.PathUUIDOrRespond(ctx, w, r, "statusPageId")
	if !ok {
		return
	}

	page, err := c.service.GetStatusPageByID(ctx, id)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, page)
}

// @Summary	Update a status page
// @Tags		status-pages
// @Accept		json
// @Produce	json
// @Param		statusPageId	path		string				true	"Status page ID"
// @Param		body			body		StatusPagePayload	true	"Status page fields"
// @Success	200				{object}	StatusPage
// @Failure	400				{object}	httpx.ErrorResponse
// @Failure	401				{object}	httpx.ErrorResponse
// @Failure	403				{object}	httpx.ErrorResponse
// @Failure	404				{object}	httpx.ErrorResponse
// @Failure	409				{object}	httpx.ErrorResponse
// @Failure	500				{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/status-pages/{statusPageId} [patch]
func (c *StatusPageAPIController) UpdateStatusPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := httpx.PathUUIDOrRespond(ctx, w, r, "statusPageId")
	if !ok {
		return
	}

	payload := StatusPagePayload{}
	if ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	page, serviceErr := c.service.UpdateStatusPage(ctx, StatusPage{
		ID:                 id,
		Slug:               payload.Slug,
		Title:              payload.Title,
		Description:        payload.Description,
		GroupMonitorsByTag: payload.GroupMonitorsByTag,
		TagIDs:             payload.TagIDs,
	})
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, page)
}

// @Summary	Delete a status page
// @Tags		status-pages
// @Produce	json
// @Param		statusPageId	path	string	true	"Status page ID"
// @Success	204
// @Failure	400	{object}	httpx.ErrorResponse
// @Failure	401	{object}	httpx.ErrorResponse
// @Failure	403	{object}	httpx.ErrorResponse
// @Failure	404	{object}	httpx.ErrorResponse
// @Failure	500	{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/status-pages/{statusPageId} [delete]
func (c *StatusPageAPIController) DeleteStatusPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := httpx.PathUUIDOrRespond(ctx, w, r, "statusPageId")
	if !ok {
		return
	}

	if err := c.service.DeleteStatusPage(ctx, id); err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
