package statuspages

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/stats"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/platform/apperr"
	"github.com/m-milek/leszmonitor/platform/audit"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/m-milek/leszmonitor/platform/util"
)

type IStatusPageService interface {
	CreateStatusPage(ctx context.Context, page StatusPage) (*StatusPage, *apperr.ServiceError)
	GetStatusPageByID(ctx context.Context, id uuid.UUID) (*StatusPage, *apperr.ServiceError)
	GetStatusPageBySlug(ctx context.Context, slug string) (*StatusPage, *apperr.ServiceError)
	GetPublicMonitors(ctx context.Context, pageID uuid.UUID, from time.Time, to time.Time) ([]PublicMonitor, *apperr.ServiceError)
	GetAllStatusPages(ctx context.Context) ([]StatusPage, *apperr.ServiceError)
	UpdateStatusPage(ctx context.Context, page StatusPage) (*StatusPage, *apperr.ServiceError)
	DeleteStatusPage(ctx context.Context, id uuid.UUID) *apperr.ServiceError
	ExportToYAML(ctx context.Context, id uuid.UUID, w io.Writer) *apperr.ServiceError
}

type StatusPageService struct {
	db    db.DB
	stats stats.IMonitorStatsService
}

type StatusPageServiceDeps struct {
	DB    db.DB
	Stats stats.IMonitorStatsService
}

func NewStatusPageService(deps StatusPageServiceDeps) *StatusPageService {
	return &StatusPageService{
		db:    deps.DB,
		stats: deps.Stats,
	}
}

func (s *StatusPageService) CreateStatusPage(ctx context.Context, page StatusPage) (*StatusPage, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameStatusPage, "CreateStatusPage")
	logger.Trace().Interface("statusPage", page).Msg("Creating status page")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewInternalError("user claims not found in context")
	}

	newPage := StatusPage{
		ID:                 uuid.New(),
		Slug:               page.Slug,
		Title:              page.Title,
		Description:        page.Description,
		GroupMonitorsByTag: page.GroupMonitorsByTag,
		RequireAuth:        page.RequireAuth,
		Enabled:            page.Enabled,
		TimeWindowSeconds:  page.TimeWindowSeconds,
		Source:             util.ResourceSourceUI,
		TagIDs:             page.TagIDs,
	}

	if err := newPage.Validate(); err != nil {
		logger.Error().Err(err).Msg("Invalid status page")
		return nil, apperr.NewBadRequestError("invalid status page: %w", err)
	}

	created, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) (*StatusPage, *audit.AuditLogParams, error) {
		p, createErr := NewStatusPageDAO(q).InsertStatusPage(ctx, newPage)
		if createErr != nil {
			return nil, nil, createErr
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &p.ID,
			Action:     audit.ActionCreateStatusPage,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Status page with ID %s created", p.ID),
			After:      p,
		}
		return p, params, nil
	})
	if txErr != nil {
		if errors.Is(txErr, db.ErrAlreadyExists) {
			return nil, apperr.NewConflictError("status page with slug %s already exists", newPage.Slug)
		}
		logger.Error().Err(txErr).Msg("Failed to create status page within transaction")
		return nil, apperr.NewInternalError("failed to create status page within transaction: %w", txErr)
	}

	logger.Debug().Str("id", created.ID.String()).Msg("Status page created")
	return created, nil
}

func (s *StatusPageService) GetStatusPageByID(ctx context.Context, id uuid.UUID) (*StatusPage, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameStatusPage, "GetStatusPageByID")
	logger.Trace().Str("id", id.String()).Msg("Retrieving status page by ID")

	page, err := NewStatusPageDAO(s.db.Querier()).GetStatusPageByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("id", id.String()).Msg("Status page not found in database")
			return nil, apperr.NewNotFoundError("status page with ID %s not found", id)
		}
		logger.Error().Err(err).Str("id", id.String()).Msg("Failed to retrieve status page from database")
		return nil, apperr.NewInternalError("failed to retrieve status page: %w", err)
	}

	logger.Debug().Str("id", id.String()).Msg("Status page retrieved successfully")
	return page, nil
}

func (s *StatusPageService) GetStatusPageBySlug(ctx context.Context, slug string) (*StatusPage, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameStatusPage, "GetStatusPageBySlug")
	logger.Trace().Str("slug", slug).Msg("Retrieving status page by slug")

	page, err := NewStatusPageDAO(s.db.Querier()).GetStatusPageBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, apperr.NewNotFoundError("status page with slug %s not found", slug)
		}
		logger.Error().Err(err).Str("slug", slug).Msg("Failed to retrieve status page from database")
		return nil, apperr.NewInternalError("failed to retrieve status page: %w", err)
	}

	return page, nil
}

