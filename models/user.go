package models

import "time"

type User struct {
	ID           uint      `gorm:"primary_key"`
	Nickname     string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	Password     string    `gorm:"not null"`
	Role         string    `gorm:"type:varchar(20);default:'user'"`
	AvatarURL    string    `gorm:"type:text"`
	RegisteredAt time.Time `gorm:"autoCreateTime"`
}
