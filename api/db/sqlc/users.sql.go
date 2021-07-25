// Code generated by sqlc. DO NOT EDIT.
// source: users.sql

package db

import (
	"context"
)

const changeUser = `-- name: ChangeUser :exec
UPDATE users
SET username=$1
WHERE id=$2
`

type ChangeUserParams struct {
	Username string `json:"username"`
	ID       int32  `json:"id"`
}

func (q *Queries) ChangeUser(ctx context.Context, arg ChangeUserParams) error {
	_, err := q.db.ExecContext(ctx, changeUser, arg.Username, arg.ID)
	return err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING id
`

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Password)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id=$1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, username, password FROM users WHERE username=$1
`

func (q *Queries) GetUser(ctx context.Context, username string) (Users, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i Users
	err := row.Scan(&i.ID, &i.Username, &i.Password)
	return i, err
}
