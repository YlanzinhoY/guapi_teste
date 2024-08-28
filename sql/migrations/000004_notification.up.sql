CREATE TABLE IF NOT EXISTS notification(
notification_id uuid NOT NULL DEFAULT uuid_generate_v4(),
created_at timestamp NOT NULL DEFAULT now(),
ping INT DEFAULT 0 NOT NULL,
is_read boolean NOT NULL DEFAULT false,

PRIMARY KEY(notification_id),
    fk_user_id uuid REFERENCES participants(participants_id) ON DELETE CASCADE,
    fk_chat_room_id uuid REFERENCES chat_room(chat_room_id) ON DELETE CASCADE,
    fk_message_id uuid REFERENCES message(message_id) ON DELETE CASCADE
);