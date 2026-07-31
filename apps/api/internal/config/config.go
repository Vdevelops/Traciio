package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server         ServerConfig
	Database       DatabaseConfig
	JWT            JWTConfig
	Cerebras       CerebrasConfig
	Storage        StorageConfig
	KPI            KPIConfig
	RateLimit      RateLimitConfig
	HSTS           HSTSConfig
	OSRM           OSRMConfig
	GoogleCalendar GoogleCalendarConfig
	Encryption     EncryptionConfig
	Geocoding      GeocodingConfig
	Redis          RedisConfig
	Kafka          KafkaConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  int // in hours
	RefreshTokenTTL int // in days
}

type CerebrasConfig struct {
	BaseURL string
	APIKey  string
	Model   string // Default model name
}

type StorageConfig struct {
	Type      string // Storage type: "local" or "r2"
	UploadDir string // Directory for uploaded files (local storage only)
	BaseURL   string // Base URL for serving files (e.g., /uploads or https://cdn.example.com)
	// R2 Configuration
	R2Endpoint        string // R2 endpoint URL (e.g., https://<account-id>.r2.cloudflarestorage.com)
	R2AccessKeyID     string // R2 Access Key ID
	R2SecretAccessKey string // R2 Secret Access Key
	R2Bucket          string // R2 Bucket name
	R2PublicURL       string // Public URL for R2 bucket (e.g., https://<bucket>.<domain>.com or custom domain)
}

// RateLimitRule defines rate limit configuration for a specific endpoint type
type RateLimitRule struct {
	Requests int // Number of requests allowed
	Window   int // Time window in seconds
}

// RateLimitConfig defines rate limit configuration for different endpoint types
type RateLimitConfig struct {
	Login   RateLimitRule // Login endpoint: 5 requests per 15 minutes (Level 1 - IP)
	Refresh RateLimitRule // Refresh token endpoint: 10 requests per hour
	Upload  RateLimitRule // File upload endpoint: 20 requests per hour
	General RateLimitRule // General API endpoints: 100 requests per minute
	Public  RateLimitRule // Public endpoints: 200 requests per minute
	// Multi-level rate limiting for login
	LoginByEmail RateLimitRule // Level 2: 10 attempts per 15 minutes per email
	LoginGlobal  RateLimitRule // Level 3: 100 attempts per minute globally
	Mutation     RateLimitRule // Write operations (POST, PUT, DELETE): 300 requests per 5 minutes
	HighVolume   RateLimitRule // High volume read operations: 3000 requests per 5 minutes
}

// HSTSConfig defines HTTP Strict Transport Security configuration
type HSTSConfig struct {
	MaxAge            int  // Max age in seconds (default: 31536000 = 1 year)
	IncludeSubDomains bool // Include subdomains in HSTS policy
	Preload           bool // Enable HSTS preload
}

// OSRMConfig defines OSRM routing service configuration
type OSRMConfig struct {
	BaseURL string // OSRM server base URL (e.g., https://router.project-osrm.org or http://localhost:5000)
	// Enable OSRM /table for distance matrix (single request) when supported.
	TableEnabled bool
	// Safety limit for /table waypoint count.
	TableMaxWaypoints int
}

// GoogleCalendarConfig defines Google Calendar OAuth2 configuration
type GoogleCalendarConfig struct {
	ClientID     string // Google OAuth2 Client ID
	ClientSecret string // Google OAuth2 Client Secret
	RedirectURL  string // OAuth2 redirect URL - backend callback (e.g. https://api.gilabs.id/api/v1/google-calendar/callback)
	Scopes       string // OAuth2 scopes (comma-separated)
	FrontendURL  string // Frontend URL for post-OAuth redirect (e.g. https://app.gilabs.id)
}

// EncryptionConfig defines encryption configuration for sensitive data
type EncryptionConfig struct {
	Key string // 32-byte key for AES-256 encryption (base64 encoded)
}

// GeocodingConfig defines geocoding service configuration
type GeocodingConfig struct {
	Provider string // "nominatim" (default, free) or "google" (requires API key)
	APIKey   string // Google Maps API key (required if provider is "google")
	Enabled  bool   // Enable/disable geocoding (default: true)
}

