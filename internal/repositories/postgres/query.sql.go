// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: query.sql

package postgres

import (
	"context"
	"database/sql"
)

const createBlog = `-- name: CreateBlog :exec
INSERT INTO blog(name, description)
VALUES ($1, $2)
`

type CreateBlogParams struct {
	Name        sql.NullString
	Description sql.NullString
}

func (q *Queries) CreateBlog(ctx context.Context, arg CreateBlogParams) error {
	_, err := q.db.Exec(ctx, createBlog, arg.Name, arg.Description)
	return err
}

const createContent = `-- name: CreateContent :exec
INSERT INTO content(user_email, name, description)
VALUES ($1, $2, $3)
`

type CreateContentParams struct {
	UserEmail   sql.NullString
	Name        sql.NullString
	Description sql.NullString
}

func (q *Queries) CreateContent(ctx context.Context, arg CreateContentParams) error {
	_, err := q.db.Exec(ctx, createContent, arg.UserEmail, arg.Name, arg.Description)
	return err
}

const createSession = `-- name: CreateSession :exec
INSERT INTO sessions(user_email, refresh_token, access_token, user_agent, ip, expires_in, created_at)
VALUES ($1, $2, $3, $4, $5, $6, DEFAULT)
`

type CreateSessionParams struct {
	UserEmail    sql.NullString
	RefreshToken sql.NullString
	AccessToken  sql.NullString
	UserAgent    sql.NullString
	Ip           sql.NullString
	ExpiresIn    int64
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) error {
	_, err := q.db.Exec(ctx, createSession,
		arg.UserEmail,
		arg.RefreshToken,
		arg.AccessToken,
		arg.UserAgent,
		arg.Ip,
		arg.ExpiresIn,
	)
	return err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users(name, password, email, role)
VALUES ($1, $2, $3, $4)
`

type CreateUserParams struct {
	Name     sql.NullString
	Password string
	Email    string
	Role     string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser,
		arg.Name,
		arg.Password,
		arg.Email,
		arg.Role,
	)
	return err
}

const deleteContent = `-- name: DeleteContent :exec
DELETE FROM content
WHERE user_email=$1 and name=$2
`

type DeleteContentParams struct {
	UserEmail sql.NullString
	Name      sql.NullString
}

func (q *Queries) DeleteContent(ctx context.Context, arg DeleteContentParams) error {
	_, err := q.db.Exec(ctx, deleteContent, arg.UserEmail, arg.Name)
	return err
}

const deleteSession = `-- name: DeleteSession :exec
DELETE FROM sessions
WHERE user_agent=$1 and ip=$2 and user_email=$3
`

type DeleteSessionParams struct {
	UserAgent sql.NullString
	Ip        sql.NullString
	UserEmail sql.NullString
}

func (q *Queries) DeleteSession(ctx context.Context, arg DeleteSessionParams) error {
	_, err := q.db.Exec(ctx, deleteSession, arg.UserAgent, arg.Ip, arg.UserEmail)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
UPDATE users SET is_deleted=true
WHERE email = $1
`

func (q *Queries) DeleteUser(ctx context.Context, email string) error {
	_, err := q.db.Exec(ctx, deleteUser, email)
	return err
}

const readBlog = `-- name: ReadBlog :one
SELECT id, name, description FROM blog
WHERE id=$1
`

func (q *Queries) ReadBlog(ctx context.Context, id int32) (Blog, error) {
	row := q.db.QueryRow(ctx, readBlog, id)
	var i Blog
	err := row.Scan(&i.ID, &i.Name, &i.Description)
	return i, err
}

const readContent = `-- name: ReadContent :one
SELECT id, user_email, name, description FROM content
WHERE user_email=$1 and id=$2
`

type ReadContentParams struct {
	UserEmail sql.NullString
	ID        int32
}

func (q *Queries) ReadContent(ctx context.Context, arg ReadContentParams) (Content, error) {
	row := q.db.QueryRow(ctx, readContent, arg.UserEmail, arg.ID)
	var i Content
	err := row.Scan(
		&i.ID,
		&i.UserEmail,
		&i.Name,
		&i.Description,
	)
	return i, err
}

const readEmailRoleFromSessions = `-- name: ReadEmailRoleFromSessions :one
SELECT email, role, s.expires_in
FROM users INNER JOIN sessions s on users.email = s.user_email
WHERE s.refresh_token = $1
`

type ReadEmailRoleFromSessionsRow struct {
	Email     string
	Role      string
	ExpiresIn int64
}

func (q *Queries) ReadEmailRoleFromSessions(ctx context.Context, refreshToken sql.NullString) (ReadEmailRoleFromSessionsRow, error) {
	row := q.db.QueryRow(ctx, readEmailRoleFromSessions, refreshToken)
	var i ReadEmailRoleFromSessionsRow
	err := row.Scan(&i.Email, &i.Role, &i.ExpiresIn)
	return i, err
}

const readSession = `-- name: ReadSession :one
SELECT id, user_email, refresh_token, access_token, user_agent, ip, expires_in, created_at FROM sessions
WHERE user_email=$1 and  user_agent=$2 and ip=$3
`

type ReadSessionParams struct {
	UserEmail sql.NullString
	UserAgent sql.NullString
	Ip        sql.NullString
}

func (q *Queries) ReadSession(ctx context.Context, arg ReadSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, readSession, arg.UserEmail, arg.UserAgent, arg.Ip)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserEmail,
		&i.RefreshToken,
		&i.AccessToken,
		&i.UserAgent,
		&i.Ip,
		&i.ExpiresIn,
		&i.CreatedAt,
	)
	return i, err
}

