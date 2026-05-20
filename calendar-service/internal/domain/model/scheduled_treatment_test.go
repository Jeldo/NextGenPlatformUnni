package model_test

import (
	"testing"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/google/uuid"
)

func pendingSchedule() *model.ScheduledTreatment {
	return &model.ScheduledTreatment{
		ID:            uuid.New(),
		UserID:        "user-1",
		RecordID:      uuid.New(),
		TreatmentID:   uuid.New(),
		TreatmentName: "보톡스",
		ScheduledDate: time.Now().Add(24 * time.Hour),
		Status:        model.StatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func TestScheduledTreatment_MarkReminded(t *testing.T) {
	t.Parallel()

	s := pendingSchedule()
	now := time.Now()
	s.MarkReminded(now)

	if s.Status != model.StatusReminded {
		t.Errorf("Status = %v, want %v", s.Status, model.StatusReminded)
	}
	if s.RemindedAt == nil || !s.RemindedAt.Equal(now) {
		t.Errorf("RemindedAt = %v, want %v", s.RemindedAt, now)
	}
	if !s.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", s.UpdatedAt, now)
	}
}

func TestScheduledTreatment_MarkCompleted(t *testing.T) {
	t.Parallel()

	s := pendingSchedule()
	now := time.Now()
	s.MarkCompleted(now)

	if s.Status != model.StatusCompleted {
		t.Errorf("Status = %v, want %v", s.Status, model.StatusCompleted)
	}
	if s.CompletedAt == nil || !s.CompletedAt.Equal(now) {
		t.Errorf("CompletedAt = %v, want %v", s.CompletedAt, now)
	}
	if !s.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", s.UpdatedAt, now)
	}
}

func TestScheduledTreatment_IsDue(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		date   time.Time
		status model.ScheduleStatus
		want   bool
	}{
		{"pending and past due", now.Add(-1 * time.Hour), model.StatusPending, true},
		{"pending and exactly now", now, model.StatusPending, true},
		{"pending but future", now.Add(1 * time.Hour), model.StatusPending, false},
		{"reminded and past due", now.Add(-1 * time.Hour), model.StatusReminded, false},
		{"completed and past due", now.Add(-1 * time.Hour), model.StatusCompleted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := pendingSchedule()
			s.ScheduledDate = tt.date
			s.Status = tt.status
			if got := s.IsDue(now); got != tt.want {
				t.Errorf("IsDue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScheduledTreatment_BuildReminderMessage(t *testing.T) {
	t.Parallel()

	s := pendingSchedule()
	s.UserID = "user-42"
	s.TreatmentName = "레이저 토닝"

	msg := s.BuildReminderMessage(90)

	if msg.UserID != "user-42" {
		t.Errorf("UserID = %v, want user-42", msg.UserID)
	}
	if msg.TreatmentName != "레이저 토닝" {
		t.Errorf("TreatmentName = %v, want 레이저 토닝", msg.TreatmentName)
	}
	if msg.DaysSince != 90 {
		t.Errorf("DaysSince = %v, want 90", msg.DaysSince)
	}
	if msg.ScheduleID != s.ID {
		t.Errorf("ScheduleID = %v, want %v", msg.ScheduleID, s.ID)
	}
}

func TestScheduledTreatment_IsOwnedBy(t *testing.T) {
	t.Parallel()

	s := pendingSchedule()
	s.UserID = "user-1"

	if !s.IsOwnedBy("user-1") {
		t.Error("IsOwnedBy should return true for owner")
	}
	if s.IsOwnedBy("user-2") {
		t.Error("IsOwnedBy should return false for non-owner")
	}
}
