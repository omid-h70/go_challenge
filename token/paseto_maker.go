package token

import (
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"time"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func (p PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayLoad(username, duration)
	if err != nil {
		return "", err
	}
	//SigningMethodHS256 is symmetric
	return p.paseto.Encrypt(p.symmetricKey, payload, nil)
}

func (p PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func NewPasetoMaker(symmetrickey string) (Maker, error) {
	if len(symmetrickey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("symmetrickey key size must be more than %d", chacha20poly1305.KeySize)
	}

	return PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetrickey),
	}, nil
}