const readUser = `-- name: ReadUser :one
SELECT id, name, password, email, role, is_deleted, verifay FROM users
WHERE email = $1
`

func (q *Queries) ReadUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, readUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Password,
		&i.Email,
		&i.Role,
		&i.IsDeleted,
		&i.Verifay,
	)
	return i, err
}

const updateContent = `-- name: UpdateContent :exec
UPDATE content SET user_email=$1, name=$3, description=$4
WHERE user_email= $1 and name =$2
`

type UpdateContentParams struct {
	UserEmail   sql.NullString
	Name        sql.NullString
	Name_2      sql.NullString
	Description sql.NullString
}

func (q *Queries) UpdateContent(ctx context.Context, arg UpdateContentParams) error {
	_, err := q.db.Exec(ctx, updateContent,
		arg.UserEmail,
		arg.Name,
		arg.Name_2,
		arg.Description,
	)
	return err
}

const updateSession = `-- name: UpdateSession :exec
UPDATE sessions SET refresh_token=$3, expires_in=$4, created_at=DEFAULT
WHERE user_email= $1 and refresh_token=$2
`

type UpdateSessionParams struct {
	UserEmail      sql.NullString
	RefreshToken   sql.NullString
	RefreshToken_2 sql.NullString
	ExpiresIn      int64
}

func (q *Queries) UpdateSession(ctx context.Context, arg UpdateSessionParams) error {
	_, err := q.db.Exec(ctx, updateSession,
		arg.UserEmail,
		arg.RefreshToken,
		arg.RefreshToken_2,
		arg.ExpiresIn,
	)
	return err
}

const updateSessionTokenOnly = `-- name: UpdateSessionTokenOnly :exec
UPDATE sessions SET refresh_token=$2, expires_in=$3, created_at=DEFAULT
WHERE refresh_token=$1
`

type UpdateSessionTokenOnlyParams struct {
	RefreshToken   sql.NullString
	RefreshToken_2 sql.NullString
	ExpiresIn      int64
}

func (q *Queries) UpdateSessionTokenOnly(ctx context.Context, arg UpdateSessionTokenOnlyParams) error {
	_, err := q.db.Exec(ctx, updateSessionTokenOnly, arg.RefreshToken, arg.RefreshToken_2, arg.ExpiresIn)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users SET name=$2, password=$3, role=$4, is_deleted=$5
WHERE email = $1
`

type UpdateUserParams struct {
	Email     string
	Name      sql.NullString
	Password  string
	Role      string
	IsDeleted sql.NullBool
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.Email,
		arg.Name,
		arg.Password,
		arg.Role,
		arg.IsDeleted,
	)
	return err
}