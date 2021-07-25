-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING id;

-- name: GetUser :one
SELECT * FROM users WHERE username=$1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id=$1;

-- name: ChangeUser :exec
UPDATE users
SET username=$1
WHERE id=$2;