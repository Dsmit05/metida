-- name: CreateUser :exec
INSERT INTO users(name, password, email, role)
VALUES ($1, $2, $3, $4);

-- name: UpdateUser :exec
UPDATE users SET name=$2, password=$3, role=$4, is_deleted=$5
WHERE email = $1;

-- name: DeleteUser :exec
UPDATE users SET is_deleted=true
WHERE email = $1;

-- name: ReadUser :one
SELECT * FROM users
WHERE email = $1;

-- name: CreateSession :exec
INSERT INTO sessions(user_email, refresh_token, access_token, user_agent, ip, expires_in, created_at)
VALUES ($1, $2, $3, $4, $5, $6, DEFAULT);

-- name: UpdateSession :exec
UPDATE sessions SET refresh_token=$3, expires_in=$4, created_at=DEFAULT
WHERE user_email= $1 and refresh_token=$2;

-- name: UpdateSessionTokenOnly :exec
UPDATE sessions SET refresh_token=$2, expires_in=$3, created_at=DEFAULT
WHERE refresh_token=$1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE user_agent=$1 and ip=$2 and user_email=$3;

-- name: ReadSession :one
SELECT * FROM sessions
WHERE user_email=$1 and  user_agent=$2 and ip=$3;

-- name: CreateContent :exec
INSERT INTO content(user_email, name, description)
VALUES ($1, $2, $3);

-- name: UpdateContent :exec
UPDATE content SET user_email=$1, name=$3, description=$4
WHERE user_email= $1 and name =$2;

-- name: DeleteContent :exec
DELETE FROM content
WHERE user_email=$1 and name=$2;

-- name: ReadContent :one
SELECT * FROM content
WHERE user_email=$1 and id=$2;

-- name: CreateBlog :exec
INSERT INTO blog(name, description)
VALUES ($1, $2);

-- name: ReadBlog :one
SELECT * FROM blog
WHERE id=$1;

-- name: ReadEmailRoleFromSessions :one
SELECT email, role, s.expires_in
FROM users INNER JOIN sessions s on users.email = s.user_email
WHERE s.refresh_token = $1;
