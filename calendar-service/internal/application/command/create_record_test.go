package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/application/command"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/service"
	"github.com/google/uuid"
)

func setupCreateHandler(recordRepo *port.MockRecordRepo, cycleClient *port.MockCycleRuleClient) *command.CreateRecordCommandHandler {
	scheduleRepo := port.NewMockScheduleRepo()
	calc := service.NewCycleCalculator()
	scheduleHandler := command.NewCalculateScheduleCommandHandler(scheduleRepo, cycleClient, calc)
	return command.NewCreateRecordCommandHandler(recordRepo, scheduleHandler)
}

func TestCreateRecordCommandHandler_Success(t *testing.T) {
	t.Parallel()

	recordRepo := port.NewMockRecordRepo()
	cycleClient := port.NewMockCycleRuleClient()
	treatmentID := uuid.New()
	cycleClient.Rules[treatmentID] = &port.CycleRule{TreatmentID: treatmentID, CycleDays: 90}

	handler := setupCreateHandler(recordRepo, cycleClient)

	cmd := command.CreateRecordCommand{
		UserID:        "user-1",
		TreatmentID:   treatmentID,
		TreatmentName: "보톡스",
		CategoryID:    uuid.New(),
		CategoryName:  "주름개선",
		TreatmentDate: time.Now(),
	}

	record, err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record == nil {
		t.Fatal("expected record, got nil")
	}
	if record.UserID != "user-1" {
		t.Errorf("UserID = %v, want user-1", record.UserID)
	}
	if record.Source != model.SourceManual {
		t.Errorf("Source = %v, want MANUAL", record.Source)
	}
	if len(recordRepo.Records) != 1 {
		t.Errorf("expected 1 saved record, got %d", len(recordRepo.Records))
	}
}

func TestCreateRecordCommandHandler_ValidationFailure(t *testing.T) {
	t.Parallel()

	recordRepo := port.NewMockRecordRepo()
	cycleClient := port.NewMockCycleRuleClient()
	handler := setupCreateHandler(recordRepo, cycleClient)

	cmd := command.CreateRecordCommand{
		UserID:        "", // missing
		TreatmentID:   uuid.New(),
		TreatmentName: "보톡스",
		CategoryID:    uuid.New(),
		CategoryName:  "주름개선",
		TreatmentDate: time.Now(),
	}

	_, err := handler.Handle(context.Background(), cmd)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	var appErr *model.AppError
	if errors.As(err, &appErr) {
		if appErr.Code != model.ErrValidation {
			t.Errorf("error code = %v, want VALIDATION_ERROR", appErr.Code)
		}
	} else {
		t.Errorf("expected AppError, got %T", err)
	}
	if len(recordRepo.Records) != 0 {
		t.Error("record should not be saved on validation failure")
	}
}

func TestCreateRecordCommandHandler_CycleRuleNotFound(t *testing.T) {
	t.Parallel()

	recordRepo := port.NewMockRecordRepo()
	cycleClient := port.NewMockCycleRuleClient() // no rules configured

	handler := setupCreateHandler(recordRepo, cycleClient)

	cmd := command.CreateRecordCommand{
		UserID:        "user-1",
		TreatmentID:   uuid.New(),
		TreatmentName: "보톡스",
		CategoryID:    uuid.New(),
		CategoryName:  "주름개선",
		TreatmentDate: time.Now(),
	}

	record, err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record == nil {
		t.Fatal("expected record even without cycle rule")
	}
	// Record saved successfully despite no cycle rule
	if len(recordRepo.Records) != 1 {
		t.Errorf("expected 1 saved record, got %d", len(recordRepo.Records))
	}
}

func TestCreateRecordCommandHandler_SaveError(t *testing.T) {
	t.Parallel()

	recordRepo := port.NewMockRecordRepo()
	recordRepo.SaveFn = func(ctx context.Context, record *model.TreatmentRecord) error {
		return errors.New("db connection failed")
	}
	cycleClient := port.NewMockCycleRuleClient()
	handler := setupCreateHandler(recordRepo, cycleClient)

	cmd := command.CreateRecordCommand{
		UserID:        "user-1",
		TreatmentID:   uuid.New(),
		TreatmentName: "보톡스",
		CategoryID:    uuid.New(),
		CategoryName:  "주름개선",
		TreatmentDate: time.Now(),
	}

	_, err := handler.Handle(context.Background(), cmd)
	if err == nil {
		t.Fatal("expected error on save failure")
	}
	var appErr *model.AppError
	if errors.As(err, &appErr) {
		if appErr.Code != model.ErrInternal {
			t.Errorf("error code = %v, want INTERNAL_ERROR", appErr.Code)
		}
	}
}
