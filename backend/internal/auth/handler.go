package auth

import (
	"backend/internal"
	errorPkg "backend/internal/error"
	"backend/internal/problem"
	"backend/internal/user"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type JWTIssuer interface {
	New(ctx context.Context, id, username string, role string) (string, error)
}

type UserStore interface {
	Create(ctx context.Context, name, password string) (user.User, error)
	GetByName(ctx context.Context, name string) (user.User, error)
}

type Handler struct {
	validator *validator.Validate
	logger    *zap.Logger
	tracer    trace.Tracer

	userStore UserStore
	jwtIssuer JWTIssuer
}

func NewHandler(validator *validator.Validate, logger *zap.Logger, userStore UserStore, jwtIssuer JWTIssuer) *Handler {
	return &Handler{
		tracer:    otel.Tracer("auth/handler"),
		validator: validator,
		logger:    logger,
		userStore: userStore,
		jwtIssuer: jwtIssuer,
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "LoginEndpoint")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	var request LoginRequest
	err := internal.ParseAndValidateRequestBody(traceCtx, h.validator, r, &request)
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	userEntity, err := h.userStore.GetByName(traceCtx, request.Username)
	if err != nil {
		logger.Warn("Failed to get user by name", zap.String("username", request.Username), zap.Error(err))

		// Prevent leaking information about whether the user exists
		problem.WriteError(traceCtx, w, fmt.Errorf("%w: %v", errorPkg.ErrCredentialInvalid, err), logger)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(request.Password))
	if err != nil {
		problem.WriteError(traceCtx, w, fmt.Errorf("%w: %v", errorPkg.ErrCredentialInvalid, err), logger)
		return
	}

	token, err := h.jwtIssuer.New(traceCtx, userEntity.ID.String(), userEntity.Name, "USER")
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	logger.Debug("User logged in", zap.String("username", request.Username), zap.String("token", token))

	response := LoginResponse{
		Token: token,
	}
	internal.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "LoginEndpoint")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	var request RegisterRequest
	err := internal.ParseAndValidateRequestBody(traceCtx, h.validator, r, &request)
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	userEntity, err := h.userStore.Create(traceCtx, request.Username, string(hashBytes))
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	logger.Info("User registered", zap.String("username", request.Username), zap.String("user_id", userEntity.ID.String()))
	w.WriteHeader(http.StatusCreated)
}
