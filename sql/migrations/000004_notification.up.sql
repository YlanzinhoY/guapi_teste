CREATE TABLE IF NOT EXISTS notification(
notification_id uuid NOT NULL DEFAULT uuid_generate_v4(),
created_at timestamp NOT NULL DEFAULT now(),
message varchar(100) NOT NULL,

PRIMARY KEY(notification_id),
    fk_chat_room_id uuid NOT NULL REFERENCES chat_room(chat_room_id) ON DELETE CASCADE,
    fk_message_id uuid NOT NULL REFERENCES message(message_id) ON DELETE CASCADE
);