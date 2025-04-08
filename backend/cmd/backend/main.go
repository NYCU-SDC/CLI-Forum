package main

import (
	"backend/internal"
	"backend/internal/auth"
	"backend/internal/comment"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/jwt"
	"backend/internal/post"
	"backend/internal/user"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"time"
)

var AppName = "no-app-name"

var Version = "no-version"

var BuildTime = "no-build-time"

var CommitHash = "no-commit-hash"

func main() {
	AppName = os.Getenv("APP_NAME")
	if AppName == "" {
		AppName = "cli-forum-dev"
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
	err := cfg.Validate()
	if err != nil {
		if errors.Is(err, config.ErrDatabaseURLRequired) {
			title := "Database URL is required"
			message := "Please set the DATABASE_URL environment variable or provide a config file with the database_url key."
			message = EarlyApplicationFailed(title, message)
			log.Fatal(message)
		} else {
			log.Fatalf("Failed to validate config: %v, exiting...", err)
		}
	}

	logger, err := initLogger(&cfg, appMetadata)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v, exiting...", err)
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

	dbPool, err := initDatabasePool(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to initialize database pool", zap.Error(err))
	}
	defer dbPool.Close()

	shutdown, err := initOpenTelemetry(AppName, Version, BuildTime, CommitHash, cfg.OtelCollectorUrl)
	if err != nil {
		logger.Fatal("Failed to initialize OpenTelemetry", zap.Error(err))
	}

	validator := internal.NewValidator()

	// initialize service
	jwtService := jwt.NewService(logger, cfg.Secret, 24*time.Hour)
	userService := user.NewService(logger, dbPool)
	commentService := comment.NewService(logger, dbPool)
	postService := post.NewService(logger, dbPool)

	// initialize middleware
	jwtMiddleware := jwt.NewMiddleware(jwtService, logger)

	// initialize handler
	authHandler := auth.NewHandler(validator, logger, userService, jwtService)
	userHandler := user.NewHandler(validator, logger, userService)
	commentHandler := comment.NewHandler(validator, logger, commentService)
	postHandler := post.NewHandler(validator, logger, postService)

	// initialize mux
	mux := http.NewServeMux()

	// set up routes
	mux.HandleFunc("POST /api/login", basicMiddleware(authHandler.LoginHandler, logger, cfg.Debug))
	// This handler duplicates the above handler intentionally for teaching clarity.
	// mux.HandleFunc("POST /api/login", internal.TraceMiddleware(internal.RecoverMiddleware(authHandler.LoginHandler, logger), logger))

	mux.HandleFunc("POST /api/register", basicMiddleware(authHandler.RegisterHandler, logger, cfg.Debug))

	mux.HandleFunc("POST /api/user", requireUserRoleMiddleware(userHandler.CreateHandler, jwtMiddleware, logger, cfg.Debug))

	mux.HandleFunc("GET /api/comments", requireUserRoleMiddleware(commentHandler.GetAllHandler, jwtMiddleware, logger, cfg.Debug))
	mux.HandleFunc("GET /api/post/{post_id}/comments", requireUserRoleMiddleware(commentHandler.GetByPostHandler, jwtMiddleware, logger, cfg.Debug))
	mux.HandleFunc("POST /api/post/{post_id}/comments", requireUserRoleMiddleware(commentHandler.CreateHandler, jwtMiddleware, logger, cfg.Debug))
	mux.HandleFunc("GET /api/comment/{id}", requireUserRoleMiddleware(commentHandler.GetByIdHandler, jwtMiddleware, logger, cfg.Debug))

  mux.HandleFunc("GET /api/posts", requireUserRoleMiddleware(postHandler.GetAllHandler, jwtMiddleware, logger, cfg.Debug))
	mux.HandleFunc("POST /api/posts", requireUserRoleMiddleware(postHandler.CreateHandler, jwtMiddleware, logger, cfg.Debug))
	mux.HandleFunc("GET /api/post/{id}", requireUserRoleMiddleware(postHandler.GetHandler, jwtMiddleware, logger, cfg.Debug))

	logger.Info("Starting listening request", zap.String("host", cfg.Host), zap.String("port", cfg.Port))
	err = http.ListenAndServe(cfg.Host+":"+cfg.Port, mux)
	if err != nil {
		logger.Fatal("Fail to start server with error", zap.Error(err))
	}

	// graceful shutdown
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			logger.Error("Failed to shutdown OpenTelemetry", zap.Error(err))
		}
	}()
}

func basicMiddleware(next http.HandlerFunc, logger *zap.Logger, debug bool) http.HandlerFunc {
	return internal.TraceMiddleware(internal.RecoverMiddleware(next, logger, debug), logger)
}

func requireUserRoleMiddleware(next http.HandlerFunc, jwtMiddleware jwt.Middleware, logger *zap.Logger, debug bool) http.HandlerFunc {
	return internal.TraceMiddleware(internal.RecoverMiddleware(jwtMiddleware.HandlerFunc(auth.Middleware(next, logger, "USER")), logger, debug), logger)
}

// initLogger create a new logger. If debug is enabled, it will create a development logger without metadata for better
// readability, otherwise it will create a production logger with metadata and json format.
func initLogger(cfg *config.Config, appMetadata []zap.Field) (*zap.Logger, error) {
	var err error
	var logger *zap.Logger
	if cfg.Debug {
		logger, err = internal.ZapDevelopmentConfig().Build()
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
func initDatabasePool(databaseURL string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// initOpenTelemetry initializes OpenTelemetry with the given app name, version, build time and commit hash. If the
// collector URL is not empty, it will create a gRPC connection to the collector and set up the trace exporter.
func initOpenTelemetry(appName, version, buildTime, commitHash, collectorUrl string) (func(context.Context) error, error) {
	ctx := context.Background()

	serviceName := semconv.ServiceNameKey.String(appName)
	serviceVersion := semconv.ServiceVersionKey.String(version)
	serviceBuildTime := semconv.DeploymentEnvironmentKey.String(buildTime)
	serviceCommitHash := semconv.DeploymentEnvironmentKey.String(commitHash)

	res, err := resource.New(ctx,
		resource.WithAttributes(
			serviceName,
			serviceVersion,
			serviceBuildTime,
			serviceCommitHash,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	options := []trace.TracerProviderOption{
		sdktrace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()),
	}

	if collectorUrl != "" {
		conn, err := initGrpcConn(collectorUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
		}

		traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, fmt.Errorf("failed to create trace exporter: %w", err)
		}

		// Use batch span processor to optimize span export by batching multiple spans together
		bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
		options = append(options, sdktrace.WithSpanProcessor(bsp))
	}

	tracerProvider := trace.NewTracerProvider(options...)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tracerProvider.Shutdown, nil
}

// initGrpcConn simply creates a gRPC connection to the given target using insecure credentials.
func initGrpcConn(target string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}

func EarlyApplicationFailed(title, action string) string {
	result := `
-----------------------------------------
Application Failed to Start
-----------------------------------------

# What's wrong?
%s

# How to fix it?
%s

`

	result = fmt.Sprintf(result, title, action)
	return result
}
