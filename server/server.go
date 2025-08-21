package main

import (
	"context"
	"develapar-server/config"
	"develapar-server/controller"
	_ "develapar-server/docs" // Import the generated docs
	"develapar-server/middleware"
	"develapar-server/repository"
	"develapar-server/service"
	"develapar-server/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// loggerAdapter adapts utils.Logger to middleware.Logger interface
type loggerAdapter struct {
	logger utils.Logger
}

func (la *loggerAdapter) Error(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	// Convert map[string]interface{} to []utils.Field
	var utilsFields []utils.Field
	for key, value := range fields {
		utilsFields = append(utilsFields, utils.Field{Key: key, Value: value})
	}
	la.logger.Error(ctx, msg, err, utilsFields...)
}

func (la *loggerAdapter) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	// Convert map[string]interface{} to []utils.Field
	var utilsFields []utils.Field
	for key, value := range fields {
		utilsFields = append(utilsFields, utils.Field{Key: key, Value: value})
	}
	la.logger.Warn(ctx, msg, utilsFields...)
}

func (la *loggerAdapter) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	// Convert map[string]interface{} to []utils.Field
	var utilsFields []utils.Field
	for key, value := range fields {
		utilsFields = append(utilsFields, utils.Field{Key: key, Value: value})
	}
	la.logger.Info(ctx, msg, utilsFields...)
}

type Server struct {
	uS          service.UserService
	cS          service.CategoryService
	aS          service.ArticleService
	bS          service.BookmarkService
	tS          service.TagService
	atS         service.ArticleTagService
	coS         service.CommentService
	lS          service.LikeService
	pS          service.ProductService
	jS          service.JwtService
	mD          middleware.AuthMiddleware
	eMD         middleware.ErrorHandler
	hC          *controller.HealthController
	mC          *controller.MetricsController
	poolManager config.ConnectionPoolManager
	engine      *gin.Engine
	portApp     string
}

func (s *Server) initiateRoute() {
	routerGroup := s.engine.Group("/api/v1")
	controller.NewUserController(s.uS, s.mD, routerGroup, s.eMD).Route()
	controller.NewCategoryController(s.cS, routerGroup, s.mD, s.eMD).Route()
	controller.NewArticleController(s.aS, s.mD, routerGroup, s.eMD).Route()
	controller.NewBookmarkController(s.bS, routerGroup, s.mD, s.eMD).Route()
	controller.NewTagController(s.tS, routerGroup, s.mD, s.eMD).Route()
	controller.NewArticleTagController(s.atS, routerGroup, s.mD, s.eMD).Route()
	controller.NewCommentController(s.coS, routerGroup, s.mD, s.eMD).Route()
	controller.NewLikeController(s.lS, routerGroup, s.mD, s.eMD).Route()
	controller.NewProductController(s.pS, routerGroup, s.mD, s.eMD).Route()

	// Health check routes (no authentication required)
	s.hC.Route(routerGroup)

	// Metrics routes (no authentication required for monitoring)
	s.mC.Route(routerGroup)

	// Swagger UI
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *Server) Start() {
	s.initiateRoute()
	s.engine.Run(s.portApp)
}

// parseLogLevel converts string log level to utils.LogLevel
func parseLogLevel(level string) utils.LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return utils.DebugLevel
	case "INFO":
		return utils.InfoLevel
	case "WARN", "WARNING":
		return utils.WarnLevel
	case "ERROR":
		return utils.ErrorLevel
	case "FATAL":
		return utils.FatalLevel
	default:
		return utils.InfoLevel // Default fallback
	}
}

