package models

import "time"

type Session struct {
    Id		string
    UserId 	string
    ExpiresAt 	time.Time
}
