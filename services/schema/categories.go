package schema

import "github.com/jmoiron/sqlx"

type Category struct {
	Base

	CategoryId   string `db:"categoryId" json:"categoryId"`
	CategoryName string `db:"categoryName" json:"categoryName"`

	Timestamp
}

func InitCategories(db *sqlx.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS categories (
			categoryId TEXT PRIMARY KEY NOT NULL,

			userId TEXT NOT NULL,
			active INTEGER DEFAULT 1,

			categoryName TEXT NOT NULL COLLATE NOCASE,

			createdAt INTEGER DEFAULT (strftime('%s','now')),
			updatedAt INTEGER DEFAULT (strftime('%s','now')),

			UNIQUE(userId, categoryName),

			FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE
		);
	`

	db.MustExec(query)
}
