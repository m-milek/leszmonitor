package tags

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/auth"
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
	if ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	if _, ok := auth.ExtractUserOrRespond(ctx, w, r); !ok {
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

func (c *TagAPIController) GetAllTagsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	allTags, err := c.service.GetAllTags(ctx)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, allTags)
}

func (c *TagAPIController) GetTagByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tagID := r.PathValue("tagId")
	if tagID == "" {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, messageTagIDIsRequired)
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
func (c *TagAPIController) UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tagID := r.PathValue("tagId")
	if tagID == "" {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, messageTagIDIsRequired)
		return
	}

	tagUUID, parseErr := uuid.Parse(tagID)
	if parseErr != nil {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid tag ID format")
		return
	}

	payload := TagPayload{}
	if ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	if _, ok := auth.ExtractUserOrRespond(ctx, w, r); !ok {
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

func (c *TagAPIController) DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tagID := r.PathValue("tagId")
	if tagID == "" {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, messageTagIDIsRequired)
		return
	}

	if _, ok := auth.ExtractUserOrRespond(ctx, w, r); !ok {
		return
	}

	if err := c.service.DeleteTag(ctx, tagID); err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondMessage(ctx, w, http.StatusOK, "Tag deleted successfully")
}
