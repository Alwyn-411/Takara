package schema

import "github.com/jmoiron/sqlx"

type Base struct {
	UserID string `db:"user_id" json:"userId"`
	Active int    `db:"active" json:"active"` // SQLite: 0 or 1
}

type Timestamp struct {
	CreatedAt int64 `db:"created_at" json:"createdAt"`
	UpdatedAt int64 `db:"updated_at" json:"updatedAt"`
}

func ForeignKeysEnabled(dbInstance *sqlx.DB) {
	command := "PRAGMA foreign_keys = ON;"

	dbInstance.MustExec(command)
}
