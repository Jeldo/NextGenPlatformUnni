package command

import (
	"context"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/google/uuid"
)

type CreateRecordCommand struct {
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

type CreateRecordCommandHandler struct {
	repo            port.TreatmentRecordRepository
	scheduleHandler *CalculateScheduleCommandHandler
}

func NewCreateRecordCommandHandler(repo port.TreatmentRecordRepository, scheduleHandler *CalculateScheduleCommandHandler) *CreateRecordCommandHandler {
	return &CreateRecordCommandHandler{repo: repo, scheduleHandler: scheduleHandler}
}

func (h *CreateRecordCommandHandler) Handle(ctx context.Context, cmd CreateRecordCommand) (*model.TreatmentRecord, error) {
	now := time.Now()
	id, _ := uuid.NewV7()

	record := &model.TreatmentRecord{
		ID:            id,
		UserID:        cmd.UserID,
		TreatmentID:   cmd.TreatmentID,
		TreatmentName: cmd.TreatmentName,
		CategoryID:    cmd.CategoryID,
		CategoryName:  cmd.CategoryName,
		TreatmentDate: cmd.TreatmentDate,
		HospitalName:  cmd.HospitalName,
		DosageValue:   cmd.DosageValue,
		DosageUnit:    cmd.DosageUnit,
		Source:        model.SourceManual,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := record.Validate(); err != nil {
		return nil, err
	}

	if err := h.repo.Save(ctx, record); err != nil {
		return nil, model.NewAppError(model.ErrInternal, "failed to save record", err)
	}

	// Calculate schedule asynchronously (best-effort)
	_, _ = h.scheduleHandler.Handle(ctx, CalculateScheduleCommand{Record: record})

	return record, nil
}
