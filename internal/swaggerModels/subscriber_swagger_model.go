package swaggerModels

import (
	"github.com/google/uuid"
)

type SubscribeSwaggerModel struct {
	FkParticipantsID uuid.UUID `json:"participant_id"`
	FkChatRoomID     uuid.UUID `json:"chat_room_id"`
}
