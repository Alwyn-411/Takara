package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"takara/services/schema"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	dbInstance *sqlx.DB
}

func NewAuthHandler(db *sqlx.DB) *AuthHandler {
	return &AuthHandler{
		dbInstance: db,
	}
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func checkPassword(Hashed string, Plaintext string) error {
	return bcrypt.CompareHashAndPassword([]byte(Hashed), []byte(Plaintext))
}

func (handler *AuthHandler) AuthorizeUserWithUserNameAndPassword(ctx *gin.Context) {
	enteredUserName := ctx.Query("userName")
	enteredUserPassword := ctx.Query("password")

	query := `
		SELECT userId, userName, password
		FROM users
		WHERE userName = ?
	`

	var user schema.User

	err := handler.dbInstance.Get(&user, query, enteredUserName)

	if errors.Is(err, sql.ErrNoRows) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = checkPassword(user.Password, enteredUserPassword)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": user.UserID,
	})
}
