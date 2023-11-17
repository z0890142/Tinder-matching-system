package config

import (
	"sync"
)

var globalConfig *Config
var configOnce sync.Once

// ResetConfig set config to Nil, used for tests
func ResetConfig() {
	globalConfig = nil
}

// GetConfig 獲取該服務相關配置
func GetConfig() *Config {
	configOnce.Do(func() {
		globalConfig = &Config{}
	})
	return globalConfig
}

// Config 該服務相關配置
type Config struct {
	Env      string  `mapstructure:"ENV"`
	Service  Service `mapstructure:"SERVICE"`
	LogLevel string  `mapstructure:"LOG_LEVEL"`
	LogFile  string  `mapstructure:"LOG_FILE"`
}

// Service defines service configuration struct.
type Service struct {
	Name string `mapstructure:"NAME"`
	Host string `mapstructure:"HOST"`
	Port int    `mapstructure:"PORT"`
	Url  string `mapstructure:"URL"`
}
