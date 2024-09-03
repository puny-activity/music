package album

import (
	"github.com/google/uuid"
	"github.com/puny-activity/music/internal/entity/song/cover"
)

type ID uuid.UUID

var DefaultID = uuid.MustParse("00000000-0000-0000-0000-000000000000")

func NewID(uuidID uuid.UUID) ID {
	return ID(uuidID)
}

func GenerateID() ID {
	return ID(uuid.New())
}

func (e ID) String() string {
	return uuid.UUID(e).String()
}

type Album struct {
	ID    *ID
	Title string
	Cover *cover.Cover
}
