package schema

import (
	"github.com/jmoiron/sqlx"
)

type User struct {
	Base
	Username string `db:"userName" json:"userName"`
	AltName  string `db:"altName" json:"altName,omitempty"`
	Email    string `db:"email" json:"email"`
	AltEmail string `db:"altEmail" json:"altEmail,omitempty"`

	Password string `db:"password" json:"password"`
	Timestamp
}

func InitUsers(dbInstance *sqlx.DB) {
	create := `
		CREATE TABLE IF NOT EXISTS users (
		userId TEXT PRIMARY KEY NOT NULL,
		active INTEGER DEFAULT 1,

		userName TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		email TEXT,
		altName TEXT,
		altEmail TEXT,

		createdAt INTEGER DEFAULT (strftime('%s','now')),
		updatedAt INTEGER DEFAULT (strftime('%s','now'))
		)
	`

	dbInstance.MustExec(create)
}
