package api

import (
	"github.com/m-milek/leszmonitor/api/controllers"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/services"
)

type Handlers struct {
	Monitor                controllers.MonitorAPIController
	MonitorResults         controllers.MonitorResultsAPIController
	MonitorStats           controllers.MonitorStatsAPIController
	AuditLog               controllers.AuditLogAPIController
	User                   controllers.UserAPIController
	Tag                    tags.TagAPIController
	InstanceMetadata       controllers.InstanceMetadataAPIController
	AuthzMiddlewareService services.IAuthzMiddlewareService
}
