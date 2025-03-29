package config

import (
	"errors"
	"flag"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
)

const DefaultSecret = "default-secret"

type Config struct {
	Debug           bool   `yaml:"debug"            envconfig:"DEBUG"`
	Host            string `yaml:"host"             envconfig:"HOST"`
	Port            string `yaml:"port"             envconfig:"PORT"`
	Secret          string `yaml:"secret"           envconfig:"SECRET"`
	DatabaseURL     string `yaml:"database_url"     envconfig:"DATABASE_URL"`
	MigrationSource string `yaml:"migration_source" envconfig:"MIGRATION_SOURCE"`
}

func Load() Config {
	config := &Config{
		Debug:           false,
		Host:            "localhost",
		Port:            "8080",
		Secret:          DefaultSecret,
		DatabaseURL:     "",
		MigrationSource: "file://internal/database/migrations",
	}

	var err error

	config, err = FromFile("config.yaml", config)
	if err != nil {
		zap.L().Warn("Failed to load config from file", zap.Error(err), zap.String("path", "config.yaml"))
	}

	config, err = FromEnv(config)
	if err != nil {
		zap.L().Warn("Failed to load config from env", zap.Error(err))
	}

	config, err = FromFlags(config)
	if err != nil {
		zap.L().Warn("Failed to load config from flags", zap.Error(err))
	}

	return *config
}

func FromFile(filePath string, config *Config) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	fileConfig := Config{}
	if err := yaml.NewDecoder(file).Decode(&fileConfig); err != nil {
		return config, err
	}

	return merge(config, &fileConfig)
}

func FromEnv(config *Config) (*Config, error) {
	if err := godotenv.Overload(); err != nil {
		return config, err
	}

	envConfig := &Config{
		Debug:           os.Getenv("DEBUG") == "true",
		Host:            os.Getenv("HOST"),
		Port:            os.Getenv("PORT"),
		Secret:          os.Getenv("SECRET"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		MigrationSource: os.Getenv("MIGRATION_SOURCE"),
	}

	return merge(config, envConfig)
}

func FromFlags(config *Config) (*Config, error) {
	flagConfig := &Config{}

	flag.BoolVar(&flagConfig.Debug, "debug", false, "debug mode")
	flag.StringVar(&flagConfig.Host, "host", "", "host")
	flag.StringVar(&flagConfig.Port, "port", "", "port")
	flag.StringVar(&flagConfig.Secret, "secret", "", "secret")
	flag.StringVar(&flagConfig.DatabaseURL, "database_url", "", "database url")
	flag.StringVar(&flagConfig.MigrationSource, "migration_source", "", "migration source")

	flag.Parse()

	return merge(config, flagConfig)
}

func merge(base, override *Config) (*Config, error) {
	if base == nil {
		return nil, errors.New("base config cannot be nil")
	}
	if override == nil {
		return base, nil
	}

	final := *base
	baseVal := reflect.ValueOf(&final).Elem()
	overrideVal := reflect.ValueOf(override).Elem()

	if baseVal.Type() != overrideVal.Type() {
		return nil, errors.New("config types do not match")
	}

	for i := 0; i < baseVal.NumField(); i++ {
		field := baseVal.Field(i)
		overrideField := overrideVal.Field(i)
		zero := reflect.Zero(field.Type()).Interface()

		if field.CanSet() && !reflect.DeepEqual(overrideField.Interface(), zero) {
			field.Set(overrideField)
		}
	}

	return &final, nil
}
