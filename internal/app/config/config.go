package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int

	JWT   JWTConfig
	Redis RedisConfig
}
type RedisConfig struct {
	Host        string
	Password    string
	Port        string
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

type JWTConfig struct {
	ExpiresIn time.Duration
	Token     string
}

const (
	envRedisHost = "REDIS_HOST"
	envRedisPort = "REDIS_PORT"
	envRedisUser = "REDIS_USER"
	envRedisPass = "REDIS_PASSWORD"
)

func NewConfig() (*Config, error) {
	var err error

	configName := "config"

	_ = godotenv.Load()

	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}
	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	cfg.Redis.Host = os.Getenv(envRedisHost)
	cfg.Redis.Port = os.Getenv(envRedisPort)
	cfg.Redis.Password = os.Getenv(envRedisPass)
	cfg.Redis.User = os.Getenv(envRedisUser)
	cfg.JWT.Token = os.Getenv("JWT_TOKEN")
	cfg.JWT.ExpiresIn = time.Hour
	log.Info("config parsed")

	return cfg, nil
}
