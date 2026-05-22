package schema

import "github.com/jmoiron/sqlx"

type Base struct {
	UserID string `db:"userId" json:"userId"`
	Active int    `db:"active" json:"active"` // SQLite: 0 or 1
}

type Timestamp struct {
	CreatedAt int64 `db:"createdAt" json:"createdAt"`
	UpdatedAt int64 `db:"updatedAt" json:"updatedAt"`
}

type ListResponse[T any] struct {
	Count   int `json:"count"`
	Records []T `json:"records"`
}

func ForeignKeysEnabled(dbInstance *sqlx.DB) {
	command := "PRAGMA foreign_keys = ON;"

	dbInstance.MustExec(command)
}
