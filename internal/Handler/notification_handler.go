package handler

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/ylanzinhoy/guapi_teste/internal/entity"
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

func (s *NotificationHandler) SendNotificationLikeUnLikeMessage(c echo.Context) error {
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

	initialLikes, err := s.dbHandler.GetMessagesLikesByChatId(c.Request().Context(), chatRoomIdParam)
	if err != nil {
		return err
	}
	var participantEntity []entity.ParticipantsEntity

	participantsSubs, err := s.dbHandler.FindAllParticipantsSubscribers(c.Request().Context(), chatRoomIdParam)
	if err != nil {
		return err
	}

	for _, participant := range participantsSubs {
		pe := entity.ParticipantsEntity{
			Id:          participant.ParticipantsID,
			IsSubscribe: participant.IsSubscribe,
			ChatRoomId:  chatRoomIdParam,
			Name:        participant.Name,
		}
		participantEntity = append(participantEntity, pe)
	}

	for {
		currentLikes, err := s.dbHandler.GetMessagesLikesByChatId(c.Request().Context(), chatRoomIdParam)
		if err != nil {
			return err
		}
		for messageID, currentLikeCount := range currentLikes {

			if currentLikeCount.LikeMessage > initialLikes[messageID].LikeMessage {
				notification := map[string]interface{}{
					"notification_message":                         "Like Message",
					"send_notification_to_participant_subscribers": &participantEntity,
					"chat_room_id":                                 chatRoomIdParam,
					"content":                                      currentLikeCount.Content,
					"participant_name":                             currentLikeCount.Name,
					"message_id":                                   currentLikeCount.MessageID,
					"like_message":                                 currentLikeCount.LikeMessage,
				}
				s.broadcastNotification(notification)

				initialLikes[messageID] = currentLikeCount
			} else if currentLikeCount.LikeMessage < initialLikes[messageID].LikeMessage {
				notification := map[string]interface{}{
					"notification_message":                         "Deslike Message",
					"send_notification_to_participant_subscribers": &participantEntity,
					"chat_room_id":                                 chatRoomIdParam,
					"content":                                      currentLikeCount.Content,
					"participant_name":                             currentLikeCount.Name,
					"message_id":                                   currentLikeCount.MessageID,
					"like_message":                                 currentLikeCount.LikeMessage,
				}

				s.broadcastNotification(notification)

				initialLikes[messageID] = currentLikeCount
			}
		}
	}
}

func (s *NotificationHandler) SendNotificationMessage(c echo.Context) error {
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

			notification := map[string]interface{}{
				"notification": "New messages available",
				"count":        messagesCount,
				"users":        participants,
			}
			s.broadcastNotification(notification)
			initialMessageCount = messagesCount
		}
		time.Sleep(2 * time.Second)
	}
}

func (s *NotificationHandler) broadcastNotification(notification map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.wsConnections {
		if err := conn.WriteJSON(notification); err != nil {
			log.Printf("error writing notification to websocket: %v", err)
			conn.Close()
			delete(s.wsConnections, conn)
		}
	}
}
