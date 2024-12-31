package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type APIKey struct {
	gorm.Model

	Key         uuid.UUID `gorm:"unique,not null,type:uuid;default:uuid_generate_v4()"`
	Description string    `gorm:"default:'No description provided'"`

	OwnerID *uint

	TimesUsed uint `gorm:"default:0"`
	MaxUsage  *int

	ExpirationDate *time.Time
	RevokeReason   *string
}

type APIKeyPermissions struct {
	KeyID uint `gorm:"unique,not null"`

	KeyAssign bool `gorm:"default:false"`
	KeyCreate bool `gorm:"default:false"`
	KeyRead   bool `gorm:"default:false"`
	KeyUpdate bool `gorm:"default:false"`
	KeyRevoke bool `gorm:"default:false"`
	KeyDelete bool `gorm:"default:false"`

	Health  bool `gorm:"default:false"`
	Metrics bool `gorm:"default:false"`

	UserAuthLogin    bool `gorm:"default:false"`
	UserAuthRegister bool `gorm:"default:false"`
	UserUpdate       bool `gorm:"default:false"`
	UserRead         bool `gorm:"default:false"`
	UserDelete       bool `gorm:"default:false"`
}
