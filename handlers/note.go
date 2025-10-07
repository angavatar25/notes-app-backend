package handlers

import (
	"database/sql"
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
	rows, err := h.DB.Query(`SELECT title, bodyText, noteColor FROM note`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.Title, &n.BodyText, &n.NoteColor); err != nil {
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
}

func (h *Handler) CreateNote(c *gin.Context) {
	var note models.Note

	if err := c.BindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payloads"})
		return
	}

	query := `insert into note (title, bodyText, notecolor) values ($1, $2, $3) returning id`

	var id string
	err := h.DB.QueryRow(query, note.Title, note.BodyText, note.NoteColor).Scan(&id)

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

	query := `update note set title=$1, bodyText=$2, notecolor=$3 where id=$4`

	_, err := h.DB.Exec(query, note.Title, note.BodyText, note.NoteColor, id)
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
