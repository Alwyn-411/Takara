package holdings

import (
	"github.com/jmoiron/sqlx"
)

type HoldingsHandler struct {
	dbInstance *sqlx.DB
}

func NewHoldingsHandler(db *sqlx.DB) *HoldingsHandler {
	return &HoldingsHandler{dbInstance: db}
}
