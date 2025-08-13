package services

import (
	"github.com/m-milek/leszmonitor/logger"
	"github.com/rs/zerolog"
)

type BaseService struct {
	logger zerolog.Logger
}

func InitializeServices() {
	TeamService.logger = logger.NewServiceLogger("TeamService")
	UserService.logger = logger.NewServiceLogger("UserService")
	MonitorService.logger = logger.NewServiceLogger("MonitorService")
}
