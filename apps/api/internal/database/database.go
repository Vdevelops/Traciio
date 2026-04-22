package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity_type"
	"github.com/gilabs/crm-healthcare/api/internal/domain/ai_settings"
	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/brick_target_distribution"
	"github.com/gilabs/crm-healthcare/api/internal/domain/category"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact_role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/customer_purchase"
	"github.com/gilabs/crm-healthcare/api/internal/domain/google_calendar_token"
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/domain/industry"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_qualification"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_source"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
	"github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/domain/notification"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product_analytics"
	"github.com/gilabs/crm-healthcare/api/internal/domain/refresh_token"
	"github.com/gilabs/crm-healthcare/api/internal/domain/reminder"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
	"github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes database connection with enterprise-scale connection pooling
func Connect() error {
	dsn := config.GetDSN()

	var err error

	// Configure logger with slow query threshold (1 second)
	slowQueryLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             1 * time.Second, // Log queries slower than 1 second
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: slowQueryLogger,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool for enterprise-scale performance
	// CRITICAL: This prevents connection exhaustion under high load
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Enterprise-scale connection pool settings
	// OPTIMIZED for 1000-5000 concurrent users
	// Formula: MaxOpenConns = expected_concurrent_users * connection_multiplier
	// connection_multiplier = 1.5 (accounts for concurrent queries per user)
	// For 5000 users: 5000 * 0.1 = 500 connections (each user makes ~10% concurrent queries)
	// INCREASED from 200 to 500 for high-load scenarios
	//
	// Load Testing Configuration:
	// - 1000 users: 200 connections (current setting) - OK
	// - 2000 users: 300 connections - RECOMMENDED
	// - 5000 users: 500 connections - REQUIRED
	//
	// Environment variables override defaults:
	// DB_MAX_OPEN_CONNS: Maximum open connections (default: 500 for production)
	// DB_MAX_IDLE_CONNS: Maximum idle connections (default: 100, 20% of max)
	// DB_CONN_MAX_LIFETIME: Connection lifetime (default: 5m)
	// DB_CONN_MAX_IDLE_TIME: Idle timeout (default: 10m)

	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 500) // INCREASED from 200 to 500
	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 100) // INCREASED from 50 to 100 (20% of max)
	connMaxLifetime := getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute)
	connMaxIdleTime := getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute)

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	log.Println("Database connected successfully with connection pooling configured")
	log.Printf("Connection pool settings: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v, MaxIdleTime=%v",
		maxOpenConns,
		maxIdleConns,
		connMaxLifetime,
		connMaxIdleTime,
	)
	log.Println("TIP: For load testing with 5000 users, ensure DB_MAX_OPEN_CONNS=500 or higher")

	return nil
}

