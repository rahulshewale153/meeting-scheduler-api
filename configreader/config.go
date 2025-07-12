package configreader

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	_cfg Config
)

// Config represents the parsed configuration from the file.
type Config struct {
	MySQL      DBConfig
	Connection HTTPServerConfig
}

// DBConfig represents the configuration for a specific database connection.
type DBConfig struct {
	Host              string
	Port              int
	Username          string
	Password          string
	Database          string
	ConnectionTimeout int
	MaxIdleConns      int
	MaxOpenConns      int
	ConnMaxLifetime   int
	ParseTime         bool
	ReadTimeout       int
}

// HTTPServerConfig represents the configuration for a host and port combination.
type HTTPServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

// ReadConfigFile reads the configuration from the specified file path and performs validation.
// It returns a pointer to the Config struct if successful, or an error if any occurred.
func ReadConfigFile(configFilePath string) (*Config, error) {
	viper.SetConfigFile(configFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	_cfg = config
	return &config, nil
}

// GetAPIServiceConfig return configuration details
func GetAPIServiceConfig() *Config {
	return &_cfg
}
