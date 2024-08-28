package handler

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	db "github.com/ylanzinhoy/guapi_teste/sql"
	"log"
	"sync"
	"time"
)

type NotificationHandler struct {
	dbHandler     *db.Queries
	wsUpgrader    *websocket.Upgrader
	wsConnections map[*websocket.Conn]bool
	mu            sync.Mutex
}

func NewNotificationHandler(dbHandler *db.Queries, wsUpgrader *websocket.Upgrader,
	wsConnections map[*websocket.Conn]bool) *NotificationHandler {
	return &NotificationHandler{
		dbHandler:     dbHandler,
		wsUpgrader:    wsUpgrader,
		wsConnections: wsConnections,
	}
}

func (s *NotificationHandler) SendNotification(c echo.Context) error {
	chatRoomIdParam := uuid.MustParse(c.Param("chat_room_id"))

	// Upgrade HTTP connection to WebSocket
	ws, err := s.wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	// Add WebSocket connection to the map
	s.mu.Lock()
	s.wsConnections[ws] = true
	s.mu.Unlock()

	defer func() {
		// Remove WebSocket connection from the map and close the connection
		s.mu.Lock()
		delete(s.wsConnections, ws)
		s.mu.Unlock()
		ws.Close()
	}()

	// Get initial message count
	initialMessageCount, err := s.dbHandler.CountMessageById(c.Request().Context(), chatRoomIdParam)
	if err != nil {
		return err
	}

	for {
		// Get the current message count
		messagesCount, err := s.dbHandler.CountMessageById(c.Request().Context(), chatRoomIdParam)
		if err != nil {
			return err
		}

		// Check if there are new messages
		if messagesCount > initialMessageCount {
			participants, err := s.dbHandler.FindAllParticipantsSubscribers(c.Request().Context(), chatRoomIdParam)
			if err != nil {
				return err
			}

			// Prepare notification payload
			notification := map[string]interface{}{
				"notification": "New messages available",
				"count":        messagesCount,
				"users":        participants,
			}

			// Send notification to all WebSocket connections
			s.mu.Lock()
			for conn := range s.wsConnections {
				if err := conn.WriteJSON(notification); err != nil {
					log.Printf("error writing notification to websocket: %v", err)
					conn.Close()
					delete(s.wsConnections, conn)
				}
			}
			s.mu.Unlock()

			// Update the initial message count
			initialMessageCount = messagesCount
		}

		// Pause before the next check
		time.Sleep(2 * time.Second)
	}
}
