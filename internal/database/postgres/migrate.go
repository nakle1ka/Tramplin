package postgres

import (
	"github.com/nakle1ka/Tramplin/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Applicant{},
		&model.Employer{},
		&model.Curator{},
		&model.Opportunity{},
		&model.Tag{},
		&model.Application{},
		&model.Contact{},
		&model.Recomendation{},
		&model.Favorites{},
	)
}
