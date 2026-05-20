package service_test

import (
	"testing"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/service"
	"github.com/google/uuid"
	"pgregory.net/rapid"
)

func TestCycleCalculator_Calculate(t *testing.T) {
	t.Parallel()

	calc := service.NewCycleCalculator()

	tests := []struct {
		name      string
		date      time.Time
		cycleDays int
		wantDate  time.Time
	}{
		{"90 days (botox)", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), 90, time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)},
		{"180 days (filler)", time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), 180, time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC)},
		{"14 days (laser)", time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), 14, time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)},
		{"365 days", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), 365, time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			record := &model.TreatmentRecord{
				ID:            uuid.New(),
				UserID:        "user-1",
				TreatmentID:   uuid.New(),
				TreatmentName: "test",
				TreatmentDate: tt.date,
			}
			result := calc.Calculate(record, tt.cycleDays)
			if !result.ScheduledDate.Equal(tt.wantDate) {
				t.Errorf("ScheduledDate = %v, want %v", result.ScheduledDate, tt.wantDate)
			}
			if result.Status != model.StatusPending {
				t.Errorf("Status = %v, want PENDING", result.Status)
			}
			if result.UserID != record.UserID {
				t.Errorf("UserID = %v, want %v", result.UserID, record.UserID)
			}
			if result.RecordID != record.ID {
				t.Errorf("RecordID = %v, want %v", result.RecordID, record.ID)
			}
		})
	}
}

func TestCycleCalculator_Properties(t *testing.T) {
	t.Parallel()

	calc := service.NewCycleCalculator()

	rapid.Check(t, func(t *rapid.T) {
		year := rapid.IntRange(2020, 2030).Draw(t, "year")
		month := rapid.IntRange(1, 12).Draw(t, "month")
		day := rapid.IntRange(1, 28).Draw(t, "day")
		cycleDays := rapid.IntRange(1, 730).Draw(t, "cycleDays")

		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		record := &model.TreatmentRecord{
			ID:            uuid.New(),
			UserID:        "user-1",
			TreatmentID:   uuid.New(),
			TreatmentName: "test",
			TreatmentDate: date,
		}

		result := calc.Calculate(record, cycleDays)

		// Property 1: result date is always after input date
		if !result.ScheduledDate.After(date) {
			t.Fatalf("ScheduledDate %v should be after input date %v (cycleDays=%d)", result.ScheduledDate, date, cycleDays)
		}

		// Property 2: cycleDays > 0 means result != input date
		if result.ScheduledDate.Equal(date) {
			t.Fatalf("ScheduledDate should not equal input date when cycleDays=%d", cycleDays)
		}

		// Property 3: status is always PENDING
		if result.Status != model.StatusPending {
			t.Fatalf("Status = %v, want PENDING", result.Status)
		}

		// Property 4: record fields are preserved
		if result.UserID != record.UserID {
			t.Fatalf("UserID mismatch")
		}
		if result.RecordID != record.ID {
			t.Fatalf("RecordID mismatch")
		}
		if result.TreatmentID != record.TreatmentID {
			t.Fatalf("TreatmentID mismatch")
		}
	})
}
