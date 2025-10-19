package monitors

import "github.com/m-milek/leszmonitor/util"

type PingMonitorResponse struct {
	baseMonitorResponse `bson:",inline"`
	Tries               int64 `json:"tries" bson:"tries"` // Number of tries made to pingType the host
}

func NewPingMonitorResponse() *PingMonitorResponse {
	return &PingMonitorResponse{
		baseMonitorResponse: baseMonitorResponse{
			Status:    success,
			Timestamp: util.GetUnixTimestamp(),
			Failures:  []string{},
			ErrorsStr: []string{},
			Errors:    []error{},
		},
		Tries: 0,
	}
}
