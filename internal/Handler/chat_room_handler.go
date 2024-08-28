package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ylanzinhoy/guapi_teste/internal/entity"
	db "github.com/ylanzinhoy/guapi_teste/sql"
)

type Error struct {
	Message string `json:"message"`
}

type ChatRoomHandler struct {
	dbHandler *db.Queries
}

func NewChatRoomHandler(dbHandler *db.Queries) *ChatRoomHandler {

	return &ChatRoomHandler{
		dbHandler: dbHandler,
	}
}

// ShowAccount godoc
// @Summary      Create Chat room
// @Description  Create Chat room
// @Tags         createRoom
// @Accept       json
// @Produce      json
// @Param        request body entity.ChatRoomEntity true "user request"
// @Success      201 
// @Failure      500  {object}  Error
// @Router       /v1/chatroom [post]
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

func (s *ChatRoomHandler) DeleteChatRoom(c echo.Context) error {
	idParam := uuid.MustParse(c.Param("id"))
	res, err := s.dbHandler.DeleteChatRoom(c.Request().Context(), idParam)

	if err != nil {
		return c.JSON(404, map[string]string{
			"error": "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"msg":  "deleted",
		"data": res,
	})

}
