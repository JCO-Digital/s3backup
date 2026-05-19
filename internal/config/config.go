package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Endpoint   string `mapstructure:"endpoint"`
	Region     string `mapstructure:"region"`
	AccessKey  string `mapstructure:"access_key"`
	SecretKey  string `mapstructure:"secret_key"`
	Bucket     string `mapstructure:"bucket"`
	TargetPath string `mapstructure:"target_path"`
	MaxAge     int    `mapstructure:"max_age"`
	KeepFree   string `mapstructure:"keep_free"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.MaxAge == 0 {
		cfg.MaxAge = 7
	}

	if cfg.KeepFree == "" {
		cfg.KeepFree = "100G"
	}

	return &cfg, nil
}
