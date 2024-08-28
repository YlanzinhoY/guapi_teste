package entity

import (
	"time"

	"github.com/google/uuid"
)

type MessageEntity struct {
	MessageId      uuid.UUID `json:"message_id"`
	Content        string    `json:"content"`
	ParticipantsId uuid.UUID `json:"participant_id"`
	ChatRoomId     uuid.UUID `json:"chat_room_id"`
	CreatedAt      time.Time `json:"created_at"`
	LikeMessage    int32     `json:"like_message"`
}
