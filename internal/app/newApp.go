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
	curatorRepo := repository.NewCuratorRepository(a.db)
	userRepo := repository.NewUserRepository(a.db)
	opportunityRepo := repository.NewOpportunityRepository(a.db)
	tagRepo := repository.NewTagRepository(a.db)
	applicationRepo := repository.NewApplicationRepository(a.db)
	contactRepo := repository.NewContactRepository(a.db)
	recomendationRepo := repository.NewRecomendationRepository(a.db)
	favoritesRepo := repository.NewFavoritesRepository(a.db)
	cacheRepo := repository.NewCacheRepository(a.cache)

	authSrv := service.NewAuthService(
		userRepo,
		applicantRepo,
		curatorRepo,
		employerRepo,
		cacheRepo,
		transactionManager,
		tokenHasher,
		passwordHasher,
		tokenManager,

		service.WithAccessExpires(a.cfg.JWT.AccessTokenLifeTime),
		service.WithRefreshExpires(a.cfg.JWT.RefreshTokenLifeTime),
	)
	applicantSrv := service.NewApplicantService(applicantRepo)
	employerSrv := service.NewEmployerService(employerRepo)
	curatorSrv := service.NewCuratorService(curatorRepo)
	opportunitySrv := service.NewOpportunityService(opportunityRepo, tagRepo, employerRepo, curatorRepo)
	applicationSrv := service.NewApplicationService(applicationRepo, opportunityRepo, applicantRepo, employerRepo)
	contactSrv := service.NewContactService(contactRepo, applicantRepo)
	recomendationSrv := service.NewRecomendationService(recomendationRepo, applicantRepo, opportunityRepo, contactRepo)
	favoritesSrv := service.NewFavoritesService(favoritesRepo, applicantRepo, opportunityRepo)

	authHnd := handler.NewAuthHandler(authSrv)
	applicantHnd := handler.NewApplicantHandler(applicantSrv)
	employerHnd := handler.NewEmployerHandler(employerSrv)
	curatorHnd := handler.NewCuratorHandler(curatorSrv)
	opportunityHnd := handler.NewOpportunityHandler(opportunitySrv)
	applicationHnd := handler.NewApplicationHandler(applicationSrv)
	contactHnd := handler.NewContactHandler(contactSrv)
	recomendationHnd := handler.NewRecommendationHandler(recomendationSrv)
	favoriteshnd := handler.NewFavoritesHandler(favoritesSrv)

	router := gin.New()
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	v1.Use(middleware.ParseJWT(tokenManager))

	protectedV1 := v1.Group("/")
	protectedV1.Use(middleware.RegisteredOnly())

	routes.SetupAuthRoutes(v1, authHnd)
	routes.SetupApplicantRoutes(v1, protectedV1, applicantHnd)
	routes.SetupEmployerRoutes(v1, protectedV1, employerHnd)
	routes.SetupCuratorRoutes(protectedV1, curatorHnd)
	routes.SetupOpportunityRoutes(v1, protectedV1, opportunityHnd)
	routes.SetupApplicationRoutes(protectedV1, applicationHnd)
	routes.SetupContactRoutes(protectedV1, contactHnd)
	routes.SetupRecomendationRoutes(protectedV1, recomendationHnd)
	routes.SetupFavoritesRoutes(protectedV1, favoriteshnd)

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
