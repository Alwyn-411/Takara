package transactions

import (
	"fmt"

	"github.com/bojanz/currency"
	"github.com/jmoiron/sqlx"
)

var validTypes = map[string]bool{
	"Debit":  true,
	"Credit": true,
}

func IsValidType(t string) bool {
	return validTypes[t]
}

func signedDelta(txType string, accountAmount currency.Amount) (currency.Amount, error) {
	switch txType {
	case "Credit":
		return accountAmount, nil
	case "Debit":
		return accountAmount.Mul("-1")
	default:
		return currency.Amount{}, fmt.Errorf("invalid transaction type: %q", txType)
	}
}

func applyToBalance(dbTx *sqlx.Tx, accountId string, accountCurrency string, txType string, accountAmountStr string, updatedAt int64) error {
	var balanceStr string
	if err := dbTx.Get(&balanceStr, `SELECT balance FROM accounts WHERE accountId = ?`, accountId); err != nil {
		return fmt.Errorf("read balance: %w", err)
	}
	if balanceStr == "" {
		balanceStr = "0"
	}

	balance, err := currency.NewAmount(balanceStr, accountCurrency)
	if err != nil {
		return fmt.Errorf("parse balance %q: %w", balanceStr, err)
	}

	txAmount, err := currency.NewAmount(accountAmountStr, accountCurrency)
	if err != nil {
		return fmt.Errorf("parse tx amount %q: %w", accountAmountStr, err)
	}

	delta, err := signedDelta(txType, txAmount)
	if err != nil {
		return err
	}

	newBalance, err := balance.Add(delta)
	if err != nil {
		return fmt.Errorf("apply delta: %w", err)
	}

	if _, err := dbTx.Exec(
		`UPDATE accounts SET balance = ?, updatedAt = ? WHERE accountId = ?`,
		newBalance.Number(), updatedAt, accountId,
	); err != nil {
		return fmt.Errorf("write balance: %w", err)
	}
	return nil
}

func reverseFromBalance(dbTx *sqlx.Tx, accountId string, accountCurrency string, txType string, accountAmountStr string, updatedAt int64) error {
	reversed := "Credit"
	if txType == "Credit" {
		reversed = "Debit"
	}
	return applyToBalance(dbTx, accountId, accountCurrency, reversed, accountAmountStr, updatedAt)
}
