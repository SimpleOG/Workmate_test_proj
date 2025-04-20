-- name: CreateUser :one
INSERT INTO users (username,
                   password)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserForLogin :one
SELECT *
FROM users
WHERE username = $1
  and password = $2;


-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

