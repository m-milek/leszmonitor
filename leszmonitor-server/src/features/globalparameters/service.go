package globalparameters

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/m-milek/leszmonitor/platform/apperr"
	"github.com/m-milek/leszmonitor/platform/audit"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/config"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
)

type IGlobalParameterService interface {
	GetAllParameters(ctx context.Context) ([]GlobalParameter, *apperr.ServiceError)
	GetParameter(ctx context.Context, key GlobalParameterKey) (*GlobalParameter, *apperr.ServiceError)
	SetParameter(ctx context.Context, key GlobalParameterKey, value any) (*GlobalParameter, *apperr.ServiceError)
}

type GlobalParameterService struct {
	db                 db.DB
	globalParameterDAO IGlobalParameterDAO
}

type GlobalParameterServiceDeps struct {
	DB db.DB
}

func NewGlobalParameterService(deps GlobalParameterServiceDeps) *GlobalParameterService {
	return &GlobalParameterService{
		db:                 deps.DB,
		globalParameterDAO: NewGlobalParameterDAO(deps.DB.Querier()),
	}
}

func (s *GlobalParameterService) GetAllParameters(ctx context.Context) ([]GlobalParameter, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameGlobalParameter, "GetAllParameters")
	logger.Trace().Msg("Retrieving all global parameters")

	params := make([]GlobalParameter, 0, len(GlobalParameterDefinitions))
	for key, definition := range GlobalParameterDefinitions {
		param, err := resolveGlobalParameter(ctx, definition, s.globalParameterDAO)
		if err != nil {
			logger.Error().Err(err).Str("key", string(key)).Msg("Failed to resolve global parameter")
			return nil, apperr.NewInternalError("failed to resolve global parameter %s: %w", key, err)
		}
		params = append(params, *param)
	}

	slices.SortFunc(params, func(a, b GlobalParameter) int {
		return strings.Compare(string(a.Key), string(b.Key))
	})

	logger.Debug().Int("count", len(params)).Msg("Global parameters retrieved")
	return params, nil
}

func (s *GlobalParameterService) GetParameter(ctx context.Context, key GlobalParameterKey) (*GlobalParameter, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameGlobalParameter, "GetParameter")
	logger.Trace().Str("key", string(key)).Msg("Retrieving global parameter")

	definition, ok := GlobalParameterDefinitions[key]
	if !ok {
		logger.Error().Str("key", string(key)).Msg("Unknown global parameter")
		return nil, apperr.NewNotFoundError("global parameter %s not found", key)
	}

	param, err := resolveGlobalParameter(ctx, definition, s.globalParameterDAO)
	if err != nil {
		logger.Error().Err(err).Str("key", string(key)).Msg("Failed to resolve global parameter")
		return nil, apperr.NewInternalError("failed to resolve global parameter %s: %w", key, err)
	}

	logger.Debug().Str("key", string(key)).Str("source", string(param.Source)).Msg("Global parameter retrieved")
	return param, nil
}

func (s *GlobalParameterService) SetParameter(ctx context.Context, key GlobalParameterKey, value any) (*GlobalParameter, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameGlobalParameter, "SetParameter")
	logger.Trace().Str("key", string(key)).Interface("value", value).Msg("Setting global parameter")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewUnauthorizedError("user claims not found in context")
	}

	definition, ok := GlobalParameterDefinitions[key]
	if !ok {
		logger.Error().Str("key", string(key)).Msg("Unknown global parameter")
		return nil, apperr.NewNotFoundError("global parameter %s not found", key)
	}

	if configValue, inConfig := config.LookupGlobalParameter(string(key)); inConfig && configValue != nil {
		logger.Error().Str("key", string(key)).Msg("Global parameter is set in config file")
		return nil, apperr.NewConflictError("global parameter %s is set in the config file and cannot be changed", key)
	}

	raw, err := serializeValue(definition, value)
	if err != nil {
		logger.Error().Err(err).Str("key", string(key)).Msg("Invalid global parameter value")
		return nil, apperr.NewBadRequestError("invalid value for global parameter %s: %w", key, err)
	}

	param, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) (*GlobalParameter, *audit.AuditLogParams, error) {
		dao := NewGlobalParameterDAO(q)

		before, resolveErr := resolveGlobalParameter(ctx, definition, dao)
		if resolveErr != nil {
			return nil, nil, fmt.Errorf("failed to resolve global parameter before update: %w", resolveErr)
		}

		if _, upsertErr := dao.UpsertParameter(ctx, GlobalParameterRecord{Key: key, Value: raw}); upsertErr != nil {
			return nil, nil, fmt.Errorf("failed to save global parameter: %w", upsertErr)
		}

		after, resolveErr := resolveGlobalParameter(ctx, definition, dao)
		if resolveErr != nil {
			return nil, nil, fmt.Errorf("failed to resolve global parameter after update: %w", resolveErr)
		}

		params := &audit.AuditLogParams{
			Username:  &userClaims.Username,
			Action:    audit.ActionUpdateGlobalParam,
			IsSuccess: true,
			Summary:   fmt.Sprintf("Global parameter %s updated", key),
			Before:    before,
			After:     after,
		}
		return after, params, nil
	})
	if txErr != nil {
		if serviceErr, isServiceErr := errors.AsType[*apperr.ServiceError](txErr); isServiceErr {
			return nil, serviceErr
		}
		logger.Error().Err(txErr).Str("key", string(key)).Msg("Failed to set global parameter within transaction")
		return nil, apperr.NewInternalError("failed to set global parameter %s within transaction: %w", key, txErr)
	}

	logger.Debug().Str("key", string(key)).Msg("Global parameter set")
	return param, nil
}

func serializeValue(definition GlobalParameterDefinition, value any) (string, error) {
	switch definition.Type {
	case GlobalParameterBoolType:
		if v, ok := value.(bool); ok {
			return strconv.FormatBool(v), nil
		}
	case GlobalParameterIntType:
		switch v := value.(type) {
		case int:
			return strconv.Itoa(v), nil
		case float64:
			if v == math.Trunc(v) {
				return strconv.Itoa(int(v)), nil
			}
		}
	case GlobalParameterStringType:
		if v, ok := value.(string); ok {
			return v, nil
		}
	default:
		return "", fmt.Errorf("unsupported parameter type %s", definition.Type)
	}

	return "", fmt.Errorf("expected %s, got %T", definition.Type, value)
}
