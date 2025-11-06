package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Http
type Http struct {
	Port int `yaml:"port"`
	// Domain   string `yaml:"domain"`
	// Protocol string `yaml:"protocol"`
}

// MySQL
type MySQL struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"db_name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	// MaxIdleConns    int    `yaml:"max_idle_conns"`
	// MaxOpenConns    int    `yaml:"max_open_conns"`
	// ConnMaxLefeTime int    `yaml:"conn_max_lefe_time"`
	Charset string `yaml:"charset"`
}

// LogConf
type LogConf struct {
	Level string `yaml:"level"`
}

type Config struct {
	Http  Http    `yaml:"http"`
	MySQL MySQL   `yaml:"mysql"`
	Log   LogConf `yaml:"log"`
}

var conf *Config

func GetConf() *Config {
	if conf == nil {

		if err := InitConfig(); err != nil {
			fmt.Printf("Warning: Failed to initialize config: %v\n", err)
		}
	}
	return conf
}

// InitConfig
func InitConfig() error {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "./config/config.yaml"
	}

	conf = &Config{}

	fmt.Printf("Loading configuration from: %s\n", configPath)
	err := getYamlConf(configPath, conf)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// 配置校验
	if err := validateConfig(conf); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return nil
}

// getYamlConf
func getYamlConf(filePath string, out interface{}) error {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, out)
	if err != nil {
		return fmt.Errorf("failed to parse yaml: %w", err)
	}

	return nil
}

// validateConfig
func validateConfig(config *Config) error {
	if config.MySQL.Host == "" {
		return fmt.Errorf("MySQL host is required")
	}
	if config.MySQL.DBName == "" {
		return fmt.Errorf("MySQL database name is required")
	}
	if config.MySQL.User == "" {
		return fmt.Errorf("MySQL username is required")
	}
	if config.Http.Port == 0 {
		config.Http.Port = 8090 // 设置默认端口
	}
	return nil
}
