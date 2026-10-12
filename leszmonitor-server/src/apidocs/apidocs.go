package apidocs

import (
	_ "embed"
	"encoding/json"
)

//go:embed swagger.json
var Spec []byte

// SpecWithVersion returns the spec with info.version set to version, or the spec unchanged if version is empty.
func SpecWithVersion(version string) ([]byte, error) {
	if version == "" {
		return Spec, nil
	}

	var spec map[string]any
	if err := json.Unmarshal(Spec, &spec); err != nil {
		return nil, err
	}
	info, _ := spec["info"].(map[string]any)
	if info == nil {
		info = map[string]any{}
		spec["info"] = info
	}
	info["version"] = version

	return json.Marshal(spec)
}
