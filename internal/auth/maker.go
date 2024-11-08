package auth

import "time"

type Maker interface {
	Create(username string, duration time.Duration) (string, error)
	Verify(tokenString string) (*Payload, error)
}
