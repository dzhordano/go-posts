package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	EnvLocal = "local"
	Prod     = "prod"
)

type (
	Config struct {
		Env  string `yaml:"env" env-default:"local"`
		AUTH AUTHConfig
		HTTP HTTPConfig
		PG   PGConfig
		SMTP SMTPConfig
	}

	AUTHConfig struct {
		JWT                    JWTConfig
		PasswordSalt           string
		VerificationCodeLength int `mapstructure:"verificationCodeLength"`
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
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
		Username string
		Password string
		SSLMode  string
	}

	SMTPConfig struct {
		Port int    `mapstructure:"port"`
		Host string `mapstructure:"host"`
		From string `mapstructure:"from"`
		Pass string
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

	if err := viper.UnmarshalKey("auth", &cfg.AUTH.JWT); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.AUTH); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("smtp", &cfg.SMTP); err != nil {
		return err
	}

	return nil
}

func setFromEnv(cfg *Config) {
	cfg.Env = os.Getenv("APP_ENV")

	cfg.AUTH.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.AUTH.JWT.SigningKey = os.Getenv("SIGNING_KEY")

	cfg.PG.Host = os.Getenv("PG_HOST")
	cfg.PG.Port = os.Getenv("PG_PORT")
	cfg.PG.Username = os.Getenv("PG_USERNAME")
	cfg.PG.Password = os.Getenv("PG_PASSWORD")

	cfg.SMTP.Pass = os.Getenv("SMTP_PASSWORD")
}
