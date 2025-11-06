package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Http web配置结构体
type Http struct {
	Port int `yaml:"port"`
	// Domain   string `yaml:"domain"`
	// Protocol string `yaml:"protocol"`
}

// MySQL 配置结构体
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
	//Collation string `yaml:"collation"`
}

// LogConf 日志相关配置
type LogConf struct {
	Level string `yaml:"level"`
}

type JWTConf struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}

type Config struct {
	Http  Http    `yaml:"http"`
	MySQL MySQL   `yaml:"mysql"`
	Log   LogConf `yaml:"log"`
	JWT   JWTConf `yaml:"jwt"`
}

var conf *Config

func GetConf() *Config {
	if conf == nil {
		// 如果配置未初始化，尝试初始化
		if err := InitConfig(); err != nil {
			fmt.Printf("Warning: Failed to initialize config: %v\n", err)
		}
	}
	return conf
}

// InitConfig 初始化配置
func InitConfig() error {
	configPath := os.Getenv("CONFIG_PATH")
	//fmt.Printf("CONFIG_PATH from env: %s\n", configPath)
	if configPath == "" {
		configPath = "./config/config.yaml"
	}
	//fmt.Printf("CONFIG_PATH: %s\n", configPath)

	// 初始化conf结构体
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

// getYamlConf 解析yaml配置
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

// validateConfig 验证配置有效性
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
