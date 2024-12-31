package repository

import (
	"github.com/foss-opensolace/api.opensolace.com/internal/api/model"
	"github.com/foss-opensolace/api.opensolace.com/internal/api/model/dto"
	"github.com/foss-opensolace/api.opensolace.com/pkg/utils"
	"gorm.io/gorm"
)

type apiKeyRepository struct {
	db *gorm.DB
}

type APIKeyRepository interface {
	Create(registerOptions ...dto.APIKeyCreate) (*dto.APIKeyLookup, error)
}

func NewAPIKeyRepository(postgres *gorm.DB) APIKeyRepository {
	return &apiKeyRepository{db: postgres}
}

func (a *apiKeyRepository) baseFetchQuery() *gorm.DB {
	return a.db.Table("api_keys ak").
		Select(`
            ak.id AS id,
            ak.created_at AS created_at,
            ak.updated_at AS updated_at,
            ak.deleted_at AS deleted_at,
            ak.key AS key,
            ak.description AS description,
            u.id as owner_id,
            u.display_name as owner_display_name,
            u.username as owner_username,
            p.key_assign as permissions_key_assign,
            p.key_create as permissions_key_create,
            p.key_read as permissions_key_read,
            p.key_update as permissions_key_update,
            p.key_revoke as permissions_key_revoke,
            p.key_delete as permissions_key_delete,
            p.health as permissions_health,
            p.metrics as permissions_metrics,
            p.user_auth_login as permissions_user_auth_login,
            p.user_auth_register as permissions_user_auth_register,
            p.user_update as permissions_user_update,
            p.user_read as permissions_user_read,
            p.user_delete as permissions_user_delete,
            ak.times_used AS times_used,
            ak.max_usage AS max_usage,
            (
                CASE WHEN ak.times_used >= ak.max_usage
                    OR ak.expiration_date >= CURRENT_DATE
                    OR ak.revoke_reason IS NOT NULL THEN FALSE ELSE TRUE END
            ) AS can_use,
            ak.expiration_date as expiration_date,
            ak.revoke_reason as revoke_reason
		`).
		Joins("LEFT JOIN users u ON ak.owner_id = u.id").
		Joins("LEFT JOIN api_key_permissions p ON ak.id = p.key_id")
}

func (a *apiKeyRepository) Create(registerOptions ...dto.APIKeyCreate) (*dto.APIKeyLookup, error) {
	transaction := a.db.Begin()

	var apiKey model.APIKey
	var apiKeyLookup dto.APIKeyLookup
	var options dto.APIKeyCreate

	if len(registerOptions) > 0 {
		options = registerOptions[0]
	}

	if options.Description != nil {
		apiKey.Description = *options.Description
	}
	if options.MaxUsage != nil {
		apiKey.MaxUsage = options.MaxUsage
	}
	if options.ExpirationDate != nil {
		apiKey.ExpirationDate = options.ExpirationDate
	}

	if options.OwnerID != nil {
		if err := transaction.Model(&model.User{}).First("id = ?", options.OwnerID).Error; err != nil {
			transaction.Rollback()
			return nil, err
		}

		apiKey.OwnerID = utils.ToPtr(uint(*options.OwnerID))
	}

	if err := transaction.Create(&apiKey).Error; err != nil {
		transaction.Rollback()
		return nil, err
	}

	if options.Permissions != nil {
		opt := model.APIKeyPermissions{
			KeyID: apiKey.ID,
		}

		if options.Permissions.KeyAssign != nil {
			opt.KeyAssign = *options.Permissions.KeyAssign
		}
		if options.Permissions.KeyCreate != nil {
			opt.KeyCreate = *options.Permissions.KeyCreate
		}
		if options.Permissions.KeyRead != nil {
			opt.KeyRead = *options.Permissions.KeyRead
		}
		if options.Permissions.KeyUpdate != nil {
			opt.KeyUpdate = *options.Permissions.KeyUpdate
		}
		if options.Permissions.KeyRevoke != nil {
			opt.KeyRevoke = *options.Permissions.KeyRevoke
		}
		if options.Permissions.KeyDelete != nil {
			opt.KeyDelete = *options.Permissions.KeyDelete
		}
		if options.Permissions.Health != nil {
			opt.Health = *options.Permissions.Health
		}
		if options.Permissions.Metrics != nil {
			opt.Metrics = *options.Permissions.Metrics
		}
		if options.Permissions.UserAuthLogin != nil {
			opt.UserAuthLogin = *options.Permissions.UserAuthLogin
		}
		if options.Permissions.UserAuthRegister != nil {
			opt.UserAuthRegister = *options.Permissions.UserAuthRegister
		}
		if options.Permissions.UserUpdate != nil {
			opt.UserUpdate = *options.Permissions.UserUpdate
		}
		if options.Permissions.UserRead != nil {
			opt.UserRead = *options.Permissions.UserRead
		}
		if options.Permissions.UserDelete != nil {
			opt.UserDelete = *options.Permissions.UserDelete
		}

		if err := transaction.Create(&opt).Error; err != nil {
			transaction.Rollback()
			return nil, err
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return nil, err
	}

	if err := a.baseFetchQuery().Where("ak.id = ?", apiKey.ID).Scan(&apiKeyLookup).Error; err != nil {
		return nil, err
	}

	if apiKeyLookup.Owner.ID == nil {
		apiKeyLookup.Owner = nil
	}

	return &apiKeyLookup, nil
}
