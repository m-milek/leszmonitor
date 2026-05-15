package authorization

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/db"
)

type ProjectAuthorization struct {
	ProjectID uuid.UUID
	Username  string
}

type Payload struct {
	MonitorID   string
	ProjectSlug string
}

func newFromConfig(ctx context.Context, config Payload) (*ProjectAuthorization, error) {
	if config.MonitorID != "" {
		return newFromMonitorID(ctx, config.MonitorID)
	} else if config.ProjectSlug != "" {
		return newFromProjectSlug(ctx, config.ProjectSlug)
	}

	return nil, fmt.Errorf("either monitor ID or project slug must be provided for authorization")
}

func NewOrRespond(ctx context.Context, w http.ResponseWriter, config Payload) (*ProjectAuthorization, bool) {
	projectAuth, err := newFromConfig(ctx, config)
	if err != nil {
		util.RespondError(ctx, w, http.StatusBadRequest, err)
		return nil, false
	}
	return projectAuth, true
}

func newFromMonitorID(ctx context.Context, monitorID string) (*ProjectAuthorization, error) {
	monitorUUID, err := uuid.Parse(monitorID)
	if err != nil {
		return nil, fmt.Errorf("invalid monitor ID format: %w", err)
	}
	monitor, err := db.Get().Monitors().GetMonitorByID(context.Background(), monitorUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get monitor: %w", err)
	}
	project, err := db.Get().Projects().GetProjectByID(context.Background(), monitor.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	username, err := GetUsernameFromRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get username from request: %w", err)
	}

	return &ProjectAuthorization{
		ProjectID: project.ID,
		Username:  *username,
	}, nil
}

func newFromProjectSlug(ctx context.Context, projectSlug string) (*ProjectAuthorization, error) {
	project, err := db.Get().Projects().GetProjectBySlug(context.Background(), projectSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	username, err := GetUsernameFromRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get username from request: %w", err)
	}

	return &ProjectAuthorization{
		ProjectID: project.ID,
		Username:  *username,
	}, nil
}
