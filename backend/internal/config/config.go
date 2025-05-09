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

var ErrDatabaseURLRequired = errors.New("database_url is required")

type Config struct {
	Debug            bool   `yaml:"debug"              envconfig:"DEBUG"`
	Host             string `yaml:"host"               envconfig:"HOST"`
	Port             string `yaml:"port"               envconfig:"PORT"`
	Secret           string `yaml:"secret"             envconfig:"SECRET"`
	DatabaseURL      string `yaml:"database_url"       envconfig:"DATABASE_URL"`
	MigrationSource  string `yaml:"migration_source"   envconfig:"MIGRATION_SOURCE"`
	OtelCollectorUrl string `yaml:"otel_collector_url" envconfig:"OTEL_COLLECTOR_URL"`
}

func (c Config) Validate() error {
	if c.DatabaseURL == "" {
		return ErrDatabaseURLRequired
	}

	return nil
}

type LogBuffer struct {
	buffer []logEntry
}

type logEntry struct {
	msg  string
	err  error
	meta map[string]string
}

func NewConfigLogger() *LogBuffer {
	return &LogBuffer{}
}

func (cl *LogBuffer) Warn(msg string, err error, meta map[string]string) {
	cl.buffer = append(cl.buffer, logEntry{msg: msg, err: err, meta: meta})
}

func (cl *LogBuffer) FlushToZap(logger *zap.Logger) {
	for _, e := range cl.buffer {
		var fields []zap.Field
		if e.err != nil {
			fields = append(fields, zap.Error(e.err))
		}
		for k, v := range e.meta {
			fields = append(fields, zap.String(k, v))
		}
		logger.Warn(e.msg, fields...)
	}
	cl.buffer = nil
}

func Load() (Config, *LogBuffer) {
	logger := NewConfigLogger()

	config := &Config{
		Debug:            false,
		Host:             "localhost",
		Port:             "8080",
		Secret:           DefaultSecret,
		DatabaseURL:      "",
		MigrationSource:  "file://internal/database/migrations",
		OtelCollectorUrl: "",
	}

	var err error

	config, err = FromFile("config.yaml", config)
	if err != nil {
		logger.Warn("Failed to load config from file", err, map[string]string{"path": "config.yaml"})
	}

	config, err = FromEnv(config, logger)
	if err != nil {
		zap.L().Warn("Failed to load config from env", zap.Error(err))
	}

	config, err = FromFlags(config)
	if err != nil {
		zap.L().Warn("Failed to load config from flags", zap.Error(err))
	}

	return *config, logger
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

func FromEnv(config *Config, logger *LogBuffer) (*Config, error) {
	if err := godotenv.Overload(); err != nil {
		if os.IsNotExist(err) {
			logger.Warn("No .env file found", err, map[string]string{"path": ".env"})
		} else {
			return nil, err
		}
	}

	envConfig := &Config{
		Debug:            os.Getenv("DEBUG") == "true",
		Host:             os.Getenv("HOST"),
		Port:             os.Getenv("PORT"),
		Secret:           os.Getenv("SECRET"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		MigrationSource:  os.Getenv("MIGRATION_SOURCE"),
		OtelCollectorUrl: os.Getenv("OTEL_COLLECTOR_URL"),
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
	flag.StringVar(&flagConfig.OtelCollectorUrl, "otel_collector_url", "", "OpenTelemetry collector URL")

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
