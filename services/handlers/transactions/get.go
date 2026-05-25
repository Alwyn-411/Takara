package transactions

import (
	"database/sql"
	"net/http"
	"takara/services/schema"

	"github.com/gin-gonic/gin"
)

type TransactionResponse struct {
	schema.Transaction
	Tags []schema.Tag `json:"tags"`
}

func (handler *TransactionsHandler) GetTransactionById(ctx *gin.Context) {
	/*
		1. Fetch Data from DB with transactionId
		2. Fetch Tags from DB related to this transaction
	*/
	transactionId := ctx.Param("transactionId")

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
	err := handler.dbInstance.Get(&tx, query, transactionId)
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
	handler.dbInstance.Select(&tags, `
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

func (handler *TransactionsHandler) GetTransactionsByAccountId(ctx *gin.Context) {
	accountId := ctx.Param("accountId")
	limit := ctx.Param("limit")
	offset := ctx.Param("offset")

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
		LIMIT ? OFFSET ?
	`

	transactions := []schema.Transaction{}
	err := handler.dbInstance.Select(&transactions, query, accountId, limit, offset)
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

/*	TODO: This might not be needed :: Think about it
	func (handler *TransactionsHandler) GetTransactionsByUserId(ctx *gin.Context) {
		userId := ctx.Param("userId")

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
		err := handler.dbInstance.Select(&transactions, query, userId)
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
*/
