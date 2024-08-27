CREATE TABLE chat_room(
chat_room_id uuid DEFAULT uuid_generate_v4(),
chat_room_name int UNIQUE NOT NULL,

PRIMARY KEY(chat_room_id)
);


CREATE TABLE participants(
participants_id uuid DEFAULT uuid_generate_v4(),
name varchar(255) NOT NULL,

PRIMARY KEY(participants_id),
chat_room_id uuid REFERENCES chat_room(chat_room_id) ON DELETE CASCADE NOT NULL
);