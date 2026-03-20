package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/config"
	"github.com/nakle1ka/Tramplin/internal/handler"
	"github.com/nakle1ka/Tramplin/internal/middleware"
	"github.com/nakle1ka/Tramplin/internal/pkg/auth"
	"github.com/nakle1ka/Tramplin/internal/pkg/hash"
	"github.com/nakle1ka/Tramplin/internal/repository"
	"github.com/nakle1ka/Tramplin/internal/routes"
	"github.com/nakle1ka/Tramplin/internal/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	cfg   *config.Config
	db    *gorm.DB
	cache *redis.Client
}

func (a *App) Run() error {
	tokenHasher := hash.NewTokenHasher()
	passwordHasher := hash.NewPasswordHasher()

	tokenManager := auth.NewTokenManager(a.cfg.JWT.SecretKey)
	transactionManager := repository.NewTransactionManager(a.db)

	applicantRepo := repository.NewApplicantRepository(a.db)
	employerRepo := repository.NewEmployerRepository(a.db)
	curatotRepo := repository.NewCuratorRepository(a.db)
	userRepo := repository.NewUserRepository(a.db)
	cacheRepo := repository.NewCacheRepository(a.cache)

	authSrv := service.NewAuthService(
		userRepo,
		applicantRepo,
		curatotRepo,
		employerRepo,
		cacheRepo,
		transactionManager,
		tokenHasher,
		passwordHasher,
		tokenManager,

		service.WithAccessExpires(a.cfg.JWT.AccessTokenLifeTime),
		service.WithRefreshExpires(a.cfg.JWT.RefreshTokenLifeTime),
	)

	authHnd := handler.NewAuthHandler(authSrv)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	v1 := router.Group("/api/v1")

	routes.SetupAuthRoutes(v1, authHnd)

	protected := v1.Group("/")
	protected.Use(middleware.JWTAuth(tokenManager))

	addr := fmt.Sprintf(":%v", a.cfg.App.Port)
	return router.Run(addr)
}

func NewApp(cfg *config.Config, db *gorm.DB, cache *redis.Client) *App {
	return &App{
		cfg:   cfg,
		db:    db,
		cache: cache,
	}
}
