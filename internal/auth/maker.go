package auth

import "time"

type Maker interface {
	Create(username string, userId int64, duration time.Duration) (string, error)
	Verify(tokenString string) (*Payload, error)
}
