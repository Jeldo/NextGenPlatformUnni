package postgres

import (
	"context"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository struct {
	pool *pgxpool.Pool
}

func NewScheduleRepository(pool *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{pool: pool}
}

func (r *ScheduleRepository) Save(ctx context.Context, schedule *model.ScheduledTreatment) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO scheduled_treatments (id, user_id, record_id, treatment_id, treatment_name, scheduled_date, status, reminded_at, completed_at, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		schedule.ID, schedule.UserID, schedule.RecordID, schedule.TreatmentID, schedule.TreatmentName, schedule.ScheduledDate, schedule.Status, schedule.RemindedAt, schedule.CompletedAt, schedule.CreatedAt, schedule.UpdatedAt)
	return err
}

func (r *ScheduleRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.ScheduledTreatment, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT seq, id, user_id, record_id, treatment_id, treatment_name, scheduled_date, status, reminded_at, completed_at, created_at, updated_at
		 FROM scheduled_treatments WHERE id = $1`, id)
	return scanSchedule(row)
}

func (r *ScheduleRepository) FindByUserIDAndDateRange(ctx context.Context, userID string, from, to time.Time) ([]*model.ScheduledTreatment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT seq, id, user_id, record_id, treatment_id, treatment_name, scheduled_date, status, reminded_at, completed_at, created_at, updated_at
		 FROM scheduled_treatments WHERE user_id = $1 AND scheduled_date >= $2 AND scheduled_date <= $3 ORDER BY scheduled_date ASC`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*model.ScheduledTreatment
	for rows.Next() {
		s, err := scanSchedule(rows)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *ScheduleRepository) FindDueSchedules(ctx context.Context, now time.Time) ([]*model.ScheduledTreatment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT seq, id, user_id, record_id, treatment_id, treatment_name, scheduled_date, status, reminded_at, completed_at, created_at, updated_at
		 FROM scheduled_treatments WHERE status = $1 AND scheduled_date <= $2`, model.StatusPending, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*model.ScheduledTreatment
	for rows.Next() {
		s, err := scanSchedule(rows)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *ScheduleRepository) FindByRecordID(ctx context.Context, recordID uuid.UUID) ([]*model.ScheduledTreatment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT seq, id, user_id, record_id, treatment_id, treatment_name, scheduled_date, status, reminded_at, completed_at, created_at, updated_at
		 FROM scheduled_treatments WHERE record_id = $1`, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*model.ScheduledTreatment
	for rows.Next() {
		s, err := scanSchedule(rows)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *ScheduleRepository) Update(ctx context.Context, schedule *model.ScheduledTreatment) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE scheduled_treatments SET status=$1, reminded_at=$2, completed_at=$3, updated_at=$4 WHERE id=$5`,
		schedule.Status, schedule.RemindedAt, schedule.CompletedAt, schedule.UpdatedAt, schedule.ID)
	return err
}

func (r *ScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM scheduled_treatments WHERE id = $1`, id)
	return err
}

func (r *ScheduleRepository) DeleteByRecordID(ctx context.Context, recordID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM scheduled_treatments WHERE record_id = $1`, recordID)
	return err
}

func (r *ScheduleRepository) FindFutureByUserAndTreatment(ctx context.Context, userID string, treatmentID uuid.UUID, after time.Time) ([]*model.ScheduledTreatment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT seq, id, user_id, record_id, treatment_id, treatment_name, scheduled_date, status, reminded_at, completed_at, created_at, updated_at
		 FROM scheduled_treatments WHERE user_id = $1 AND treatment_id = $2 AND scheduled_date > $3 AND status = $4`, userID, treatmentID, after, model.StatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*model.ScheduledTreatment
	for rows.Next() {
		s, err := scanSchedule(rows)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func scanSchedule(row scannable) (*model.ScheduledTreatment, error) {
	var s model.ScheduledTreatment
	err := row.Scan(&s.Seq, &s.ID, &s.UserID, &s.RecordID, &s.TreatmentID, &s.TreatmentName, &s.ScheduledDate, &s.Status, &s.RemindedAt, &s.CompletedAt, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
