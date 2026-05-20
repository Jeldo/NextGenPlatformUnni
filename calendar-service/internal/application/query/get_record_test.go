package query_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/application/query"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/google/uuid"
)

func TestGetRecordQueryHandler_Success(t *testing.T) {
	t.Parallel()

	repo := port.NewMockRecordRepo()
	recordID := uuid.New()
	repo.Records[recordID] = &model.TreatmentRecord{
		ID:            recordID,
		UserID:        "user-1",
		TreatmentID:   uuid.New(),
		TreatmentName: "보톡스",
		CategoryID:    uuid.New(),
		CategoryName:  "주름개선",
		TreatmentDate: time.Now(),
		Source:        model.SourceManual,
	}

	handler := query.NewGetRecordQueryHandler(repo)

	result, err := handler.Handle(context.Background(), query.GetRecordQuery{
		ID:     recordID,
		UserID: "user-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != recordID {
		t.Errorf("ID = %v, want %v", result.ID, recordID)
	}
}

func TestGetRecordQueryHandler_NotFound(t *testing.T) {
	t.Parallel()

	repo := port.NewMockRecordRepo()
	handler := query.NewGetRecordQueryHandler(repo)

	_, err := handler.Handle(context.Background(), query.GetRecordQuery{
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

func TestGetRecordQueryHandler_Unauthorized(t *testing.T) {
	t.Parallel()

	repo := port.NewMockRecordRepo()
	recordID := uuid.New()
	repo.Records[recordID] = &model.TreatmentRecord{
		ID:            recordID,
		UserID:        "user-1",
		TreatmentID:   uuid.New(),
		TreatmentName: "보톡스",
		CategoryID:    uuid.New(),
		CategoryName:  "주름개선",
		TreatmentDate: time.Now(),
		Source:        model.SourceManual,
	}

	handler := query.NewGetRecordQueryHandler(repo)

	_, err := handler.Handle(context.Background(), query.GetRecordQuery{
		ID:     recordID,
		UserID: "user-other",
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
}
