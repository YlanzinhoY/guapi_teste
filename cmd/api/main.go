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

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v1/
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
	e.GET("/ws/message/:chatRoomId/:participantId", messageHandler.CreateMessageWS)
	e.PATCH("/v1/message/:messageId", messageHandler.LikeMessage)
	e.DELETE("/v1/message/:messageId", messageHandler.RemoveLikeMessage)

	subscribeHandler := handler.NewSubscribeHandler(r.DbHandler())
	e.POST("/v1/subscribe", subscribeHandler.CreateSubscribeInChatRoom)

	notificationHandler := handler.NewNotificationHandler(r.DbHandler(), &ws, make(map[*websocket.Conn]bool))

	e.GET("/ws/notifier/:chat_room_id", notificationHandler.SendNotification)

	e.Logger.Fatal(e.Start(":9001"))
}
