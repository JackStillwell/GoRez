package base

import (
	"fmt"

	"go.uber.org/zap"
)

type Service interface {
	GetLogger() *zap.Logger
	ShutdownLogger() error
}

type service struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) Service {
	if logger == nil {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			panic(fmt.Errorf("initing logger: %w", err))
		}
	}
	return &service{
		logger: logger,
	}
}

func (s *service) GetLogger() *zap.Logger {
	return s.logger
}

func (s *service) ShutdownLogger() error {
	return s.logger.Sync()
}
