package entity

import "github.com/google/uuid"

type ParticipantsEntity struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"participant_name"`
	ChatRoomId uuid.UUID `json:"chat_room_id"`
}
