package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"todo-list/models"
	"todo-list/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) UserLogin(c *gin.Context) {
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := h.DB.QueryRow(
		"SELECT id, name, email, passwords FROM users WHERE email=$1",
		input.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Passwords)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		fmt.Println("Invalid email or password in sql no rows")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Passwords), []byte(input.Passwords))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		fmt.Println("Invalid email or password in compare hash")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func (h *Handler) UserRegister(c *gin.Context) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing string
	err := h.DB.QueryRow("SELECT email FROM users WHERE email=$1", input.Email).Scan(&existing)
	if err != sql.ErrNoRows {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Passwords), bcrypt.DefaultCost)

	_, err = h.DB.Exec(
		"INSERT INTO users (name, email, passwords) VALUES ($1, $2, $3)",
		input.Name, input.Email, string(hashedPassword),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
