package entity

import (
	"github.com/google/uuid"
	"time"
)

type SubscribeEntity struct {
	SubscriberID     uuid.UUID `json:"subscriber_id"`
	SubscribedAt     time.Time `json:"subscribed_at"`
	FkParticipantsID uuid.UUID `json:"participants_id"`
	FkChatRoomID     uuid.UUID `json:"chat_room_id"`
}
