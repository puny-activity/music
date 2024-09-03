package user

import (
	"github.com/google/uuid"
	"github.com/puny-activity/music/internal/entity/user/email"
)

type ID uuid.UUID

func NewID(uuidID uuid.UUID) ID {
	return ID(uuidID)
}

func GenerateID() ID {
	return ID(uuid.New())
}

type User struct {
	ID       *ID
	Nickname string
	Email    email.Email
}