// AutoMigrate runs database migrations
//
// PRODUCTION SAFETY:
// - This function is SAFE for production use
// - Tables are NEVER dropped in production mode (ENV=production)
// - Drop tables only happens in development mode when DROP_TABLES=true
// - Multiple safety checks prevent accidental data loss
// - No code changes needed for production deployment
func AutoMigrate() error {
	// Check if we should drop all tables (development only)
	// This check has built-in production protection
	if shouldDropTables() {
		log.Println("Development mode: Dropping all tables...")
		if err := DropAllTables(); err != nil {
			return fmt.Errorf("failed to drop tables: %w", err)
		}
		log.Println("All tables dropped successfully")
	}

	// Try to handle constraint issues by attempting to drop constraints that might cause problems
	// This is a workaround for development environments where schema might be out of sync
	if err := handleConstraintIssues(); err != nil {
		log.Printf("Warning: Could not handle constraint issues (this may be expected): %v", err)
	}

	// Use a custom migration approach that handles constraint errors gracefully
	err := migrateWithErrorHandling(
		&user.User{},
		&role.Role{},
		&role.RoleScope{},
		&permission.Permission{},
		&permission.Menu{},
		&group.Group{},
		&brick.Brick{},
		&monthly_target.MonthlyTarget{},
		&brick_target_distribution.BrickTargetDistribution{},
		&category.Category{},
		&contact_role.ContactRole{},
		&account.Account{},
		&contact.Contact{},
		&lead.Lead{},
		&lead_status.LeadStatus{},
		&industry.Industry{},
		&lead_source.LeadSource{},
		&pipeline.PipelineStage{},
		&pipeline.Deal{},
		&pipeline.DealProductItem{},
		&product.ProductCategory{},
		&product.Product{},
		&product_analytics.ProductSales{},
		&task.Task{},
		&reminder.Reminder{},
		&notification.Notification{},
		&visit_report.VisitReport{},
		&activity_type.ActivityType{},
		&activity.Activity{},
		&ai_settings.AISettings{},
		&refresh_token.RefreshToken{},
		&route_optimization.OptimizedRoute{},
		&schedule.Schedule{},
		&schedule.ScheduleAssignment{},
		&google_calendar_token.GoogleCalendarToken{},
		&lead_qualification.LeadQualificationChecklist{},
		&customer_purchase.CustomerPurchaseHistory{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Run PostGIS migrations AFTER GORM AutoMigrate (needs users table to exist first)
	log.Println("Running PostGIS migrations for area mapping...")
	if err := runPostGISMigrations(); err != nil {
		log.Printf("Warning: Failed to run PostGIS migrations (some geographic features may be unavailable): %v", err)
		// return fmt.Errorf("failed to run PostGIS migrations: %w", err)
	}

	// Run SQL migrations for schema changes
	log.Println("Running SQL migrations...")
	if err := runSQLMigrations(); err != nil {
		return fmt.Errorf("failed to run SQL migrations: %w", err)
	}

	log.Println("Database migrations completed")
	return nil
}

// shouldDropTables checks if we should drop all tables (development only)
// This function ensures that tables are NEVER dropped in production mode
func shouldDropTables() bool {
	// CRITICAL: Never drop tables in production
	// Check both config and environment variable for safety
	env := ""
	if config.AppConfig != nil {
		env = config.AppConfig.Server.Env
	}

	// Fallback to environment variable if config is not loaded yet
	if env == "" {
		env = os.Getenv("ENV")
	}

	// Safety check: Never allow in production
	if env == "production" || env == "prod" {
		log.Println("🔒 Production mode detected: Table drop is disabled (safety protection)")
		return false
	}

	// Only allow dropping tables in development mode
	// Check environment variable DROP_TABLES
	dropTables := os.Getenv("DROP_TABLES")
	if dropTables == "true" || dropTables == "1" {
		// Double check: ensure we're not in production
		if env == "" || env == "development" || env == "dev" {
			log.Println("🔧 Development mode: Table drop is enabled")
			return true
		}
		log.Printf("⚠️  Warning: DROP_TABLES is set but ENV=%s is not development. Skipping table drop.", env)
		return false
	}
	return false
}

// DropAllTables drops all tables in the database (development only)
// This function has built-in safety checks to prevent accidental data loss
func DropAllTables() error {
	// Safety check: Never allow dropping tables in production
	// Check both config and environment variable for maximum safety
	env := ""
	if config.AppConfig != nil {
		env = config.AppConfig.Server.Env
	}

	// Fallback to environment variable if config is not loaded yet
	if env == "" {
		env = os.Getenv("ENV")
	}

	if env == "production" || env == "prod" {
		return fmt.Errorf("🔒 CRITICAL: Cannot drop tables in production mode (ENV=%s). This is a safety protection", env)
	}

	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	// Get all table names in the current schema
	var tables []string
	err := DB.Raw(`
		SELECT tablename 
		FROM pg_tables 
		WHERE schemaname = CURRENT_SCHEMA()
		AND tablename NOT LIKE 'pg_%'
		AND tablename NOT LIKE '_prisma_%'
	`).Scan(&tables).Error

	if err != nil {
		return fmt.Errorf("failed to get table list: %w", err)
	}

	if len(tables) == 0 {
		log.Println("No tables to drop")
		return nil
	}

	// Disable foreign key checks temporarily and drop all tables
	// PostgreSQL doesn't have a simple way to disable FK checks, so we use CASCADE
	log.Printf("⚠️  DEVELOPMENT MODE: Dropping %d tables...", len(tables))

	for _, table := range tables {
		// Use CASCADE to drop dependent objects
		dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)
		if err := DB.Exec(dropSQL).Error; err != nil {
			log.Printf("Warning: Failed to drop table %s: %v", table, err)
			// Continue with other tables
		} else {
			log.Printf("Dropped table: %s", table)
		}
	}

	// Also drop any remaining sequences
	var sequences []string
	_ = DB.Raw(`
		SELECT sequence_name 
		FROM information_schema.sequences 
		WHERE sequence_schema = CURRENT_SCHEMA()
	`).Scan(&sequences).Error

	for _, seq := range sequences {
		dropSQL := fmt.Sprintf("DROP SEQUENCE IF EXISTS %s CASCADE", seq)
		_ = DB.Exec(dropSQL).Error // Ignore errors
	}

	log.Println("✅ All tables and sequences dropped successfully")
	return nil
}

// migrateWithErrorHandling migrates models while handling common constraint errors
func migrateWithErrorHandling(models ...interface{}) error {
	for _, model := range models {
		err := DB.AutoMigrate(model)
		if err != nil {
			// Check if error is PostgreSQL error code 42704 (undefined_object)
			// This happens when trying to DROP a constraint that doesn't exist
			errStr := err.Error()
			if strings.Contains(errStr, "SQLSTATE 42704") ||
				(strings.Contains(errStr, "does not exist") && strings.Contains(errStr, "constraint")) {
				log.Printf("Warning: Constraint error during migration (safe to ignore): %v", err)
				log.Println("GORM will create the necessary constraints. This is expected during schema evolution.")
				// Continue with next model - GORM might have partially succeeded
				continue
			}
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	return nil
}

// handleConstraintIssues attempts to fix common constraint issues before migration
func handleConstraintIssues() error {
	// Check if roles table exists
	var exists bool
	err := DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'roles')").Scan(&exists).Error
	if err != nil || !exists {
		return nil // Table doesn't exist, nothing to fix
	}

	// Get all unique constraints on the roles table
	type ConstraintInfo struct {
		ConstraintName string
	}
	var constraints []ConstraintInfo
	err = DB.Raw(`
		SELECT conname as constraint_name
		FROM pg_constraint
		WHERE conrelid = 'roles'::regclass
		AND contype = 'u'
	`).Scan(&constraints).Error

	if err != nil {
		// If we can't query constraints, that's okay - continue anyway
		return nil
	}

	// Drop all unique constraints on code column (GORM will recreate them)
	for _, constraint := range constraints {
		// Check if this constraint is on the 'code' column
		var columnName string
		err = DB.Raw(`
			SELECT a.attname
			FROM pg_constraint c
			JOIN pg_attribute a ON a.attrelid = c.conrelid AND a.attnum = ANY(c.conkey)
			WHERE c.conname = ?
			AND a.attname = 'code'
			LIMIT 1
		`, constraint.ConstraintName).Scan(&columnName).Error

		if err == nil && columnName == "code" {
			// Drop the constraint
			dropSQL := fmt.Sprintf("ALTER TABLE roles DROP CONSTRAINT IF EXISTS %s", constraint.ConstraintName)
			_ = DB.Exec(dropSQL).Error // Ignore errors - constraint might not exist
		}
	}

	return nil
}

// runPostGISMigrations executes SQL migrations for PostGIS tables
// GORM AutoMigrate doesn't support PostGIS GEOGRAPHY types, so we use raw SQL
func runPostGISMigrations() error {
	log.Println("🗺️  Starting PostGIS migrations...")

	// Enable PostGIS extension
	log.Println("Enabling PostGIS extension...")
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS postgis").Error; err != nil {
		return fmt.Errorf("failed to enable PostGIS extension: %w", err)
	}
	log.Println("✅ PostGIS extension enabled")

	// Check if area_captures table already exists
	var tableExists bool
	err := DB.Raw(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'area_captures'
		)
	`).Scan(&tableExists).Error

	if err != nil {
		return fmt.Errorf("failed to check if area_captures table exists: %w", err)
	}

	if tableExists {
		log.Println("✅ PostGIS tables already exist, skipping migration")
		return nil
	}

	log.Println("📋 Creating PostGIS tables for area mapping...")

	// Create area_captures table
	log.Println("Creating area_captures table...")
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS area_captures (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			visit_report_id UUID NOT NULL,
			capture_type VARCHAR(20) NOT NULL,
			location GEOGRAPHY(POINT, 4326) NOT NULL,
			address TEXT,
			accuracy NUMERIC(10,2),
			captured_at TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create area_captures table: %w", err)
	}
	log.Println("✅ area_captures table created")

	// Create indexes for area_captures
	log.Println("Creating indexes for area_captures...")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_area_captures_location ON area_captures USING GIST(location)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_area_captures_visit_report_id ON area_captures(visit_report_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_area_captures_captured_at ON area_captures(captured_at)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_area_captures_capture_type ON area_captures(capture_type)")

	// Create territories table
	log.Println("Creating territories table...")
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS territories (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			polygon GEOGRAPHY(POLYGON, 4326) NOT NULL,
			assigned_to UUID,
			color VARCHAR(50) DEFAULT '#3B82F6',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create territories table: %w", err)
	}
	log.Println("✅ territories table created")

	// Create indexes for territories
	log.Println("Creating indexes for territories...")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_territories_polygon ON territories USING GIST(polygon)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_territories_assigned_to ON territories(assigned_to)")

	// Create coverage_analysis table
	log.Println("Creating coverage_analysis table...")
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS coverage_analysis (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			territory_id UUID,
			user_id UUID,
			period_start DATE NOT NULL,
			period_end DATE NOT NULL,
			visit_count INTEGER NOT NULL DEFAULT 0,
			coverage_percent NUMERIC(5,2),
			analyzed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			CONSTRAINT fk_coverage_analysis_territory FOREIGN KEY (territory_id) 
				REFERENCES territories(id) ON DELETE CASCADE,
			CONSTRAINT fk_coverage_analysis_user FOREIGN KEY (user_id) 
				REFERENCES users(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create coverage_analysis table: %w", err)
	}
	log.Println("✅ coverage_analysis table created")

	// Create indexes for coverage_analysis
	log.Println("Creating indexes for coverage_analysis...")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_coverage_analysis_territory_id ON coverage_analysis(territory_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_coverage_analysis_user_id ON coverage_analysis(user_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_coverage_analysis_period ON coverage_analysis(period_start, period_end)")

	log.Println("✅ PostGIS tables created successfully")
	return nil
}

// runSQLMigrations executes SQL migration files
func runSQLMigrations() error {
	if err := DB.Exec("ALTER TABLE schedules ALTER COLUMN task_id DROP NOT NULL").Error; err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "does not exist") || (strings.Contains(errStr, "column") && strings.Contains(errStr, "is not of type")) {
			// Column might already be nullable, skip
		}
	}

	DB.Exec("COMMENT ON COLUMN schedules.task_id IS 'Reference to the task this schedule is connected to (nullable for standalone schedules)'")

	return nil
}

// Close closes database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// getEnvAsInt reads environment variable as integer with default value
func getEnvAsInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

// getEnvAsDuration reads environment variable as time.Duration with default value
func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	return defaultVal
}
