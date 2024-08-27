/*
Chat Room
*/

-- name: CreateChatRoom :exec
INSERT INTO chat_room(
chat_room_id,
chat_room_name 
) values (
	$1, $2
)
returning *;

-- name: DeleteChatRoom :one
DELETE FROM chat_room
WHERE chat_room_id = $1
AND NOT EXISTS (
    SELECT 1 FROM message WHERE message.fk_chat_room_id = chat_room.chat_room_id
)
RETURNING *;

/*
    Participants
*/
-- name: CreateParticipants :exec
INSERT INTO participants(
participants_id,
name,
chat_room_id
) VALUES (
$1,
$2,
$3
) RETURNING participants_id, name, chat_room_id;

/*
 Messages
*/
-- name: CreateMessage :exec
INSERT INTO message(
message_id,
fk_participants_id,
fk_chat_room_id,
content
) VALUES(
    $1,
    $2,
    $3,
    $4
) RETURNING message_id, content, created_at, fk_chat_room_id, fk_participants_id;