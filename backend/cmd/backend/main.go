package main

import (
	"backend/internal"
	"backend/internal/config"
	"backend/internal/database"
	"go.uber.org/zap"
)

func main() {
	logger, err := internal.ZapProductionConfig().Build()
	if err != nil {
		zap.S().Fatalw("Failed to initialize logger", zap.Error(err))
	}

	defer func() {
		err := logger.Sync()
		if err != nil {
			zap.S().Errorw("Failed to sync logger", zap.Error(err))
		}
	}()

	cfg := config.Load()
	if cfg.Debug {
		logger, err = zap.NewDevelopment()
		if err != nil {
			logger.Fatal("Failed to initialize logger with development config", zap.Error(err))
		}
	}

	logger.Info("Running mode", zap.Bool("debug", cfg.Debug), zap.String("host", cfg.Host), zap.String("port", cfg.Port))

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
