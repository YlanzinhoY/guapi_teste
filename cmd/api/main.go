package main

import (
	"database/sql"
	"log"

	handler "github.com/ylanzinhoy/guapi_teste/internal/Handler"
	db "github.com/ylanzinhoy/guapi_teste/sql"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {

	e := echo.New()

	connString := "host=localhost port=5432 user=postgres password=postgres dbname=guapi_teste sslmode=disable"

	dbConn, err := sql.Open("postgres", connString)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer func(dbConn *sql.DB) {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}(dbConn)

	query := db.New(dbConn)

	chatRoomHandler := handler.NewChatRoomHandler(query)
	e.POST("/v1/chatroom", chatRoomHandler.CreateChatRoom)

	participatnsHandler := handler.NewParticipantsHandler(query)
	e.POST("/v1/participants", participatnsHandler.CreateParticipants)

	e.Logger.Fatal(e.Start(":9001"))
}
