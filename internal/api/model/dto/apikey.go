package dto

import (
	"time"

	"gorm.io/gorm"
)

type APIKeyPermissions struct {
	KeyAssign *bool `validate:"omitnil" json:"key_assign"`
	KeyCreate *bool `validate:"omitnil" json:"key_create"`
	KeyRead   *bool `validate:"omitnil" json:"key_read"`
	KeyUpdate *bool `validate:"omitnil" json:"key_update"`
	KeyRevoke *bool `validate:"omitnil" json:"key_revoke"`
	KeyDelete *bool `validate:"omitnil" json:"key_delete"`

	Health  *bool `validate:"omitnil" json:"health"`
	Metrics *bool `validate:"omitnil" json:"metrics"`

	UserAuthLogin    *bool `validate:"omitnil" json:"user_auth_login"`
	UserAuthRegister *bool `validate:"omitnil" json:"user_auth_register"`
	UserUpdate       *bool `validate:"omitnil" json:"user_update"`
	UserRead         *bool `validate:"omitnil" json:"user_read"`
	UserDelete       *bool `validate:"omitnil" json:"user_delete"`
}

type APIKeyCreate struct {
	OwnerID *int `validate:"omitnil,min=0" json:"owner_id"`
	APIKeyUpdate
}

type APIKeyUpdate struct {
	Description *string `validate:"omitnil,max=128" json:"description"`

	Permissions *APIKeyPermissions `gorm:"embedded;embeddedPrefix:permissions_" validate:"omitnil,dive" json:"permissions"`

	MaxUsage       *int       `validate:"omitnil,min=0,max=1000000" json:"max_usage"`
	ExpirationDate *time.Time `validate:"omitnil" json:"expiration_date"`
}

type APIKeyRevoke struct {
	Reason string `validate:"required,max=1000" json:"reason"`
}

type APIKeyLookup struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	Key         string `json:"key"`
	Description string `json:"description"`

	Owner *struct {
		ID          *uint   `json:"id"`
		DisplayName *string `json:"display_name"`
		Username    *string `json:"username"`
	} `gorm:"embedded;embeddedPrefix:owner_" json:"owner"`

	Permissions struct {
		KeyAssign bool `json:"key_assign"`
		KeyCreate bool `json:"key_create"`
		KeyRead   bool `json:"key_read"`
		KeyUpdate bool `json:"key_update"`
		KeyRevoke bool `json:"key_revoke"`
		KeyDelete bool `json:"key_delete"`

		Health  bool `json:"health"`
		Metrics bool `json:"metrics"`

		UserAuthLogin    bool `json:"user_auth_login"`
		UserAuthRegister bool `json:"user_auth_register"`
		UserUpdate       bool `json:"user_update"`
		UserRead         bool `json:"user_read"`
		UserDelete       bool `json:"user_delete"`
	} `gorm:"embedded;embeddedPrefix:permissions_" json:"permissions"`

	TimesUsed uint  `json:"times_used"`
	MaxUsage  *uint `json:"max_usage"`

	// CanUse displays, in a developer friendly way, is the key has reached its max usage, has expired or revoked.
	CanUse         bool       `json:"can_use"`
	ExpirationDate *time.Time `json:"expiration_date"`
	RevokeReason   *string    `json:"revoke_reason"`
}
