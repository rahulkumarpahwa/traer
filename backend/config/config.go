package config

import (
	"github.com/rahulkumarpahwa/traer/bootstrap"
)

type Config struct {
	App struct {
		Name string `yaml:"name" env:"APP_NAME" env-default:"traer-backend"`
		Env  string `yaml:"env" env:"APP_ENV" env-default:"development"`
	} `yaml:"app"`

	HTTPServer struct {
		Address      string `yaml:"address" env:"HTTP_ADDRESS" env-default:"localhost:8080"`
		ReadTimeout  int    `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"20"`
		WriteTimeout int    `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"20"`
		IdleTimeout  int    `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"20"`
	} `yaml:"http_server"`

	Workers struct {
		Instances int `yaml:"instances" env:"WORKER_INSTANCES" env-default:"2"`
		QueueSize int `yaml:"queue_size" env:"QUEUE_SIZE" env-default:"100"`
	} `yaml:"workers"`

	Jobs struct {
		MaxDuration     string `yaml:"max_duration" env:"JOB_MAX_DURATION" env-default:"30m"`
		CleanupInterval string `yaml:"cleanup_interval" env:"JOB_CLEANUP_INTERVAL" env-default:"10m"`
		RetentionTime   string `yaml:"retention_time" env:"JOB_RETENTION_TIME" env-default:"1h"`
	} `yaml:"jobs"`

	Storage struct {
		DownloadPath string `yaml:"download_path" env:"DOWNLOAD_PATH" env-default:"downloads/"`
		MaxFileSize  int    `yaml:"max_file_size_mb" env:"MAX_FILE_SIZE_MB" env-default:"1024"`
		AutoCleanup  bool   `yaml:"auto_cleanup" env:"AUTO_CLEANUP" env-default:"true"`
	} `yaml:"storage"`

	// Nextcloud holds the credentials and destination used by the music
	// pipeline to upload finished MP3 files over WebDAV. These values are read
	// from environment variables (or a git-ignored local config) and are never
	// exposed to the frontend or written to logs.
	Nextcloud struct {
		URL         string `yaml:"url" env:"NEXTCLOUD_URL"`
		Username    string `yaml:"username" env:"NEXTCLOUD_USERNAME"`
		AppPassword string `yaml:"app_password" env:"NEXTCLOUD_APP_PASSWORD"`
		BasePath    string `yaml:"base_path" env:"NEXTCLOUD_BASE_PATH" env-default:"Music/Downloads"`
	} `yaml:"nextcloud"`

	Executables struct {
		Required []string `yaml:"required"`
	} `yaml:"executables"`

	Database struct {
		Path string `yaml:"path" env:"DB_PATH" env-default:"database/database.db"`
	} `yaml:"database"`

	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
		Format string `yaml:"format" env:"LOG_FORMAT" env-default:"text"`
	} `yaml:"logging"`

	Jwt struct {
		Secret string `yaml:"secret" env:"JWT_SECRET"`
	} `yaml:"jwt"`
}

// Configured reports whether the Nextcloud integration has the minimum required
// settings to attempt an upload.
func (c *Config) NextcloudConfigured() bool {
	return c.Nextcloud.URL != "" && c.Nextcloud.Username != "" && c.Nextcloud.AppPassword != ""
}

// ResolveExecutables resolves all required executables and returns a map of tool names to paths.
// It uses the DependencyManager to find tools in the system PATH and known install locations.
func (c *Config) ResolveExecutables() (map[string]string, error) {
	depManager := bootstrap.NewDependencyManager()

	if err := depManager.Ensure(c.Executables.Required); err != nil {
		return nil, err
	}

	return depManager.Required, nil
}
