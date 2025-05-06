package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/spf13/viper"
)

type Config struct {
	PrimaryDBHost       string `mapstructure:"PRIMARY_DB_HOST"`
	PrimaryDBPort       string `mapstructure:"PRIMARY_DB_PORT"`
	PrimaryDBUser       string `mapstructure:"PRIMARY_DB_USER"`
	PrimaryDBPassword   string `mapstructure:"PRIMARY_DB_PASSWORD"`
	PrimaryDBName       string `mapstructure:"PRIMARY_DB_NAME"`
	PrimaryDBSslMode    string `mapstructure:"PRIMARY_DB_SSL_MODE"`
	SecondaryDBHost     string `mapstructure:"SECONDARY_DB_HOST"`
	SecondaryDBPort     string `mapstructure:"SECONDARY_DB_PORT"`
	SecondaryDBUser     string `mapstructure:"SECONDARY_DB_USER"`
	SecondaryDBPassword string `mapstructure:"SECONDARY_DB_PASSWORD"`
	SecondaryDBName     string `mapstructure:"SECONDARY_DB_NAME"`
	SecondaryDBSslMode  string `mapstructure:"SECONDARY_DB_SSL_MODE"`
	SSOPort             string `mapstructure:"SSO_PORT"`
	RedisURI            string `mapstructure:"REDIS_URI"`
	RootURL             string `mapstructure:"ROOT_URL"`
	ClientID            string `mapstructure:"CLIENT_ID"`
	CallbackURL         string `mapstructure:"CALLBACK_URL"`
	SignOutURL          string `mapstructure:"SIGN_OUT_URL"`
	PayloadURL          string `mapstructure:"PAYLOAD_URL"`
	PublicKeyAPI        string `mapstructure:"PUBLIC_KEY_API"`
	SSOUsername         string `mapstructure:"SSO_USERNAME"`
	SSOPassword         string `mapstructure:"SSO_PASSWORD"`
	SessionKey          string `mapstructure:"SESSION_KEY"`
}

func LoadConfig() (Config, error) {
	v := viper.New()
	env := os.Getenv("APP_ENV")
	envsWithEnvVars := []string{"preview", "staging", "prod"}
	if slices.Contains(envsWithEnvVars, env) {
		// Read in environment variables that match
		v.BindEnv("SSO_URL")
		v.BindEnv("CLIENT_ID")
		v.BindEnv("CALLBACK_URL")
		v.BindEnv("SIGN_OUT_URL")
		v.BindEnv("PAYLOAD_URL")
		v.BindEnv("PUBLIC_KEY_API")
		v.BindEnv("SSO_USERNAME")
		v.BindEnv("SSO_PASSWORD")
		v.BindEnv("REDIS_URI")
		v.BindEnv("ROOT_URL")
		v.BindEnv("SERVER_PORT")
		v.BindEnv("MODE")
		v.BindEnv("DB_HOST")
		v.BindEnv("DB_PORT")
		v.BindEnv("DB_USER")
		v.BindEnv("DB_PASSWORD")
		v.BindEnv("DB_NAME")
		v.BindEnv("DB_SSL_MODE")
		v.BindEnv("SESSION_KEY")
	} else {
		v.SetDefault("SERVER_PORT", "8080")
		v.SetDefault("MODE", "debug")

		v.SetConfigName(configFile())
		v.SetConfigType("env")
		v.AddConfigPath(configPath())

		// Read the config file
		if err := v.ReadInConfig(); err != nil {
			return Config{}, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func configFile() string {
	env := os.Getenv("APP_ENV")
	if env == "test" {
		return "config.test"
	}
	return "config.dev"
}

func configPath() string {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "..")
	fmt.Println(projectRoot)
	return filepath.Join(projectRoot, "config")
}
