package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	ValidToken(token string) (*Payload, error)
}
