package artistrepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/pkg/queryer"
)

type createAllEntity struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (r *Repository) CreateAll(ctx context.Context, artists []artist.Base) error {
	return r.createAll(ctx, r.db, artists)
}

func (r *Repository) CreateAllTx(ctx context.Context, tx *sqlx.Tx, artists []artist.Base) error {
	return r.createAll(ctx, tx, artists)
}

func (r *Repository) createAll(ctx context.Context, queryer queryer.Queryer, artists []artist.Base) error {
	if len(artists) == 0 {
		return nil
	}

	query := `
INSERT INTO artists(id, name) 
VALUES (:id, :name)
`

	artistsRepo := make([]createAllEntity, len(artists))
	for i, artistItem := range artists {
		artistsRepo[i] = createAllEntity{
			ID:   uuid.UUID(*artistItem.ID),
			Name: artistItem.Name,
		}
	}

	_, err := queryer.NamedExecContext(ctx, query, artistsRepo)
	if err != nil {
		return err
	}

	return nil
}
