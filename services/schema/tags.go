package schema

import "github.com/jmoiron/sqlx"

type Tag struct {
	Base

	TagId   string `db:"tagId" json:"tagId"`
	TagName string `db:"tagName" json:"tagName"`

	Timestamp
}

func InitTags(db *sqlx.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS tags (
			tagId TEXT PRIMARY KEY NOT NULL,

			userId TEXT NOT NULL,
			active INTEGER DEFAULT 1,

			tagName TEXT NOT NULL COLLATE NOCASE,

			createdAt INTEGER DEFAULT (strftime('%s','now')),
			updatedAt INTEGER DEFAULT (strftime('%s','now')),

			UNIQUE(userId, tagName),

			FOREIGN KEY(userId)
				REFERENCES users(userId)
				ON DELETE CASCADE
		);
	`

	db.MustExec(query)
}
