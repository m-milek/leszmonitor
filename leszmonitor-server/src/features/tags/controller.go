package tags

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/httpx"
)

type TagAPIController struct {
	service ITagService
}

func NewTagAPIController(service ITagService) TagAPIController {
	return TagAPIController{
		service: service,
	}
}

// TagPayload is the client-writable part of a tag. The ID and timestamps are server-side.
type TagPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ColorHex    string `json:"colorHex"`
}

// CreateTagHandler handles the creation of a new tag.
//
//	@Summary	Create a tag
//	@Tags		tags
//	@Accept		json
//	@Produce	json
//	@Param		body	body		TagPayload	true	"Tag fields"
//	@Success	201		{object}	Tag
//	@Failure	400		{object}	httpx.ErrorResponse
//	@Failure	401		{object}	httpx.ErrorResponse
//	@Failure	403		{object}	httpx.ErrorResponse
//	@Failure	409		{object}	httpx.ErrorResponse
//	@Failure	500		{object}	httpx.ErrorResponse
//	@Security	BearerAuth
//	@Router		/tags [post]
func (c *TagAPIController) CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload := TagPayload{}
	if ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	tag, serviceErr := c.service.CreateTag(ctx, Tag{
		Name:        payload.Name,
		Description: payload.Description,
		ColorHex:    payload.ColorHex,
	})
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusCreated, tag)
}

// @Summary	List tags
// @Tags		tags
// @Produce	json
// @Success	200	{array}		Tag
// @Failure	401	{object}	httpx.ErrorResponse
// @Failure	403	{object}	httpx.ErrorResponse
// @Failure	500	{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/tags [get]
func (c *TagAPIController) GetAllTagsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	allTags, err := c.service.GetAllTags(ctx)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, allTags)
}

// @Summary	Get a tag
// @Tags		tags
// @Produce	json
// @Param		tagId	path		string	true	"Tag ID"
// @Success	200		{object}	Tag
// @Failure	400		{object}	httpx.ErrorResponse
// @Failure	401		{object}	httpx.ErrorResponse
// @Failure	403		{object}	httpx.ErrorResponse
// @Failure	404		{object}	httpx.ErrorResponse
// @Failure	500		{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/tags/{tagId} [get]
func (c *TagAPIController) GetTagByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tagID, ok := httpx.PathUUIDOrRespond(ctx, w, r, "tagId")
	if !ok {
		return
	}

	tag, err := c.service.GetTagByID(ctx, tagID)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, tag)
}

// UpdateTagHandler handles the update of an existing tag.
//
//	@Summary	Update a tag
//	@Tags		tags
//	@Accept		json
//	@Produce	json
//	@Param		tagId	path		string		true	"Tag ID"
//	@Param		body	body		TagPayload	true	"Tag fields"
//	@Success	200		{object}	Tag
//	@Failure	400		{object}	httpx.ErrorResponse
//	@Failure	401		{object}	httpx.ErrorResponse
//	@Failure	403		{object}	httpx.ErrorResponse
//	@Failure	404		{object}	httpx.ErrorResponse
//	@Failure	409		{object}	httpx.ErrorResponse
//	@Failure	500		{object}	httpx.ErrorResponse
//	@Security	BearerAuth
//	@Router		/tags/{tagId} [patch]
func (c *TagAPIController) UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tagUUID, ok := httpx.PathUUIDOrRespond(ctx, w, r, "tagId")
	if !ok {
		return
	}

	payload := TagPayload{}
	if ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	tag, serviceErr := c.service.UpdateTag(ctx, Tag{
		ID:          tagUUID,
		Name:        payload.Name,
		Description: payload.Description,
		ColorHex:    payload.ColorHex,
	})
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, tag)
}

// @Summary	Delete a tag
// @Tags		tags
// @Produce	json
// @Param		tagId	path	string	true	"Tag ID"
// @Success	204
// @Failure	400	{object}	httpx.ErrorResponse
// @Failure	401	{object}	httpx.ErrorResponse
// @Failure	403	{object}	httpx.ErrorResponse
// @Failure	404	{object}	httpx.ErrorResponse
// @Failure	409	{object}	httpx.ErrorResponse
// @Failure	500	{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/tags/{tagId} [delete]
func (c *TagAPIController) DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tagID, ok := httpx.PathUUIDOrRespond(ctx, w, r, "tagId")
	if !ok {
		return
	}

	if err := c.service.DeleteTag(ctx, tagID); err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
