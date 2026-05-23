package handlers

import (
	"fmt"
	"net/http"
	"takara/services/forex"
	"takara/services/schema"
	"time"

	"github.com/bojanz/currency"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// type Transaction struct {
// 	Base

// 	AccountId     string `db:"accountId" json:"accountId"`
// 	TransactionId string `db:"transactionId" json:"transactionId"`

// 	Type string `db:"type" json:"type"` // Debit | Credit

// 	// other party billed or settled amount
// 	SettledAmount   string `db:"settledAmount" json:"settledAmount"`
// 	SettledCurrency string `db:"settledCurrency" json:"settledCurrency"`

// 	// source account amount
// 	AccountAmount   string `db:"accountAmount" json:"accountAmount"`     // Cost in account currency
// 	AccountCurrency string `db:"accountCurrency" json:"accountCurrency"` // ISO 4217 e.g. "INR"

// 	ExchangeRate string `db:"exchangeRate" json:"exchangeRate"` // Conversion rate at transaction time

// 	Merchant      string `db:"merchant" json:"merchant"` // Name of the Merchant
// 	CategoryId    string `db:"categoryId" json:"categoryId"`
// 	Description   string `db:"description" json:"description,omitempty"`
// 	TransactionAt int64  `db:"transactionAt" json:"transactionAt"` // Transaction Timestamp

// 	Timestamp
// }

type CreateTransactionRequest struct {
	UserId          string   `json:"userId" binding:"required"`
	AccountId       string   `json:"accountId" binding:"required"`
	Type            string   `json:"type" binding:"required"`            // Debit | Credit
	SettledAmount   string   `json:"settledAmount" binding:"required"`   // "30.00"
	SettledCurrency string   `json:"settledCurrency" binding:"required"` // "USD"
	Merchant        string   `json:"merchant" binding:"required"`
	TransactionAt   int64    `json:"transactionAt" binding:"required"`
	CategoryId      string   `json:"categoryId,omitempty"`
	Description     string   `json:"description,omitempty"`
	TagIds          []string `json:"tagIds,omitempty"`
}

func CreateTransaction(ctx *gin.Context, dbInstance *sqlx.DB, forexService *forex.Frankfurter) {
	var data CreateTransactionRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	settledAmount, err := currency.NewAmount(data.SettledAmount, data.SettledCurrency)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount or currency"})
		return
	}

	getAccountQuery := `
		SELECT userId, accountId, type, name, accountNumber, description, currency, interest, balance, active
		FROM accounts
		WHERE accountId = ? AND active = 1
	`

	var account schema.Account
	err = dbInstance.Get(&account, getAccountQuery, data.AccountId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("No Account | Error: %s", err.Error()))
		return
	}

	var accountAmount currency.Amount
	var rateUsed string

	if data.SettledCurrency != account.Currency {
		rate, err := forexService.GetRate(data.SettledCurrency, account.Currency)
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
		Merchant:        data.Merchant,
		CategoryId:      data.CategoryId,
		Description:     data.Description,
		TransactionAt:   data.TransactionAt,
	}

	dbTx, err := dbInstance.Beginx()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start db transaction"})
		return
	}

	_, err = dbTx.NamedExec(`
		INSERT INTO transactions (
			transactionId, userId, accountId, type,
			settledAmount, settledCurrency,
			accountAmount, accountCurrency,
			exchangeRate, merchant, categoryId,
			description, transactionAt
		) VALUES (
			:transactionId, :userId, :accountId, :type,
			:settledAmount, :settledCurrency,
			:accountAmount, :accountCurrency,
			:exchangeRate, :merchant, :categoryId,
			:description, :transactionAt
		)`, args)

	if err != nil {
		dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert transaction"})
		return
	}

	// Insert tags if provided
	if len(data.TagIds) > 0 {
		for _, tagId := range data.TagIds {
			_, err = dbTx.Exec(
				`INSERT INTO transaction_tags (transactionId, tagId) VALUES (?, ?)`,
				transactionId, tagId,
			)
			if err != nil {
				dbTx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to link tag: " + tagId})
				return
			}
		}
	}

	sign := "-"
	if args.Type == "Credit" {
		sign = "+"
	}
	_, err = dbTx.Exec(
		fmt.Sprintf(`
			UPDATE accounts
			SET balance = printf('%%.2f', CAST(balance AS REAL) %s CAST(? AS REAL)),
				updatedAt = ?
			WHERE accountId = ?`, sign),
		args.AccountAmount, time.Now().Unix(), args.AccountId,
	)
	if err != nil {
		dbTx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update balance"})
		return
	}

	if err := dbTx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit"})
		return
	}

	ctx.JSON(http.StatusCreated, args)
}
