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

-- name: GetChatRoomById :one
SELECT * FROM chat_room
WHERE chat_room_id = $1;


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
participants_id,
chat_room_id,
content
) VALUES(
    $1,
    $2,
    $3,
    $4
) RETURNING message_id, content, created_at, chat_room_id, participants_id;