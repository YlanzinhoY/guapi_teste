package main

import (
	"github.com/gorilla/websocket"
	handler "github.com/ylanzinhoy/guapi_teste/internal/Handler"
	"github.com/ylanzinhoy/guapi_teste/internal/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	ws websocket.Upgrader
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	r := repository.NewDatabaseRepository("host=localhost port=5432 user=postgres password=postgres dbname=guapi_teste sslmode=disable")
	defer r.DatabaseConn()

	chatRoomHandler := handler.NewChatRoomHandler(r.DbHandler())
	e.POST("/v1/chatroom", chatRoomHandler.CreateChatRoom)
	e.DELETE("/v1/chatroom/:id", chatRoomHandler.DeleteChatRoom)

	participatnsHandler := handler.NewParticipantsHandler(r.DbHandler())
	e.POST("/v1/participants", participatnsHandler.CreateParticipants)

	messageHandler := handler.NewMessageHandler(r.DbHandler(), &ws, make(map[*websocket.Conn]bool))
	e.GET("/ws/:chatRoomId/:participantId", messageHandler.CreateMessageWS)
	e.PATCH("/v1/message/:messageId", messageHandler.LikeMessage)
	e.DELETE("/v1/message/:messageId", messageHandler.RemoveLikeMessage)
	e.Logger.Fatal(e.Start(":9001"))

}
