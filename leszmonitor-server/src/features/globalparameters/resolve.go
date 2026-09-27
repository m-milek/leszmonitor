package globalparameters

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/m-milek/leszmonitor/platform/config"
	"github.com/m-milek/leszmonitor/platform/db"
)

func ValidateConfigFileParameters() error {
	for _, definition := range GlobalParameterDefinitions {
		if _, err := resolveFromConfigFile(definition); err != nil {
			return err
		}
	}
	return nil
}

func resolveFromConfigFile(paramDefinition GlobalParameterDefinition) (*GlobalParameter, error) {
	configFileValue, ok := config.LookupGlobalParameter(string(paramDefinition.Key))
	if !ok || configFileValue == nil {
		return nil, nil
	}

	param := &GlobalParameter{
		GlobalParameterDefinition: paramDefinition,
		Value:                     configFileValue,
		Source:                    GlobalParameterSourceConfig,
	}
	if err := param.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config file value for %s: %w", paramDefinition.Key, err)
	}
	return param, nil
}

func resolveGlobalParameter(ctx context.Context, paramDefinition GlobalParameterDefinition, dao IGlobalParameterDAO) (*GlobalParameter, error) {
	fromConfig, err := resolveFromConfigFile(paramDefinition)
	if err != nil {
		return nil, err
	}
	if fromConfig != nil {
		return fromConfig, nil
	}

	record, err := dao.GetParameterByKey(ctx, paramDefinition.Key)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return nil, err
	}

	if record != nil {
		value, parseErr := parseStoredValue(paramDefinition, record.Value)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid stored value for %s: %w", paramDefinition.Key, parseErr)
		}

		param := &GlobalParameter{
			GlobalParameterDefinition: paramDefinition,
			Value:                     value,
			Source:                    GlobalParameterSourceUI,
		}
		if err := param.Validate(); err != nil {
			return nil, fmt.Errorf("invalid stored value for %s: %w", paramDefinition.Key, err)
		}
		return param, nil
	}

	return &GlobalParameter{
		GlobalParameterDefinition: paramDefinition,
		Value:                     paramDefinition.DefaultValue,
		Source:                    GlobalParameterSourceDefault,
	}, nil
}

func parseStoredValue(paramDefinition GlobalParameterDefinition, raw string) (any, error) {
	switch paramDefinition.Type {
	case GlobalParameterBoolType:
		return strconv.ParseBool(raw)
	case GlobalParameterIntType:
		return strconv.Atoi(raw)
	case GlobalParameterStringType, GlobalParameterEnumType:
		return raw, nil
	default:
		return nil, fmt.Errorf("unsupported parameter type %s", paramDefinition.Type)
	}
}
