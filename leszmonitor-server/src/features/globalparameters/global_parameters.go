package globalparameters

import (
	"fmt"
	"time"
)

type GlobalParameterType string

const (
	GlobalParameterStringType GlobalParameterType = "string"
	GlobalParameterIntType    GlobalParameterType = "int"
	GlobalParameterBoolType   GlobalParameterType = "bool"
)

type GlobalParameterKey string

const (
	GlobalParameterAllowSelfRegistration GlobalParameterKey = "accounts.allow_self_registration"
)

type GlobalParameterSource string

const (
	GlobalParameterSourceConfig  GlobalParameterSource = "config"
	GlobalParameterSourceUI      GlobalParameterSource = "ui"
	GlobalParameterSourceDefault GlobalParameterSource = "default"
)

type GlobalParameterDefinition struct {
	Key          GlobalParameterKey  `json:"key"`
	DisplayName  string              `json:"displayName"`
	Description  string              `json:"description"`
	Type         GlobalParameterType `json:"type"`
	DefaultValue any                 `json:"defaultValue"`
}

var GlobalParameterDefinitions = map[GlobalParameterKey]GlobalParameterDefinition{
	GlobalParameterAllowSelfRegistration: {
		Key:          GlobalParameterAllowSelfRegistration,
		DisplayName:  "Allow Self-Registration",
		Description:  "Whether to allow users to create their own accounts. If disabled, only admins/owners can create accounts.",
		Type:         GlobalParameterBoolType,
		DefaultValue: true,
	},
}

type GlobalParameter struct {
	GlobalParameterDefinition
	Value  any                   `json:"value"`
	Source GlobalParameterSource `json:"source"`
}

func (p *GlobalParameter) Validate() error {
	var ok bool
	switch p.Type {
	case GlobalParameterBoolType:
		_, ok = p.Value.(bool)
	case GlobalParameterIntType:
		_, ok = p.Value.(int)
	case GlobalParameterStringType:
		_, ok = p.Value.(string)
	default:
		return fmt.Errorf("unsupported parameter type %s", p.Type)
	}

	if !ok {
		return fmt.Errorf("expected %s, got %T", p.Type, p.Value)
	}
	return nil
}

type GlobalParameterRecord struct {
	Key       GlobalParameterKey `db:"key"`
	Value     string             `db:"value"`
	UpdatedAt time.Time          `db:"updated_at"`
}
