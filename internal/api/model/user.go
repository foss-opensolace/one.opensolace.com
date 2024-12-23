package model

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	DisplayName *string `gorm:"index:,composite:name"`
	Username    string  `gorm:"index:,not null,unique,composite:name"`
	Email       string  `gorm:"index:,not null,unique"`
	Password    string  `gorm:"not null"`
}

func (u *User) ToSafe() fiber.Map {
	return fiber.Map{
		"id":           u.ID,
		"created_at":   u.CreatedAt,
		"updated_at":   u.UpdatedAt,
		"display_name": u.DisplayName,
		"username":     u.Username,
	}
}
