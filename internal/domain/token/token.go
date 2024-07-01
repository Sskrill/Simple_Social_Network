package domain

import "time"

type RefreshToken struct {
	Id        int
	UserID    int
	Token     string
	ExpiresAt time.Time
}
