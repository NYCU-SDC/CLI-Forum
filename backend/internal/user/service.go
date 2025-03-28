package user

import (
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	query  *Queries
}

func NewService(logger *zap.Logger, db DBTX) *Service {
	return &Service{
		logger: logger,
		query:  New(db),
	}
}
