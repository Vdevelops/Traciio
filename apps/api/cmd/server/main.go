package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	areamappinghandler "github.com/gilabs/crm-healthcare/api/internal/api/area_mapping"
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/api/routes"
	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/database"
	domainevents "github.com/gilabs/crm-healthcare/api/internal/domain/events"
	"github.com/gilabs/crm-healthcare/api/internal/hub"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	accountrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/account"
	activityrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/activity"
	activitytyperepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/activity_type"
	aisettingsrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/ai_settings"
	areamappingrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/area_mapping"
	"github.com/gilabs/crm-healthcare/api/internal/repository/postgres/auth"
	brickrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/brick"
	bricktargetdistributionrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/brick_target_distribution"
	categoryrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/category"
	contactrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/contact"
	contactrolerepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/contact_role"
	customerpurchaserepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/customer_purchase_history"
	dealrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/deal"
	dealhistoryrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/deal_history"
	dealproductitemrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/deal_product_item"
	googlecalendartokenrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/google_calendar_token"
	grouprepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/group"
	industryrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/industry"
	leadrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/lead"
	leadqualificationrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/lead_qualification"
	leadsourcerepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/lead_source"
	leadstatusrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/lead_status"
	monthlytargetrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/monthly_target"
	notificationrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/notification"
	permissionrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/permission"
	pipelinerepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/pipeline"
	productrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/product"
	productanalyticsrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/product_analytics"
	productcategoryrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/product_category"
	refreshtokenrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/refresh_token"
	reminderrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/reminder"
	rolerepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/role"
	routeoptimizationrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/route_optimization"
	salesoverviewrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/sales_overview"
	schedulerepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/schedule"
	taskrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/task"
	userrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/user"
	visitreportrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/visit_report"
	accountservice "github.com/gilabs/crm-healthcare/api/internal/service/account"
	activityservice "github.com/gilabs/crm-healthcare/api/internal/service/activity"
	activitytypeservice "github.com/gilabs/crm-healthcare/api/internal/service/activity_type"
	aiservice "github.com/gilabs/crm-healthcare/api/internal/service/ai"
	aisettingsservice "github.com/gilabs/crm-healthcare/api/internal/service/ai_settings"
	areamappingservice "github.com/gilabs/crm-healthcare/api/internal/service/area_mapping"
	authservice "github.com/gilabs/crm-healthcare/api/internal/service/auth"
	brickservice "github.com/gilabs/crm-healthcare/api/internal/service/brick"
	brickanalyticsservice "github.com/gilabs/crm-healthcare/api/internal/service/brick_analytics"
	bricktargetdistributionservice "github.com/gilabs/crm-healthcare/api/internal/service/brick_target_distribution"
	categoryservice "github.com/gilabs/crm-healthcare/api/internal/service/category"
	contactservice "github.com/gilabs/crm-healthcare/api/internal/service/contact"
	contactroleservice "github.com/gilabs/crm-healthcare/api/internal/service/contact_role"
	customerpurchaseservice "github.com/gilabs/crm-healthcare/api/internal/service/customer_purchase"
	dashboardservice "github.com/gilabs/crm-healthcare/api/internal/service/dashboard"
	fileservice "github.com/gilabs/crm-healthcare/api/internal/service/file"
	googlecalendartokenservice "github.com/gilabs/crm-healthcare/api/internal/service/google_calendar_token"
	groupservice "github.com/gilabs/crm-healthcare/api/internal/service/group"
	industryservice "github.com/gilabs/crm-healthcare/api/internal/service/industry"
	leadservice "github.com/gilabs/crm-healthcare/api/internal/service/lead"
	leadqualificationservice "github.com/gilabs/crm-healthcare/api/internal/service/lead_qualification"
	leadsourceservice "github.com/gilabs/crm-healthcare/api/internal/service/lead_source"
	leadstatusservice "github.com/gilabs/crm-healthcare/api/internal/service/lead_status"
	monthlytargetservice "github.com/gilabs/crm-healthcare/api/internal/service/monthly_target"
	notificationservice "github.com/gilabs/crm-healthcare/api/internal/service/notification"
	permissionservice "github.com/gilabs/crm-healthcare/api/internal/service/permission"
	pipelineservice "github.com/gilabs/crm-healthcare/api/internal/service/pipeline"
	productservice "github.com/gilabs/crm-healthcare/api/internal/service/product"
	productanalyticsservice "github.com/gilabs/crm-healthcare/api/internal/service/product_analytics"
	reportservice "github.com/gilabs/crm-healthcare/api/internal/service/report"
	roleservice "github.com/gilabs/crm-healthcare/api/internal/service/role"
	routeoptimizationservice "github.com/gilabs/crm-healthcare/api/internal/service/route_optimization"
	salesoverviewservice "github.com/gilabs/crm-healthcare/api/internal/service/sales_overview"
	scheduleservice "github.com/gilabs/crm-healthcare/api/internal/service/schedule"
	taskservice "github.com/gilabs/crm-healthcare/api/internal/service/task"
	userservice "github.com/gilabs/crm-healthcare/api/internal/service/user"
	visitreportservice "github.com/gilabs/crm-healthcare/api/internal/service/visit_report"
	"github.com/gilabs/crm-healthcare/api/internal/worker"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/cerebras"
	"github.com/gilabs/crm-healthcare/api/pkg/events"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gilabs/crm-healthcare/api/pkg/logger"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gilabs/crm-healthcare/api/seeders"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Seed data (ONLY in development mode when explicitly enabled)
	// PRODUCTION SAFETY: General seeders are disabled by default
	// Set ENABLE_SEEDERS=true to run seeders (development only)
	if shouldRunSeeders() {
		log.Println("🌱 Running database seeders...")
		if os.Getenv("ONLY_ACCOUNT") == "true" {
			if err := seeders.SeedAccountOnly(); err != nil {
				log.Fatal("Failed to seed account data:", err)
			}
		} else {
			if err := seeders.SeedAll(database.DB); err != nil {
				log.Fatal("Failed to seed data:", err)
			}
		}
		log.Println("✅ Database seeding completed")
	} else {
		log.Println("📦 General seeders skipped (ENABLE_SEEDERS is not set or ENV=production)")
	}

	// SPECIAL CASE: Default Admin Seeder for Production
	// This allows creating the initial admin account safely in production
	if os.Getenv("SEED_DEFAULT_ADMIN") == "true" {
		log.Println("🔑 Running Default Admin Seeder...")
		if err := seeders.SeedDefaultAdmin(); err != nil {
			log.Printf("❌ Failed to seed default admin: %v", err)
		}
	}

	// Initialize Redis cache with high availability support
	// Initialize Redis cache with high availability support
	var redisCache *cache.Cache
	var redisClient redis.UniversalClient

	if config.AppConfig.Redis.Enabled {
		var err error
		redisCache, err = cache.New(&config.AppConfig.Redis)
		if err != nil {
			log.Printf("Warning: Redis initialization failed: %v. Caching disabled.", err)
		} else if redisCache.IsEnabled() {
			defer redisCache.Close()
			redisClient = redisCache.GetClient()
			log.Println("Redis cache initialized successfully")
		}
	} else {
		log.Println("Redis cache is disabled")
	}

	// Initialize rate limiter with Redis support (falls back to in-memory if Redis unavailable)
	if redisClient != nil {
		middleware.InitRateLimiter(redisClient)
	} else {
		middleware.InitRateLimiter(nil)
		log.Println("Warning: Rate limiting using in-memory storage (single-instance only)")
	}

	// Initialize event producer for Kafka-ready architecture
	// When Kafka is disabled (default), events are logged via NoOpProducer
	// When Kafka is enabled, events will be published to Kafka topics
	var eventProducer events.EventProducer
	if config.AppConfig.Kafka.Enabled {
		// TODO: Implement KafkaProducer when Kafka is installed
		// brokers := strings.Split(config.AppConfig.Kafka.Brokers, ",")
		// kafkaProducer, err := events.NewKafkaProducer(brokers, config.AppConfig.Kafka.TopicPrefix)
		// if err != nil {
		// 	log.Fatalf("Failed to initialize Kafka producer: %v", err)
		// }
		// defer kafkaProducer.Close()
		// eventProducer = kafkaProducer
		// log.Println("Kafka event producer initialized")
		log.Println("WARNING: Kafka is enabled but KafkaProducer is not implemented yet. Using NoOpProducer.")
		eventProducer = events.NewNoOpProducer()
	} else {
		eventProducer = events.NewNoOpProducer()
		log.Println("Event producer initialized: NoOpProducer (Kafka disabled)")
	}

	// Create event helper for easy event emission
	eventHelper := domainevents.NewHelper(eventProducer)

	// Setup JWT Manager
	jwtManager := jwt.NewJWTManager(
		config.AppConfig.JWT.SecretKey,
		time.Duration(config.AppConfig.JWT.AccessTokenTTL)*time.Hour,
		time.Duration(config.AppConfig.JWT.RefreshTokenTTL)*24*time.Hour,
	)

	// Setup repositories
	authRepo := auth.NewRepository(database.DB)
	refreshTokenRepo := refreshtokenrepo.NewRepository(database.DB)
	userRepo := userrepo.NewRepository(database.DB)
	roleRepo := rolerepo.NewRepository(database.DB)
	basePermissionRepo := permissionrepo.NewRepository(database.DB)
	menuRepo := permissionrepo.NewMenuRepository(database.DB)

	// Wrap permission repository with cache if Redis is available
	var permissionRepo interfaces.PermissionRepository
	permissionRepo = basePermissionRepo

	groupRepo := grouprepo.NewRepository(database.DB)
	monthlyTargetRepo := monthlytargetrepo.NewRepository(database.DB)
	categoryRepo := categoryrepo.NewRepository(database.DB)
	contactRoleRepo := contactrolerepo.NewRepository(database.DB)
	accountRepo := accountrepo.NewRepository(database.DB)
	contactRepo := contactrepo.NewRepository(database.DB)
	pipelineRepo := pipelinerepo.NewRepository(database.DB)
	dealRepo := dealrepo.NewRepository(database.DB)
	dealProductItemRepo := dealproductitemrepo.NewRepository(database.DB)
	dealHistoryRepo := dealhistoryrepo.NewRepository(database.DB)
	leadRepo := leadrepo.NewRepository(database.DB)
	leadStatusRepo := leadstatusrepo.NewRepository(database.DB)
	industryRepo := industryrepo.NewRepository(database.DB)
	leadSourceRepo := leadsourcerepo.NewRepository(database.DB)
	visitReportRepo := visitreportrepo.NewRepository(database.DB)
	activityRepo := activityrepo.NewRepository(database.DB)
	activityTypeRepo := activitytyperepo.NewRepository(database.DB)
	productCategoryRepo := productcategoryrepo.NewRepository(database.DB)
	productRepo := productrepo.NewRepository(database.DB)
	productAnalyticsRepo := productanalyticsrepo.NewRepository(database.DB)
	taskRepo := taskrepo.NewRepository(database.DB)
	reminderRepo := reminderrepo.NewRepository(database.DB)
	notificationRepo := notificationrepo.NewRepository(database.DB)
	aiSettingsRepo := aisettingsrepo.NewRepository(database.DB)
	routeOptimizationRepo := routeoptimizationrepo.NewRepository(database.DB)
	scheduleRepo := schedulerepo.NewRepository(database.DB)
	googleCalendarTokenRepo := googlecalendartokenrepo.NewRepository(database.DB)
	salesOverviewRepo := salesoverviewrepo.NewRepository(database.DB)
	areaMappingRepo := areamappingrepo.NewRepository(database.DB)
	brickRepo := brickrepo.NewRepository(database.DB)
	brickTargetDistributionRepo := bricktargetdistributionrepo.NewRepository(database.DB)
	leadQualificationRepo := leadqualificationrepo.NewLeadQualificationRepository(database.DB)
	customerPurchaseRepo := customerpurchaserepo.NewCustomerPurchaseHistoryRepository(database.DB)

	// Setup services
	permissionService := permissionservice.NewService(permissionRepo, roleRepo, userRepo, redisClient)
	authService := authservice.NewService(authRepo, refreshTokenRepo, jwtManager, permissionService)
	userService := userservice.NewService(userRepo, roleRepo, groupRepo, brickRepo, monthlyTargetRepo, redisCache)
	profileService := userservice.NewProfileService(userRepo, activityRepo, dealRepo, visitReportRepo, taskRepo)
	settingsService := userservice.NewSettingsService(userRepo, dealRepo, visitReportRepo, taskRepo, profileService)
	roleService := roleservice.NewService(roleRepo, userRepo)
	groupService := groupservice.NewService(groupRepo, userRepo)
	monthlyTargetService := monthlytargetservice.NewService(monthlyTargetRepo, groupRepo, userRepo, brickRepo)
	brickService := brickservice.NewService(brickRepo, userRepo, monthlyTargetRepo, database.DB)
	brickTargetDistributionService := bricktargetdistributionservice.NewService(brickTargetDistributionRepo, brickRepo, monthlyTargetRepo, userRepo)

	// Initialize brick helper for auto-populating brick_id
	brickHelper := brickservice.NewBrickHelper(userRepo, brickRepo, accountRepo)
	categoryService := categoryservice.NewService(categoryRepo)
	contactRoleService := contactroleservice.NewService(contactRoleRepo)
	accountService := accountservice.NewService(accountRepo, categoryRepo, brickHelper)
	contactService := contactservice.NewService(contactRepo, accountRepo, contactRoleRepo)
	pipelineService := pipelineservice.NewService(database.DB, pipelineRepo, dealRepo, accountRepo, productRepo, dealProductItemRepo, dealHistoryRepo, taskRepo, visitReportRepo, leadRepo, customerPurchaseRepo, brickHelper, eventHelper)
	leadService := leadservice.NewService(database.DB, leadRepo, dealRepo, pipelineRepo, accountRepo, contactRepo, categoryRepo, contactRoleRepo, userRepo, activityRepo, visitReportRepo, taskRepo, dealHistoryRepo, leadStatusRepo, eventHelper)
	leadQualificationService := leadqualificationservice.NewService(leadQualificationRepo, leadRepo)
	customerPurchaseService := customerpurchaseservice.NewService(customerPurchaseRepo, accountRepo)
	leadStatusService := leadstatusservice.NewService(leadStatusRepo, database.DB)
	industryService := industryservice.NewService(industryRepo, database.DB)
	leadSourceService := leadsourceservice.NewService(leadSourceRepo, database.DB)
	activityService := activityservice.NewService(activityRepo, activityTypeRepo, accountRepo, contactRepo, userRepo, database.DB, eventHelper)
	activityTypeService := activitytypeservice.NewService(activityTypeRepo)
	visitReportService := visitreportservice.NewService(visitReportRepo, accountRepo, contactRepo, userRepo, activityRepo, activityTypeRepo, leadRepo, taskRepo, notificationRepo, brickHelper, database.DB)
	dashboardService := dashboardservice.NewService(visitReportRepo, accountRepo, activityRepo, userRepo, dealRepo, taskRepo, pipelineRepo, leadRepo, roleRepo, monthlyTargetRepo, brickRepo, scheduleRepo)
	salesOverviewService := salesoverviewservice.NewService(salesOverviewRepo, monthlyTargetRepo)
	areaMappingService := areamappingservice.NewService(areaMappingRepo)

	// Setup file service with storage provider
	var storageProvider fileservice.StorageProvider
	storageConfig := config.AppConfig.Storage

	if storageConfig.Type == "r2" {
		// Initialize R2 storage
		r2Storage, err := fileservice.NewR2Storage(
			storageConfig.R2Endpoint,
			storageConfig.R2AccessKeyID,
			storageConfig.R2SecretAccessKey,
			storageConfig.R2Bucket,
			storageConfig.R2PublicURL,
			storageConfig.BaseURL,
		)
		if err != nil {
			log.Fatalf("Failed to initialize R2 storage: %v", err)
		}
		storageProvider = r2Storage
	} else {
		// Initialize local storage (default)
		storageProvider = fileservice.NewLocalStorage(
			storageConfig.UploadDir,
			storageConfig.BaseURL,
		)
	}

	fileService := fileservice.NewService(storageProvider)
	reportService := reportservice.NewService(visitReportRepo, accountRepo, activityRepo, userRepo, dealRepo)
	productService := productservice.NewService(productRepo, productCategoryRepo)
	productAnalyticsService := productanalyticsservice.NewService(productAnalyticsRepo)
	taskService := taskservice.NewService(taskRepo, reminderRepo, userRepo, accountRepo, contactRepo, dealRepo, leadRepo, eventHelper)
	routeOptimizationService := routeoptimizationservice.NewService(routeOptimizationRepo, userRepo)

	// Setup Google Calendar token service
	googleCalendarTokenService := googlecalendartokenservice.NewService(
		googleCalendarTokenRepo,
		&config.AppConfig.GoogleCalendar,
		config.AppConfig.Encryption.Key,
	)

	scheduleService := scheduleservice.NewService(
		scheduleRepo,
		taskRepo,
		userRepo,
		googleCalendarTokenService,
	)

	// Connect schedule service to task service for auto-creating schedules
	taskService.SetScheduleService(scheduleService)

	// Connect Google Calendar token service to task service for auto-syncing tasks
	taskService.SetGoogleCalendarTokenService(googleCalendarTokenService)

	// Setup WebSocket hub
	notificationHub := hub.NewNotificationHub()
	go notificationHub.Run()

	// Setup notification service with hub
	notificationService := notificationservice.NewService(notificationRepo, eventHelper)
	notificationService.SetHub(notificationHub)

	// Setup role service with hub for permission updates broadcast
	roleService.SetNotificationHub(notificationHub)

	// Set permission service on role service for cache invalidation
	roleService.SetPermissionService(permissionService)

	// Setup Cerebras AI Client
	cerebrasClient := cerebras.NewClient(
		config.AppConfig.Cerebras.BaseURL,
		config.AppConfig.Cerebras.APIKey,
		config.AppConfig.Cerebras.Model,
	)

	// Setup AI Settings Service
	aiSettingsService := aisettingsservice.NewService(aiSettingsRepo)

	// Setup AI Service
	aiService := aiservice.NewService(
		cerebrasClient,
		visitReportRepo,
		accountRepo,
		contactRepo,
		dealRepo,
		leadRepo,
		activityRepo,
		taskRepo,
		productRepo,
		pipelineRepo,
		aiSettingsRepo,
		permissionService,
		dashboardService,         // For analytics data
		routeOptimizationService, // For creating real routes from AI
		leadService,
		taskService,
		pipelineService,
		scheduleService,
		config.AppConfig.Cerebras.APIKey,
	)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, profileService, settingsService)
	roleHandler := handlers.NewRoleHandler(roleService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)
	menuHandler := handlers.NewMenuHandler(menuRepo)
	groupHandler := handlers.NewGroupHandler(groupService)
	monthlyTargetHandler := handlers.NewMonthlyTargetHandler(monthlyTargetService)
	brickHandler := handlers.NewBrickHandler(brickService)
	brickTargetDistributionHandler := handlers.NewBrickTargetDistributionHandler(brickTargetDistributionService)

	// Initialize brick analytics service and handler
	brickAnalyticsService := brickanalyticsservice.NewService(
		database.DB,
		brickRepo,
		dealRepo,
		visitReportRepo,
		accountRepo,
		monthlyTargetRepo,
		brickTargetDistributionRepo,
		userRepo,
	)
	brickAnalyticsHandler := handlers.NewBrickAnalyticsHandler(brickAnalyticsService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	contactRoleHandler := handlers.NewContactRoleHandler(contactRoleService)
	accountHandler := handlers.NewAccountHandler(accountService)
	contactHandler := handlers.NewContactHandler(contactService)
	pipelineHandler := handlers.NewPipelineHandler(pipelineService)
	dealHandler := handlers.NewDealHandler(pipelineService, visitReportService, activityService, accountService, contactService, productService)
	leadHandler := handlers.NewLeadHandler(leadService, visitReportService, activityService, leadQualificationService)
	customerPurchaseHandler := handlers.NewCustomerPurchaseHandler(customerPurchaseService)
	leadStatusHandler := handlers.NewLeadStatusHandler(leadStatusService)
	industryHandler := handlers.NewIndustryHandler(industryService)
	leadSourceHandler := handlers.NewLeadSourceHandler(leadSourceService)
	activityHandler := handlers.NewActivityHandler(activityService)
	activityTypeHandler := handlers.NewActivityTypeHandler(activityTypeService)
	visitReportHandler := handlers.NewVisitReportHandler(visitReportService, fileService, accountService, contactService, pipelineService, leadService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	reportHandler := handlers.NewReportHandler(reportService)
	productHandler := handlers.NewProductHandler(productService)
	productAnalyticsHandler := handlers.NewProductAnalyticsHandler(productAnalyticsService)
	fileHandler := handlers.NewFileHandler(fileService)
	taskHandler := handlers.NewTaskHandler(taskService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	aiHandler := handlers.NewAIHandler(aiService)
	aiSettingsHandler := handlers.NewAISettingsHandler(aiSettingsService)
	routeOptimizationHandler := handlers.NewRouteOptimizationHandler(routeOptimizationService)
	geocodingHandler := handlers.NewGeocodingHandler()
	scheduleHandler := handlers.NewScheduleHandler(scheduleService)
	googleCalendarAuthHandler := handlers.NewGoogleCalendarAuthHandler(googleCalendarTokenService, &config.AppConfig.GoogleCalendar)
	salesOverviewHandler := handlers.NewSalesOverviewHandler(salesOverviewService, userRepo, brickRepo)
	areaMappingHandler := areamappinghandler.NewHandler(areaMappingService)

	// Setup health handler for monitoring
	healthHandler := handlers.NewHealthHandler()

	// Setup WebSocket handler
	wsHandler := handlers.NewWebSocketHandler(notificationHub, jwtManager)

	// Setup reminder worker
	reminderWorker := worker.NewReminderWorker(
		reminderRepo,
		notificationService,
		notificationHub,
		1*time.Minute, // Run every 1 minute
	)
	reminderWorker.Start()

	// Setup refresh token cleanup worker
	// Run every 24 hours to clean up expired refresh tokens
	refreshTokenCleanupWorker := worker.NewRefreshTokenCleanupWorker(
		refreshTokenRepo,
		24*time.Hour, // Run every 24 hours
	)
	refreshTokenCleanupWorker.Start()

	// Create RBAC scope middleware instance for data-scoped routes
	scopeMiddleware := middleware.ScopeMiddleware(permissionService, roleRepo, userRepo, brickRepo, redisClient)

	// Setup router
	router := setupRouter(
		jwtManager,
		authHandler,
		userHandler,
		roleHandler,
		permissionHandler,
		menuHandler,
		groupHandler,
		monthlyTargetHandler,
		brickHandler,
		brickTargetDistributionHandler,
		brickAnalyticsHandler,
		categoryHandler,
		contactRoleHandler,
		accountHandler,
		contactHandler,
		pipelineHandler,
		dealHandler,
		leadHandler,
		leadStatusHandler,
		industryHandler,
		leadSourceHandler,
		activityHandler,
		activityTypeHandler,
		visitReportHandler,
		dashboardHandler,
		reportHandler,
		productHandler,
		productAnalyticsHandler,
		fileHandler,
		taskHandler,
		notificationHandler,
		wsHandler,
		aiHandler,
		aiSettingsHandler,
		routeOptimizationHandler,
		geocodingHandler,
		scheduleHandler,
		googleCalendarAuthHandler,
		customerPurchaseHandler,
		areaMappingHandler,
		salesOverviewHandler,
		healthHandler,
		permissionService,
		scopeMiddleware,
	)

	// Configure HTTP server with production-ready timeouts
	port := config.AppConfig.Server.Port
	addr := ":" + port

	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   35 * time.Second, // Must be > request timeout to allow proper completion
		MaxHeaderBytes: 1 << 20,          // 1MB max header size
	}

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling
	go func() {
		log.Printf("🚀 Server starting in %s mode on %s", config.AppConfig.Server.Env, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ listen error: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so no need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("⏳ Shutting down server...")

	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️  Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server shutdown completed")
}

func setupRouter(
	jwtManager *jwt.JWTManager,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	roleHandler *handlers.RoleHandler,
	permissionHandler *handlers.PermissionHandler,
	menuHandler *handlers.MenuHandler,
	groupHandler *handlers.GroupHandler,
	monthlyTargetHandler *handlers.MonthlyTargetHandler,
	brickHandler *handlers.BrickHandler,
	brickTargetDistributionHandler *handlers.BrickTargetDistributionHandler,
	brickAnalyticsHandler *handlers.BrickAnalyticsHandler,
	categoryHandler *handlers.CategoryHandler,
	contactRoleHandler *handlers.ContactRoleHandler,
	accountHandler *handlers.AccountHandler,
	contactHandler *handlers.ContactHandler,
	pipelineHandler *handlers.PipelineHandler,
	dealHandler *handlers.DealHandler,
	leadHandler *handlers.LeadHandler,
	leadStatusHandler *handlers.LeadStatusHandler,
	industryHandler *handlers.IndustryHandler,
	leadSourceHandler *handlers.LeadSourceHandler,
	activityHandler *handlers.ActivityHandler,
	activityTypeHandler *handlers.ActivityTypeHandler,
	visitReportHandler *handlers.VisitReportHandler,
	dashboardHandler *handlers.DashboardHandler,
	reportHandler *handlers.ReportHandler,
	productHandler *handlers.ProductHandler,
	productAnalyticsHandler *handlers.ProductAnalyticsHandler,
	fileHandler *handlers.FileHandler,
	taskHandler *handlers.TaskHandler,
	notificationHandler *handlers.NotificationHandler,
	wsHandler *handlers.WebSocketHandler,
	aiHandler *handlers.AIHandler,
	aiSettingsHandler *handlers.AISettingsHandler,
	routeOptimizationHandler *handlers.RouteOptimizationHandler,
	geocodingHandler *handlers.GeocodingHandler,
	scheduleHandler *handlers.ScheduleHandler,
	googleCalendarAuthHandler *handlers.GoogleCalendarAuthHandler,
	customerPurchaseHandler *handlers.CustomerPurchaseHandler,
	areaMappingHandler *areamappinghandler.Handler,
	salesOverviewHandler *handlers.SalesOverviewHandler,
	healthHandler *handlers.HealthHandler,
	permissionService *permissionservice.Service,
	scopeMiddleware gin.HandlerFunc,
) *gin.Engine {
	// Set Gin mode
	if config.AppConfig.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Configure max body size (50MB for file uploads)
	router.MaxMultipartMemory = 50 << 20 // 50 MB

	// Global middleware
	router.Use(middleware.ConcurrencyLimiterMiddleware(500)) // Prevent goroutine/connection exhaustion under extreme load
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.HSTSMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())             // Enhanced security headers
	router.Use(middleware.InputSanitizationMiddleware())           // XSS & SQL injection prevention
	router.Use(middleware.CSRFMiddleware())                        // CSRF protection for authenticated endpoints
	router.Use(middleware.MaxBodySizeMiddleware(10 * 1024 * 1024)) // 10MB max body size (DoS protection)

	// A+ SECURITY & PERFORMANCE MIDDLEWARE
	// -------------------------------------------------------------------------
	// 1. Audit Logging: Tracks all sensitive actions
	router.Use(middleware.AuditLogMiddleware())

	// 2. Request Timeout: Prevents slow requests from locking the system.
	// Default is 5 seconds, but route optimization needs longer due to OSRM calls.
	// Note: Repositories must use db.WithContext(ctx) for this to kill DB queries.
	router.Use(middleware.TimeoutMiddlewareByPath(5*time.Second, map[string]time.Duration{
		"/api/v1/route-optimization/":        30 * time.Second,
		"/api/v1/mobile/route-optimization/": 30 * time.Second,
	}))

	// 3. Prometheus Metrics: Observability for latency and errors
	router.Use(middleware.MetricsMiddleware())
	// -------------------------------------------------------------------------

	// Health check endpoints
	router.GET("/health", healthHandler.GetSystemHealth)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Prometheus Metrics Endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Detailed health endpoints
	router.GET("/health/cache", healthHandler.GetCacheHealth)
	router.GET("/health/cache/metrics", healthHandler.GetCacheMetrics)
	router.GET("/health/circuit-breakers", healthHandler.GetCircuitBreakerStats)
	router.GET("/health/runtime", healthHandler.GetRuntimeStats)

	// pprof endpoints for profiling (load testing)
	router.GET("/debug/pprof/*any", gin.WrapH(healthHandler.PprofGroup()))

	// Serve uploaded files statically (only for local storage)
	if config.AppConfig.Storage.Type != "r2" {
		router.Static(config.AppConfig.Storage.BaseURL, config.AppConfig.Storage.UploadDir)
	}

	// Google Calendar OAuth2 callback route (no auth required - called by Google)
	// Registered at root level before v1 group to ensure proper route matching
	router.GET("/api/v1/google-calendar/callback", googleCalendarAuthHandler.HandleCallback)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			response.SuccessResponse(c, gin.H{
				"message": "API v1",
				"version": "1.0.0",
			}, nil)
		})

		// Auth routes
		routes.SetupAuthRoutes(v1, authHandler, permissionHandler, userHandler, jwtManager)

		// User routes
		routes.SetupUserRoutes(v1, userHandler, permissionHandler, jwtManager, scopeMiddleware)

		// Role routes
		routes.SetupRoleRoutes(v1, roleHandler, jwtManager)

		// Permission routes
		routes.SetupPermissionRoutes(v1, permissionHandler, jwtManager)

		// Menu routes
		routes.SetupMenuRoutes(v1, menuHandler, jwtManager)

		// Group routes
		routes.SetupGroupRoutes(v1, groupHandler, jwtManager)

		// Monthly Target routes
		routes.SetupMonthlyTargetRoutes(v1, monthlyTargetHandler, jwtManager, scopeMiddleware)

		// Brick routes
		routes.SetupBrickRoutes(v1, brickHandler, brickTargetDistributionHandler, brickAnalyticsHandler, jwtManager)

		// Category routes
		routes.SetupCategoryRoutes(v1, categoryHandler, jwtManager)

		// Contact Role routes
		routes.SetupContactRoleRoutes(v1, contactRoleHandler, jwtManager)

		// Account routes
		routes.SetupAccountRoutes(v1, accountHandler, jwtManager)

		// Contact routes
		routes.SetupContactRoutes(v1, contactHandler, jwtManager)

		// Visit Report routes
		routes.SetupVisitReportRoutes(v1, visitReportHandler, activityTypeHandler, jwtManager, scopeMiddleware, permissionService)

		// Activity routes
		routes.SetupActivityRoutes(v1, activityHandler, jwtManager, scopeMiddleware)

		// Pipeline & Deals routes
		routes.SetupPipelineRoutes(v1, pipelineHandler, dealHandler, jwtManager, scopeMiddleware)

		// Lead routes
		routes.SetupLeadRoutes(v1, leadHandler, jwtManager, scopeMiddleware)

		// Lead Status routes
		routes.SetupLeadStatusRoutes(v1, leadStatusHandler, jwtManager)

		// Industry routes
		routes.SetupIndustryRoutes(v1, industryHandler, jwtManager)

		// Lead Source routes
		routes.SetupLeadSourceRoutes(v1, leadSourceHandler, jwtManager)

		// Dashboard routes
		routes.SetupDashboardRoutes(v1, dashboardHandler, jwtManager, scopeMiddleware)

		// Report routes
		routes.SetupReportRoutes(v1, reportHandler, jwtManager, permissionService)

		// Master Data routes

		// Product routes
		routes.SetupProductRoutes(v1, productHandler, jwtManager)

		// Customer Purchase routes
		routes.SetupCustomerPurchaseRoutes(v1, customerPurchaseHandler, jwtManager)

		// Product Analytics routes
		routes.SetupProductAnalyticsRoutes(v1, productAnalyticsHandler, jwtManager, scopeMiddleware)

		// File upload routes
		routes.SetupFileRoutes(v1, fileHandler, jwtManager)

		// Task & Reminder routes
		routes.SetupTaskRoutes(v1, taskHandler, jwtManager, scopeMiddleware)

		// Notification routes
		routes.SetupNotificationRoutes(v1, notificationHandler, wsHandler, jwtManager)

		// AI routes
		routes.SetupAIRoutes(v1, aiHandler, aiSettingsHandler, jwtManager, scopeMiddleware)

		// Route Optimization routes
		routes.SetupRouteOptimizationRoutes(v1, routeOptimizationHandler, jwtManager)

		// Geocoding routes
		routes.SetupGeocodingRoutes(v1, geocodingHandler, jwtManager)

		// Schedule routes
		routes.SetupScheduleRoutes(v1, scheduleHandler, googleCalendarAuthHandler, jwtManager, scopeMiddleware)

		// Sales Overview routes
		routes.SetupSalesOverviewRoutes(v1, salesOverviewHandler, jwtManager, scopeMiddleware)

		// Area Mapping routes
		routes.SetupAreaMappingRoutes(v1, areaMappingHandler, jwtManager)
	}

	return router
}

// shouldRunSeeders checks if we should run database seeders
// This function ensures that seeders are NEVER run in production mode
//
// PRODUCTION SAFETY:
// - Seeders are disabled by default (safe for production)
// - Must be explicitly enabled with ENABLE_SEEDERS=true
// - Production mode (ENV=production) blocks seeders even if ENABLE_SEEDERS=true
// - No code changes needed for production deployment
func shouldRunSeeders() bool {
	// CRITICAL: Never run seeders in production
	env := os.Getenv("ENV")
	if env == "production" || env == "prod" {
		log.Println("🔒 Production mode detected: Seeders are disabled (safety protection)")
		return false
	}

	// Only run seeders in development mode when explicitly enabled
	enableSeeders := os.Getenv("ENABLE_SEEDERS")
	if enableSeeders == "true" || enableSeeders == "1" {
		if env == "" || env == "development" || env == "dev" {
			log.Println("🔧 Development mode: Seeders are enabled")
			return true
		}
		log.Printf("⚠️  Warning: ENABLE_SEEDERS is set but ENV=%s is not development. Skipping seeders.", env)
	}
	return false
}
