package configreader

import (
	"fmt"

	"github.com/spf13/viper"
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

func ReadConfigFileOrEnv(configFilePath string) (*Config, error) {
	// If a config file path is provided, read the configuration from the file, for local development or testing.
	if configFilePath != "" {
		return ReadConfigFile(configFilePath)
	}

	// Read the configuration from environment variables, for production or containerized environments.
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	config := &Config{
		MySQL: DBConfig{
			Host:              viper.GetString("MYSQL_HOST"),
			Port:              viper.GetInt("MYSQL_PORT"),
			Username:          viper.GetString("MYSQL_USERNAME"),
			Password:          viper.GetString("MYSQL_PASSWORD"),
			Database:          viper.GetString("MYSQL_DATABASE"),
			ConnectionTimeout: viper.GetInt("MYSQL_CONNECTION_TIMEOUT"),
			MaxIdleConns:      viper.GetInt("MYSQL_MAX_IDLE_CONNS"),
			MaxOpenConns:      viper.GetInt("MYSQL_MAX_OPEN_CONNS"),
			ConnMaxLifetime:   viper.GetInt("MYSQL_CONN_MAX_LIFETIME"),
			ParseTime:         viper.GetBool("MYSQL_PARSE_TIME"),
			ReadTimeout:       viper.GetInt("MYSQL_READ_TIMEOUT"),
		},
		Connection: HTTPServerConfig{
			Host:         viper.GetString("HTTP_HOST"),
			Port:         viper.GetInt("HTTP_PORT"),
			ReadTimeout:  viper.GetInt("HTTP_READ_TIMEOUT"),
			WriteTimeout: viper.GetInt("HTTP_WRITE_TIMEOUT"),
			IdleTimeout:  viper.GetInt("HTTP_IDLE_TIMEOUT"),
		},
	}
	return config, nil

}

// ReadConfigFile reads the configuration from the specified file path and performs validation.
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
	return &config, nil
}
