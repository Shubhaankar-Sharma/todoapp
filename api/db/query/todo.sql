-- name: CreateToDo :exec
INSERT INTO todo (body, end_date, user_id, done)
VALUES ($1, $2, $3, $4);

-- name: GetToDoAll :many
SELECT *
FROM todo
WHERE user_id = $1
ORDER BY end_date;

-- name: DeleteToDo :exec
DELETE
FROM todo
WHERE id = $1;

-- name: PatchToDo :exec
UPDATE todo
SET body=$1,
    end_date=$2
WHERE id = $3;