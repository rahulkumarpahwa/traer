package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rahulkumarpahwa/traer/bootstrap"
)

type Config struct {
	AppName string
	AppEnv  string

	AdminEmail    string
	AdminPassword string

	HTTPAddress      string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration

	WorkerInstances int
	WorkerQueueSize int

	JobMaxDuration     time.Duration
	JobCleanupInterval time.Duration
	JobRetentionTime   time.Duration

	StorageDownloadPath string
	StorageMaxFileSize  int
	StorageAutoCleanup  bool

	Executables string

	DatabasePath string

	LoggingLevel  string
	LoggingFormat string

	JWTSecret string

	NextcloudURL         string
	NextcloudUsername    string
	NextcloudAppPassword string
	NextcloudBasePath    string
}

func getInt(key string, fallback int) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return v
}

func getBool(key string, fallback bool) bool {
	v, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return v
}

func getDuration(key string, fallback time.Duration) time.Duration {
	v, err := time.ParseDuration(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return v
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using system environment variables")
	}

	return &Config{
		AppName: os.Getenv("APP_NAME"),
		AppEnv:  os.Getenv("APP_ENV"),

		AdminEmail:    os.Getenv("ADMIN_USERNAME"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),

		HTTPAddress:      os.Getenv("HTTP_SERVER_ADDRESS"),
		HTTPReadTimeout:  time.Duration(getInt("HTTP_SERVER_READ_TIMEOUT", 60)) * time.Second,
		HTTPWriteTimeout: time.Duration(getInt("HTTP_SERVER_WRITE_TIMEOUT", 3600)) * time.Second,
		HTTPIdleTimeout:  time.Duration(getInt("HTTP_SERVER_IDLE_TIMEOUT", 120)) * time.Second,

		WorkerInstances: getInt("WORKERS_INSTANCES", 2),
		WorkerQueueSize: getInt("WORKERS_QUEUE_SIZE", 100),

		JobMaxDuration:     getDuration("JOBS_MAX_DURATION", 30*time.Minute),
		JobCleanupInterval: getDuration("JOBS_CLEANUP_INTERVAL", 10*time.Minute),
		JobRetentionTime:   getDuration("JOBS_RETENTION_TIME", time.Hour),

		StorageDownloadPath: os.Getenv("STORAGE_DOWNLOAD_PATH"),
		StorageMaxFileSize:  getInt("STORAGE_MAX_FILE_SIZE_MB", 1024),
		StorageAutoCleanup:  getBool("STORAGE_AUTO_CLEANUP", true),

		Executables: os.Getenv("EXECUTABLES_REQUIRED"),

		DatabasePath: os.Getenv("DATABASE_PATH"),

		LoggingLevel:  os.Getenv("LOGGING_LEVEL"),
		LoggingFormat: os.Getenv("LOGGING_FORMAT"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		NextcloudURL:         os.Getenv("NEXTCLOUD_URL"),
		NextcloudUsername:    os.Getenv("NEXTCLOUD_USERNAME"),
		NextcloudAppPassword: os.Getenv("NEXTCLOUD_APP_PASSWORD"),
		NextcloudBasePath:    os.Getenv("NEXTCLOUD_BASE_PATH"),
	}
}

// Configured reports whether the Nextcloud integration has the minimum required
// settings to attempt an upload.
func (c *Config) NextcloudConfigured() bool {
	return c.NextcloudURL != "" && c.NextcloudUsername != "" && c.NextcloudAppPassword != ""
}

// ResolveExecutables resolves all required executables and returns a map of tool names to paths.
// It uses the DependencyManager to find tools in the system PATH and known install locations.
func (c *Config) ResolveExecutables() (map[string]string, error) {
	depManager := bootstrap.NewDependencyManager()

	if err := depManager.Ensure(c.Executables); err != nil {
		return nil, err
	}

	return depManager.Required, nil
}
