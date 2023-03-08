package repositories

import "gitlab.privy.id/privypass/privypass-oauth2-core-se/pkg/mariadb"

type storer struct {
	dbTx mariadb.Transaction
}

//  NewStore create ne instance database transaction
func NewStore(dbTx mariadb.Transaction) *storer {
	return &storer{dbTx: dbTx}
}
