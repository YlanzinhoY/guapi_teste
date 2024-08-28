CREATE TABLE IF NOT EXISTS subscriber(
subscriber_id uuid NOT NULL DEFAULT uuid_generate_v4(),
subscribed_at timestamp NOT NULL DEFAULT now(),


PRIMARY KEY(subscriber_id),
fk_participants_id uuid REFERENCES participants(participants_id) ON DELETE CASCADE,
fk_chat_room_id uuid  REFERENCES chat_room(chat_room_id) ON DELETE CASCADE,
UNIQUE(fk_participants_id, fk_chat_room_id)
);