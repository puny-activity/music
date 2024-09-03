package txmanager

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Transactor interface {
	Transaction(ctx context.Context, f func(context.Context, *sqlx.Tx) error) error
}
