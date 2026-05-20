package port

import (
	"context"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/google/uuid"
)

type CycleRule struct {
	TreatmentID uuid.UUID
	CycleDays   int
}

type CycleRuleClient interface {
	GetByTreatmentID(ctx context.Context, treatmentID uuid.UUID) (*CycleRule, error)
}

type NotificationClient interface {
	SendReminder(ctx context.Context, msg model.ReminderMessage) error
}
