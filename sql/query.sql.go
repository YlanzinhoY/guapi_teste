// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package sql

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countMessageById = `-- name: CountMessageById :one
SELECT COUNT(*) AS message_count
FROM message
WHERE fk_chat_room_id = $1
`

func (q *Queries) CountMessageById(ctx context.Context, fkChatRoomID uuid.UUID) (int64, error) {
	row := q.db.QueryRowContext(ctx, countMessageById, fkChatRoomID)
	var message_count int64
	err := row.Scan(&message_count)
	return message_count, err
}

const createChatRoom = `-- name: CreateChatRoom :exec
/*
Chat Room
*/

INSERT INTO chat_room(
chat_room_id,
chat_room_name 
) values (
	$1, $2
)
returning chat_room_id, chat_room_name
`

type CreateChatRoomParams struct {
	ChatRoomID   uuid.UUID
	ChatRoomName int32
}

func (q *Queries) CreateChatRoom(ctx context.Context, arg CreateChatRoomParams) error {
	_, err := q.db.ExecContext(ctx, createChatRoom, arg.ChatRoomID, arg.ChatRoomName)
	return err
}

const createMessage = `-- name: CreateMessage :exec
/*
 Messages
*/
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
) RETURNING message_id, content, like_message ,created_at, fk_chat_room_id, fk_participants_id
`

type CreateMessageParams struct {
	MessageID        uuid.UUID
	FkParticipantsID uuid.UUID
	FkChatRoomID     uuid.UUID
	Content          string
	LikeMessage      int32
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) error {
	_, err := q.db.ExecContext(ctx, createMessage,
		arg.MessageID,
		arg.FkParticipantsID,
		arg.FkChatRoomID,
		arg.Content,
		arg.LikeMessage,
	)
	return err
}

const createParticipants = `-- name: CreateParticipants :exec
/*
    Participants
*/
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
) RETURNING participants_id, name, chat_room_id, is_subscribe
`

type CreateParticipantsParams struct {
	ParticipantsID uuid.UUID
	Name           string
	ChatRoomID     uuid.UUID
	IsSubscribe    bool
}

func (q *Queries) CreateParticipants(ctx context.Context, arg CreateParticipantsParams) error {
	_, err := q.db.ExecContext(ctx, createParticipants,
		arg.ParticipantsID,
		arg.Name,
		arg.ChatRoomID,
		arg.IsSubscribe,
	)
	return err
}

const createSubscribe = `-- name: CreateSubscribe :exec
/*
 Subscribe
 */


 INSERT INTO subscriber(
    subscriber_id,
    fk_participants_id,
    fk_chat_room_id
 ) VALUES($1, $2, $3)
 ON CONFLICT (fk_participants_id, fk_chat_room_id) DO NOTHING
 RETURNING subscriber_id, subscribed_at, fk_chat_room_id, fk_participants_id
`

type CreateSubscribeParams struct {
	SubscriberID     uuid.UUID
	FkParticipantsID uuid.UUID
	FkChatRoomID     uuid.UUID
}

func (q *Queries) CreateSubscribe(ctx context.Context, arg CreateSubscribeParams) error {
	_, err := q.db.ExecContext(ctx, createSubscribe, arg.SubscriberID, arg.FkParticipantsID, arg.FkChatRoomID)
	return err
}

const deleteChatRoom = `-- name: DeleteChatRoom :one
DELETE FROM chat_room
WHERE chat_room_id = $1
AND NOT EXISTS (
    SELECT 1 FROM message WHERE message.fk_chat_room_id = chat_room.chat_room_id
)
RETURNING chat_room_id, chat_room_name
`

func (q *Queries) DeleteChatRoom(ctx context.Context, chatRoomID uuid.UUID) (ChatRoom, error) {
	row := q.db.QueryRowContext(ctx, deleteChatRoom, chatRoomID)
	var i ChatRoom
	err := row.Scan(&i.ChatRoomID, &i.ChatRoomName)
	return i, err
}

