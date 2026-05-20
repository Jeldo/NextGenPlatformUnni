package service

import (
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/google/uuid"
)

type CycleCalculator struct{}

func NewCycleCalculator() *CycleCalculator {
	return &CycleCalculator{}
}

func (c *CycleCalculator) Calculate(record *model.TreatmentRecord, cycleDays int) *model.ScheduledTreatment {
	scheduledDate := record.TreatmentDate.AddDate(0, 0, cycleDays)
	now := time.Now()
	id, _ := uuid.NewV7()
	return &model.ScheduledTreatment{
		ID:            id,
		UserID:        record.UserID,
		RecordID:      record.ID,
		TreatmentID:   record.TreatmentID,
		TreatmentName: record.TreatmentName,
		ScheduledDate: scheduledDate,
		Status:        model.StatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
