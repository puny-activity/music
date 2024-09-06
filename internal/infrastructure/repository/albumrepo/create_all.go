package albumrepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/pkg/queryer"
)

type createAllEntity struct {
	ID    uuid.UUID `db:"id"`
	Title string    `db:"title"`
}

func (r *Repository) CreateAll(ctx context.Context, albums []album.Base) error {
	return r.createAll(ctx, r.db, albums)
}

func (r *Repository) CreateAllTx(ctx context.Context, tx *sqlx.Tx, albums []album.Base) error {
	return r.createAll(ctx, tx, albums)
}

func (r *Repository) createAll(ctx context.Context, queryer queryer.Queryer, albums []album.Base) error {
	if len(albums) == 0 {
		return nil
	}

	query := `
INSERT INTO albums(id, title) 
VALUES (:id, :title)
`

	albumsRepo := make([]createAllEntity, len(albums))
	for i, albumItem := range albums {
		albumsRepo[i] = createAllEntity{
			ID:    uuid.UUID(*albumItem.ID),
			Title: albumItem.Title,
		}
	}

	_, err := queryer.NamedExecContext(ctx, query, albumsRepo)
	if err != nil {
		return err
	}

	return nil
}
