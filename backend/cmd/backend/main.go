package main

import (
	"backend/internal"
	"backend/internal/config"
	"backend/internal/database"
	"go.uber.org/zap"
	"time"
)

var Version = "no-version"

var BuildTime = "no-build-time"

var CommitHash = "no-commit-hash"

func main() {
	if BuildTime == "no-build-time" {
		now := time.Now()
		BuildTime = "not provided (now: " + now.Format(time.RFC3339) + ")"
	}

	appMetadata := []zap.Field{
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

	logger.Info("Application initialization", zap.Bool("debug", cfg.Debug), zap.String("host", cfg.Host), zap.String("port", cfg.Port))

	logger.Info("Starting database migration...")
	err = database.MigrationUp(cfg.MigrationSource, cfg.DatabaseURL, logger)
	if err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	//// initialize database
	//dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	//if err != nil {
	//	panic("Unable to create connection pool: " + err.Error())
	//}
	//defer dbpool.Close()
	//
	//// initialize jwt service
	//jwtService := jwt.NewService(logger, []byte(os.Getenv("BACKEND_SECRET_KEY")), 15)
	//// initialize auth service
	//authService := auth.NewService(logger, dbpool, jwtService)
	//authHandler := auth.NewHandler(authService)
	//
	//// initialize mux
	//mux := http.NewServeMux()
	//
	//// set up routes
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

// initLogger create a new logger. If debug is enabled, it will create a development logger without metadata for better readability,
// otherwise it will create a production logger with metadata and json format.
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
