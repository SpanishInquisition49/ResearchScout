-- name: CreateCall :one
INSERT OR IGNORE INTO calls (title, deadline, requirements, apply_module)
VALUES (
    sqlc.arg(title),
    sqlc.arg(deadline),
    sqlc.arg(requirements),
    sqlc.arg(apply_module)
) RETURNING id;

-- name: CreateUser :one
INSERT INTO bot_users (chat_id, first_interaction)
VALUES (
    sqlc.arg(chat_id),
    sqlc.arg(first_interaction)
)
ON CONFLICT(chat_id) DO UPDATE SET
    first_interaction = excluded.first_interaction,
    is_active = 1
RETURNING chat_id;

-- name: DeactivateUser :one
UPDATE bot_users
SET is_active = false
WHERE chat_id = sqlc.arg(chat_id)
RETURNING chat_id;

-- Create the relationship between a user and a call
-- if the user has already received the call, it will not be added again
-- name: SendCallToUser :many
INSERT INTO users_calls (user_chat_id, call_id)
VALUES (
    sqlc.arg(user_chat_id),
    sqlc.arg(call_id)
) RETURNING user_chat_id, call_id;

-- Get all the calls that a user has not yet received
-- name: GetCallsToSend :many
SELECT c.id, c.title, c.deadline, c.requirements, c.apply_module
FROM calls c
WHERE
    c.id
    NOT IN (SELECT call_id FROM users_calls WHERE user_chat_id = sqlc.arg(user_chat_id))
;

-- name: GetUsers :many
SELECT chat_id, first_interaction
FROM bot_users
WHERE is_active = 1
;

-- Remove all the calls with a past deadline
-- name: DeleteOlderCalls :many
DELETE
FROM calls
WHERE date('now') > date(deadline)
RETURNING id;

-- Remove all the relationship between users and all the past calls
-- name: CleanNotificationsHistory :exec
DELETE FROM users_calls
WHERE call_id IN (
  SELECT id
  FROM calls
  WHERE date('now') > date(deadline)
);
