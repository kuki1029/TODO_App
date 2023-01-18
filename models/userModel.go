package models

import (
	"gorm.io/gorm"
)

// The email is already validated through HTML
type User struct {
	gorm.Model
	ID       uint   `json:"id"`
	Email    string `json:"email`
	Password string `json:"password"`
	Tasks    []Task `json:"tasks"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Tasks []Task `json:"tasks"`
}
