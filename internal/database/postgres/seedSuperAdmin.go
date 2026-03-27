package postgres

import (
	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/pkg/hash"
	"gorm.io/gorm"
)

func SeedSuperAdmin(
	db *gorm.DB,
	th hash.Hasher,
	email, password string,
) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		hash, err := th.Hash([]byte(password))
		if err != nil {
			return err
		}

		user := &model.User{
			ID:           uuid.New(),
			Email:        email,
			PasswordHash: string(hash),
			Role:         model.RoleCurator,
			IsVerified:   true,
		}
		usrErr := tx.Create(user).Error
		if usrErr != nil {
			return usrErr
		}

		crtErr := tx.Create(&model.Curator{
			ID:           uuid.New(),
			UserID:       user.ID,
			IsSuperAdmin: true,
		}).Error
		if crtErr != nil {
			return crtErr
		}

		return nil
	})

	return err
}
