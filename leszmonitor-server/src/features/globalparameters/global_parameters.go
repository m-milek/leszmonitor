package globalparameters

import (
	"fmt"
	"slices"
	"time"

	"github.com/m-milek/leszmonitor/platform/auth"
)

type GlobalParameterType string

const (
	GlobalParameterStringType GlobalParameterType = "string"
	GlobalParameterIntType    GlobalParameterType = "int"
	GlobalParameterBoolType   GlobalParameterType = "bool"
	GlobalParameterEnumType   GlobalParameterType = "enum"
)

type GlobalParameterKey string

const (
	GlobalParameterAllowSelfRegistration GlobalParameterKey = "accounts.allow_self_registration"
	GlobalParameterDefaultRole           GlobalParameterKey = "accounts.default_role"
	GlobalParameterInstanceName          GlobalParameterKey = "instance.name"
	GlobalParameterInstancePublicURL     GlobalParameterKey = "instance.public_url"
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
	Options      []string            `json:"options,omitempty"`
}

var GlobalParameterDefinitions = map[GlobalParameterKey]GlobalParameterDefinition{
	GlobalParameterAllowSelfRegistration: {
		Key:          GlobalParameterAllowSelfRegistration,
		DisplayName:  "Allow Self-Registration",
		Description:  "Whether to allow users to create their own accounts. If disabled, only admins/owners can create accounts.",
		Type:         GlobalParameterBoolType,
		DefaultValue: true,
	},
	GlobalParameterDefaultRole: {
		Key:          GlobalParameterDefaultRole,
		DisplayName:  "Default Role",
		Description:  "Role assigned to newly registered users.",
		Type:         GlobalParameterEnumType,
		DefaultValue: string(auth.RoleViewer),
		Options:      []string{string(auth.RoleViewer), string(auth.RoleWriter), string(auth.RoleAdmin)},
	},
	GlobalParameterInstanceName: {
		Key:          GlobalParameterInstanceName,
		DisplayName:  "Instance Name",
		Description:  "Name of this Leszmonitor instance.",
		Type:         GlobalParameterStringType,
		DefaultValue: "Leszmonitor",
	},
	GlobalParameterInstancePublicURL: {
		Key:          GlobalParameterInstancePublicURL,
		DisplayName:  "Public URL",
		Description:  "URL under which this instance is reachable, used in links pointing to it.",
		Type:         GlobalParameterStringType,
		DefaultValue: "",
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
	case GlobalParameterEnumType:
		var value string
		value, ok = p.Value.(string)
		if ok && !slices.Contains(p.Options, value) {
			return fmt.Errorf("value %s is not one of %v", value, p.Options)
		}
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
