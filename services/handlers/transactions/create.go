package transactions

import (
	"fmt"
	"net/http"
	"takara/services/schema"
	"time"

	"github.com/bojanz/currency"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func (handler *TransactionsHandler) CreateTransaction(ctx *gin.Context) {
	/*
		1. Parse Data from Body
		2. Query From Account Remaining Balance
		3. Convert Amount to account currency with forex
		4. Begin transaction
		5. Create Row with NEW data in transactions table
		6. Bind Tags to Transactions
		7. Update Account Balance Value
		8. Wrap up
	*/
	var data CreateTransactionRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	settledAmount, err := currency.NewAmount(data.SettledAmount, data.SettledCurrency)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	getAccountQuery := `
		SELECT userId, accountId, type, name, accountNumber, description, currency, interest, balance, active
		FROM accounts
		WHERE accountId = ? AND active = 1
	`

	var account schema.Account
	err = handler.dbInstance.Get(&account, getAccountQuery, data.AccountId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("No Account | Error: %s", err.Error()))
		return
	}

	var accountAmount currency.Amount
	var rateUsed string

	if data.SettledCurrency != account.Currency {
		rate, err := handler.forEx.GetRate(data.SettledCurrency, account.Currency)
		if err != nil {
			ctx.AbortWithError(http.StatusBadGateway, err)
			return
		}

		accountAmount, err = settledAmount.Convert(account.Currency, rate)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
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

	dbTx, err := handler.dbInstance.Beginx()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
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
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(data.TagIds) > 0 {
		for _, tagId := range data.TagIds {
			_, err = dbTx.Exec(
				`INSERT INTO transaction_tags (transactionId, tagId) VALUES (?, ?)`,
				transactionId, tagId,
			)
			if err != nil {
				dbTx.Rollback()
				ctx.AbortWithError(http.StatusInternalServerError, err)
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
			SET balance = printf('%%.2f', CAST(balance AS REAL) %s CAST(? AS REAL)), updatedAt = ?
			WHERE accountId = ?`, sign),
		args.AccountAmount, time.Now().Unix(), args.AccountId,
	)
	if err != nil {
		dbTx.Rollback()
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := dbTx.Commit(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": transactionId,
	})
}
