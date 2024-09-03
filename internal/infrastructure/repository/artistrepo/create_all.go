package artistrepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/util"
)

type createAllEntity struct {
	ID      uuid.UUID  `db:"id"`
	Name    string     `db:"name"`
	CoverID *uuid.UUID `db:"cover_id"`
}

func (r *Repository) CreateAll(ctx context.Context, artists []artist.Artist) error {
	return r.createAll(ctx, r.db, artists)
}

func (r *Repository) CreateAllTx(ctx context.Context, tx *sqlx.Tx, artists []artist.Artist) error {
	return r.createAll(ctx, tx, artists)
}

func (r *Repository) createAll(ctx context.Context, queryer queryer.Queryer, artists []artist.Artist) error {
	if len(artists) == 0 {
		return nil
	}

	query := `
INSERT INTO artists(id, name, cover_id) 
VALUES (:id, :name, :cover_id)
`

	artistsRepo := make([]createAllEntity, len(artists))
	for i, artistItem := range artists {
		var coverID *uuid.UUID
		if artistItem.Cover != nil {
			coverID = util.ToPointer(uuid.UUID(*artistItem.Cover.ID))
		}
		artistsRepo[i] = createAllEntity{
			ID:      uuid.UUID(*artistItem.ID),
			Name:    artistItem.Name,
			CoverID: coverID,
		}
	}

	_, err := queryer.NamedExecContext(ctx, query, artistsRepo)
	if err != nil {
		return err
	}

	return nil
}
