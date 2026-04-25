package workers

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/auth"
	"github.com/m-milek/leszmonitor/log"
)

type wsAuthMessage struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type wsAuthResponse struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

func closeUnauthorized(conn *websocket.Conn, reason string) {
	_ = conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.ClosePolicyViolation, reason),
		time.Now().Add(2*time.Second),
	)
}

func readAuthFrame(conn *websocket.Conn) (*wsAuthMessage, error) {
	messageType, message, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	if messageType != websocket.TextMessage {
		return nil, errors.New("first message must be auth")
	}

	var authMsg wsAuthMessage
	if err := json.Unmarshal(message, &authMsg); err != nil {
		return nil, errors.New("invalid auth message")
	}
	if authMsg.Type != "auth" || strings.TrimSpace(authMsg.Token) == "" {
		return nil, errors.New("invalid auth message")
	}

	return &authMsg, nil
}

func writeAuthAck(conn *websocket.Conn) error {
	authAck, err := json.Marshal(wsAuthResponse{Type: "auth", Status: "ok"})
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, authAck)
}

func authenticateConnection(ctx context.Context, conn *websocket.Conn) (context.Context, bool) {
	authMsg, err := readAuthFrame(conn)
	if err != nil {
		log.Api.Error().Err(err).Msg("Failed to read auth message from WebSocket connection")
		closeUnauthorized(conn, err.Error())
		return ctx, false
	}

	userClaims, err := auth.ValidateJwt(authMsg.Token)
	if err != nil {
		log.Api.Error().Err(err).Msg("Failed to validate JWT token from WebSocket connection")
		closeUnauthorized(conn, "invalid token")
		return ctx, false
	}

	if err := writeAuthAck(conn); err != nil {
		log.Api.Error().Err(err).Msg("Failed to write auth acknowledgment to WebSocket connection")
		return ctx, false
	}

	log.Api.Info().Any("username", userClaims.Username).Msg("WebSocket connection authenticated successfully")

	return middleware.SetUserContext(ctx, userClaims), true
}

func RunWebSocketWorker(ctx context.Context, conn *websocket.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	log.Api.Info().Any("remoteAddr", conn.RemoteAddr()).Msg("WebSocket connection established")

	authCtx, ok := authenticateConnection(ctx, conn)
	if !ok {
		return
	}

	for {
		select {
		case <-authCtx.Done():
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
