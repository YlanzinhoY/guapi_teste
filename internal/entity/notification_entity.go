package entity

import (
	"time"

	"github.com/google/uuid"
)

type NotificationEntity struct {
	NotificationID  uuid.UUID `json:"notification_id"`
	CreatedAt       time.Time `json:"created_at"`
	Message         string    `json:"message_notification"`
	FKParticipantID uuid.UUID `json:"participant_id"`
	FkChatRoomID    uuid.UUID `json:"chat_room_id"`
	FkMessageID     uuid.UUID `json:"message_id"`
}
