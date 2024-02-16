package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Env  string `yaml:"env" env-default:"local"`
		HTTP HTTPConfig
		PG   PGConfig
	}

	HTTPConfig struct {
		Host    string
		Port    string
		Timeout time.Duration
	}

	PGConfig struct {
		Host     string
		Port     string
		DBName   string
		SSLMode  string
		Username string
		Password string
	}
)

func MustLoad(cfgPath string) (*Config, error) {
	godotenv.Load()
	viper.AddConfigPath(cfgPath)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	viper.SetConfigName(os.Getenv("APP_ENV"))

	viper.MergeInConfig()

	var cfg Config
	if err := unmarshalVals(&cfg); err != nil {
		log.Fatalf("error unmarshalling cfg: %v", err)
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func unmarshalVals(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgres", &cfg.PG); err != nil {
		return err
	}

	return nil
}

func setFromEnv(cfg *Config) {
	cfg.Env = os.Getenv("APP_ENV")

	cfg.PG.Host = os.Getenv("PG_HOST")
	cfg.PG.Port = os.Getenv("PG_PORT")
	cfg.PG.Username = os.Getenv("PG_USERNAME")
	cfg.PG.Password = os.Getenv("PG_PASSWORD")
}
