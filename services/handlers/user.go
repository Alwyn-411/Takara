package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"takara/services/schema"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GetUserById(ctx *gin.Context, dbInstance *sqlx.DB) {
	id := ctx.Param("id")

	query := `
		SELECT 
			userId, active, userName,
			altName, email, altEmail,
			password, createdAt, updatedAt
		FROM users
		WHERE userId = ? AND active = ?
	`

	var user schema.User
	err := dbInstance.Get(&user, query, id, true)

	if err == sql.ErrNoRows {
		ctx.AbortWithError(http.StatusNotFound, sql.ErrNoRows)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func CreateUser(ctx *gin.Context, dbInstance *sqlx.DB) {
	id := uuid.New().String()

	var data schema.User

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	data.UserID = id

	// Hash password
	hashedPassword, err := HashPassword(data.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	data.Password = string(hashedPassword)

	insertQuery := `
		INSERT INTO users (
			userId, userName, altName,
			email, altEmail, password
		) VALUES (
			:userId, :userName, :altName,
			:email, :altEmail, :password
		)
	`

	_, err = dbInstance.NamedExec(insertQuery, data)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"error": "username or email already exists",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": data.UserID,
	})
}

func UpdateUserById(ctx *gin.Context, dbInstance *sqlx.DB) {
	id := ctx.Param("id")

	var data map[string]any

	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	set := []string{}
	args := []any{}

	allowedFields := map[string]bool{
		"userName": true,
		"altName":  true,
		"email":    true,
		"altEmail": true,
		"password": true,
	}

	for key, value := range data {
		if !allowedFields[key] {
			continue
		}

		if key == "password" {
			password, ok := value.(string)
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
					"error": "invalid password",
				})
				return
			}

			hashedPassword, err := HashPassword(password)
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			value = string(hashedPassword)
		}

		set = append(set, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}

	if len(set) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "no valid fields provided",
		})
		return
	}

	set = append(set, "updatedAt = ?")
	args = append(args, time.Now().Unix())

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE userId = ?
	`, strings.Join(set, ", "))

	result, err := dbInstance.Exec(query, args...)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("user with userId = %s does not exist", id),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func DeleteUserById(ctx *gin.Context, dbInstance *sqlx.DB) {
	id := ctx.Param("id")

	query := `
		UPDATE users
		SET active = :active
		WHERE userId = :userId
	`

	result, err := dbInstance.NamedExec(query, gin.H{
		"userId": id,
		"active": false,
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func AuthorizeUserWithUserNameAndPassword(ctx *gin.Context, dbInstance *sqlx.DB) {
	enteredUserName := ctx.Query("userName")
	enteredUserPassword := ctx.Query("password")

	query := `
		SELECT 
			userId, userName, password
		FROM users
		WHERE userName = ?
	`

	var user schema.User

	err := dbInstance.Get(&user, query, enteredUserName)

	if errors.Is(err, sql.ErrNoRows) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = CheckPassword(user.Password, enteredUserPassword)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": user.UserID,
	})
}