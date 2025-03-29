package main

import (
	"backend/internal"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/jwt"
	"backend/internal/user"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

var AppName = "no-app-name"

var Version = "no-version"

var BuildTime = "no-build-time"

var CommitHash = "no-commit-hash"

func main() {
	if AppName == "no-app-name" {
		AppName = "cli-forum-dev-" + uuid.New().String()
	}

	if BuildTime == "no-build-time" {
		now := time.Now()
		BuildTime = "not provided (now: " + now.Format(time.RFC3339) + ")"
	}

	appMetadata := []zap.Field{
		zap.String("app_name", AppName),
		zap.String("version", Version),
		zap.String("build_time", BuildTime),
		zap.String("commit_hash", CommitHash),
	}

	cfg := config.Load()

	logger, err := initLogger(&cfg, appMetadata)
	if err != nil {
		zap.L().Warn("Critical error occurred, exiting...", appMetadata...)
		zap.L().Fatal("Failed to initialize logger", zap.Error(err))
	}

	if cfg.Secret == config.DefaultSecret && !cfg.Debug {
		logger.Warn("Default secret detected in production environment, replace it with a secure random string")
		cfg.Secret = uuid.New().String()
	}

	logger.Info("Application initialization", zap.Bool("debug", cfg.Debug), zap.String("host", cfg.Host), zap.String("port", cfg.Port))

	logger.Info("Starting database migration...")
	err = database.MigrationUp(cfg.MigrationSource, cfg.DatabaseURL, logger)
	if err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	dbPool, err := initDatabasePool(&cfg)
	if err != nil {
		logger.Fatal("Failed to initialize database pool", zap.Error(err))
	}
	defer dbPool.Close()

	// initialize service
	_ = jwt.NewService(logger, cfg.Secret, 24*time.Hour)
	userService := user.NewService(logger, dbPool)

	// initialize handler
	userHandler := user.NewHandler(logger, userService)

	// initialize auth service
	//authService := auth.NewService(logger, dbpool, jwtService)
	//authHandler := auth.NewHandler(authService)
	//
	// initialize mux
	mux := http.NewServeMux()

	// set up routes
	mux.HandleFunc("POST /api/user", userHandler.CreateHandler)

	//mux.HandleFunc("POST /login", authHandler.LoginHandler)
	//mux.HandleFunc("POST /register", authHandler.RegisterHandler)
	//
	//// start server on port 8090
	//logger.Info("Server starting on localhost:8090")
	//err = http.ListenAndServe("localhost:8090", mux)
	//if err != nil {
	//	logger.Fatal("Fail to start server with error : ", zap.Error(err))
	//}
}

// initLogger create a new logger. If debug is enabled, it will create a development logger without metadata for better
// readability, otherwise it will create a production logger with metadata and json format.
func initLogger(cfg *config.Config, appMetadata []zap.Field) (*zap.Logger, error) {
	var err error
	var logger *zap.Logger
	if cfg.Debug {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}

		logger.Info("Running in debug mode", appMetadata...)
	} else {
		logger, err = internal.ZapProductionConfig().Build()
		if err != nil {
			return nil, err
		}

		logger = logger.With(appMetadata...)
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			zap.S().Errorw("Failed to sync logger", zap.Error(err))
		}
	}()

	return logger, nil
}

// initDatabasePool creates a new pgxpool.Pool with the given database URL in the config, it uses the default config
// provided by pgxpool.ParseConfig:
//
//   - pool_max_conns: 4
//   - pool_min_conns: 0
//   - pool_max_conn_lifetime: 1 hour
//   - pool_max_conn_idle_time: 30 minutes
//   - pool_health_check_period: 1 minute
//   - pool_max_conn_lifetime_jitter: 0
func initDatabasePool(cfg *config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
