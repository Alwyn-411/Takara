package transactions

import (
	"database/sql"
	"net/http"
	"takara/services/schema"
	"time"

	"github.com/gin-gonic/gin"
)

func (handler *TransactionsHandler) DeleteTransaction(ctx *gin.Context) {
	transactionId := ctx.Param("transactionId")

	existing := schema.Transaction{}
	if err := handler.dbInstance.Get(&existing, `
		SELECT transactionId, accountId, type, accountAmount, accountCurrency
		FROM transactions
		WHERE transactionId = ? AND active = 1
	`, transactionId); err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now().Unix()

	dbTx, err := handler.dbInstance.Beginx()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start db transaction"})
		return
	}

	// Soft-delete the transaction
	if _, err = dbTx.Exec(`
		UPDATE transactions SET active = 0, updatedAt = ? WHERE transactionId = ?
	`, now, transactionId); err != nil {
		dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete transaction"})
		return
	}

	// Clean up tag bindings (hard delete — junction rows have no meaning without the transaction)
	if _, err = dbTx.Exec(`DELETE FROM transaction_tags WHERE transactionId = ?`, transactionId); err != nil {
		dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clean up tags"})
		return
	}

	// Reverse the balance effect
	if err = reverseFromBalance(dbTx, existing.AccountId, existing.AccountCurrency, existing.Type, existing.AccountAmount, now); err != nil {
		dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reverse balance"})
		return
	}

	if err := dbTx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
