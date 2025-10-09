package handlers

import (
	"net/http"
	"todo-list/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserData(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	userID, ok := userIDVal.(string)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	rows, err := h.DB.Query(`SELECT id, name, email FROM users WHERE id=$1`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if rows.Next() {
		var user models.UserProfile
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
		return
	}
}
