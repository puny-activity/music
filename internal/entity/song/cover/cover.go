package cover

import (
	"github.com/google/uuid"
	"github.com/puny-activity/music/internal/entity/remotefile"
)

type ID uuid.UUID

func NewID(uuidID uuid.UUID) ID {
	return ID(uuidID)
}

func GenerateID() ID {
	return ID(uuid.New())
}

func (e ID) String() string {
	return uuid.UUID(e).String()
}

type Cover struct {
	ID     *ID
	Width  string
	Height string
	File   remotefile.File
}
