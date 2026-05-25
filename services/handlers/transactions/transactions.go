package transactions

import (
	"takara/services/forex"

	"github.com/jmoiron/sqlx"
)

type TransactionsHandler struct {
	dbInstance *sqlx.DB
	forEx      *forex.ForEx
}

func NewTransactionsHandler(db *sqlx.DB, forEx *forex.ForEx) *TransactionsHandler {
	return &TransactionsHandler{
		dbInstance: db,
		forEx:      forEx,
	}
}
