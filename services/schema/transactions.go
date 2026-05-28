package schema

import (
	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	Base

	AccountId     string `db:"accountId" json:"accountId"`
	TransactionId string `db:"transactionId" json:"transactionId"`

	Type string `db:"type" json:"type"` // Debit | Credit

	// other party billed or settled amount
	SettledAmount   string `db:"settledAmount" json:"settledAmount"`
	SettledCurrency string `db:"settledCurrency" json:"settledCurrency"`

	// source account amount
	AccountAmount   string `db:"accountAmount" json:"accountAmount"`     // Cost in account currency
	AccountCurrency string `db:"accountCurrency" json:"accountCurrency"` // ISO 4217 e.g. "INR"

	ExchangeRate string `db:"exchangeRate" json:"exchangeRate"` // Conversion rate at transaction time

	MerchantId    string `db:"merchantId" json:"merchantId"` // Name of the MerchantId
	CategoryId    string `db:"categoryId" json:"categoryId"`
	Description   string `db:"description" json:"description,omitempty"`
	TransactionAt int64  `db:"transactionAt" json:"transactionAt"` // Transaction Timestamp

	Timestamp
}

// Join Transaction Tag Table for : many to one on tagId
type TransactionTag struct {
	TransactionId string `db:"transactionId" json:"transactionId"`
	TagId         string `db:"tagId" json:"tagId"`
}

func InitTransactions(db *sqlx.DB) {
	query := `
			CREATE TABLE IF NOT EXISTS transactions (
			transactionId TEXT PRIMARY KEY NOT NULL,
			userId TEXT NOT NULL,
			active INTEGER DEFAULT 1,
			accountId TEXT NOT NULL,
			type TEXT NOT NULL,
 
			settledAmount TEXT NOT NULL,
			settledCurrency TEXT NOT NULL,
			accountAmount TEXT NOT NULL,
			accountCurrency TEXT NOT NULL,
			exchangeRate TEXT DEFAULT '1',
 
			merchantId TEXT,
			categoryId TEXT,
			description TEXT,
			transactionAt INTEGER NOT NULL,
 
			createdAt INTEGER DEFAULT (strftime('%s','now')),
			updatedAt INTEGER DEFAULT (strftime('%s','now')),
 
			FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE,
			FOREIGN KEY(accountId) REFERENCES accounts(accountId) ON DELETE CASCADE,
			FOREIGN KEY(categoryId) REFERENCES categories(categoryId) ON DELETE SET NULL,
			FOREIGN KEY(merchantId) REFERENCES merchants(merchantId) ON DELETE SET NULL
		);
	`

	db.MustExec(query)
}

func InitTransactionTags(db *sqlx.DB) {
	query := `
        CREATE TABLE IF NOT EXISTS transaction_tags (
        transactionId TEXT NOT NULL,
        tagId TEXT NOT NULL,
        PRIMARY KEY (transactionId, tagId),
        FOREIGN KEY(transactionId) REFERENCES transactions(transactionId) ON DELETE CASCADE,
        FOREIGN KEY(tagId) REFERENCES tags(tagId) ON DELETE CASCADE
    );
    `

	db.MustExec(query)
}

func InitIndexTransactions(db *sqlx.DB) {
	db.MustExec(`
		CREATE INDEX IF NOT EXISTS idx_transactions_user_date
			ON transactions(userId, transactionAt);
  
		CREATE INDEX IF NOT EXISTS idx_transactions_user_type
			ON transactions(userId, type, active);
 
		CREATE INDEX IF NOT EXISTS idx_transactions_category
			ON transactions(userId, categoryId);
 
		CREATE INDEX IF NOT EXISTS idx_transactions_merchantId
			ON transactions(userId, merchantId);
 
		CREATE INDEX IF NOT EXISTS idx_transaction_tags_tag
			ON transaction_tags(tagId);
	`)
}
