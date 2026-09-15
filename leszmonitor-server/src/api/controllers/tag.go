package controllers

import (
	"net/http"

	"github.com/google/uuid"
	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/services"
)

type TagAPIController struct {
	service services.ITagService
}

func NewTagAPIController(service services.ITagService) TagAPIController {
	return TagAPIController{
		service: service,
	}
}

const messageTagIDIsRequired = "Tag ID is required"

// TagPayload is the client-writable part of a tag. The ID and timestamps are server-side.
type TagPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ColorHex    string `json:"colorHex"`
}

// CreateTagHandler handles the creation of a new tag.
func (c *TagAPIController) CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload := TagPayload{}
	if ok := util.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	if _, ok := authorization.ExtractUserOrRespond(ctx, w, r); !ok {
		return
	}

	tag, serviceErr := c.service.CreateTag(ctx, models.Tag{
		Name:        payload.Name,
		Description: payload.Description,
		ColorHex:    payload.ColorHex,
	})
	if serviceErr != nil {
		util.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusCreated, tag)
}

func (c *TagAPIController) GetAllTagsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	allTags, err := c.service.GetAllTags(ctx)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, allTags)
}

func (c *TagAPIController) GetTagByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tagID := r.PathValue("tagId")
	if tagID == "" {
		util.RespondMessage(ctx, w, http.StatusBadRequest, messageTagIDIsRequired)
		return
	}

	tag, err := c.service.GetTagByID(ctx, tagID)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, tag)
}

// UpdateTagHandler handles the update of an existing tag.
func (c *TagAPIController) UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tagID := r.PathValue("tagId")
	if tagID == "" {
		util.RespondMessage(ctx, w, http.StatusBadRequest, messageTagIDIsRequired)
		return
	}

	tagUUID, parseErr := uuid.Parse(tagID)
	if parseErr != nil {
		util.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid tag ID format")
		return
	}

	payload := TagPayload{}
	if ok := util.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	if _, ok := authorization.ExtractUserOrRespond(ctx, w, r); !ok {
		return
	}

	tag, serviceErr := c.service.UpdateTag(ctx, models.Tag{
		ID:          tagUUID,
		Name:        payload.Name,
		Description: payload.Description,
		ColorHex:    payload.ColorHex,
	})
	if serviceErr != nil {
		util.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, tag)
}

func (c *TagAPIController) DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tagID := r.PathValue("tagId")
	if tagID == "" {
		util.RespondMessage(ctx, w, http.StatusBadRequest, messageTagIDIsRequired)
		return
	}

	if _, ok := authorization.ExtractUserOrRespond(ctx, w, r); !ok {
		return
	}

	if err := c.service.DeleteTag(ctx, tagID); err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondMessage(ctx, w, http.StatusOK, "Tag deleted successfully")
}
