package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrTokenGotExpired = errors.New("token got expired")
	ErrInvalidToken    = errors.New("invalid token")
)

type Payload struct {
	ID        uuid.UUID `json:"uuid"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (p Payload) Valid() error {
	if time.Now().After(p.ExpiresAt) {
		return ErrTokenGotExpired
	}
	return nil
}

func NewPayLoad(username string, duration time.Duration) (*Payload, error) {
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        randomUUID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}
