package workers

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/m-milek/leszmonitor/log"
)

func RunWebSocketWorker(ctx context.Context, conn *websocket.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	log.Api.Info().Any("remoteAddr", conn.RemoteAddr()).Msg("WebSocket connection established")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				return
			}

			if messageType == websocket.TextMessage && string(message) == "ping" {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("pong")); err != nil {
					return
				}
				continue
			}

			if err := conn.WriteMessage(messageType, message); err != nil {
				return
			}
		}
	}
}
