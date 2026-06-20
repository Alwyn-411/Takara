package schema

import (
	"github.com/jmoiron/sqlx"
)

// Holding describes a thing the user owns (asset) or owes (liability).
// It carries identity and lifecycle only — its monetary value lives in the
// holding_valuations time series, never on this row. "Current value" is always
// derived as the latest valuation, so there is no single field to drift.
type Holding struct {
	Base

	HoldingID string `db:"holdingId" json:"holdingId"`

	Kind        string `db:"kind" json:"kind"`
	Type        string `db:"type" json:"type"` // app-validated: 'cash','brokerage','mortgage',...
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description,omitempty"`
	Currency    string `db:"currency" json:"currency"`

	OpenedAt int64  `db:"openedAt" json:"openedAt"`           // entered your financial life
	ClosedAt *int64 `db:"closedAt" json:"closedAt,omitempty"` // nil = still held

	Timestamp
}

// HoldingValuation is one observation of a holding's worth at a point in time.
type HoldingValuation struct {
	Base

	ValuationID string `db:"valuationId" json:"valuationId"`
	HoldingID   string `db:"holdingId" json:"holdingId"`

	Value     string  `db:"value" json:"value"`                   // authoritative magnitude, in the holding's currency
	Quantity  *string `db:"quantity" json:"quantity,omitempty"`   // optional: units held (shares, grams)
	UnitPrice *string `db:"unitPrice" json:"unitPrice,omitempty"` // optional: price per unit at observedAt

	ObservedAt int64   `db:"observedAt" json:"observedAt"` // VALID time: when this value was true
	Note       *string `db:"note" json:"note,omitempty"`

	Timestamp // createdAt = TRANSACTION time: when the row was recorded
}

func InitHolding(dbInstance *sqlx.DB) {
	command := `
		CREATE TABLE IF NOT EXISTS holdings (
			holdingId   TEXT PRIMARY KEY NOT NULL,
			userId      TEXT NOT NULL,
			active      INTEGER DEFAULT 1,
			kind        TEXT NOT NULL,
			type        TEXT NOT NULL CHECK(type IN ('Asset','Liability')),
			name        TEXT NOT NULL,
			description TEXT,
			currency    TEXT NOT NULL DEFAULT 'INR',
			openedAt    INTEGER NOT NULL,
			closedAt    INTEGER,
			createdAt   INTEGER DEFAULT (strftime('%s','now')),
			updatedAt   INTEGER DEFAULT (strftime('%s','now')),
			FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE
		);`
	dbInstance.MustExec(command)
}

func InitHoldingValuation(dbInstance *sqlx.DB) {
	command := `
		CREATE TABLE IF NOT EXISTS holding_valuations (
			valuationId TEXT PRIMARY KEY NOT NULL,
			userId      TEXT NOT NULL,
			holdingId   TEXT NOT NULL,
			active      INTEGER DEFAULT 1,
			value       TEXT NOT NULL,
			quantity    TEXT,
			unitPrice   TEXT,
			observedAt  INTEGER NOT NULL,
			note        TEXT,
			createdAt   INTEGER DEFAULT (strftime('%s','now')),
			updatedAt   INTEGER DEFAULT (strftime('%s','now')),
			FOREIGN KEY(userId)    REFERENCES users(userId)       ON DELETE CASCADE,
			FOREIGN KEY(holdingId) REFERENCES holdings(holdingId) ON DELETE CASCADE
		);`
	dbInstance.MustExec(command)
}

func InitIndexHolding(dbInstance *sqlx.DB) {
	command := `
		CREATE INDEX IF NOT EXISTS idx_holdings_userId
			ON holdings(userId, active);`
	dbInstance.MustExec(command)
}

func InitIndexHoldingValuation(dbInstance *sqlx.DB) {
	command := `
		CREATE INDEX IF NOT EXISTS idx_holding_valuations_holding
			ON holding_valuations(holdingId, active, observedAt);`
	dbInstance.MustExec(command)
}
