package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ScheduledTreatment struct {
	Seq           int64
	ID            uuid.UUID
	UserID        string
	RecordID      uuid.UUID
	TreatmentID   uuid.UUID
	TreatmentName string
	ScheduledDate time.Time
	Status        ScheduleStatus
	RemindedAt    *time.Time
	CompletedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (s *ScheduledTreatment) MarkReminded(now time.Time) {
	s.Status = StatusReminded
	s.RemindedAt = &now
	s.UpdatedAt = now
}

func (s *ScheduledTreatment) MarkCompleted(now time.Time) {
	s.Status = StatusCompleted
	s.CompletedAt = &now
	s.UpdatedAt = now
}

func (s *ScheduledTreatment) IsOwnedBy(userID string) bool {
	return s.UserID == userID
}

func (s *ScheduledTreatment) IsDue(now time.Time) bool {
	return !now.Before(s.ScheduledDate) && s.Status == StatusPending
}

func (s *ScheduledTreatment) BuildReminderMessage(daysSince int) ReminderMessage {
	return ReminderMessage{
		UserID:        s.UserID,
		TreatmentName: s.TreatmentName,
		DaysSince:     daysSince,
		ScheduleID:    s.ID,
	}
}

func (s *ScheduledTreatment) ReminderText() string {
	return fmt.Sprintf("%s 맞은 지 %d일이 됐어요", s.TreatmentName, int(time.Since(s.CreatedAt).Hours()/24))
}