// RedisConfig defines Redis cache configuration
type RedisConfig struct {
	Host     string // Redis host
	Port     string // Redis port
	Password string // Redis password (empty for no auth)
	DB       int    // Redis database number
	Enabled  bool   // Enable/disable Redis cache
	// Cache TTL settings (in seconds)
	PermissionsTTL int // TTL for permission cache (default: 300 = 5 minutes)
	SessionTTL     int // TTL for session cache (default: 1800 = 30 minutes)
	GeneralTTL     int // TTL for general cache (default: 300 = 5 minutes)
}

// KafkaConfig defines Kafka event bus configuration
type KafkaConfig struct {
	Enabled       bool   // Enable/disable Kafka event publishing
	Brokers       string // Kafka broker addresses (comma-separated, e.g., "localhost:9092,localhost:9093")
	TopicPrefix   string // Topic prefix for all events (e.g., "crm")
	ConsumerGroup string // Consumer group ID for event consumers
}

var AppConfig *Config

func Load() error {
	// Load .env file if exists (for local development only)
	// Skip .env loading in production to use Docker environment variables
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "crm_healthcare"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey:       getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			AccessTokenTTL:  getEnvAsInt("JWT_ACCESS_TTL", 24), // 24 hours
			RefreshTokenTTL: getEnvAsInt("JWT_REFRESH_TTL", 7), // 7 days
		},
		Cerebras: CerebrasConfig{
			BaseURL: getEnv("CEREBRAS_BASE_URL", "https://api.cerebras.ai"),
			APIKey:  getEnv("CEREBRAS_API_KEY", ""),
			Model:   getEnv("CEREBRAS_MODEL", "gpt-oss-120b"), // Default model
		},
		Storage: StorageConfig{
			Type:              getEnv("STORAGE_TYPE", "local"), // "local" or "r2"
			UploadDir:         getEnv("STORAGE_UPLOAD_DIR", "./uploads"),
			BaseURL:           getEnv("STORAGE_BASE_URL", "/uploads"),
			R2Endpoint:        getEnv("R2_ENDPOINT", ""),
			R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
			R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
			R2Bucket:          getEnv("R2_BUCKET", ""),
			R2PublicURL:       getEnv("R2_PUBLIC_URL", ""),
		},
		KPI: loadKPIConfig(),
		RateLimit: RateLimitConfig{
			Login: RateLimitRule{
				Requests: getEnvAsInt("RATE_LIMIT_LOGIN_REQUESTS", 5), // 5 requests per 15 minutes (Level 1 - IP)
				Window:   getEnvAsInt("RATE_LIMIT_LOGIN_WINDOW", 900), // 15 minutes (900 seconds)
			},
			Refresh: RateLimitRule{
				Requests: getEnvAsInt("RATE_LIMIT_REFRESH_REQUESTS", 10), // 10 requests
				Window:   getEnvAsInt("RATE_LIMIT_REFRESH_WINDOW", 3600), // 1 hour (3600 seconds)
			},
			Upload: RateLimitRule{
				Requests: getEnvAsInt("RATE_LIMIT_UPLOAD_REQUESTS", 20), // 20 requests
				Window:   getEnvAsInt("RATE_LIMIT_UPLOAD_WINDOW", 3600), // 1 hour (3600 seconds)
			},
			General: RateLimitRule{
				Requests: getEnvAsInt("RATE_LIMIT_GENERAL_REQUESTS", 100), // 100 requests
				Window:   getEnvAsInt("RATE_LIMIT_GENERAL_WINDOW", 60),    // 1 minute (60 seconds)
			},
			Public: RateLimitRule{
				Requests: getEnvAsInt("RATE_LIMIT_PUBLIC_REQUESTS", 200), // 200 requests
				Window:   getEnvAsInt("RATE_LIMIT_PUBLIC_WINDOW", 60),    // 1 minute (60 seconds)
			},
			// Level 2: Rate limit by email/username (prevents brute force even if IP changes)
			LoginByEmail: RateLimitRule{
				Requests: getEnvAsInt("RATE_LIMIT_LOGIN_BY_EMAIL_REQUESTS", 10), // 10 requests per 15 minutes per email
				Window:   getEnvAsInt("RATE_LIMIT_LOGIN_BY_EMAIL_WINDOW", 900),  // 15 minutes (900 seconds)
			},
			// Level 3: Global rate limit (prevents DOS on entire system)
			LoginGlobal: RateLimitRule{
				Requests: getEnvAsInt("RATE_LIMIT_LOGIN_GLOBAL_REQUESTS", 100), // 100 requests per minute globally
				Window:   getEnvAsInt("RATE_LIMIT_LOGIN_GLOBAL_WINDOW", 60),    // 1 minute (60 seconds)
			},
			Mutation: RateLimitRule{
				Requests: getEnvAsInt("RATE_LIMIT_MUTATION_REQUESTS", 300), // 300 requests
				Window:   getEnvAsInt("RATE_LIMIT_MUTATION_WINDOW", 300),   // 5 minutes (300 seconds)
			},
			HighVolume: RateLimitRule{
				Requests: getEnvAsInt("RATE_LIMIT_HIGH_VOLUME_REQUESTS", 3000), // 3000 requests
				Window:   getEnvAsInt("RATE_LIMIT_HIGH_VOLUME_WINDOW", 300),    // 5 minutes (300 seconds)
			},
		},
		HSTS: HSTSConfig{
			MaxAge:            getEnvAsInt("HSTS_MAX_AGE", 31536000), // 1 year in seconds
			IncludeSubDomains: getEnv("HSTS_INCLUDE_SUBDOMAINS", "true") == "true",
			Preload:           getEnv("HSTS_PRELOAD", "true") == "true",
		},
		OSRM: OSRMConfig{
			BaseURL:           getEnv("OSRM_BASE_URL", "https://router.project-osrm.org"), // Default to public OSRM instance
			TableEnabled:      getEnv("OSRM_TABLE_ENABLED", "true") == "true",             // Enable by default; public OSRM supports /table
			TableMaxWaypoints: getEnvAsInt("OSRM_TABLE_MAX_WAYPOINTS", 25),
		},
		GoogleCalendar: GoogleCalendarConfig{
			ClientID:     getEnv("GOOGLE_CALENDAR_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CALENDAR_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_CALENDAR_REDIRECT_URL", ""),
			Scopes:       getEnv("GOOGLE_CALENDAR_SCOPES", "https://www.googleapis.com/auth/calendar.events"),
			FrontendURL:  getEnv("GOOGLE_CALENDAR_FRONTEND_URL", "http://localhost:3000"),
		},
		Encryption: EncryptionConfig{
			Key: getEnv("ENCRYPTION_KEY", ""), // Must be 32 bytes, base64 encoded
		},
		Geocoding: GeocodingConfig{
			Provider: getEnv("GEOCODING_PROVIDER", "nominatim"), // "nominatim" or "google"
			APIKey:   getEnv("GEOCODING_API_KEY", ""),           // Google Maps API key (if using Google)
			Enabled:  getEnv("GEOCODING_ENABLED", "true") == "true",
		},
		Redis: RedisConfig{
			Host:           getEnv("REDIS_HOST", "localhost"),
			Port:           getEnv("REDIS_PORT", "6379"),
			Password:       getEnv("REDIS_PASSWORD", ""),
			DB:             getEnvAsInt("REDIS_DB", 0),
			Enabled:        getEnv("REDIS_ENABLED", "true") == "true",
			PermissionsTTL: getEnvAsInt("REDIS_PERMISSIONS_TTL", 300), // 5 minutes
			SessionTTL:     getEnvAsInt("REDIS_SESSION_TTL", 1800),    // 30 minutes
			GeneralTTL:     getEnvAsInt("REDIS_GENERAL_TTL", 300),     // 5 minutes
		},
		Kafka: KafkaConfig{
			Enabled:       getEnv("KAFKA_ENABLED", "false") == "true", // Disabled by default until Kafka is installed
			Brokers:       getEnv("KAFKA_BROKERS", "localhost:9092"),
			TopicPrefix:   getEnv("KAFKA_TOPIC_PREFIX", "crm"),
			ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "crm-api"),
		},
	}

	// SECURITY HARDENING: Enforce strong JWT secrets
	// This prevents the application from starting with weak or default secrets in production
	if len(AppConfig.JWT.SecretKey) < 32 {
		return fmt.Errorf("SECURITY ERROR: JWT_SECRET must be at least 32 characters long. Current length: %d", len(AppConfig.JWT.SecretKey))
	}

	if AppConfig.JWT.SecretKey == "your-secret-key-change-in-production" {
		// Only allow default secret in development environment
		if AppConfig.Server.Env == "production" {
			return fmt.Errorf("SECURITY CRITICAL: Default JWT_SECRET detected in production environment. You MUST change this immediately.")
		}
		// Warn in development
		fmt.Println("⚠️  WARNING: Using default JWT_SECRET. This is unsafe for production!")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetDSN() string {
	db := AppConfig.Database
	// Ensure sslmode is set, default to disable if empty
	sslmode := db.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, sslmode,
	)
}
