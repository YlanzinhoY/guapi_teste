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
