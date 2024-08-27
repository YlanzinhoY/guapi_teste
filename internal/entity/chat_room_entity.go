package entity

import "github.com/google/uuid"

type ChatRoomEntity struct {
	Id   uuid.UUID `json:"chat_room_id"`
	Name int32     `json:"chat_room_name"`
}
