package transactions

import (
	"database/sql"
	"net/http"
	"takara/services/handlers"
	"takara/services/schema"
	"time"

	"github.com/bojanz/currency"
	"github.com/gin-gonic/gin"
)

type UpdateTransactionRequest struct {
	SettledAmount   *string  `json:"settledAmount"`
	SettledCurrency *string  `json:"settledCurrency"`
	MerchantName    *string  `json:"merchantName"`
	CategoryName    *string  `json:"categoryName"`
	Description     *string  `json:"description"`
	TransactionAt   *int64   `json:"transactionAt"`
	TagNames        []string `json:"tagNames"`
}

func (handler *TransactionsHandler) UpdateTransaction(ctx *gin.Context) {
	transactionId := ctx.Param("transactionId")

	var req UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Load existing transaction
	existing := schema.Transaction{}
	if err := handler.dbInstance.Get(&existing, `
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

	oldAccountAmount := existing.AccountAmount
	now := time.Now().Unix()

	// Detect if the financial fields changed
	amountChanged := req.SettledAmount != nil && *req.SettledAmount != existing.SettledAmount
	currencyChanged := req.SettledCurrency != nil && *req.SettledCurrency != existing.SettledCurrency
	needsRebalance := amountChanged || currencyChanged

	// Prepare new settled values
	newSettledAmountStr := existing.SettledAmount
	newSettledCurrency := existing.SettledCurrency
	if req.SettledAmount != nil {
		newSettledAmountStr = *req.SettledAmount
	}
	if req.SettledCurrency != nil {
		newSettledCurrency = *req.SettledCurrency
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.TransactionAt != nil {
		existing.TransactionAt = *req.TransactionAt
	}

	// Begin DB transaction
	dbTx, err := handler.dbInstance.Beginx()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start db transaction"})
		return
	}

	// Resolve merchant name → ID
	if req.MerchantName != nil {
		merchantId, err := handlers.ResolveOrCreateMerchant(dbTx, existing.UserID, *req.MerchantName)
		if err != nil {
			dbTx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve merchant"})
			return
		}
		existing.MerchantId = merchantId
	}

	// Resolve category name → ID
	if req.CategoryName != nil {
		categoryId, err := handlers.ResolveOrCreateCategory(dbTx, existing.UserID, *req.CategoryName)
		if err != nil {
			dbTx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve category"})
			return
		}
		existing.CategoryId = categoryId
	}

	// Recompute account amount if settled amount/currency changed
	if needsRebalance {
		settledAmount, err := currency.NewAmount(newSettledAmountStr, newSettledCurrency)
		if err != nil {
			dbTx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount or currency"})
			return
		}
		existing.SettledAmount = settledAmount.Number()
		existing.SettledCurrency = newSettledCurrency

		if newSettledCurrency == existing.AccountCurrency {
			existing.AccountAmount = settledAmount.Number()
			existing.ExchangeRate = "1"
		} else {
			rate, err := handler.forEx.GetRate(newSettledCurrency, existing.AccountCurrency)
			if err != nil {
				dbTx.Rollback()
				ctx.JSON(http.StatusBadGateway, gin.H{"error": "could not fetch exchange rate"})
				return
			}
			converted, err := settledAmount.Convert(existing.AccountCurrency, rate)
			if err != nil {
				dbTx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "conversion failed"})
				return
			}
			existing.AccountAmount = converted.Round().Number()
			existing.ExchangeRate = rate
		}
	}

	// Update the transaction row
	if _, err = dbTx.NamedExec(`
		UPDATE transactions SET
			settledAmount   = :settledAmount,
			settledCurrency = :settledCurrency,
			accountAmount   = :accountAmount,
			exchangeRate    = :exchangeRate,
			merchantId      = :merchantId,
			categoryId      = :categoryId,
			description     = :description,
			transactionAt   = :transactionAt,
			updatedAt       = :updatedAt
		WHERE transactionId = :transactionId
	`, gin.H{
		"transactionId":   transactionId,
		"settledAmount":   existing.SettledAmount,
		"settledCurrency": existing.SettledCurrency,
		"accountAmount":   existing.AccountAmount,
		"exchangeRate":    existing.ExchangeRate,
		"merchantId":      existing.MerchantId,
		"categoryId":      existing.CategoryId,
		"description":     existing.Description,
		"transactionAt":   existing.TransactionAt,
		"updatedAt":       now,
	}); err != nil {
		dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update transaction"})
		return
	}

	// Rebalance if amount changed
	if needsRebalance {
		if err = reverseFromBalance(dbTx, existing.AccountId, existing.AccountCurrency, existing.Type, oldAccountAmount, now); err != nil {
			dbTx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reverse old balance"})
			return
		}
		if err = applyToBalance(dbTx, existing.AccountId, existing.AccountCurrency, existing.Type, existing.AccountAmount, now); err != nil {
			dbTx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to apply new balance"})
			return
		}
	}

	// Replace tags if tagNames was provided (even if empty — that clears all tags)
	if req.TagNames != nil {
		// Remove existing tag bindings
		if _, err = dbTx.Exec(`DELETE FROM transaction_tags WHERE transactionId = ?`, transactionId); err != nil {
			dbTx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear tags"})
			return
		}

		// Resolve and bind new tags
		for _, tagName := range req.TagNames {
			tagId, err := handlers.ResolveOrCreateTag(dbTx, existing.UserID, tagName)
			if err != nil {
				dbTx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve tag: " + tagName})
				return
			}
			if _, err = dbTx.Exec(
				`INSERT INTO transaction_tags (transactionId, tagId) VALUES (?, ?)`,
				transactionId, tagId,
			); err != nil {
				dbTx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to bind tag"})
				return
			}
		}
	}

	if err := dbTx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit"})
		return
	}

	ctx.JSON(http.StatusOK, existing)
}
