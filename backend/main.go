package main

import (
	"github.com/m-milek/leszmonitor/api"
)

func main() {
	var serverConfig = api.DefaultServerConfig()
	api.StartServer(serverConfig)
}
