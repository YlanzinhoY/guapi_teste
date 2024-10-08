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
type requestBodySwagger struct {
	ChatRoomName int32 `json:"chat_room_name"`
}

type ChatRoomHandler struct {
	DbHandler *db.Queries
}

func NewChatRoomHandler(dbHandler *db.Queries) *ChatRoomHandler {

	return &ChatRoomHandler{
		DbHandler: dbHandler,
	}
}

// ShowAccount godoc
// @Summary      Create Chat room
// @Description  Create Chat room
// @Tags         Room
// @Accept       json
// @Produce      json
// @Param        request body requestBodySwagger true "user request"
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

	err := s.DbHandler.CreateChatRoom(c.Request().Context(), db.CreateChatRoomParams{
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

// ShowAccount godoc
// @Summary      Delete Chat Room
// @Description  Delete Chat Room
// @Tags         Room
// @Accept       json
// @Produce      json
// @Param   	id path string true "chat_room_id"
// @Success      200
// @Failure      500  {object}  Error
// @Router       /v1/chatroom/{id} [delete]
func (s *ChatRoomHandler) DeleteChatRoom(c echo.Context) error {
	idParam := uuid.MustParse(c.Param("id"))
	res, err := s.DbHandler.DeleteChatRoom(c.Request().Context(), idParam)

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
