// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: session.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (
    session_uuid,
    user_name,
    user_agent,
    refresh_token,
    client_ip,
    is_blocked,
    expires_at
) VALUES ( $1, $2, $3, $4, $5, $6, $7 ) RETURNING session_uuid, user_name, user_agent, refresh_token, client_ip, is_blocked, expires_at, created_at
`

type CreateSessionParams struct {
	SessionUuid  uuid.UUID `json:"session_uuid"`
	UserName     string    `json:"user_name"`
	UserAgent    string    `json:"user_agent"`
	RefreshToken string    `json:"refresh_token"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.SessionUuid,
		arg.UserName,
		arg.UserAgent,
		arg.RefreshToken,
		arg.ClientIp,
		arg.IsBlocked,
		arg.ExpiresAt,
	)
	var i Session
	err := row.Scan(
		&i.SessionUuid,
		&i.UserName,
		&i.UserAgent,
		&i.RefreshToken,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSession = `-- name: GetSession :one
SELECT session_uuid, user_name, user_agent, refresh_token, client_ip, is_blocked, expires_at, created_at FROM sessions WHERE session_uuid = $1 LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, sessionUuid uuid.UUID) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSession, sessionUuid)
	var i Session
	err := row.Scan(
		&i.SessionUuid,
		&i.UserName,
		&i.UserAgent,
		&i.RefreshToken,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}
