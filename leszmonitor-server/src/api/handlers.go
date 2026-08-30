package api

import (
	"github.com/m-milek/leszmonitor/api/controllers"
	"github.com/m-milek/leszmonitor/services"
)

type Handlers struct {
	Project                controllers.ProjectAPIController
	Monitor                controllers.MonitorAPIController
	MonitorResults         controllers.MonitorResultsAPIController
	MonitorStats           controllers.MonitorStatsAPIController
	AuditLog               controllers.AuditLogAPIController
	User                   controllers.UserAPIController
	InstanceMetadata       controllers.InstanceMetadataAPIController
	AuthzMiddlewareService services.IAuthzMiddlewareService
}
