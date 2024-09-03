package queue

import (
	"github.com/google/uuid"
	"github.com/puny-activity/music/internal/entity/song"
)

type ID uuid.UUID

func NewID(uuidID uuid.UUID) ID {
	return ID(uuidID)
}

func GenerateID() ID {
	return ID(uuid.New())
}

type Queue struct {
	ID    *ID
	Songs []song.Song
}
