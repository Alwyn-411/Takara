package handlers

import (
	"database/sql"
	"net/http"
	"takara/services/forex"
	"takara/services/schema"
	"time"

	"github.com/bojanz/currency"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type UpdateTransactionRequest struct {
	SettledAmount   string `json:"settledAmount"`
	SettledCurrency string `json:"settledCurrency"`
	Merchant        string `json:"merchant"`
	CategoryId      string `json:"categoryId"`
	Description     string `json:"description"`
	TransactionAt   int64  `json:"transactionAt"`
}

func UpdateTransaction(transactionId string, ctx *gin.Context, dbInstance *sqlx.DB, forexService *forex.Frankfurter) {
	var req UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch existing transaction
	existing := schema.Transaction{}
	err := dbInstance.Get(&existing, `
		SELECT * FROM transactions WHERE transactionId = ? AND active = 1
	`, transactionId)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Apply updates — keep existing values if not provided
	if req.Merchant != "" {
		existing.Merchant = req.Merchant
	}
	if req.CategoryId != "" {
		existing.CategoryId = req.CategoryId
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.TransactionAt != 0 {
		existing.TransactionAt = req.TransactionAt
	}

	// If amount or currency changed, recalculate conversion
	amountChanged := req.SettledAmount != "" && req.SettledAmount != existing.SettledAmount
	currencyChanged := req.SettledCurrency != "" && req.SettledCurrency != existing.SettledCurrency

	if amountChanged || currencyChanged {
		newSettledAmount := existing.SettledAmount
		newSettledCurrency := existing.SettledCurrency
		if req.SettledAmount != "" {
			newSettledAmount = req.SettledAmount
		}
		if req.SettledCurrency != "" {
			newSettledCurrency = req.SettledCurrency
		}

		settledAmount, err := currency.NewAmount(newSettledAmount, newSettledCurrency)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount or currency"})
			return
		}

		existing.SettledAmount = settledAmount.Number()
		existing.SettledCurrency = newSettledCurrency

		if newSettledCurrency == existing.AccountCurrency {
			existing.AccountAmount = settledAmount.Number()
			existing.ExchangeRate = "1"
		} else {
			rate, err := forexService.GetRate(newSettledCurrency, existing.AccountCurrency)
			if err != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"error": "could not fetch exchange rate"})
				return
			}
			converted, err := settledAmount.Convert(existing.AccountCurrency, rate)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "conversion failed"})
				return
			}
			existing.AccountAmount = converted.Round().Number()
			existing.ExchangeRate = rate
		}

		// Recalculate balance: reverse old amount, apply new amount
		dbTx, err := dbInstance.Beginx()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start db transaction"})
			return
		}

		_, err = dbTx.NamedExec(`
			UPDATE transactions SET
				settledAmount = :settledAmount,
				settledCurrency = :settledCurrency,
				accountAmount = :accountAmount,
				exchangeRate = :exchangeRate,
				merchant = :merchant,
				categoryId = :categoryId,
				description = :description,
				transactionAt = :transactionAt,
				updatedAt = :updatedAt
			WHERE transactionId = :transactionId
		`, map[string]interface{}{
			"transactionId":   transactionId,
			"settledAmount":   existing.SettledAmount,
			"settledCurrency": existing.SettledCurrency,
			"accountAmount":   existing.AccountAmount,
			"exchangeRate":    existing.ExchangeRate,
			"merchant":        existing.Merchant,
			"categoryId":      existing.CategoryId,
			"description":     existing.Description,
			"transactionAt":   existing.TransactionAt,
			"updatedAt":       time.Now().Unix(),
		})
		if err != nil {
			dbTx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update transaction"})
			return
		}

		if err := dbTx.Commit(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit"})
			return
		}
	} else {
		// Simple update — no amount change, no balance recalculation needed
		_, err = dbInstance.NamedExec(`
			UPDATE transactions SET
				merchant = :merchant,
				categoryId = :categoryId,
				description = :description,
				transactionAt = :transactionAt,
				updatedAt = :updatedAt
			WHERE transactionId = :transactionId
		`, map[string]interface{}{
			"transactionId": transactionId,
			"merchant":      existing.Merchant,
			"categoryId":    existing.CategoryId,
			"description":   existing.Description,
			"transactionAt": existing.TransactionAt,
			"updatedAt":     time.Now().Unix(),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update transaction"})
			return
		}
	}

	ctx.JSON(http.StatusOK, existing)
}
