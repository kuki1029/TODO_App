package models

import (
	"gorm.io/gorm"
)

// The email is already validated through HTML
type User struct {
	gorm.Model
	Email    string `json:"email`
	Password string `json:"password"`
	Tasks    []Task `json:"tasks"`
}

type UserResponse struct {
	Email string `json:"email"`
	Tasks []Task `json:"tasks"`
}
