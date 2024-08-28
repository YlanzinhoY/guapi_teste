package handler

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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
	ws, err := s.wsUpgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		return err
	}

	s.wsConnections[ws] = true

	defer ws.Close()

	go func() {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
		}
	}()

	notificationMessage := "Notificação de teste"
	for client := range s.wsConnections {
		err := client.WriteMessage(websocket.TextMessage, []byte(notificationMessage))
		if err != nil {
			log.Println("Error sending message:", err)
			client.Close()
			delete(s.wsConnections, client)
		}
	}

	return nil

}
