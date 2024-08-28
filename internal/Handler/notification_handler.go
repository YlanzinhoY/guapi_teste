package handler

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/ylanzinhoy/guapi_teste/internal/entity"
	db "github.com/ylanzinhoy/guapi_teste/sql"
)

type NotificationHandler struct {
	dbHandler     *db.Queries
	wsUpgrader    *websocket.Upgrader
	wsConnections map[*websocket.Conn]bool
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
	chatRoomId, err := uuid.Parse(c.Param("chat_room_id"))
	if err != nil {
		return c.JSON(500, err.Error())
	}
	ws, err := s.wsUpgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		return err
	}

	s.wsConnections[ws] = true

	defer ws.Close()

	for {
		var notificationEntity entity.NotificationEntity

		err := ws.ReadJSON(&notificationEntity)
		if err != nil {
			log.Printf("error reading message: %v", err)
			delete(s.wsConnections, ws)
			break
		}

		notificationEntity.NotificationID = uuid.New()
		notificationEntity.CreatedAt = time.Now()
		notificationEntity.FkChatRoomID = chatRoomId
		notificationEntity.Message = "Nova Mensagem!"

		err = s.dbHandler.CreateNotificationForSubscribers(c.Request().Context(), db.CreateNotificationForSubscribersParams{
			Message:      notificationEntity.Message,
			FkMessageID:  notificationEntity.FkMessageID,
			FkChatRoomID: notificationEntity.FkChatRoomID,
		})

		if err != nil {
			log.Printf("error saving message: %v", err)
			continue
		}

		for conn := range s.wsConnections {
			if err := conn.WriteJSON(notificationEntity.Message); err != nil {
				log.Printf("error writing message to websocket: %v", err)
				err := conn.Close()
				if err != nil {
					return err
				}
				delete(s.wsConnections, conn)
			}
		}
	}
	return nil
}
