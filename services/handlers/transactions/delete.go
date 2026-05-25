package transactions

import (
	"database/sql"
	"fmt"
	"net/http"
	"takara/services/schema"
	"time"

	"github.com/gin-gonic/gin"
)

func (handler *TransactionsHandler) DeleteTransaction(ctx *gin.Context) {
	/*
		1. Get Account Amount
		2. Begin Transaction
		3. Soft Delete Transaction
		4. Reverse Transaction on the Account Amount
		5. Wrap up
	*/
	transactionId := ctx.Param("transactionId")

	existing := schema.Transaction{}
	err := handler.dbInstance.Get(&existing, `
		SELECT transactionId, accountId, type, accountAmount
		FROM transactions
		WHERE transactionId = ? AND active = 1
	`, transactionId)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dbTx, err := handler.dbInstance.Beginx()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start db transaction"})
		return
	}

	// Soft delete the transaction
	_, err = dbTx.Exec(`
		UPDATE transactions SET active = 0, updatedAt = ? WHERE transactionId = ?
	`, time.Now().Unix(), transactionId)
	if err != nil {
		dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete transaction"})
		return
	}

	// Reverse the balance change
	sign := "+"
	if existing.Type == "Credit" {
		sign = "-"
	}
	_, err = dbTx.Exec(
		fmt.Sprintf(`
			UPDATE accounts
			SET balance = printf('%%.2f', CAST(balance AS REAL) %s CAST(? AS REAL)),
				updatedAt = ?
			WHERE accountId = ?`, sign),
		existing.AccountAmount, time.Now().Unix(), existing.AccountId,
	)
	if err != nil {
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
