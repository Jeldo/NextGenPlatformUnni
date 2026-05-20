package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/application/command"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/google/uuid"
)

func TestDeleteRecordCommandHandler_Success(t *testing.T) {
	t.Parallel()

	recordRepo := port.NewMockRecordRepo()
	scheduleRepo := port.NewMockScheduleRepo()

	recordID := uuid.New()
	recordRepo.Records[recordID] = &model.TreatmentRecord{
		ID:            recordID,
		UserID:        "user-1",
		TreatmentID:   uuid.New(),
		TreatmentName: "보톡스",
		CategoryID:    uuid.New(),
		CategoryName:  "주름개선",
		TreatmentDate: time.Now(),
		Source:        model.SourceManual,
	}

	handler := command.NewDeleteRecordCommandHandler(recordRepo, scheduleRepo)

	err := handler.Handle(context.Background(), command.DeleteRecordCommand{
		ID:     recordID,
		UserID: "user-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := recordRepo.Records[recordID]; exists {
		t.Error("record should be deleted")
	}
}

func TestDeleteRecordCommandHandler_NotFound(t *testing.T) {
	t.Parallel()

	recordRepo := port.NewMockRecordRepo()
	scheduleRepo := port.NewMockScheduleRepo()
	handler := command.NewDeleteRecordCommandHandler(recordRepo, scheduleRepo)

	err := handler.Handle(context.Background(), command.DeleteRecordCommand{
		ID:     uuid.New(),
		UserID: "user-1",
	})
	if err == nil {
		t.Fatal("expected error for non-existent record")
	}
	var appErr *model.AppError
	if errors.As(err, &appErr) {
		if appErr.Code != model.ErrNotFound {
			t.Errorf("error code = %v, want NOT_FOUND", appErr.Code)
		}
	}
}

func TestDeleteRecordCommandHandler_Unauthorized(t *testing.T) {
	t.Parallel()

	recordRepo := port.NewMockRecordRepo()
	scheduleRepo := port.NewMockScheduleRepo()

	recordID := uuid.New()
	recordRepo.Records[recordID] = &model.TreatmentRecord{
		ID:            recordID,
		UserID:        "user-1",
		TreatmentID:   uuid.New(),
		TreatmentName: "보톡스",
		CategoryID:    uuid.New(),
		CategoryName:  "주름개선",
		TreatmentDate: time.Now(),
		Source:        model.SourceManual,
	}

	handler := command.NewDeleteRecordCommandHandler(recordRepo, scheduleRepo)

	err := handler.Handle(context.Background(), command.DeleteRecordCommand{
		ID:     recordID,
		UserID: "user-other", // different user
	})
	if err == nil {
		t.Fatal("expected forbidden error")
	}
	var appErr *model.AppError
	if errors.As(err, &appErr) {
		if appErr.Code != model.ErrForbidden {
			t.Errorf("error code = %v, want FORBIDDEN", appErr.Code)
		}
	}
	// Record should still exist
	if _, exists := recordRepo.Records[recordID]; !exists {
		t.Error("record should not be deleted when unauthorized")
	}
}
