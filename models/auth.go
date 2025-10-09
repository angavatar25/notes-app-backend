package models

type User struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name" binding:"required"`
	Email     string `json:"email" db:"email" binding:"required,email"`
	Passwords string `json:"passwords" db:"passwords" binding:"required"`
}

type UserProfile struct {
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

type LoginInput struct {
	Email     string `json:"email" binding:"required,email"`
	Passwords string `json:"passwords" binding:"required"`
}

type RegisterInput struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Passwords string `json:"passwords" binding:"required"`
}
