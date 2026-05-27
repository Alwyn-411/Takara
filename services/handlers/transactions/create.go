package transactions

import (
	"fmt"
	"net/http"
	"takara/services/handlers"
	"takara/services/schema"
	"time"

	"github.com/bojanz/currency"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTransactionRequest struct {
	UserId          string   `json:"userId" binding:"required"`
	AccountId       string   `json:"accountId" binding:"required"`
	Type            string   `json:"type" binding:"required"`
	SettledAmount   string   `json:"settledAmount" binding:"required"`
	SettledCurrency string   `json:"settledCurrency" binding:"required"`
	MerchantName    string   `json:"merchantName,omitempty"`
	CategoryName    string   `json:"categoryName,omitempty"`
	Description     string   `json:"description,omitempty"`
	TagNames        []string `json:"tagNames,omitempty"`
	TransactionAt   int64    `json:"transactionAt" binding:"required"`
}

func (handler *TransactionsHandler) CreateTransaction(ctx *gin.Context) {
	var data CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if !IsValidType(data.Type) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid type %q: must be Debit or Credit", data.Type)})
		return
	}

	settledAmount, err := currency.NewAmount(data.SettledAmount, data.SettledCurrency)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount or currency"})
		return
	}

	// Load the target account
	var account schema.Account
	if err := handler.dbInstance.Get(&account, `
		SELECT userId, accountId, type, name, accountNumber, description, currency, interest, balance, active
		FROM accounts
		WHERE accountId = ? AND userId = ? AND active = 1
	`, data.AccountId, data.UserId); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	// Convert settled amount to account currency
	var accountAmount currency.Amount
	var rateUsed string

	if data.SettledCurrency != account.Currency {
		rate, err := handler.forEx.GetRate(data.SettledCurrency, account.Currency)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "could not fetch exchange rate"})
			return
		}
		accountAmount, err = settledAmount.Convert(account.Currency, rate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "conversion failed"})
			return
		}
		accountAmount = accountAmount.Round()
		rateUsed = rate
	} else {
		accountAmount = settledAmount
		rateUsed = "1"
	}

	// Begin DB transaction — everything below is atomic
	dbTx, err := handler.dbInstance.Beginx()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start db transaction"})
		return
	}

	// Resolve names → IDs (create if new)
	merchantId, err := handlers.ResolveOrCreateMerchant(dbTx, data.UserId, data.MerchantName)
	if err != nil {
		dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve merchant"})
		return
	}

	categoryId, err := handlers.ResolveOrCreateCategory(dbTx, data.UserId, data.CategoryName)
	if err != nil {
		dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve category"})
		return
	}

	tagIds := make([]string, 0, len(data.TagNames))
	for _, tagName := range data.TagNames {
		tagId, err := handlers.ResolveOrCreateTag(dbTx, data.UserId, tagName)
		if err != nil {
			dbTx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve tag: " + tagName})
			return
		}
		tagIds = append(tagIds, tagId)
	}

	now := time.Now().Unix()
	transactionId := uuid.New().String()

	args := schema.Transaction{
		Base:            schema.Base{UserID: data.UserId, Active: 1},
		TransactionId:   transactionId,
		AccountId:       data.AccountId,
		Type:            data.Type,
		SettledAmount:   settledAmount.Number(),
		SettledCurrency: data.SettledCurrency,
		AccountAmount:   accountAmount.Number(),
		AccountCurrency: account.Currency,
		ExchangeRate:    rateUsed,
		MerchantId:      merchantId,
		CategoryId:      categoryId,
		Description:     data.Description,
		TransactionAt:   data.TransactionAt,
	}

	// Insert transaction
	if _, err = dbTx.NamedExec(`
		INSERT INTO transactions (
			transactionId, userId, accountId, type,
			settledAmount, settledCurrency,
			accountAmount, accountCurrency,
			exchangeRate, merchantId, categoryId,
			description, transactionAt
		) VALUES (
			:transactionId, :userId, :accountId, :type,
			:settledAmount, :settledCurrency,
			:accountAmount, :accountCurrency,
			:exchangeRate, :merchantId, :categoryId,
			:description, :transactionAt
		)`, args); err != nil {
		dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert transaction"})
		return
	}

	// Bind tags
	for _, tagId := range tagIds {
		if _, err = dbTx.Exec(
			`INSERT INTO transaction_tags (transactionId, tagId) VALUES (?, ?)`,
			transactionId, tagId,
		); err != nil {
			dbTx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to bind tag"})
			return
		}
	}

	// Apply to account balance
	if err = applyToBalance(dbTx, args.AccountId, args.AccountCurrency, args.Type, args.AccountAmount, now); err != nil {
		dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update balance"})
		return
	}

	if err := dbTx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": transactionId})
}
