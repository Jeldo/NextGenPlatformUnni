package query

import (
	"context"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/google/uuid"
)

type GetScheduleQuery struct {
	ID     uuid.UUID
	UserID string
}

type GetScheduleQueryHandler struct {
	repo port.ScheduledTreatmentRepository
}

func NewGetScheduleQueryHandler(repo port.ScheduledTreatmentRepository) *GetScheduleQueryHandler {
	return &GetScheduleQueryHandler{repo: repo}
}

func (h *GetScheduleQueryHandler) Handle(ctx context.Context, q GetScheduleQuery) (*model.ScheduledTreatment, error) {
	schedule, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		return nil, model.NewAppError(model.ErrNotFound, "schedule not found", err)
	}
	if !schedule.IsOwnedBy(q.UserID) {
		return nil, model.NewAppError(model.ErrForbidden, "not authorized", nil)
	}
	return schedule, nil
}
