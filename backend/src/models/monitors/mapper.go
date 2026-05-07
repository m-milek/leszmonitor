package monitors

import (
	"encoding/json"
	"fmt"

	consts "github.com/m-milek/leszmonitor/models/consts"
)

func mapProbeType(kind consts.ProbeType) Probe {
	switch kind {
	case consts.HttpConfigType:
		return &HttpProbe{}
	case consts.TCPConfigType:
		return &TCPProbe{}
	default:
		return nil
	}
}

func ProbeFromJSON(probeConfig string, probeType consts.ProbeType) (Probe, error) {
	// Map the monitor type to the appropriate config type
	probe := mapProbeType(probeType)
	if probe == nil {
		return nil, fmt.Errorf("unknown monitor type: %s", probeType)
	}

	// unmarshal the raw data into a probe instance
	if err := json.Unmarshal([]byte(probeConfig), &probe); err != nil {
		return nil, fmt.Errorf("failed to parse monitor config: %w: %s", err, probeConfig)
	}

	return probe, nil
}

func UnmarshalProbeFromBytes(kind consts.ProbeType, data []byte) (Probe, error) {
	config := mapProbeType(kind)
	if config == nil {
		return nil, fmt.Errorf("unknown monitor config type: %s", kind)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal monitor config: %w", err)
	}
	return config, nil
}
