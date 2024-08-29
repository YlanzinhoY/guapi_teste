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
chat_room_id,
is_subscribe
) VALUES (
$1,
$2,
$3,
$4
) RETURNING participants_id, name, chat_room_id, is_subscribe;

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
    set like_message = like_message + 1
WHERE message_id = $1
RETURNING message_id, content, like_message ,created_at, fk_chat_room_id, fk_participants_id;

-- name: DeleteLike :one
UPDATE message
SET like_message = like_message - 1
WHERE message_id = $1
AND like_message > 0
RETURNING message_id, content, like_message ,created_at, fk_chat_room_id, fk_participants_id;

-- name: GetMessageLikes :one
SELECT like_message FROM message;
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

-- name: UpdateParticipantSubscription :exec
UPDATE participants
    set is_subscribe = true
WHERE participants_id = $1;



/*
 Notification queries
 */

-- name: FindAllParticipantsSubscribers :many
 SELECT * FROM participants
 where is_subscribe = true AND chat_room_id = $1;

-- name: CountMessageById :one
SELECT COUNT(*) AS message_count
FROM message
WHERE fk_chat_room_id = $1;


-- name: GetMessagesLikesByChatId :many
SELECT
    m.message_id,
    m.like_message,
    m.content,
    p.name
FROM
    message AS m
        JOIN
    participants AS p
    ON
        m.fk_participants_id = p.participants_id
WHERE
    m.fk_chat_room_id = $1;