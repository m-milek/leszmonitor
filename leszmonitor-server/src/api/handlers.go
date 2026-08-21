package api

import (
	"github.com/m-milek/leszmonitor/api/controllers"
	"github.com/m-milek/leszmonitor/services"
)

type Handlers struct {
	Project                controllers.ProjectAPIController
	Monitor                controllers.MonitorAPIController
	MonitorResults         controllers.MonitorResultsAPIController
	AuditLog               controllers.AuditLogAPIController
	User                   controllers.UserAPIController
	AuthzMiddlewareService services.IAuthzMiddlewareService
}
