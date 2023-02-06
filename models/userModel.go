package models

import (
	"gorm.io/gorm"
)

// The email is already validated through HTML
type User struct {
	gorm.Model
	Name     string `json:"name`
	Email    string `json:"email`
	Password string `json:"password"`
	Tasks    []Task `json:"tasks"`
}

type UserResponse struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Email string `json:"email"`
	Tasks []Task `json:"tasks"`
}
