package api

import "github.com/m-milek/leszmonitor/api/controllers"

type Handlers struct {
	Project        controllers.ProjectAPIController
	Monitor        controllers.MonitorAPIController
	MonitorResults controllers.MonitorResultsAPIController
	AuditLog       controllers.AuditLogAPIController
	User           controllers.UserAPIController
}
