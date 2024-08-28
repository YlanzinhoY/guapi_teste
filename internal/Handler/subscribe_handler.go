package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ylanzinhoy/guapi_teste/internal/entity"
	db "github.com/ylanzinhoy/guapi_teste/sql"
)

type SubscribeHandler struct {
	dbHandler *db.Queries
}

func NewSubscribeHandler(dbHandler *db.Queries) *SubscribeHandler {
	return &SubscribeHandler{dbHandler: dbHandler}
}

func (s *SubscribeHandler) CreateSubscribeInChatRoom(c echo.Context) error {

	var subscribeEntity entity.SubscribeEntity

	participantId, err := uuid.Parse(c.Param("participant_id"))
	if err != nil {
		c.JSON(404, err.Error())
	}
	chatRoomId, err := uuid.Parse(c.Param("chat_room_id"))
	if err != nil {
		c.JSON(404, err.Error())
	}

	if err := c.Bind(&subscribeEntity); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	subscribeEntity.SubscriberID = uuid.New()
	subscribeEntity.SubscribedAt = time.Now()
	subscribeEntity.FkChatRoomID = chatRoomId
	subscribeEntity.FkParticipantsID = participantId

	err = s.dbHandler.CreateSubscribe(c.Request().Context(), db.CreateSubscribeParams{
		SubscriberID:     subscribeEntity.SubscriberID,
		FkChatRoomID:     subscribeEntity.FkChatRoomID,
		FkParticipantsID: subscribeEntity.FkParticipantsID,
	})

	if err != nil {
		return c.JSON(400, err.Error())
	}

	return c.JSON(http.StatusCreated, subscribeEntity)
}
