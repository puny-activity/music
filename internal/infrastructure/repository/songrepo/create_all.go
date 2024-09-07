package songrepo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/util"
)

type createAllDTO struct {
	ID             string  `db:"id"`
	FileID         string  `db:"file_id"`
	Title          string  `db:"title"`
	DurationNs     int64   `db:"duration_ns"`
	CoverID        *string `db:"cover_id"`
	GenreID        string  `db:"genre_id"`
	AlbumID        string  `db:"album_id"`
	ArtistID       string  `db:"artist_id"`
	Year           *int    `db:"year"`
	Number         *int    `db:"number"`
	Comment        *string `db:"comment"`
	Channels       int     `db:"channels"`
	BitrateKbps    int     `db:"bitrate_kbps"`
	SampleRateKbps int     `db:"sample_rate_hz"`
	MD5            string  `db:"md5"`
}

func (r *Repository) CreateAll(ctx context.Context, songsToCreate []song.Song) error {
	return r.createAll(ctx, r.db, songsToCreate)
}

func (r *Repository) CreateAllTx(ctx context.Context, tx *sqlx.Tx, songsToCreate []song.Song) error {
	return r.createAll(ctx, tx, songsToCreate)
}

func (r *Repository) createAll(ctx context.Context, queryer queryer.Queryer, songsToCreate []song.Song) error {
	query := `
INSERT INTO songs(id, file_id, title, duration_ns, cover_id, genre_id, album_id, artist_id, year, number, comment, channels, bitrate_kbps, sample_rate_hz, md5) 
VALUES (:id, :file_id, :title, :duration_ns, :cover_id, :genre_id, :album_id, :artist_id, :year, :number, :comment, :channels, :bitrate_kbps, :sample_rate_hz, :md5)
`

	songsDTO := make([]createAllDTO, len(songsToCreate))
	for i, songToCreate := range songsToCreate {
		var coverID *string = nil
		if songToCreate.Cover != nil {
			coverID = util.ToPointer(songToCreate.Cover.ID.String())
		}
		songsDTO[i] = createAllDTO{
			ID:             songToCreate.ID.String(),
			FileID:         songToCreate.FileID.String(),
			Title:          songToCreate.Title,
			DurationNs:     songToCreate.Duration.Nanoseconds(),
			CoverID:        coverID,
			GenreID:        songToCreate.Genre.ID.String(),
			AlbumID:        songToCreate.Album.ID.String(),
			ArtistID:       songToCreate.Artist.ID.String(),
			Year:           songToCreate.Year,
			Number:         songToCreate.Number,
			Comment:        songToCreate.Comment,
			Channels:       songToCreate.Channels,
			BitrateKbps:    songToCreate.BitrateKbps,
			SampleRateKbps: songToCreate.SampleRateHz,
			MD5:            songToCreate.MD5,
		}
	}

	_, err := queryer.NamedExecContext(ctx, query, songsDTO)
	if err != nil {
		return err
	}

	return nil
}
