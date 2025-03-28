package config

import (
	"flag"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
)

type Config struct {
	Debug           bool   `yaml:"debug"            envconfig:"DEBUG"`
	Host            string `yaml:"host"             envconfig:"HOST"`
	Port            string `yaml:"port"             envconfig:"PORT"`
	DatabaseURL     string `yaml:"database_url"     envconfig:"DATABASE_URL"`
	MigrationSource string `yaml:"migration_source" envconfig:"MIGRATION_SOURCE"`
}

// Load merges config from file, env and cli flags
func Load() Config {
	// Default config
	config := &Config{
		Debug:           false,
		Host:            "localhost",
		Port:            "8080",
		DatabaseURL:     "",
		MigrationSource: "file://internal/database/migrations",
	}

	config, err := FromFile("config.yaml", config)
	if err != nil {
		zap.L().Warn("Failed to load config from file", zap.Error(err), zap.String("path", "config.yaml"))
	}

	_, err = FromEnv(config)
	if err != nil {
		zap.L().Warn("Failed to load config from env", zap.Error(err))
	}

	FromFlags(config)

	return *config
}

func FromFile(filePath string, config *Config) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	var configFile Config
	err = yaml.NewDecoder(file).Decode(&configFile)
	if err != nil {
		return nil, err
	}

	return merge(&configFile, config), nil
}

func FromEnv(config *Config) (*Config, error) {
	err := godotenv.Overload()
	if err != nil {
		return nil, err
	}

	envConfig := Config{
		Debug:       os.Getenv("DEBUG") == "true",
		Host:        os.Getenv("HOST"),
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	return merge(&envConfig, config), nil
}

func FromFlags(config *Config) *Config {
	flagConfig := &Config{}

	flag.BoolVar(&flagConfig.Debug, "debug", false, "debug mode")
	flag.StringVar(&flagConfig.Host, "host", "localhost", "host")
	flag.StringVar(&flagConfig.Port, "port", "8080", "port")
	flag.StringVar(&flagConfig.DatabaseURL, "database_url", "", "database url")

	flag.Parse()

	return merge(flagConfig, config)
}

// merge merges two config structs, overriding the base with the override when the field is not the zero value
func merge(base, override *Config) *Config {
	final := *base

	baseVal := reflect.ValueOf(&final).Elem()
	overrideVal := reflect.ValueOf(override).Elem()

	for i := 0; i < baseVal.NumField(); i++ {
		field := baseVal.Field(i)
		overrideField := overrideVal.Field(i)

		zero := reflect.Zero(field.Type()).Interface()
		if !reflect.DeepEqual(overrideField.Interface(), zero) {
			field.Set(overrideField)
		}
	}

	return &final
}
