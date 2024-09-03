package device

import "github.com/google/uuid"

type ID uuid.UUID

func NewID(uuidID uuid.UUID) ID {
	return ID(uuidID)
}

func GenerateID() ID {
	return ID(uuid.New())
}

type Device struct {
	ID   *ID
	Name string
}
