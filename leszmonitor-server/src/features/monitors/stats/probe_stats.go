package stats

import (
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
)

func mapToProbeStats(probeType kind.ProbeType, monitorResults []results.IMonitorResult) (interface{}, error) {
	switch probeType {
	case kind.HTTPConfigType:
		return getHTTPProbeStats(monitorResults), nil
	default:
		return nil, nil
	}
}

func getHTTPProbeStats(monitorResults []results.IMonitorResult) HTTPProbeStats {
	httpProbeStats := HTTPProbeStats{
		HTTPCodeToCount: make(map[int]int),
	}

	for _, result := range monitorResults {
		resultDetails, ok := result.GetDetails().(*results.HTTPResultDetails)
		if !ok {
			continue
		}
		httpProbeStats.HTTPCodeToCount[resultDetails.StatusCode]++
	}

	return httpProbeStats
}
