package port

import (
	"context"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/google/uuid"
)

// MockRecordRepo implements TreatmentRecordRepository for testing.
type MockRecordRepo struct {
	Records map[uuid.UUID]*model.TreatmentRecord
	SaveFn  func(ctx context.Context, record *model.TreatmentRecord) error
}

func NewMockRecordRepo() *MockRecordRepo {
	return &MockRecordRepo{Records: make(map[uuid.UUID]*model.TreatmentRecord)}
}

func (m *MockRecordRepo) Save(ctx context.Context, record *model.TreatmentRecord) error {
	if m.SaveFn != nil {
		return m.SaveFn(ctx, record)
	}
	m.Records[record.ID] = record
	return nil
}

func (m *MockRecordRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.TreatmentRecord, error) {
	r, ok := m.Records[id]
	if !ok {
		return nil, model.NewAppError(model.ErrNotFound, "not found", nil)
	}
	return r, nil
}

func (m *MockRecordRepo) FindByUserIDAndDateRange(ctx context.Context, userID string, from, to time.Time) ([]*model.TreatmentRecord, error) {
	var result []*model.TreatmentRecord
	for _, r := range m.Records {
		if r.UserID == userID {
			result = append(result, r)
		}
	}
	return result, nil
}

func (m *MockRecordRepo) FindByReservationID(ctx context.Context, reservationID string) (*model.TreatmentRecord, error) {
	for _, r := range m.Records {
		if r.ReservationID != nil && *r.ReservationID == reservationID {
			return r, nil
		}
	}
	return nil, model.NewAppError(model.ErrNotFound, "not found", nil)
}

func (m *MockRecordRepo) Update(ctx context.Context, record *model.TreatmentRecord) error {
	m.Records[record.ID] = record
	return nil
}

func (m *MockRecordRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.Records, id)
	return nil
}

func (m *MockRecordRepo) GetStatsByUserID(ctx context.Context, userID string) ([]model.TreatmentStat, error) {
	return nil, nil
}

// MockScheduleRepo implements ScheduledTreatmentRepository for testing.
type MockScheduleRepo struct {
	Schedules map[uuid.UUID]*model.ScheduledTreatment
	SaveFn    func(ctx context.Context, schedule *model.ScheduledTreatment) error
}

func NewMockScheduleRepo() *MockScheduleRepo {
	return &MockScheduleRepo{Schedules: make(map[uuid.UUID]*model.ScheduledTreatment)}
}

func (m *MockScheduleRepo) Save(ctx context.Context, schedule *model.ScheduledTreatment) error {
	if m.SaveFn != nil {
		return m.SaveFn(ctx, schedule)
	}
	m.Schedules[schedule.ID] = schedule
	return nil
}

func (m *MockScheduleRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.ScheduledTreatment, error) {
	s, ok := m.Schedules[id]
	if !ok {
		return nil, model.NewAppError(model.ErrNotFound, "not found", nil)
	}
	return s, nil
}

func (m *MockScheduleRepo) FindByUserIDAndDateRange(ctx context.Context, userID string, from, to time.Time) ([]*model.ScheduledTreatment, error) {
	var result []*model.ScheduledTreatment
	for _, s := range m.Schedules {
		if s.UserID == userID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *MockScheduleRepo) FindDueSchedules(ctx context.Context, now time.Time) ([]*model.ScheduledTreatment, error) {
	var result []*model.ScheduledTreatment
	for _, s := range m.Schedules {
		if s.IsDue(now) {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *MockScheduleRepo) FindByRecordID(ctx context.Context, recordID uuid.UUID) ([]*model.ScheduledTreatment, error) {
	var result []*model.ScheduledTreatment
	for _, s := range m.Schedules {
		if s.RecordID == recordID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *MockScheduleRepo) Update(ctx context.Context, schedule *model.ScheduledTreatment) error {
	m.Schedules[schedule.ID] = schedule
	return nil
}

func (m *MockScheduleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.Schedules, id)
	return nil
}

func (m *MockScheduleRepo) DeleteByRecordID(ctx context.Context, recordID uuid.UUID) error {
	for id, s := range m.Schedules {
		if s.RecordID == recordID {
			delete(m.Schedules, id)
		}
	}
	return nil
}

func (m *MockScheduleRepo) FindFutureByUserAndTreatment(ctx context.Context, userID string, treatmentID uuid.UUID, after time.Time) ([]*model.ScheduledTreatment, error) {
	var result []*model.ScheduledTreatment
	for _, s := range m.Schedules {
		if s.UserID == userID && s.TreatmentID == treatmentID && s.ScheduledDate.After(after) {
			result = append(result, s)
		}
	}
	return result, nil
}
