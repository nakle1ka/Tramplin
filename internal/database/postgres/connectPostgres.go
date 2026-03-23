package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(opts ...opt) (*gorm.DB, error) {
	postgesCfg := defaultPostgresConfig()

	for _, opt := range opts {
		opt(&postgesCfg)
	}

	sslmode := "dasable"
	if postgesCfg.Sslmode {
		sslmode = "prefer"
	}

	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
		postgesCfg.Host,
		postgesCfg.User,
		postgesCfg.Password,
		postgesCfg.DBname,
		postgesCfg.Port,
		sslmode,
	)

	postges, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		return nil, err
	}

	return postges, nil
}
