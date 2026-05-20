package command

import (
	"context"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/google/uuid"
)

type UpdateRecordCommand struct {
	ID            uuid.UUID
	UserID        string
	TreatmentID   uuid.UUID
	TreatmentName string
	CategoryID    uuid.UUID
	CategoryName  string
	TreatmentDate time.Time
	HospitalName  *string
	DosageValue   *float64
	DosageUnit    *model.DosageUnit
}

type UpdateRecordCommandHandler struct {
	repo            port.TreatmentRecordRepository
	scheduleRepo    port.ScheduledTreatmentRepository
	scheduleHandler *CalculateScheduleCommandHandler
}

func NewUpdateRecordCommandHandler(repo port.TreatmentRecordRepository, scheduleRepo port.ScheduledTreatmentRepository, scheduleHandler *CalculateScheduleCommandHandler) *UpdateRecordCommandHandler {
	return &UpdateRecordCommandHandler{repo: repo, scheduleRepo: scheduleRepo, scheduleHandler: scheduleHandler}
}

func (h *UpdateRecordCommandHandler) Handle(ctx context.Context, cmd UpdateRecordCommand) (*model.TreatmentRecord, error) {
	record, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, model.NewAppError(model.ErrNotFound, "record not found", err)
	}

	if !record.IsOwnedBy(cmd.UserID) {
		return nil, model.NewAppError(model.ErrForbidden, "not authorized", nil)
	}

	needsReschedule := record.IsDateOrCategoryChanged(cmd.TreatmentID, cmd.TreatmentDate)

	record.ApplyUpdate(cmd.TreatmentID, cmd.TreatmentName, cmd.CategoryID, cmd.CategoryName, cmd.TreatmentDate, cmd.HospitalName, cmd.DosageValue, cmd.DosageUnit)
	record.UpdatedAt = time.Now()

	if err := record.Validate(); err != nil {
		return nil, err
	}

	if err := h.repo.Update(ctx, record); err != nil {
		return nil, model.NewAppError(model.ErrInternal, "failed to update record", err)
	}

	if needsReschedule {
		_ = h.scheduleRepo.DeleteByRecordID(ctx, record.ID)
		_, _ = h.scheduleHandler.Handle(ctx, CalculateScheduleCommand{Record: record})
	}

	return record, nil
}
