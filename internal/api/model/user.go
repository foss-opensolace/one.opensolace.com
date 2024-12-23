package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	DisplayName *string `gorm:"index:,composite:name"`
	Username    string  `gorm:"index:,not null,unique,composite:name"`
	Email       string  `gorm:"index:,not null,unique"`
	Password    string  `gorm:"not null"`
}
