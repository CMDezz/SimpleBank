-- name: CreateVerifyEmail :one
INSERT INTO verify_user_email (
    username,
    email,
    secret_code
)   VALUES (
    $1,$2,$3
) RETURNING *;