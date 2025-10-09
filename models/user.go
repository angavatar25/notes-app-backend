package models

type UserSettings struct {
	DarkMode bool `json:"darkmode" db:"darkmode"`
}
