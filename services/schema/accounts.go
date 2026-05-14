package schema

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Account struct {
	Base
	AccountID string `db:"account_id" json:"accountId"`

	Type        string `db:"type" json:"type,omitempty"`
	Name        string `db:"name" json:"name"`
	AltName     string `db:"alt_name" json:"altName,omitempty"`
	Description string `db:"description" json:"description,omitempty"`
	Currency    string `db:"currency" json:"currency"`

	Interest float64 `db:"interest" json:"interest,omitempty"`
	Balance  float64 `db:"balance" json:"balance,omitempty"`

	Timestamp
}

func InitAccounts(dbInstance *sqlx.DB) {
	command := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS accounts (
		user_id TEXT NOT NULL,
		active INTEGER DEFAULT 1,

		account_id TEXT PRIMARY KEY NOT NULL,
		type TEXT,
		name TEXT NOT NULL,
		alt_name TEXT,
		description TEXT,
		currency TEXT NOT NULL,

		interest REAL DEFAULT 0,
		balance REAL DEFAULT 0,

		created_at INTEGER DEFAULT (strftime('%s','now')),
		updated_at INTEGER DEFAULT (strftime('%s','now')),

		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`)

	dbInstance.MustExec(command)
}
