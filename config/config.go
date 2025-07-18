package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

//var appConfig = flag.String("config", "config/app.yaml", "application config path")

type Config struct {
	Server       ServerConfig           `yaml:"server"`
	DB           DBConfig               `yaml:"db"`
	Redis        RedisConfig            `yaml:"redis"`
	Storage      StorageConfig          `yaml:"storage"`
	Logger       LoggerConfig           `yaml:"logger"`
	JWTConfig    JWTConfig              `yaml:"jwt"`
	OAuthConfigs map[string]OAuthConfig `yaml:"oauth"` // 支持多种 OAuth 配置
}

type ServerConfig struct {
	ENV                    string            `yaml:"env"`
	Address                string            `yaml:"address"`
	Port                   int               `yaml:"port"`
	GracefulShutdownPeriod int               `yaml:"gracefulShutdownPeriod"`
	RateLimitsConfigs      RateLimitsConfigs `yaml:"rateLimits"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int64  `yaml:"expire"` // Token过期时间，单位为秒
	Issuer string `yaml:"issuer"`
}

type RateLimitsConfigs struct {
	FillInterval time.Duration `yaml:"fillInterval"`
	Cap          int64         `yaml:"cap"`
	Quantum      int64         `yaml:"quantum"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Migrate  bool   `yaml:"migrate"`
}

type StorageConfig struct {
	AccessKey  string `yaml:"accessKey"`
	SecretKey  string `yaml:"secretKey"`
	Bucket     string `yaml:"bucket"`
	StorageUrl string `yaml:"storageUrl"`
}

type RedisConfig struct {
	Enable   bool   `yaml:"enable"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type LoggerConfig struct {
	Level      string `yaml:"level" json:"level"`
	Filename   string `yaml:"filename" json:"filename"`
	MaxSize    int    `yaml:"maxsize" json:"maxsize"`
	MaxAge     int    `yaml:"maxage" json:"maxage"`
	MaxBackups int    `yaml:"maxbackups" json:"maxbackups"`
}

type OAuthConfig struct {
	AuthType     string `yaml:"authType"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}

// LoadConfig 从指定的配置文件加载应用程序配置
func LoadConfig(appconfig string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(appconfig)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
