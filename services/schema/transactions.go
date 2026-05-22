package schema

import (
	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	Base

	AccountId     string `db:"accountId" json:"accountId"`
	TransactionId string `db:"transactionId" json:"transactionId"`

	Type             string  `db:"type" json:"type"`         // Debit | Credit
	Amount           float64 `db:"amount" json:"amount"`     // Merchant-side amount
	Merchant         string  `db:"merchant" json:"merchant"` // Name of the Merchant
	MerchantCurrency string  `db:"merchantCurrency" json:"merchantCurrency"`

	ExchangeRate float64 `db:"exchangeRate" json:"exchangeRate"` // Conversion rate at transaction time
	BaseAmount   float64 `db:"baseAmount" json:"baseAmount"`     // Cost in account currency

	CategoryId  string `db:"categoryId" json:"categoryId"`
	Description string `db:"description" json:"description,omitempty"`

	TransactionAt int64 `db:"transactionAt" json:"transactionAt"` // Transaction Timestamp

	Timestamp
}

func InitTransactions(db *sqlx.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS transactions (
			transactionId TEXT PRIMARY KEY NOT NULL,

			userId TEXT NOT NULL,
			active INTEGER DEFAULT 1,
			accountId TEXT NOT NULL,

			type TEXT NOT NULL,
			amount REAL NOT NULL,

			merchantCurrency TEXT NOT NULL,
			exchangeRate REAL DEFAULT 1,
			baseAmount REAL NOT NULL,

			merchant TEXT NOT NULL,

			categoryId TEXT,
			description TEXT,

			transactionAt INTEGER NOT NULL,

			createdAt INTEGER DEFAULT (strftime('%s','now')),
			updatedAt INTEGER DEFAULT (strftime('%s','now')),

			FOREIGN KEY(userId)
				REFERENCES users(userId)
				ON DELETE CASCADE,

			FOREIGN KEY(accountId)
				REFERENCES accounts(accountId)
				ON DELETE CASCADE,

			FOREIGN KEY(categoryId)
				REFERENCES categories(categoryId)
				ON DELETE SET NULL
		);
	`

	db.MustExec(query)
}
