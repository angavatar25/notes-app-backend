package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"todo-list/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) GetNoteList(c *gin.Context) {
	label := c.Query("label")
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	fmt.Println(label, "label")

	if label == "All" || label == "" {
		rows, err := h.DB.Query(
			`SELECT id, title, bodyText, noteColor, labelname, created_at FROM note where userid=$1 ORDER BY created_at DESC`, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var notes []models.Note
		for rows.Next() {
			var n models.Note
			if err := rows.Scan(&n.ID, &n.Title, &n.BodyText, &n.NoteColor, &n.LabelName, &n.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			notes = append(notes, n)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, notes)
		return
	}

	rows, err := h.DB.Query(`select id, title, bodytext, notecolor, labelname, created_at from note where labelname like $1 and userid=$2`, label, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.BodyText, &n.NoteColor, &n.LabelName, &n.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		notes = append(notes, n)
	}
	c.JSON(http.StatusOK, notes)
}

func (h *Handler) GetNoteByID(c *gin.Context) {
	id := c.Param("id")

	rows, err := h.DB.Query(`SELECT id, title, bodyText, noteColor, labelname FROM note WHERE id=$1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.BodyText, &n.NoteColor, &n.LabelName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, n)
		return
	}
}

func (h *Handler) CreateNote(c *gin.Context) {
	var note models.Note

	userIDVal, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	if err := c.BindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payloads"})
		return
	}

	userID, ok := userIDVal.(string)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	query := `insert into note (title, bodyText, notecolor, labelname, userid) values ($1, $2, $3, $4, $5) returning id`

	var id string
	err := h.DB.QueryRow(query, note.Title, note.BodyText, note.NoteColor, note.LabelName, userID).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	note.ID = id
	c.JSON(http.StatusOK, gin.H{"message": "Note created successfully", "id": id})
}

func (h *Handler) UpdateNote(c *gin.Context) {
	var note models.Note
	id := c.Param("id")

	if err := c.BindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := `update note set title=$1, bodyText=$2, notecolor=$3, labelname=$4 where id=$5`

	_, err := h.DB.Exec(query, note.Title, note.BodyText, note.NoteColor, note.LabelName, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note updated successfully"})
}

func (h *Handler) DeleteNote(c *gin.Context) {
	id := c.Param("id")

	query := `delete from note where id=$1`
	_, err := h.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func (h *Handler) GetLabelList(c *gin.Context) {
	rows, err := h.DB.Query(`SELECT id, labelname FROM labels`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var labels []models.Labels

	for rows.Next() {
		var l models.Labels
		if err := rows.Scan(&l.ID, &l.LabelName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		labels = append(labels, l)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, labels)
}

func (h *Handler) GetNotesByLabel(c *gin.Context) {
	label := c.Param("label")
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	rows, err := h.DB.Query(`select id, title, bodytext, notecolor, labelname, created_at from note where labelname like $1 and userid=$2`, label, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.BodyText, &n.NoteColor, &n.LabelName, &n.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		notes = append(notes, n)
	}
	c.JSON(http.StatusOK, notes)
}
