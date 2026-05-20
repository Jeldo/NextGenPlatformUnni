package command

import (
	"context"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/google/uuid"
)

type CompleteScheduleCommand struct {
	ID     uuid.UUID
	UserID string
}

type CompleteScheduleCommandHandler struct {
	repo port.ScheduledTreatmentRepository
}

func NewCompleteScheduleCommandHandler(repo port.ScheduledTreatmentRepository) *CompleteScheduleCommandHandler {
	return &CompleteScheduleCommandHandler{repo: repo}
}

func (h *CompleteScheduleCommandHandler) Handle(ctx context.Context, cmd CompleteScheduleCommand) error {
	schedule, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return model.NewAppError(model.ErrNotFound, "schedule not found", err)
	}

	if !schedule.IsOwnedBy(cmd.UserID) {
		return model.NewAppError(model.ErrForbidden, "not authorized", nil)
	}

	schedule.MarkCompleted(time.Now())

	if err := h.repo.Update(ctx, schedule); err != nil {
		return model.NewAppError(model.ErrInternal, "failed to update schedule", err)
	}

	return nil
}
