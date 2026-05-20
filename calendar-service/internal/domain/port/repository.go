package port

import (
	"context"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/google/uuid"
)

type TreatmentRecordRepository interface {
	Save(ctx context.Context, record *model.TreatmentRecord) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.TreatmentRecord, error)
	FindByUserIDAndDateRange(ctx context.Context, userID string, from, to time.Time) ([]*model.TreatmentRecord, error)
	FindByReservationID(ctx context.Context, reservationID string) (*model.TreatmentRecord, error)
	Update(ctx context.Context, record *model.TreatmentRecord) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetStatsByUserID(ctx context.Context, userID string) ([]model.TreatmentStat, error)
}

type ScheduledTreatmentRepository interface {
	Save(ctx context.Context, schedule *model.ScheduledTreatment) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.ScheduledTreatment, error)
	FindByUserIDAndDateRange(ctx context.Context, userID string, from, to time.Time) ([]*model.ScheduledTreatment, error)
	FindDueSchedules(ctx context.Context, now time.Time) ([]*model.ScheduledTreatment, error)
	FindByRecordID(ctx context.Context, recordID uuid.UUID) ([]*model.ScheduledTreatment, error)
	Update(ctx context.Context, schedule *model.ScheduledTreatment) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByRecordID(ctx context.Context, recordID uuid.UUID) error
	FindFutureByUserAndTreatment(ctx context.Context, userID string, treatmentID uuid.UUID, after time.Time) ([]*model.ScheduledTreatment, error)
}
