package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ylanzinhoy/guapi_teste/internal/entity"
	db "github.com/ylanzinhoy/guapi_teste/sql"
)

type ParticipantsHandler struct {
	dbHandler *db.Queries
}

func NewParticipantsHandler(dbHandler *db.Queries) *ParticipantsHandler {
	return &ParticipantsHandler{
		dbHandler: dbHandler,
	}
}

func (s *ParticipantsHandler) CreateParticipants(c echo.Context) error {
	var participantsEntity *entity.ParticipantsEntity

	if err := c.Bind(&participantsEntity); err != nil {
		return c.JSON(400, map[string]string{
			"error": err.Error(),
		})
	}

	participantsEntity.Id = uuid.New()
	participantsEntity.IsSubscribe = false

	err := s.dbHandler.CreateParticipants(c.Request().Context(), db.CreateParticipantsParams{
		ParticipantsID: participantsEntity.Id,
		Name:           participantsEntity.Name,
		ChatRoomID:     participantsEntity.ChatRoomId,
		IsSubscribe:    participantsEntity.IsSubscribe,
	})

	if err != nil {

		return c.JSON(500, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &participantsEntity)

}
