package transactions

import (
	"database/sql"
	"net/http"
	"takara/services/schema"

	"github.com/gin-gonic/gin"
)

type MerchantInfo struct {
	MerchantId   string `db:"merchantId" json:"merchantId"`
	MerchantName string `db:"merchantName" json:"merchantName"`
}

type CategoryInfo struct {
	CategoryId   string `db:"categoryId" json:"categoryId"`
	CategoryName string `db:"categoryName" json:"categoryName"`
}

type TransactionResponse struct {
	schema.Transaction
	Merchant *MerchantInfo `json:"merchant,omitempty"`
	Category *CategoryInfo `json:"category,omitempty"`
	Tags     []schema.Tag  `json:"tags"`
}

func (handler *TransactionsHandler) GetTransactionById(ctx *gin.Context) {
	transactionId := ctx.Param("transactionId")

	tx := schema.Transaction{}
	if err := handler.dbInstance.Get(&tx, `
		SELECT transactionId, userId, accountId, type,
			settledAmount, settledCurrency,
			accountAmount, accountCurrency,
			exchangeRate, merchantId, categoryId,
			description, transactionAt, active,
			createdAt, updatedAt
		FROM transactions
		WHERE transactionId = ? AND active = 1
	`, transactionId); err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := TransactionResponse{Transaction: tx, Tags: []schema.Tag{}}

	// Resolve merchant name
	if tx.MerchantId != "" {
		var m MerchantInfo
		if err := handler.dbInstance.Get(&m, `
			SELECT merchantId, merchantName FROM merchants WHERE merchantId = ?
		`, tx.MerchantId); err == nil {
			resp.Merchant = &m
		}
	}

	// Resolve category name
	if tx.CategoryId != "" {
		var c CategoryInfo
		if err := handler.dbInstance.Get(&c, `
			SELECT categoryId, categoryName FROM categories WHERE categoryId = ?
		`, tx.CategoryId); err == nil {
			resp.Category = &c
		}
	}

	// Fetch tags
	if err := handler.dbInstance.Select(&resp.Tags, `
		SELECT t.tagId, t.tagName, t.userId, t.active, t.createdAt, t.updatedAt
		FROM tags t
		JOIN transaction_tags tt ON t.tagId = tt.tagId
		WHERE tt.transactionId = ?
	`, transactionId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

type TransactionListItem struct {
	schema.Transaction
	MerchantName string `db:"merchantName" json:"merchantName"`
	CategoryName string `db:"categoryName" json:"categoryName"`
}

func (handler *TransactionsHandler) GetTransactionsByAccountId(ctx *gin.Context) {
	accountId := ctx.Param("accountId")
	limit := ctx.DefaultQuery("limit", "20")
	offset := ctx.DefaultQuery("offset", "0")

	// Use LEFT JOINs so we get merchant/category names in one query
	transactions := []TransactionListItem{}
	if err := handler.dbInstance.Select(&transactions, `
		SELECT
			t.transactionId, t.userId, t.accountId, t.type,
			t.settledAmount, t.settledCurrency,
			t.accountAmount, t.accountCurrency,
			t.exchangeRate, t.merchantId, t.categoryId,
			t.description, t.transactionAt, t.active,
			t.createdAt, t.updatedAt,
			COALESCE(m.merchantName, '') AS merchantName,
			COALESCE(c.categoryName, '') AS categoryName
		FROM transactions t
		LEFT JOIN merchants m ON t.merchantId = m.merchantId
		LEFT JOIN categories c ON t.categoryId = c.categoryId
		WHERE t.accountId = ? AND t.active = 1
		ORDER BY t.transactionAt DESC
		LIMIT ? OFFSET ?
	`, accountId, limit, offset); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, schema.ListResponse[TransactionListItem]{
		Count:   len(transactions),
		Records: transactions,
	})
}
