package schema

import (
	"github.com/jmoiron/sqlx"
)

type Account struct {
	Base
	AccountID string `db:"accountId" json:"accountId"`

	Type          string `db:"type" json:"type,omitempty"`
	Name          string `db:"name" json:"name"`
	AccountNumber string `db:"accountNumber" json:"accountNumber"`
	Description   string `db:"description" json:"description,omitempty"`
	Currency      string `db:"currency" json:"currency"`

	Interest float64 `db:"interest" json:"interest,omitempty"`
	Balance  float64 `db:"balance" json:"balance,omitempty"`

	Timestamp
}

func InitAccounts(dbInstance *sqlx.DB) {
	command := `
		CREATE TABLE IF NOT EXISTS accounts (
		userId TEXT NOT NULL,
		active INTEGER DEFAULT 1,

		accountId TEXT PRIMARY KEY NOT NULL,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		accountNumber TEXT NOT NULL,
		description TEXT,
		currency TEXT NOT NULL,

		interest REAL DEFAULT 0,
		balance REAL DEFAULT 0,

		createdAt INTEGER DEFAULT (strftime('%s','now')),
		updatedAt INTEGER DEFAULT (strftime('%s','now')),

		FOREIGN KEY(userId) REFERENCES users(id) ON DELETE CASCADE
	);`

	dbInstance.MustExec(command)
}
