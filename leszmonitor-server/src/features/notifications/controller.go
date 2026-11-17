package notifications

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/httpx"
)

type NotificationChannelAPIController struct {
	service INotificationChannelService
}

func NewNotificationChannelAPIController(service INotificationChannelService) NotificationChannelAPIController {
	return NotificationChannelAPIController{
		service: service,
	}
}

// NotificationChannelPayload is the client-writable part of a notification channel.
// Secret config fields left empty on update keep their stored value.
type NotificationChannelPayload struct {
	Name    string                  `json:"name"`
	Kind    NotificationChannelKind `json:"kind"`
	Enabled bool                    `json:"enabled"`
	Config  ChannelConfig           `json:"config" swaggertype:"object"`
}

// @Summary	Create a notification channel
// @Tags		notification-channels
// @Accept		json
// @Produce	json
// @Param		body	body		NotificationChannelPayload	true	"Notification channel fields"
// @Success	201		{object}	NotificationChannel
// @Failure	400		{object}	httpx.ErrorResponse
// @Failure	401		{object}	httpx.ErrorResponse
// @Failure	403		{object}	httpx.ErrorResponse
// @Failure	409		{object}	httpx.ErrorResponse
// @Failure	500		{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/notification-channels [post]
func (c *NotificationChannelAPIController) CreateChannelHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload := NotificationChannelPayload{}
	if ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	channel, serviceErr := c.service.CreateChannel(ctx, NotificationChannel{
		Name:    payload.Name,
		Kind:    payload.Kind,
		Enabled: payload.Enabled,
		Config:  payload.Config,
	})
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusCreated, channel)
}

// @Summary	List notification channels
// @Tags		notification-channels
// @Produce	json
// @Success	200	{array}		NotificationChannel
// @Failure	401	{object}	httpx.ErrorResponse
// @Failure	403	{object}	httpx.ErrorResponse
// @Failure	500	{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/notification-channels [get]
func (c *NotificationChannelAPIController) GetAllChannelsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	channels, err := c.service.GetAllChannels(ctx)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, channels)
}

// @Summary	Get a notification channel
// @Tags		notification-channels
// @Produce	json
// @Param		channelId	path		string	true	"Notification channel ID"
// @Success	200			{object}	NotificationChannel
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	404			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/notification-channels/{channelId} [get]
func (c *NotificationChannelAPIController) GetChannelByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	channelID, ok := httpx.PathUUIDOrRespond(ctx, w, r, "channelId")
	if !ok {
		return
	}

	channel, err := c.service.GetChannelByID(ctx, channelID)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, channel)
}

// @Summary	Update a notification channel
// @Tags		notification-channels
// @Accept		json
// @Produce	json
// @Param		channelId	path		string						true	"Notification channel ID"
// @Param		body		body		NotificationChannelPayload	true	"Notification channel fields"
// @Success	200			{object}	NotificationChannel
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	404			{object}	httpx.ErrorResponse
// @Failure	409			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/notification-channels/{channelId} [patch]
func (c *NotificationChannelAPIController) UpdateChannelHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	channelID, ok := httpx.PathUUIDOrRespond(ctx, w, r, "channelId")
	if !ok {
		return
	}

	payload := NotificationChannelPayload{}
	if ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload); !ok {
		return
	}

	channel, serviceErr := c.service.UpdateChannel(ctx, NotificationChannel{
		ID:      channelID,
		Name:    payload.Name,
		Kind:    payload.Kind,
		Enabled: payload.Enabled,
		Config:  payload.Config,
	})
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, channel)
}

// @Summary	Delete a notification channel
// @Tags		notification-channels
// @Produce	json
// @Param		channelId	path	string	true	"Notification channel ID"
// @Success	204
// @Failure	400	{object}	httpx.ErrorResponse
// @Failure	401	{object}	httpx.ErrorResponse
// @Failure	403	{object}	httpx.ErrorResponse
// @Failure	404	{object}	httpx.ErrorResponse
// @Failure	500	{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/notification-channels/{channelId} [delete]
func (c *NotificationChannelAPIController) DeleteChannelHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	channelID, ok := httpx.PathUUIDOrRespond(ctx, w, r, "channelId")
	if !ok {
		return
	}

	if err := c.service.DeleteChannel(ctx, channelID); err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary	Export a notification channel as YAML config
// @Tags		notification-channels
// @Produce	plain
// @Param		channelId	path		string	true	"Notification channel ID"
// @Success	200			{string}	string	"Notification channel config in YAML, secrets replaced with ${ENV} placeholders"
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	404			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/notification-channels/{channelId}/export [get]
func (c *NotificationChannelAPIController) ExportToYAMLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	channelID, ok := httpx.PathUUIDOrRespond(ctx, w, r, "channelId")
	if !ok {
		return
	}

	if err := c.service.ExportToYAML(ctx, channelID, w); err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
	}
}
