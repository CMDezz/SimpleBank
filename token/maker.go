package token

import "time"

//Maker create interface for token
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	ValidToken(token string) (*Payload, error)
}
