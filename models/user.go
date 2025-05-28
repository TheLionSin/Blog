package models

import "time"

type User struct {
	ID           uint      `gorm:"primary_key"`
	Nickname     string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	Password     string    `gorm:"not null"`
	RegisteredAt time.Time `gorm:"autoCreateTime"`
}
