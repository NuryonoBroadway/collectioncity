package repositories

import "collection-squad/collection/collection-city/pkg/mariadb"

type storer struct {
	dbTx mariadb.Transaction
}

// NewStore create ne instance database transaction
func NewStore(dbTx mariadb.Transaction) *storer {
	return &storer{dbTx: dbTx}
}
