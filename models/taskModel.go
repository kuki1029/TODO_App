package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID        uint
	TaskName  string `validate:"omitempty,ascii"`
	Assignee  string
	CreatedAt time.Time
	IsDone    bool `gorm:"default:false" json:"isDone"`
	UserID    uint
}
type TaskResponse struct {
	ID        uint
	TaskName  string `validate:"omitempty,ascii"`
	Assignee  string
	CreatedAt time.Time
	IsDone    bool `gorm:"default:false" json:"isDone"`
	UserID    uint
}