// CORSMiddleware adalah middleware yang akan menangani CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("CORS Middleware hit:", c.Request.Method, c.Request.URL.Path)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func NewServer() *Server {
	co, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create context for database initialization
	ctx := context.Background()

	// Initialize connection pool manager with context support
	poolManager, err := config.NewConnectionPoolManager(ctx, co.DbConfig, co.PoolConfig)
	if err != nil {
		log.Fatalf("failed to create connection pool manager: %v", err)
	}

	// Get database connection from pool manager
	db, err := poolManager.GetConnection(ctx)
	if err != nil {
		log.Fatalf("failed to get database connection: %v", err)
	}

	// Perform initial health check
	if err := poolManager.HealthCheck(ctx); err != nil {
		log.Fatalf("database health check failed: %v", err)
	}

	log.Printf("Database connection pool initialized successfully")
	stats := poolManager.GetStats(ctx)
	log.Printf("Connection pool stats: Open=%d, InUse=%d, Idle=%d",
		stats.OpenConnections, stats.InUseConnections, stats.IdleConnections)

	engine := gin.Default()

	// Initialize middleware components
	contextManager := middleware.NewContextManager()
	contextMiddleware := middleware.NewContextMiddleware(contextManager)
	errorHandler := middleware.NewErrorHandler(nil) // Using default logger

	// Initialize logger factory with config-based log level
	logLevel := parseLogLevel(co.LoggingConfig.Level)
	loggerFactory := utils.NewLoggerFactory(logLevel)

	// Initialize metrics service and related components
	metricsLogger := loggerFactory.GetLogger("metrics")
	metricsService := service.NewMetricsService(metricsLogger)
	metricsController := controller.NewMetricsController(metricsService)
	metricsMiddleware := middleware.NewMetricsMiddleware(metricsService, metricsLogger)

	// Initialize rate limiting components
	rateLimitUtilsLogger := loggerFactory.GetLogger("rate_limiter")

	// Create logger adapter to convert utils.Logger to middleware.Logger
	rateLimitLogger := &loggerAdapter{logger: rateLimitUtilsLogger}

	rateLimitStore := middleware.NewInMemoryStore(rateLimitLogger)
	rateLimiter := middleware.NewSlidingWindowRateLimiter(rateLimitStore, rateLimitLogger)

	// Create rate limit configuration from config
	rateLimitConfig := &middleware.RateLimitConfig{
		DefaultLimit:        co.RateLimitConfig.RequestsPerMinute,
		DefaultWindow:       co.RateLimitConfig.WindowSize,
		AuthenticatedLimit:  co.RateLimitConfig.AuthenticatedRPM,
		AuthenticatedWindow: co.RateLimitConfig.WindowSize,
		AnonymousLimit:      co.RateLimitConfig.AnonymousRPM,
		AnonymousWindow:     co.RateLimitConfig.WindowSize,
		SkipPaths:           []string{"/health", "/metrics", "/swagger"},
		IncludeHeaders:      true,
		KeyStrategy:         "ip",
	}

	// Create monitored rate limiting middleware
	rateLimitMiddleware := middleware.NewMonitoredRateLimitMiddleware(rateLimiter, rateLimitConfig, rateLimitLogger)

	// Start periodic cleanup for rate limiting
	cleanupMiddleware := rateLimitMiddleware.CleanupMiddleware(co.RateLimitConfig.CleanupInterval)

	// Start periodic monitoring for rate limiting
	rateLimitMonitor := rateLimitMiddleware.GetMonitor()
	rateLimitMonitor.StartPeriodicLogging(ctx, 10*time.Minute) // Log stats every 10 minutes

	// Initialize request logging middleware with context support
	requestLogger := loggerFactory.GetLogger("request")
	requestLoggingMiddleware := middleware.NewRequestLoggerWithMetrics(requestLogger)
	log.Printf("Request logging middleware initialized with context support and metrics collection")

	// Configure middleware order for proper context propagation
	// 1. Recovery middleware (should be first to catch panics)
	engine.Use(middleware.RecoveryMiddleware(errorHandler))

	// 2. CORS middleware
	engine.Use(CORSMiddleware())

	// 3. Context middleware (injects request ID, user ID, start time)
	engine.Use(contextMiddleware.InjectContext())

	// 4. Request logging middleware (after context injection for proper context correlation)
	engine.Use(requestLoggingMiddleware.LogRequestsWithMetrics())

	// 5. Rate limiting middleware (after context and logging, before metrics)
	if co.RateLimitConfig.Enabled {
		engine.Use(rateLimitMiddleware.MonitoredMiddleware())
		engine.Use(cleanupMiddleware)
		log.Printf("Rate limiting enabled: %d RPM for anonymous, %d RPM for authenticated users",
			co.RateLimitConfig.AnonymousRPM, co.RateLimitConfig.AuthenticatedRPM)
	} else {
		log.Printf("Rate limiting disabled")
	}

	// 6. Metrics middleware (collects request metrics with context)
	engine.Use(metricsMiddleware.CollectMetrics())

	// 7. System metrics collection (background collection)
	engine.Use(metricsMiddleware.CollectSystemMetrics())

	// 8. Error handling middleware (should be last to catch all errors)
	engine.Use(middleware.ErrorMiddleware(errorHandler))

	portApp := co.AppPort

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)
	tagRepo := repository.NewTagRepository(db)
	articleTagRepo := repository.NewArticleTagRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	productRepo := repository.NewProductRepository(db)

	passwordHasher := utils.NewPasswordHasher()
	jwtService := service.NewJwtService(co.SecurityConfig)

	// Initialize error wrapper and validation service for pagination
	errorWrapper := utils.NewErrorWrapper()
	validationService := service.NewValidationService(errorWrapper)
	paginationService := service.NewPaginationService(validationService, errorWrapper)

	userService := service.NewUserservice(userRepo, jwtService, passwordHasher, paginationService, validationService)
	categoryService := service.NewCategoryService(categoryRepo, validationService)
	articleTagService := service.NewArticleTagService(tagRepo, articleTagRepo, validationService)
	articleService := service.NewArticleService(articleRepo, articleTagService, paginationService, validationService)
	bookmarkService := service.NewBookmarkService(bookmarkRepo, validationService)
	tagService := service.NewTagService(tagRepo, validationService)
	commentService := service.NewCommentService(commentRepo, validationService)
	likeService := service.NewLikeService(likeRepo, validationService)
	productService := service.NewProductService(productRepo, validationService, paginationService)

	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	healthController := controller.NewHealthController(poolManager)

	return &Server{
		cS:          categoryService,
		uS:          userService,
		aS:          articleService,
		bS:          bookmarkService,
		tS:          tagService,
		jS:          jwtService,
		atS:         articleTagService,
		coS:         commentService,
		lS:          likeService,
		pS:          productService,
		mD:          authMiddleware,
		eMD:         errorHandler,
		hC:          healthController,
		mC:          metricsController,
		poolManager: poolManager,
		portApp:     portApp,
		engine:      engine,
	}
}