const deleteLike = `-- name: DeleteLike :one
UPDATE message
SET like_message = like_message - 1
WHERE message_id = $1
AND like_message > 0
RETURNING message_id, content, like_message ,created_at, fk_chat_room_id, fk_participants_id
`

type DeleteLikeRow struct {
	MessageID        uuid.UUID
	Content          string
	LikeMessage      int32
	CreatedAt        time.Time
	FkChatRoomID     uuid.UUID
	FkParticipantsID uuid.UUID
}

func (q *Queries) DeleteLike(ctx context.Context, messageID uuid.UUID) (DeleteLikeRow, error) {
	row := q.db.QueryRowContext(ctx, deleteLike, messageID)
	var i DeleteLikeRow
	err := row.Scan(
		&i.MessageID,
		&i.Content,
		&i.LikeMessage,
		&i.CreatedAt,
		&i.FkChatRoomID,
		&i.FkParticipantsID,
	)
	return i, err
}

const findAllParticipantsSubscribers = `-- name: FindAllParticipantsSubscribers :many
/*
 Notification queries
 */

 SELECT participants_id, name, is_subscribe, chat_room_id FROM participants
 where is_subscribe = true AND chat_room_id = $1
`

func (q *Queries) FindAllParticipantsSubscribers(ctx context.Context, chatRoomID uuid.UUID) ([]Participant, error) {
	rows, err := q.db.QueryContext(ctx, findAllParticipantsSubscribers, chatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Participant
	for rows.Next() {
		var i Participant
		if err := rows.Scan(
			&i.ParticipantsID,
			&i.Name,
			&i.IsSubscribe,
			&i.ChatRoomID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMessageLikes = `-- name: GetMessageLikes :one
SELECT like_message FROM message
`

func (q *Queries) GetMessageLikes(ctx context.Context) (int32, error) {
	row := q.db.QueryRowContext(ctx, getMessageLikes)
	var like_message int32
	err := row.Scan(&like_message)
	return like_message, err
}

const getMessagesLikesByChatId = `-- name: GetMessagesLikesByChatId :many
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
    m.fk_chat_room_id = $1
`

type GetMessagesLikesByChatIdRow struct {
	MessageID   uuid.UUID
	LikeMessage int32
	Content     string
	Name        string
}

func (q *Queries) GetMessagesLikesByChatId(ctx context.Context, fkChatRoomID uuid.UUID) ([]GetMessagesLikesByChatIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesLikesByChatId, fkChatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMessagesLikesByChatIdRow
	for rows.Next() {
		var i GetMessagesLikesByChatIdRow
		if err := rows.Scan(
			&i.MessageID,
			&i.LikeMessage,
			&i.Content,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const patchLikeMessage = `-- name: PatchLikeMessage :one
UPDATE message
    set like_message = like_message + 1
WHERE message_id = $1
RETURNING message_id, content, like_message ,created_at, fk_chat_room_id, fk_participants_id
`

type PatchLikeMessageRow struct {
	MessageID        uuid.UUID
	Content          string
	LikeMessage      int32
	CreatedAt        time.Time
	FkChatRoomID     uuid.UUID
	FkParticipantsID uuid.UUID
}

func (q *Queries) PatchLikeMessage(ctx context.Context, messageID uuid.UUID) (PatchLikeMessageRow, error) {
	row := q.db.QueryRowContext(ctx, patchLikeMessage, messageID)
	var i PatchLikeMessageRow
	err := row.Scan(
		&i.MessageID,
		&i.Content,
		&i.LikeMessage,
		&i.CreatedAt,
		&i.FkChatRoomID,
		&i.FkParticipantsID,
	)
	return i, err
}

const updateParticipantSubscription = `-- name: UpdateParticipantSubscription :exec
UPDATE participants
    set is_subscribe = true
WHERE participants_id = $1
`

func (q *Queries) UpdateParticipantSubscription(ctx context.Context, participantsID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, updateParticipantSubscription, participantsID)
	return err
}
