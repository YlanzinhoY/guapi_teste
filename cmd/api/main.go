package main

import (
	"github.com/gorilla/websocket"
	handler "github.com/ylanzinhoy/guapi_teste/internal/Handler"
	"github.com/ylanzinhoy/guapi_teste/internal/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger" // echo-swagger middleware
	_ "github.com/ylanzinhoy/guapi_teste/docs"
)

var (
	ws websocket.Upgrader
)

// @title Guapi Teste API
// @version 1.0
func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	r := repository.NewDatabaseRepository("host=localhost port=5432 user=postgres password=postgres dbname=guapi_teste sslmode=disable")
	defer r.DatabaseConn()
	ws.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	chatRoomHandler := handler.NewChatRoomHandler(r.DbHandler())
	e.POST("/v1/chatroom", chatRoomHandler.CreateChatRoom)
	e.DELETE("/v1/chatroom/:id", chatRoomHandler.DeleteChatRoom)

	participatnsHandler := handler.NewParticipantsHandler(r.DbHandler())
	e.POST("/v1/participants", participatnsHandler.CreateParticipants)

	messageHandler := handler.NewMessageHandler(r.DbHandler(), &ws, make(map[*websocket.Conn]bool))
	e.GET("/ws/notifier/message/:chatRoomId/:participantId", messageHandler.CreateMessageWS)
	e.PATCH("/v1/message/like/:messageId", messageHandler.LikeMessage)
	e.DELETE("/v1/message/dislike/:messageId", messageHandler.RemoveLikeMessage)

	subscribeHandler := handler.NewSubscribeHandler(r.DbHandler())
	e.POST("/v1/subscribe", subscribeHandler.CreateSubscribeInChatRoom)

	notificationHandler := handler.NewNotificationHandler(r.DbHandler(), &ws, make(map[*websocket.Conn]bool))

	e.GET("/ws/notifier/message/:chat_room_id", notificationHandler.SendNotificationMessage)
	e.GET("/ws/notifier/likesUnlike/:chat_room_id", notificationHandler.SendNotificationLikeUnLikeMessage)

	e.Logger.Fatal(e.Start(":9001"))
}
