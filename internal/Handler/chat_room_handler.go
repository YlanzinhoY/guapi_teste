package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ylanzinhoy/guapi_teste/internal/entity"
	db "github.com/ylanzinhoy/guapi_teste/sql"
)

type ChatRoomHandler struct {
	dbHandler *db.Queries
}

func NewChatRoomHandler(dbHandler *db.Queries) *ChatRoomHandler {
	return &ChatRoomHandler{
		dbHandler: dbHandler,
	}
}

func (s *ChatRoomHandler) CreateChatRoom(c echo.Context) error {

	var chatRoomEntity *entity.ChatRoomEntity

	if err := c.Bind(&chatRoomEntity); err != nil {
		return c.JSON(500, map[string]string{
			"error": err.Error(),
		})
	}

	chatRoomEntity.Id = uuid.New()

	err := s.dbHandler.CreateChatRoom(c.Request().Context(), db.CreateChatRoomParams{
		ChatRoomID:   chatRoomEntity.Id,
		ChatRoomName: chatRoomEntity.Name,
	})

	if err != nil {
		return c.JSON(400, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, chatRoomEntity)
}
