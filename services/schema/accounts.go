package schema

import (
	"github.com/jmoiron/sqlx"
)

type Account struct {
	Base
	AccountID string `db:"account_id" json:"accountId"`

	Type          string `db:"type" json:"type,omitempty"`
	Name          string `db:"name" json:"name"`
	AccountNumber string `db:"account_number" json:"account_number"`
	Description   string `db:"description" json:"description,omitempty"`
	Currency      string `db:"currency" json:"currency"`

	Interest float64 `db:"interest" json:"interest,omitempty"`
	Balance  float64 `db:"balance" json:"balance,omitempty"`

	Timestamp
}

func InitAccounts(dbInstance *sqlx.DB) {
	command := `
		CREATE TABLE IF NOT EXISTS accounts (
		user_id TEXT NOT NULL,
		active INTEGER DEFAULT 1,

		account_id TEXT PRIMARY KEY NOT NULL,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		account_number TEXT NOT NULL,
		description TEXT,
		currency TEXT NOT NULL,

		interest REAL DEFAULT 0,
		balance REAL DEFAULT 0,

		created_at INTEGER DEFAULT (strftime('%s','now')),
		updated_at INTEGER DEFAULT (strftime('%s','now')),

		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	dbInstance.MustExec(command)
}
