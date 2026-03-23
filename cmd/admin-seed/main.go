package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/nakle1ka/Tramplin/internal/config"
	"github.com/nakle1ka/Tramplin/internal/database/postgres"
	"github.com/nakle1ka/Tramplin/internal/pkg/hash"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		slog.Error("init db: failed load config", "err", err)
		os.Exit(1)
	}

	db, err := postgres.NewPostgresDB(
		postgres.WithDbname(cfg.Database.Name),
		postgres.WithHost(cfg.Database.Host),
		postgres.WithPort(cfg.Database.Port),
		postgres.WithUser(cfg.Database.Username),
		postgres.WithPassword(cfg.Database.Password),
		postgres.WithSslmode(cfg.Database.SSLMode),
	)
	if err != nil {
		slog.Error("init db: failed connect to postgres", "err", err)
		os.Exit(1)
	}

	migrationErr := postgres.Migrate(db)
	if migrationErr != nil {
		slog.Error("init db: failed migrate db", "err", migrationErr)
		os.Exit(1)
	}

	th := hash.NewPasswordHasher()

	seedErr := postgres.SeedSuperAdmin(
		db,
		th,
		cfg.Database.SuperAdminEmail,
		cfg.Database.SuperAdminPassword,
	)

	if seedErr != nil && !errors.Is(seedErr, gorm.ErrDuplicatedKey) {
		slog.Error("failed seed super admin", "err", seedErr)
		os.Exit(1)
	}
}
