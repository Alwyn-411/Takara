package handlers

import (
	"database/sql"
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
		SELECT user_id, altEmail, altName, email, name, createdTimeStamp, updatedTimeStamp, active
		FROM users
		WHERE user_id = ? AND active = ?
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
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	data.UserID = id

	query := `INSERT 
	INTO users (
		user_id, altEmail, altName, email, 
		name, passHash, createdTimeStamp, 
		updatedTimeStamp
	) VALUES (
		:user_id, :altEmail, :altName, :email, 
		:name, :passHash, :createdTimeStamp, 
		:updatedTimeStamp
	 )`

	_, err = dbInstance.NamedExec(query, data)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": data.UserID})
}

func UpdateUserById(ctx *gin.Context, dbInstance *sqlx.DB) {
	id := ctx.Param("id")

	var data map[string]any

	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	set := []string{}
	args := []any{}

	allowedFields := map[string]bool{
		"altEmail": true,
		"altName":  true,
		"email":    true,
		"name":     true,
		"passHash": true,
	}

	for key, value := range data {
		if !allowedFields[key] {
			continue
		}
		set = append(set, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}

	set = append(set, "updatedTimeStamp = ?")
	args = append(args, time.Now())

	args = append(args, id)

	query := fmt.Sprintf(`UPDATE users SET %s WHERE user_id = ?`, strings.Join(set, ", "))

	result, err := dbInstance.Exec(query, args...)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if rows == 0 {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("User with user_id = %s does not Exist", id))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func DeleteUserById(ctx *gin.Context, dbInstance *sqlx.DB) {
	id := ctx.Param("id")

	query := `
		UPDATE users SET active = :active WHERE user_id = :id
	`
	_, err := dbInstance.NamedExec(query, gin.H{"user_id": id, "active": false})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
