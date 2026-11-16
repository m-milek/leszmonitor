package statuspages

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/apperr"
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
	RequireAuth        bool        `json:"requireAuth"`
	Enabled            bool        `json:"enabled"`
	TimeWindowSeconds  int         `json:"timeWindowSeconds"`
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
		RequireAuth:        payload.RequireAuth,
		Enabled:            payload.Enabled,
		TimeWindowSeconds:  payload.TimeWindowSeconds,
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

// @Summary	Get a public status page
// @Tags		status-pages
// @Produce	json
// @Param		slug	path		string	true	"Status page slug"
// @Success	200		{object}	StatusPage
// @Failure	401		{object}	httpx.ErrorResponse
// @Failure	403		{object}	httpx.ErrorResponse
// @Failure	404		{object}	httpx.ErrorResponse
// @Failure	500		{object}	httpx.ErrorResponse
// @Router		/public/status-pages/{slug} [get]
func (c *StatusPageAPIController) GetPublicStatusPageHandler(
	requireReader func(http.HandlerFunc) http.HandlerFunc,
) http.HandlerFunc {
	return c.withPublicStatusPage(requireReader, func(w http.ResponseWriter, r *http.Request, page *StatusPage) {
		httpx.RespondJSON(r.Context(), w, http.StatusOK, page)
	})
}

// @Summary	List monitors on a public status page
// @Tags		status-pages
// @Produce	json
// @Param		slug	path		string	true	"Status page slug"
// @Success	200		{array}		PublicMonitor
// @Failure	401		{object}	httpx.ErrorResponse
// @Failure	403		{object}	httpx.ErrorResponse
// @Failure	404		{object}	httpx.ErrorResponse
// @Failure	500		{object}	httpx.ErrorResponse
// @Router		/public/status-pages/{slug}/monitors [get]
func (c *StatusPageAPIController) GetPublicMonitorsHandler(
	requireReader func(http.HandlerFunc) http.HandlerFunc,
) http.HandlerFunc {
	return c.withPublicStatusPage(requireReader, func(w http.ResponseWriter, r *http.Request, page *StatusPage) {
		ctx := r.Context()

		to := time.Now().UTC()
		from := to.Add(-time.Duration(page.TimeWindowSeconds) * time.Second)

		publicMonitors, err := c.service.GetPublicMonitors(ctx, page.ID, from, to)
		if err != nil {
			httpx.RespondError(ctx, w, err.Code, err.Err)
			return
		}

		httpx.RespondJSON(ctx, w, http.StatusOK, publicMonitors)
	})
}

func (c *StatusPageAPIController) withPublicStatusPage(
	requireReader func(http.HandlerFunc) http.HandlerFunc,
	next func(w http.ResponseWriter, r *http.Request, page *StatusPage),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		page, err := c.service.GetStatusPageBySlug(ctx, r.PathValue("slug"))
		if err != nil {
			httpx.RespondError(ctx, w, err.Code, err.Err)
			return
		}

		if !page.Enabled {
			notFound := apperr.NewNotFoundError("status page with slug %s not found", page.Slug)
			httpx.RespondError(ctx, w, notFound.Code, notFound.Err)
			return
		}

		handle := func(w http.ResponseWriter, r *http.Request) {
			next(w, r, page)
		}

		if page.RequireAuth {
			requireReader(handle)(w, r)
			return
		}

		handle(w, r)
	}
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
		RequireAuth:        payload.RequireAuth,
		Enabled:            payload.Enabled,
		TimeWindowSeconds:  payload.TimeWindowSeconds,
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

// @Summary	Export a status page as YAML config
// @Tags		status-pages
// @Produce	plain
// @Param		statusPageId	path		string	true	"Status page ID"
// @Success	200				{string}	string	"Status page config in YAML"
// @Failure	400				{object}	httpx.ErrorResponse
// @Failure	401				{object}	httpx.ErrorResponse
// @Failure	403				{object}	httpx.ErrorResponse
// @Failure	404				{object}	httpx.ErrorResponse
// @Failure	500				{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/status-pages/{statusPageId}/export [get]
func (c *StatusPageAPIController) ExportToYAMLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := httpx.PathUUIDOrRespond(ctx, w, r, "statusPageId")
	if !ok {
		return
	}

	if svcErr := c.service.ExportToYAML(ctx, id, w); svcErr != nil {
		httpx.RespondError(ctx, w, svcErr.Code, svcErr.Err)
		return
	}
}
