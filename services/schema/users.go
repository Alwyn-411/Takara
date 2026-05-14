package schema

import (
	"github.com/jmoiron/sqlx"
)

type User struct {
	Base
	Username string `db:"username" json:"userName"`
	AltName  string `db:"alt_name" json:"altName,omitempty"`
	Email    string `db:"email" json:"email"`
	AltEmail string `db:"alt_email" json:"altEmail,omitempty"`

	PasswordHash string `db:"password_hash" json:"-"`
	Timestamp
}

func InitUsers(dbInstance *sqlx.DB) {
	create := `
		CREATE TABLE IF NOT EXISTS users (
		user_id TEXT PRIMARY KEY NOT NULL,
		active INTEGER DEFAULT 1,

		username TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,

		created_at INTEGER DEFAULT (strftime('%s','now')),
		updated_at INTEGER DEFAULT (strftime('%s','now'))
		)
	`

	dbInstance.MustExec(create)
}
