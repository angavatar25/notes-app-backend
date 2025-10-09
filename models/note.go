package models

type Note struct {
	ID        string `json:"id" db:"id"`
	Title     string `json:"title" db:"title" binding:"required"`
	BodyText  string `json:"bodyText" db:"bodyText" binding:"required"`
	NoteColor string `json:"notecolor" db:"notecolor"`
	LabelName string `json:"labelname" db:"labelname"`
	UserID    string `json:"userid" db:"user_id"`
}

type Labels struct {
	ID        string `json:"id" db:"id"`
	LabelName string `json:"labelname" db:"labelname"`
}
