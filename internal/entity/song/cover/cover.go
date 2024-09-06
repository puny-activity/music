package cover

import (
	"github.com/google/uuid"
)

type ID uuid.UUID

func NewID(uuidID uuid.UUID) ID {
	return ID(uuidID)
}

func ParseID(id string) (ID, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return ID{}, err
	}
	return ID(uuidID), nil
}

func GenerateID() ID {
	return ID(uuid.New())
}

func (e ID) String() string {
	return uuid.UUID(e).String()
}

type Cover struct {
	ID     *ID
	Width  int
	Height int
}
