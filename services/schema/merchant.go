package schema

import "github.com/jmoiron/sqlx"

type Merchant struct {
	Base

	MerchantId   string `db:"merchantId" json:"merchantId"`
	MerchantName string `db:"merchantName" json:"merchantName"`

	Timestamp
}

func InitMerchant(db *sqlx.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS merchants (
			merchantId TEXT PRIMARY KEY NOT NULL,

			userId TEXT NOT NULL,
			active INTEGER DEFAULT 1,

			merchantName TEXT NOT NULL,

			createdAt INTEGER DEFAULT (strftime('%s','now')),
			updatedAt INTEGER DEFAULT (strftime('%s','now')),

			UNIQUE(userId, merchantName),

			FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE
		);
	`

	db.MustExec(query)
}
