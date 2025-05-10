package monitors

type BaseMonitor struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Interval    int    `json:"interval" bson:"interval"` // in seconds
	Timeout     int    `json:"timeout" bson:"timeout"`   // in seconds
	OwnerId     string `json:"owner_id" bson:"owner_id"`
	Type        string `json:"type" bson:"type"`
}

type IMonitor interface {
	Run() error
	GetName() string
	GetDescription() string
	GetInterval() int
	GetTimeout() int
}

const (
	MonitorTypeHttp = "http"
)
