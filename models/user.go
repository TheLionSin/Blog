package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint      `gorm:"primary_key"`
	Nickname     string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	Password     string    `gorm:"not null"`
	Role         string    `gorm:"type:varchar(20);default:'user'"`
	AvatarURL    string    `gorm:"type:text"`
	RegisteredAt time.Time `gorm:"autoCreateTime"`

	gorm.DeletedAt `gorm:"index"`
}

type AuditLog struct {
	ID        uint      `gorm:"primary_key"`
	UserID    uint      // кто выполнил
	Action    string    // тип действия (delete_user, update_email и т.д.)
	Object    string    // над каким типом сущности (user, avatar, team)
	ObjectID  uint      // ID объекта
	Timestamp time.Time `gorm:"autoCreateTime"`

	IP        string
	UserAgent string
	Metadata  string
}
