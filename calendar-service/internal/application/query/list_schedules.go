package query

import (
	"context"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
)

type ListSchedulesQuery struct {
	UserID string
	From   time.Time
	To     time.Time
}

type ListSchedulesQueryHandler struct {
	repo port.ScheduledTreatmentRepository
}

func NewListSchedulesQueryHandler(repo port.ScheduledTreatmentRepository) *ListSchedulesQueryHandler {
	return &ListSchedulesQueryHandler{repo: repo}
}

func (h *ListSchedulesQueryHandler) Handle(ctx context.Context, q ListSchedulesQuery) ([]*model.ScheduledTreatment, error) {
	schedules, err := h.repo.FindByUserIDAndDateRange(ctx, q.UserID, q.From, q.To)
	if err != nil {
		return nil, model.NewAppError(model.ErrInternal, "failed to list schedules", err)
	}
	return schedules, nil
}
