package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	TaskName string `validate:"omitempty,ascii"`
	IsDone   bool   `gorm:"default:false" json:"isDone"`
	UserID   uint
}
type TaskResponse struct {
	TaskName string `validate:"omitempty,ascii"`
	IsDone   bool   `gorm:"default:false" json:"isDone"`
	UserID   uint
}
