package globalparameters

import (
	"context"
	"errors"
	"fmt"
	"maps"
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
	SetParameters(ctx context.Context, values map[GlobalParameterKey]any) ([]GlobalParameter, *apperr.ServiceError)
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

func (s *GlobalParameterService) SetParameters(ctx context.Context, values map[GlobalParameterKey]any) ([]GlobalParameter, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameGlobalParameter, "SetParameters")
	logger.Trace().Interface("values", values).Msg("Setting global parameters")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewUnauthorizedError("user claims not found in context")
	}

	if len(values) == 0 {
		logger.Error().Msg("No global parameters to set")
		return nil, apperr.NewBadRequestError("no global parameters to set")
	}

	keys := slices.Sorted(maps.Keys(values))
	definitions := make([]GlobalParameterDefinition, 0, len(keys))
	records := make([]GlobalParameterRecord, 0, len(keys))
	resets := make([]GlobalParameterKey, 0, len(keys))
	for _, key := range keys {
		definition, ok := GlobalParameterDefinitions[key]
		if !ok {
			logger.Error().Str("key", string(key)).Msg("Unknown global parameter")
			return nil, apperr.NewBadRequestError("unknown global parameter %s", key)
		}

		if configValue, inConfig := config.LookupGlobalParameter(string(key)); inConfig && configValue != nil {
			logger.Error().Str("key", string(key)).Msg("Global parameter is set in config file")
			return nil, apperr.NewConflictError("global parameter %s is set in the config file and cannot be changed", key)
		}

		definitions = append(definitions, definition)

		if values[key] == nil {
			resets = append(resets, key)
			continue
		}

		raw, err := serializeValue(definition, values[key])
		if err != nil {
			logger.Error().Err(err).Str("key", string(key)).Msg("Invalid global parameter value")
			return nil, apperr.NewBadRequestError("invalid value for global parameter %s: %w", key, err)
		}

		records = append(records, GlobalParameterRecord{Key: key, Value: raw})
	}

	params, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) ([]GlobalParameter, *audit.AuditLogParams, error) {
		dao := NewGlobalParameterDAO(q)

		before, resolveErr := resolveGlobalParameters(ctx, definitions, dao)
		if resolveErr != nil {
			return nil, nil, fmt.Errorf("failed to resolve global parameters before update: %w", resolveErr)
		}

		for _, record := range records {
			if _, upsertErr := dao.UpsertParameter(ctx, record); upsertErr != nil {
				return nil, nil, fmt.Errorf("failed to save global parameter %s: %w", record.Key, upsertErr)
			}
		}

		for _, key := range resets {
			if _, deleteErr := dao.DeleteParameterByKey(ctx, key); deleteErr != nil && !errors.Is(deleteErr, db.ErrNotFound) {
				return nil, nil, fmt.Errorf("failed to reset global parameter %s: %w", key, deleteErr)
			}
		}

		after, resolveErr := resolveGlobalParameters(ctx, definitions, dao)
		if resolveErr != nil {
			return nil, nil, fmt.Errorf("failed to resolve global parameters after update: %w", resolveErr)
		}

		params := &audit.AuditLogParams{
			Username:  &userClaims.Username,
			Action:    audit.ActionUpdateGlobalParam,
			IsSuccess: true,
			Summary:   fmt.Sprintf("Global parameters updated: %s", joinKeys(keys)),
			Before:    before,
			After:     after,
		}
		return after, params, nil
	})
	if txErr != nil {
		if serviceErr, isServiceErr := errors.AsType[*apperr.ServiceError](txErr); isServiceErr {
			return nil, serviceErr
		}
		logger.Error().Err(txErr).Msg("Failed to set global parameters within transaction")
		return nil, apperr.NewInternalError("failed to set global parameters within transaction: %w", txErr)
	}

	logger.Debug().Int("count", len(params)).Msg("Global parameters set")
	return params, nil
}

func resolveGlobalParameters(ctx context.Context, definitions []GlobalParameterDefinition, dao IGlobalParameterDAO) ([]GlobalParameter, error) {
	params := make([]GlobalParameter, 0, len(definitions))
	for _, definition := range definitions {
		param, err := resolveGlobalParameter(ctx, definition, dao)
		if err != nil {
			return nil, err
		}
		params = append(params, *param)
	}
	return params, nil
}

func joinKeys(keys []GlobalParameterKey) string {
	parts := make([]string, len(keys))
	for i, key := range keys {
		parts[i] = string(key)
	}
	return strings.Join(parts, ", ")
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
	case GlobalParameterEnumType:
		if v, ok := value.(string); ok {
			if !slices.Contains(definition.Options, v) {
				return "", fmt.Errorf("value %s is not one of %v", v, definition.Options)
			}
			return v, nil
		}
	default:
		return "", fmt.Errorf("unsupported parameter type %s", definition.Type)
	}

	return "", fmt.Errorf("expected %s, got %T", definition.Type, value)
}
