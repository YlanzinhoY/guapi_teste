package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/ylanzinhoy/guapi_teste/internal/entity"
	db "github.com/ylanzinhoy/guapi_teste/sql"
)

type SubscribeHandler struct {
	dbHandler *db.Queries
}

func NewSubscribeHandler(dbHandler *db.Queries) *SubscribeHandler {
	return &SubscribeHandler{dbHandler: dbHandler}
}

// ShowAccount godoc
// @Summary      Create subscription
// @Description  Create subscription
// @Tags         createSubscribe
// @Accept       json
// @Produce      json
// @Param        request body swaggerModels.SubscribeSwaggerModel true "user request"
// @Success      201
// @Failure      500  {object}  Error
// @Router       /v1/subscribe [post]
func (s *SubscribeHandler) CreateSubscribeInChatRoom(c echo.Context) error {

	var subscribeEntity entity.SubscribeEntity

	if err := c.Bind(&subscribeEntity); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	subscribeEntity.SubscriberID = uuid.New()
	subscribeEntity.SubscribedAt = time.Now()

	err := s.dbHandler.CreateSubscribe(c.Request().Context(), db.CreateSubscribeParams{
		SubscriberID:     subscribeEntity.SubscriberID,
		FkChatRoomID:     subscribeEntity.FkChatRoomID,
		FkParticipantsID: subscribeEntity.FkParticipantsID,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = s.dbHandler.UpdateParticipantSubscription(c.Request().Context(), subscribeEntity.FkParticipantsID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, subscribeEntity)
}
