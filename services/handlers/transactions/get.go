package handlers

import (
	"database/sql"
	"net/http"
	"takara/services/schema"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type TransactionResponse struct {
	schema.Transaction
	Tags []schema.Tag `json:"tags"`
}

func GetTransactionById(transactionId string, ctx *gin.Context, dbInstance *sqlx.DB) {
	query := `
		SELECT transactionId, userId, accountId, type,
			settledAmount, settledCurrency,
			accountAmount, accountCurrency,
			exchangeRate, merchant, categoryId,
			description, transactionAt, active,
			createdAt, updatedAt
		FROM transactions
		WHERE transactionId = ? AND active = 1
	`

	tx := schema.Transaction{}
	err := dbInstance.Get(&tx, query, transactionId)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch associated tags
	tags := []schema.Tag{}
	dbInstance.Select(&tags, `
		SELECT t.tagId, t.tagName, t.userId, t.active, t.createdAt, t.updatedAt
		FROM tags t
		JOIN transaction_tags tt ON t.tagId = tt.tagId
		WHERE tt.transactionId = ?
	`, transactionId)

	response := TransactionResponse{
		Transaction: tx,
		Tags:        tags,
	}

	ctx.JSON(http.StatusOK, response)
}

func GetTransactionsByUserId(userId string, ctx *gin.Context, dbInstance *sqlx.DB) {
	query := `
		SELECT transactionId, userId, accountId, type,
			settledAmount, settledCurrency,
			accountAmount, accountCurrency,
			exchangeRate, merchant, categoryId,
			description, transactionAt, active,
			createdAt, updatedAt
		FROM transactions
		WHERE userId = ? AND active = 1
		ORDER BY transactionAt DESC
		LIMIT 50
	`

	transactions := []schema.Transaction{}
	err := dbInstance.Select(&transactions, query, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(transactions) == 0 {
		ctx.JSON(http.StatusOK, schema.ListResponse[schema.Transaction]{Count: 0, Records: []schema.Transaction{}})
		return
	}

	response := schema.ListResponse[schema.Transaction]{
		Count:   len(transactions),
		Records: transactions,
	}
	ctx.JSON(http.StatusOK, response)
}

func GetTransactionsByAccountId(accountId string, ctx *gin.Context, dbInstance *sqlx.DB) {
	query := `
		SELECT transactionId, userId, accountId, type,
			settledAmount, settledCurrency,
			accountAmount, accountCurrency,
			exchangeRate, merchant, categoryId,
			description, transactionAt, active,
			createdAt, updatedAt
		FROM transactions
		WHERE accountId = ? AND active = 1
		ORDER BY transactionAt DESC
		LIMIT 50
	`

	transactions := []schema.Transaction{}
	err := dbInstance.Select(&transactions, query, accountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(transactions) == 0 {
		ctx.JSON(http.StatusOK, schema.ListResponse[schema.Transaction]{Count: 0, Records: []schema.Transaction{}})
		return
	}

	response := schema.ListResponse[schema.Transaction]{
		Count:   len(transactions),
		Records: transactions,
	}
	ctx.JSON(http.StatusOK, response)
}
