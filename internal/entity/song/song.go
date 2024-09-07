package song

import (
	"github.com/google/uuid"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/internal/entity/song/cover"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/pkg/werr"
	"time"
)

type ID uuid.UUID

func ParseID(id string) (ID, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return ID{}, werr.WrapSE("failed to parse uuid", err)
	}
	return ID(idUUID), nil
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

const (
	PaginationNumber     = "song_number"
	PaginationTitle      = "song_title"
	PaginationYear       = "song_year"
	PaginationDuration   = "song_duration"
	PaginationBitrate    = "song_bitrate"
	PaginationSampleRate = "song_sample_rate"
)

const (
	FilterGenre  = "song_genre_filter"
	FilterAlbum  = "song_album_filter"
	FilterArtist = "song_artist_filter"
)
