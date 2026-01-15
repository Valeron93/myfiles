package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Username       string `gorm:"unique"`
	HashedPassword []byte
	CreatedAt      time.Time
}

type UserSession struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UserID    uuid.UUID
	User      User
}
