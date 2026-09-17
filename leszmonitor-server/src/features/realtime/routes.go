package realtime

import "net/http"

// RegisterRoutes registers the websocket route on the public router.
func RegisterRoutes(publicRouter *http.ServeMux) {
	publicRouter.HandleFunc("GET /api/ws", WebSocketConnectionHandler)
}
