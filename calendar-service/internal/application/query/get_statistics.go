package query

import (
	"context"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
)

type GetStatisticsQuery struct {
	UserID string
}

type GetStatisticsQueryHandler struct {
	repo port.TreatmentRecordRepository
}

func NewGetStatisticsQueryHandler(repo port.TreatmentRecordRepository) *GetStatisticsQueryHandler {
	return &GetStatisticsQueryHandler{repo: repo}
}

func (h *GetStatisticsQueryHandler) Handle(ctx context.Context, q GetStatisticsQuery) ([]model.TreatmentStat, error) {
	stats, err := h.repo.GetStatsByUserID(ctx, q.UserID)
	if err != nil {
		return nil, model.NewAppError(model.ErrInternal, "failed to get statistics", err)
	}
	return stats, nil
}
