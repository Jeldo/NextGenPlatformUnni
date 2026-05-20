package command

import (
	"context"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/google/uuid"
)

type HandleReservationFixedCommand struct {
	Event model.ReservationFixedEvent
}

type HandleReservationFixedCommandHandler struct {
	repo            port.TreatmentRecordRepository
	scheduleHandler *CalculateScheduleCommandHandler
}

func NewHandleReservationFixedCommandHandler(repo port.TreatmentRecordRepository, scheduleHandler *CalculateScheduleCommandHandler) *HandleReservationFixedCommandHandler {
	return &HandleReservationFixedCommandHandler{repo: repo, scheduleHandler: scheduleHandler}
}

func (h *HandleReservationFixedCommandHandler) Handle(ctx context.Context, cmd HandleReservationFixedCommand) error {
	event := cmd.Event
	if event.TreatmentType == model.TypeSurgery {
		return nil // skip surgery
	}

	// Idempotency check
	existing, _ := h.repo.FindByReservationID(ctx, event.ReservationID)
	if existing != nil {
		return nil
	}

	now := time.Now()
	id, _ := uuid.NewV7()
	treatmentID, _ := uuid.Parse(event.TreatmentID)
	categoryID, _ := uuid.Parse(event.CategoryID)

	record := &model.TreatmentRecord{
		ID:            id,
		UserID:        event.UserID,
		TreatmentID:   treatmentID,
		TreatmentName: event.TreatmentName,
		CategoryID:    categoryID,
		CategoryName:  event.CategoryName,
		TreatmentDate: event.TreatmentDate,
		HospitalName:  &event.HospitalName,
		DosageValue:   event.DosageValue,
		DosageUnit:    event.DosageUnit,
		Source:        model.SourceAuto,
		ReservationID: &event.ReservationID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := record.Validate(); err != nil {
		return err
	}

	if err := h.repo.Save(ctx, record); err != nil {
		return model.NewAppError(model.ErrInternal, "failed to save record", err)
	}

	_, _ = h.scheduleHandler.Handle(ctx, CalculateScheduleCommand{Record: record})
	return nil
}
