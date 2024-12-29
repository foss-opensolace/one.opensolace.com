package dto

import (
	"time"

	"gorm.io/gorm"
)

type APIKeyLookup struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	Key         string `json:"key"`
	Description string `json:"description"`

	Owner *struct {
		ID          uint    `json:"id"`
		DisplayName *string `json:"display_name"`
		Username    string  `json:"username"`
    } `gorm:"-" json:"owner"`

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
	} `gorm:"-" json:"permissions"`

	TimesUsed uint  `json:"times_used"`
	MaxUsage  *uint `json:"max_usage"`

	// IsOkay displays, in a developer friendly way, is the key has reached its max usage, has expired or revoked.
	IsOkay         bool       `json:"is_active"`
	ExpirationDate *time.Time `json:"expiration_date"`
	RevokeReason   *string    `json:"revoke_reason"`
}
