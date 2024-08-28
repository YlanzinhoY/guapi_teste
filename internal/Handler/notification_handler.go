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

	ws, err := s.wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.wsConnections[ws] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.wsConnections, ws)
		s.mu.Unlock()
		ws.Close()
	}()

	initialMessageCount, err := s.dbHandler.CountMessageById(c.Request().Context(), chatRoomIdParam)
	if err != nil {
		return err
	}

	for {
		messagesCount, err := s.dbHandler.CountMessageById(c.Request().Context(), chatRoomIdParam)
		if err != nil {
			return err
		}

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

			s.mu.Lock()
			for conn := range s.wsConnections {
				if err := conn.WriteJSON(notification); err != nil {
					log.Printf("error writing notification to websocket: %v", err)
					conn.Close()
					delete(s.wsConnections, conn)
				}
			}
			s.mu.Unlock()

			initialMessageCount = messagesCount
		}

		time.Sleep(2 * time.Second)
	}
}
