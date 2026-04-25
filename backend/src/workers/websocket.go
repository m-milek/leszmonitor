package workers

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/m-milek/leszmonitor/log"
)

func RunWebSocketWorker(ctx context.Context, conn *websocket.Conn) {
	defer conn.Close()

	log.Api.Info().Any("remoteAddr", conn.RemoteAddr()).Msg("WebSocket connection established")

	// First message auth
	//_, message, err := conn.ReadMessage()
	//if err != nil {
	//	log.Api.Error().Err(err).Msg("Failed to read authentication message from WebSocket connection")
	//	return
	//}

	//authToken := string(message)

	log.Api.Info().Any("remoteAddr", conn.RemoteAddr()).Msg("WebSocket client authenticated successfully")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}
