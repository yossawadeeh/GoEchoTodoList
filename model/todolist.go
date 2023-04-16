package model

import (
	"time"

	"gorm.io/gorm"
)

type Status struct {
	gorm.Model

	StatusName string `json:"status_name" binding:"required"`
}

type ToDoList struct {
	gorm.Model

	Discription string `json:"discription" binding:"required"`
	CreateDate  time.Time
	Status      Status `gorm:"foreignKey:StatusId"`
	StatusId    uint   `json:"status_id" binding:"required"`
	User        User   `gorm:"foreignKey:UserId"`
	UserId      uint   `json:"user_id" binding:"required"`
}
