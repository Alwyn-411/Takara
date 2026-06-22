package schema

import (
	"github.com/jmoiron/sqlx"
)

type Preferences struct {
	Base

	Currency string `db:"currency" json:"currency"`
	Theme    string `db:"theme" json:"theme"`

	Timestamp
}

func InitPreferences(dbInstance *sqlx.DB) {
	command := `
		CREATE TABLE IF NOT EXISTS userPrefs (
		userId TEXT PRIMARY KEY NOT NULL,
		active INTEGER DEFAULT 1,
		theme TEXT,
		currency TEXT NOT NULL CHECK(length(currency) = 3),
		createdAt INTEGER DEFAULT (strftime('%s','now')),
		updatedAt INTEGER DEFAULT (strftime('%s','now')),
		FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE
	);`
	dbInstance.MustExec(command)
}
