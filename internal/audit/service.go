package audit

import (
	"context"
	"errors"

	"github.com/edustack/go-boilerplate/pkg"
	"gorm.io/gorm"
)

// Service mendefinisikan use-case untuk audit Logs.
type Service interface {
	GetLogs(ctx context.Context, params QueryParams) ([]Logs, int64, error)
	GetLogByID(ctx context.Context, id int64) (*Logs, error)
}

type service struct {
	repo *Repository
}

// NewService membuat instance audit service.
func NewService(repo *Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetLogs(ctx context.Context, params QueryParams) ([]Logs, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}
	return s.repo.FindAll(ctx, params)
}

func (s *service) GetLogByID(ctx context.Context, id int64) (*Logs, error) {
	logEntry, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrLogNotFound
		}
		return nil, pkg.ErrInternal
	}
	return logEntry, nil
}
