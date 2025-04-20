-- name: CreateTask :one
INSERT INTO tasks (
                   id,
                   status,
                   payload
                   )
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTask :one
SELECT *
FROM tasks
WHERE id = $1;

-- name: UpdateTaskStatus :exec
UPDATE tasks
SET
    status  = $2
WHERE id = $1;

-- name: UpdateTaskResult :exec
UPDATE tasks
SET
    result  = $2
WHERE id = $1;

-- name: DeleteTask :exec
DELETE from tasks where id=$1;