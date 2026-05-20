package command

import (
	"context"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/google/uuid"
)

type DeleteScheduleCommand struct {
	ID     uuid.UUID
	UserID string
}

type DeleteScheduleCommandHandler struct {
	repo port.ScheduledTreatmentRepository
}

func NewDeleteScheduleCommandHandler(repo port.ScheduledTreatmentRepository) *DeleteScheduleCommandHandler {
	return &DeleteScheduleCommandHandler{repo: repo}
}

func (h *DeleteScheduleCommandHandler) Handle(ctx context.Context, cmd DeleteScheduleCommand) error {
	schedule, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return model.NewAppError(model.ErrNotFound, "schedule not found", err)
	}

	if !schedule.IsOwnedBy(cmd.UserID) {
		return model.NewAppError(model.ErrForbidden, "not authorized", nil)
	}

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return model.NewAppError(model.ErrInternal, "failed to delete schedule", err)
	}

	return nil
}