func (s *StatusPageService) GetPublicMonitors(ctx context.Context, pageID uuid.UUID, from time.Time, to time.Time) ([]PublicMonitor, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameStatusPage, "GetPublicMonitors")
	logger.Trace().Str("id", pageID.String()).Msg("Retrieving public monitors for status page")

	publicMonitors, err := NewStatusPageDAO(s.db.Querier()).GetPublicMonitors(ctx, pageID)
	if err != nil {
		logger.Error().Err(err).Str("id", pageID.String()).Msg("Failed to retrieve public monitors from database")
		return nil, apperr.NewInternalError("failed to retrieve status page monitors: %w", err)
	}

	for i := range publicMonitors {
		monitorStats, statsErr := s.stats.GetStatsByMonitorID(ctx, publicMonitors[i].ID, from, to)
		if statsErr != nil {
			return nil, statsErr
		}
		publicMonitors[i].Stats = monitorStats
	}

	return publicMonitors, nil
}

func (s *StatusPageService) GetAllStatusPages(ctx context.Context) ([]StatusPage, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameStatusPage, "GetAllStatusPages")
	logger.Trace().Msg("Retrieving all status pages")

	pages, err := NewStatusPageDAO(s.db.Querier()).GetAllStatusPages(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve status pages from database")
		return nil, apperr.NewInternalError("failed to retrieve status pages: %w", err)
	}

	logger.Debug().Int("count", len(pages)).Msg("Status pages retrieved successfully")
	return pages, nil
}

func (s *StatusPageService) UpdateStatusPage(ctx context.Context, page StatusPage) (*StatusPage, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameStatusPage, "UpdateStatusPage")
	logger.Trace().Interface("statusPage", page).Msg("Updating status page")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewInternalError("user claims not found in context")
	}

	if err := page.Validate(); err != nil {
		logger.Error().Err(err).Str("id", page.ID.String()).Msg("Invalid status page")
		return nil, apperr.NewBadRequestError("invalid status page: %w", err)
	}

	updated, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) (*StatusPage, *audit.AuditLogParams, error) {
		dao := NewStatusPageDAO(q)
		existing, err := dao.GetStatusPageByID(ctx, page.ID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return nil, nil, apperr.NewNotFoundError("status page with ID %s not found", page.ID)
			}
			return nil, nil, fmt.Errorf("failed to retrieve existing status page for update: %w", err)
		}

		if existing.Source == util.ResourceSourceConfig {
			return nil, nil, apperr.NewConflictError("status page %s is managed by config file", page.ID)
		}

		p, updateErr := dao.UpdateStatusPage(ctx, page)
		if updateErr != nil {
			return nil, nil, updateErr
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &p.ID,
			Action:     audit.ActionUpdateStatusPage,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Status page with ID %s updated", p.ID),
			Before:     existing,
			After:      p,
		}
		return p, params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*apperr.ServiceError](txErr); ok2 {
			return nil, serviceErr
		}
		if errors.Is(txErr, db.ErrAlreadyExists) {
			return nil, apperr.NewConflictError("status page with slug %s already exists", page.Slug)
		}
		logger.Error().Err(txErr).Str("id", page.ID.String()).Msg("Failed to update status page within transaction")
		return nil, apperr.NewInternalError("failed to update status page within transaction: %w", txErr)
	}

	logger.Debug().Str("id", updated.ID.String()).Msg("Status page updated")
	return updated, nil
}

func (s *StatusPageService) DeleteStatusPage(ctx context.Context, id uuid.UUID) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameStatusPage, "DeleteStatusPage")
	logger.Trace().Str("id", id.String()).Msg("Deleting status page")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return apperr.NewInternalError("user claims not found in context")
	}

	txErr := audit.WithAuditedVoidTx(ctx, s.db, func(q db.Querier) (*audit.AuditLogParams, error) {
		dao := NewStatusPageDAO(q)
		existing, err := dao.GetStatusPageByID(ctx, id)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return nil, apperr.NewNotFoundError("status page with ID %s not found", id)
			}
			return nil, fmt.Errorf("failed to retrieve status page before deletion: %w", err)
		}

		if existing.Source == util.ResourceSourceConfig {
			return nil, apperr.NewConflictError("status page %s is managed by config file", id)
		}

		if _, delErr := dao.DeleteStatusPageByID(ctx, id); delErr != nil {
			return nil, delErr
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &id,
			Action:     audit.ActionDeleteStatusPage,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Status page with ID %s deleted", id),
			Before:     existing,
		}
		return params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*apperr.ServiceError](txErr); ok2 {
			return serviceErr
		}
		logger.Error().Err(txErr).Str("id", id.String()).Msg("Failed to delete status page")
		return apperr.NewInternalError("failed to delete status page: %w", txErr)
	}

	logger.Debug().Str("id", id.String()).Msg("Status page deleted")
	return nil
}

func (s *StatusPageService) ExportToYAML(ctx context.Context, id uuid.UUID, w io.Writer) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameStatusPage, "ExportToYAML")
	logger.Trace().Str("id", id.String()).Msg("Exporting to YAML")

	page, err := NewStatusPageDAO(s.db.Querier()).GetStatusPageByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apperr.NewNotFoundError("status page with ID %s not found", id)
		}
		return apperr.NewInternalError("failed to retrieve status page: %w", err)
	}

	if err = exportToYAML(ctx, *page, tags.NewTagDAO(s.db.Querier()), w); err != nil {
		return apperr.NewInternalError("failed to export status page %s to YAML", page.Slug)
	}

	return nil
}
