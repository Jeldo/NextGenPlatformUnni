package command

import (
	"context"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/google/uuid"
)

type DeleteRecordCommand struct {
	ID     uuid.UUID
	UserID string
}

type DeleteRecordCommandHandler struct {
	repo         port.TreatmentRecordRepository
	scheduleRepo port.ScheduledTreatmentRepository
}

func NewDeleteRecordCommandHandler(repo port.TreatmentRecordRepository, scheduleRepo port.ScheduledTreatmentRepository) *DeleteRecordCommandHandler {
	return &DeleteRecordCommandHandler{repo: repo, scheduleRepo: scheduleRepo}
}

func (h *DeleteRecordCommandHandler) Handle(ctx context.Context, cmd DeleteRecordCommand) error {
	record, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return model.NewAppError(model.ErrNotFound, "record not found", err)
	}

	if !record.IsOwnedBy(cmd.UserID) {
		return model.NewAppError(model.ErrForbidden, "not authorized", nil)
	}

	_ = h.scheduleRepo.DeleteByRecordID(ctx, cmd.ID)

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return model.NewAppError(model.ErrInternal, "failed to delete record", err)
	}

	return nil
}
