package main

import (
	"database/sql"
	"fmt"
	"log"
	"todo-list/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB
var server = 5004

func main() {
	dsn := "postgres://postgres@localhost:5432/note-app?sslmode=disable"
	var err error

	db, err = sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal("❌ Error opening DB: ", err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	fmt.Println("Connected to the database successfully!")

	h := handlers.NewHandler(db)

	r := gin.Default()

	r.Use(cors.Default())

	notes := r.Group("/notes")
	{
		notes.GET("/lists", h.GetNoteList)
		notes.GET("/get/:id", h.GetNoteByID)
		notes.GET("/labels", h.GetLabelList)
		notes.POST("/create", h.CreateNote)
		notes.PUT("/update/:id", h.UpdateNote)
		notes.DELETE("/delete/:id", h.DeleteNote)
	}

	fmt.Printf("🚀 Server running on http://localhost:%d\n", server)

	r.Run(fmt.Sprintf(":%d", server))
}
