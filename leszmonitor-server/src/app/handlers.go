package app

import (
	"github.com/m-milek/leszmonitor/features/auditlog"
	"github.com/m-milek/leszmonitor/features/instance"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/features/monitors/stats"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/features/users"
)

type Handlers struct {
	Monitor                monitors.MonitorAPIController
	MonitorResults         results.MonitorResultsAPIController
	MonitorStats           stats.MonitorStatsAPIController
	AuditLog               auditlog.AuditLogAPIController
	User                   users.UserAPIController
	Tag                    tags.TagAPIController
	InstanceMetadata       instance.InstanceMetadataAPIController
	AuthzMiddlewareService users.IAuthzMiddlewareService
}
