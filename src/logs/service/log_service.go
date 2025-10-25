package service

import (
	"context"

	"gobackend/shared/pagination"
	"gobackend/src/logs/dao"
	"gobackend/src/logs/dto"
	loginterfaces "gobackend/src/logs/interfaces"
)

var _ loginterfaces.Service = (*LogService)(nil)

// LogService provides read operations for user logs.
type LogService struct {
	repo loginterfaces.Repository
}

// NewLogService constructs a new LogService.
func NewLogService(repo loginterfaces.Repository) *LogService {
	return &LogService{repo: repo}
}

// ListLogs fetches logs and maps them to DTOs.
func (s *LogService) ListLogs(ctx context.Context, params pagination.Params) ([]dto.Log, int64, error) {
	logs, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.Log, 0, len(logs))
	for _, entry := range logs {
		result = append(result, dto.Log{
			UserName:  entry.UserName,
			Action:    entry.Action,
			Detail:    entry.Detail,
			CreatedAt: entry.CreatedAt,
		})
	}

	return result, total, nil
}

// Record stores a new log entry.
func (s *LogService) Record(ctx context.Context, entry dto.NewLog) error {
	daoEntry := dao.Log{
		UserID: entry.UserID,
		Action: entry.Action,
		Detail: entry.Detail,
	}

	return s.repo.Create(ctx, daoEntry)
}
