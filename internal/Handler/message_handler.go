package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/ylanzinhoy/guapi_teste/internal/entity"
	db "github.com/ylanzinhoy/guapi_teste/sql"
)

type MessageHandler struct {
	dbHandler     *db.Queries
	wsUpgrader    *websocket.Upgrader
	wsConnections map[*websocket.Conn]bool
}

func NewMessageHandler(dbHandler *db.Queries, wsUpgrader *websocket.Upgrader,
	wsConnections map[*websocket.Conn]bool) *MessageHandler {
	return &MessageHandler{
		dbHandler:     dbHandler,
		wsUpgrader:    wsUpgrader,
		wsConnections: wsConnections,
	}

}

func (s *MessageHandler) CreateMessageWS(c echo.Context) error {
	chatRoomID := c.Param("chatRoomId")
	participantID := c.Param("participantId")

	ws, err := s.wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer ws.Close()

	s.wsConnections[ws] = true

	for {
		var messageEntity entity.MessageEntity

		err := ws.ReadJSON(&messageEntity)
		if err != nil {
			log.Printf("error reading message: %v", err)
			delete(s.wsConnections, ws)
			break
		}

		messageEntity.MessageId = uuid.New()
		messageEntity.ChatRoomId = uuid.MustParse(chatRoomID)
		messageEntity.ParticipantsId = uuid.MustParse(participantID)
		messageEntity.CreatedAt = time.Now()

		err = s.dbHandler.CreateMessage(c.Request().Context(), db.CreateMessageParams{
			MessageID:        messageEntity.MessageId,
			FkParticipantsID: messageEntity.ParticipantsId,
			FkChatRoomID:     messageEntity.ChatRoomId,
			Content:          messageEntity.Content,
			LikeMessage:      messageEntity.LikeMessage,
		})
		if err != nil {
			log.Printf("error saving message: %v", err)
			continue
		}

		for conn := range s.wsConnections {
			if err := conn.WriteJSON(messageEntity); err != nil {
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

func (s *MessageHandler) LikeMessage(c echo.Context) error {
	messageId, err := uuid.Parse(c.Param("messageId"))

	if err != nil {
		c.Logger().Fatal(err)
		return err
	}

	var messageEntity entity.MessageEntity

	messageEntity.MessageId = messageId

	if err := c.Bind(&messageEntity); err != nil {
		return c.JSON(500, map[string]string{
			"error": err.Error(),
		})
	}

	msg, err := s.dbHandler.PatchLikeMessage(c.Request().Context(), db.PatchLikeMessageParams{
		MessageID:   messageEntity.MessageId,
		LikeMessage: messageEntity.LikeMessage,
	})

	if err != nil {
		return c.JSON(500, map[string]string{
			"error": err.Error(),
		})
	}

	messageEntity = entity.MessageEntity{
		MessageId:      msg.MessageID,
		LikeMessage:    msg.LikeMessage,
		Content:        msg.Content,
		ParticipantsId: msg.FkParticipantsID,
		ChatRoomId:     msg.FkChatRoomID,
		CreatedAt:      msg.CreatedAt,
	}

	return c.JSON(200, messageEntity)

}

func (s *MessageHandler) RemoveLikeMessage(c echo.Context) error {
	messageId, err := uuid.Parse(c.Param("messageId"))

	if err != nil {
		c.Logger().Fatal(err)
		return err
	}

	var messageEntity entity.MessageEntity

	messageEntity.MessageId = messageId

	if err := c.Bind(&messageEntity); err != nil {
		return c.JSON(500, map[string]string{
			"error": err.Error(),
		})
	}

	res, err := s.dbHandler.DeleteLike(c.Request().Context(), messageId)

	if err != nil {
		return c.JSON(500, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(200, res)

}
