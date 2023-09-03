-- name: CreateUser :one
INSERT INTO users (
    user_name,
    hashed_password,
    full_name,
    email
) VALUES ( $1, $2, $3, $4 ) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE user_name = $1 LIMIT 1;

-- name: UpdateUser1 :one
UPDATE
    users
SET
    hashed_password = $1,
    full_name = $2,
    email = $3
WHERE
    user_name =  $4
RETURNING *;

-- @ is the same as sqlc.arg - @ used for named parameters - $ used for positional parameters
-- @hashed_password is the one from outside - hashed_password (without @) is database field
-- name: UpdateUser2 :one
UPDATE
    users
SET
    hashed_password = CASE
        WHEN @set_hashed_password::boolean = TRUE THEN @hashed_password
        ELSE hashed_password
    END,
    full_name = CASE
        WHEN @set_full_name::boolean = TRUE THEN @full_name
        ELSE full_name
    END,
    email = CASE
        WHEN  @set_email::boolean = TRUE THEN @email
        ELSE email
    END
WHERE
    user_name =  @user_name
RETURNING *;

-- name: UpdateUser3 :one
UPDATE
    users
SET
    hashed_password = CASE
        WHEN @set_hashed_password::boolean = TRUE THEN @hashed_password
        ELSE hashed_password
    END,
    full_name = CASE
        WHEN @set_full_name::boolean = TRUE THEN @full_name
        ELSE full_name
    END,
    email = CASE
        WHEN  @set_email::boolean = TRUE THEN @email
        ELSE email
    END
WHERE
    user_name =  @user_name
RETURNING *;

-- using nullable types
-- name: UpdateUser :one
UPDATE
    users
SET
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    email = COALESCE(sqlc.narg(email), email),
    is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified)
WHERE
    user_name =  sqlc.arg(user_name)
RETURNING *;