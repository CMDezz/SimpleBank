-- name: CreateUser :one
INSERT INTO users (
    username,
    hased_password,
    full_name,
    email
) VALUES (
    $1,$2,$3,$4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE username = $1 LIMIT 1; 

-- name: UpdateUser :one
UPDATE users
SET 
    hased_password = COALESCE(sqlc.narg(hased_password),hased_password),
    full_name =COALESCE(sqlc.narg(full_name),full_name),
    email = COALESCE(sqlc.narg(email),email)
WHERE 
    username = sqlc.arg(username)
RETURNING *;