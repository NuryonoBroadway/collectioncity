package repositories

import (
	"context"
	"database/sql"
)

// DBTransaction contract database transaction
type DBTransaction interface {
	ExecTX(ctx context.Context, options *sql.TxOptions, fn func(context.Context, StoreTX) (int64, error)) (int64, error)
}

// StoreTX data store transaction contract
type StoreTX interface {
	// Create your function contract here
}
