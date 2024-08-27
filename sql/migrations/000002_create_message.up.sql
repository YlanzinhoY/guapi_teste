CREATE TABLE IF NOT EXISTS message(
message_id uuid DEFAULT uuid_generate_v4(),
content text NOT NULL,
created_at timestamp NOT NULL DEFAULT now(),

PRIMARY KEY(message_id),
participants_id uuid REFERENCES participants(participants_id) ON DELETE CASCADE,
chat_room_id uuid REFERENCES chat_room(chat_room_id) ON DELETE CASCADE
);