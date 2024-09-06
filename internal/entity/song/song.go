package song

import (
	"github.com/google/uuid"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/internal/entity/song/cover"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"time"
)

type ID uuid.UUID

func NewID(uuidID uuid.UUID) ID {
	return ID(uuidID)
}

func GenerateID() ID {
	return ID(uuid.New())
}

func (i ID) String() string {
	return uuid.UUID(i).String()
}

type Song struct {
	ID           *ID
	FileID       *remotefile.ID
	Title        string
	Duration     time.Duration
	Cover        *cover.Cover
	Genre        genre.Base
	Album        album.Base
	Artist       artist.Base
	Year         *int
	Number       *int
	Comment      *string
	Channels     int
	BitrateKbps  int
	SampleRateHz int
	MD5          string
}
