package postgres

import (
	"context"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RecordRepository struct {
	pool *pgxpool.Pool
}

func NewRecordRepository(pool *pgxpool.Pool) *RecordRepository {
	return &RecordRepository{pool: pool}
}

func (r *RecordRepository) Save(ctx context.Context, record *model.TreatmentRecord) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO treatment_records (id, user_id, treatment_id, treatment_name, category_id, category_name, treatment_date, hospital_name, dosage_value, dosage_unit, source, reservation_id, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		record.ID, record.UserID, record.TreatmentID, record.TreatmentName, record.CategoryID, record.CategoryName, record.TreatmentDate, record.HospitalName, record.DosageValue, record.DosageUnit, record.Source, record.ReservationID, record.CreatedAt, record.UpdatedAt)
	return err
}

func (r *RecordRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.TreatmentRecord, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT seq, id, user_id, treatment_id, treatment_name, category_id, category_name, treatment_date, hospital_name, dosage_value, dosage_unit, source, reservation_id, created_at, updated_at
		 FROM treatment_records WHERE id = $1`, id)
	return scanRecord(row)
}

func (r *RecordRepository) FindByUserIDAndDateRange(ctx context.Context, userID string, from, to time.Time) ([]*model.TreatmentRecord, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT seq, id, user_id, treatment_id, treatment_name, category_id, category_name, treatment_date, hospital_name, dosage_value, dosage_unit, source, reservation_id, created_at, updated_at
		 FROM treatment_records WHERE user_id = $1 AND treatment_date >= $2 AND treatment_date <= $3 ORDER BY treatment_date DESC`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*model.TreatmentRecord
	for rows.Next() {
		rec, err := scanRecordFromRows(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r *RecordRepository) FindByReservationID(ctx context.Context, reservationID string) (*model.TreatmentRecord, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT seq, id, user_id, treatment_id, treatment_name, category_id, category_name, treatment_date, hospital_name, dosage_value, dosage_unit, source, reservation_id, created_at, updated_at
		 FROM treatment_records WHERE reservation_id = $1`, reservationID)
	return scanRecord(row)
}

func (r *RecordRepository) Update(ctx context.Context, record *model.TreatmentRecord) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE treatment_records SET treatment_id=$1, treatment_name=$2, category_id=$3, category_name=$4, treatment_date=$5, hospital_name=$6, dosage_value=$7, dosage_unit=$8, updated_at=$9 WHERE id=$10`,
		record.TreatmentID, record.TreatmentName, record.CategoryID, record.CategoryName, record.TreatmentDate, record.HospitalName, record.DosageValue, record.DosageUnit, record.UpdatedAt, record.ID)
	return err
}

func (r *RecordRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM treatment_records WHERE id = $1`, id)
	return err
}

func (r *RecordRepository) GetStatsByUserID(ctx context.Context, userID string) ([]model.TreatmentStat, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT treatment_id, treatment_name, category_name, COUNT(*) as count
		 FROM treatment_records WHERE user_id = $1 GROUP BY treatment_id, treatment_name, category_name ORDER BY count DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []model.TreatmentStat
	for rows.Next() {
		var s model.TreatmentStat
		if err := rows.Scan(&s.TreatmentID, &s.TreatmentName, &s.CategoryName, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanRecord(row scannable) (*model.TreatmentRecord, error) {
	var rec model.TreatmentRecord
	err := row.Scan(&rec.Seq, &rec.ID, &rec.UserID, &rec.TreatmentID, &rec.TreatmentName, &rec.CategoryID, &rec.CategoryName, &rec.TreatmentDate, &rec.HospitalName, &rec.DosageValue, &rec.DosageUnit, &rec.Source, &rec.ReservationID, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRecordFromRows(row rowScanner) (*model.TreatmentRecord, error) {
	return scanRecord(row)
}
