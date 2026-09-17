package api

import (
	"github.com/m-milek/leszmonitor/api/controllers"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/features/users"
)

type Handlers struct {
	Monitor                monitors.MonitorAPIController
	MonitorResults         monitors.MonitorResultsAPIController
	MonitorStats           monitors.MonitorStatsAPIController
	AuditLog               controllers.AuditLogAPIController
	User                   users.UserAPIController
	Tag                    tags.TagAPIController
	InstanceMetadata       controllers.InstanceMetadataAPIController
	AuthzMiddlewareService users.IAuthzMiddlewareService
}
