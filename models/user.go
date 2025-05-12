package models

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	Nickname     string    `json:"nickname" gorm:"unique" validate:"required,min=4"`
	Email        string    `json:"email" gorm:"unique" validate:"required,email"`
	RegisteredAt time.Time `json:"registered_at" gorm:"autoCreateTime"`
}
