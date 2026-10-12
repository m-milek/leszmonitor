package realtime

import "net/http"

// RegisterRoutes registers the websocket route.
func RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /api/v1/ws", WebSocketConnectionHandler)
}
