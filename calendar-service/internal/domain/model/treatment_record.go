package model

import (
	"time"

	"github.com/google/uuid"
)

type TreatmentRecord struct {
	Seq           int64
	ID            uuid.UUID
	UserID        string
	TreatmentID   uuid.UUID
	TreatmentName string
	CategoryID    uuid.UUID
	CategoryName  string
	TreatmentDate time.Time
	HospitalName  *string
	DosageValue   *float64
	DosageUnit    *DosageUnit
	Source        RecordSource
	ReservationID *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (r *TreatmentRecord) Validate() *AppError {
	if r.UserID == "" {
		return NewAppError(ErrValidation, "user_id is required", nil)
	}
	if r.TreatmentID == uuid.Nil {
		return NewAppError(ErrValidation, "treatment_id is required", nil)
	}
	if r.TreatmentName == "" {
		return NewAppError(ErrValidation, "treatment_name is required", nil)
	}
	if r.CategoryID == uuid.Nil {
		return NewAppError(ErrValidation, "category_id is required", nil)
	}
	if r.CategoryName == "" {
		return NewAppError(ErrValidation, "category_name is required", nil)
	}
	if r.TreatmentDate.IsZero() {
		return NewAppError(ErrValidation, "treatment_date is required", nil)
	}
	if r.DosageValue != nil && r.DosageUnit == nil {
		return NewAppError(ErrValidation, "dosage_unit is required when dosage_value is set", nil)
	}
	return nil
}

func (r *TreatmentRecord) ApplyUpdate(treatmentID uuid.UUID, treatmentName string, categoryID uuid.UUID, categoryName string, treatmentDate time.Time, hospitalName *string, dosageValue *float64, dosageUnit *DosageUnit) {
	r.TreatmentID = treatmentID
	r.TreatmentName = treatmentName
	r.CategoryID = categoryID
	r.CategoryName = categoryName
	r.TreatmentDate = treatmentDate
	r.HospitalName = hospitalName
	r.DosageValue = dosageValue
	r.DosageUnit = dosageUnit
}

func (r *TreatmentRecord) IsDateOrCategoryChanged(newTreatmentID uuid.UUID, newDate time.Time) bool {
	return r.TreatmentID != newTreatmentID || !r.TreatmentDate.Equal(newDate)
}

func (r *TreatmentRecord) IsOwnedBy(userID string) bool {
	return r.UserID == userID
}
