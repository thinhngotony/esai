// pkg/config/config.go
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey     string `mapstructure:"API_KEY"`
	ModelName  string `mapstructure:"TEXT_MODEL"`
	ImageModel string `mapstructure:"IMAGE_MODEL"`
}

func NewConfig() (*Config, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("TEXT_MODEL", "gemini-pro")
	v.SetDefault("IMAGE_MODEL", "gemini-1.5-flash")

	// Configure Viper to read from .env file
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("../")    // Look for config in the parent directory
	v.AddConfigPath("../../") // Look for config two directories up

	// Make Viper read from environment variables as well
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Try to read from .env file
	if err := v.ReadInConfig(); err != nil {
		// It's okay if we can't find the .env file, we'll use environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate required fields
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API_KEY is required in environment or .env file")
	}

	return &cfg, nil
}
