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
content,
like_message
) VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING message_id, content, like_message ,created_at, fk_chat_room_id, fk_participants_id;

-- name: PatchLikeMessage :one
UPDATE message
    set like_message = $2
WHERE message_id = $1
RETURNING message_id, content, like_message ,created_at, fk_chat_room_id, fk_participants_id;

-- name: DeleteLike :one
UPDATE message
SET like_message = like_message - 1
WHERE message_id = $1
AND like_message > 0
RETURNING message_id, content, like_message ,created_at, fk_chat_room_id, fk_participants_id;

/*
 Subscribe
 */

-- name: CreateSubscribe :exec
 INSERT INTO subscriber(
    subscriber_id,
    fk_participants_id,
    fk_chat_room_id
 ) VALUES($1, $2, $3)
 ON CONFLICT (fk_participants_id, fk_chat_room_id) DO NOTHING
 RETURNING subscriber_id, subscribed_at, fk_chat_room_id, fk_participants_id;


 /*
    notification
 */

-- name: CreateNotification :exec

INSERT INTO notification(
notification_id,
fk_chat_room_id,
fk_participant_id,
ping,
fk_message_id
) VALUES(
$1, $2, $3, $4, $5
) RETURNING notification_id, fk_chat_room_id, fk_participant_id, ping, fk_message_id;