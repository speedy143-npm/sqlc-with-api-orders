-- name: CreateMessage :one
INSERT INTO message (thread, sender, content)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetMessageByID :one
SELECT * FROM message
WHERE id = $1;

-- name: GetMessagesByThread :many
SELECT * FROM message
WHERE thread = $1
ORDER BY created_at DESC;


-- creating a thread
-- name: CreateThread :one
INSERT INTO thread (thread)
VALUES ($1)
RETURNING *;