package model_test

import (
	"testing"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/google/uuid"
)

func validRecord() *model.TreatmentRecord {
	return &model.TreatmentRecord{
		ID:            uuid.New(),
		UserID:        "user-1",
		TreatmentID:   uuid.New(),
		TreatmentName: "보톡스",
		CategoryID:    uuid.New(),
		CategoryName:  "주름개선",
		TreatmentDate: time.Now(),
		Source:        model.SourceManual,
	}
}

func TestTreatmentRecord_Validate(t *testing.T) {
	t.Parallel()

	unit := model.UnitShot
	tests := []struct {
		name    string
		modify  func(r *model.TreatmentRecord)
		wantErr bool
		errCode model.ErrorCode
	}{
		{"valid record", func(r *model.TreatmentRecord) {}, false, ""},
		{"missing user_id", func(r *model.TreatmentRecord) { r.UserID = "" }, true, model.ErrValidation},
		{"missing treatment_id", func(r *model.TreatmentRecord) { r.TreatmentID = uuid.Nil }, true, model.ErrValidation},
		{"missing treatment_name", func(r *model.TreatmentRecord) { r.TreatmentName = "" }, true, model.ErrValidation},
		{"missing category_id", func(r *model.TreatmentRecord) { r.CategoryID = uuid.Nil }, true, model.ErrValidation},
		{"missing category_name", func(r *model.TreatmentRecord) { r.CategoryName = "" }, true, model.ErrValidation},
		{"missing treatment_date", func(r *model.TreatmentRecord) { r.TreatmentDate = time.Time{} }, true, model.ErrValidation},
		{"dosage_value without unit", func(r *model.TreatmentRecord) {
			v := 10.0
			r.DosageValue = &v
			r.DosageUnit = nil
		}, true, model.ErrValidation},
		{"dosage_value with unit", func(r *model.TreatmentRecord) {
			v := 10.0
			r.DosageValue = &v
			r.DosageUnit = &unit
		}, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := validRecord()
			tt.modify(r)
			err := r.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errCode != "" && err.Code != tt.errCode {
				t.Errorf("Validate() code = %v, want %v", err.Code, tt.errCode)
			}
		})
	}
}

func TestTreatmentRecord_ApplyUpdate(t *testing.T) {
	t.Parallel()

	r := validRecord()
	newTreatmentID := uuid.New()
	newCategoryID := uuid.New()
	newDate := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	hospital := "서울성형외과"
	dosage := 20.0
	unit := model.UnitShot

	r.ApplyUpdate(newTreatmentID, "필러", newCategoryID, "볼륨", newDate, &hospital, &dosage, &unit)

	if r.TreatmentID != newTreatmentID {
		t.Errorf("TreatmentID = %v, want %v", r.TreatmentID, newTreatmentID)
	}
	if r.TreatmentName != "필러" {
		t.Errorf("TreatmentName = %v, want 필러", r.TreatmentName)
	}
	if r.CategoryID != newCategoryID {
		t.Errorf("CategoryID = %v, want %v", r.CategoryID, newCategoryID)
	}
	if r.CategoryName != "볼륨" {
		t.Errorf("CategoryName = %v, want 볼륨", r.CategoryName)
	}
	if !r.TreatmentDate.Equal(newDate) {
		t.Errorf("TreatmentDate = %v, want %v", r.TreatmentDate, newDate)
	}
	if *r.HospitalName != hospital {
		t.Errorf("HospitalName = %v, want %v", *r.HospitalName, hospital)
	}
	if *r.DosageValue != dosage {
		t.Errorf("DosageValue = %v, want %v", *r.DosageValue, dosage)
	}
	if *r.DosageUnit != unit {
		t.Errorf("DosageUnit = %v, want %v", *r.DosageUnit, unit)
	}
}

func TestTreatmentRecord_IsOwnedBy(t *testing.T) {
	t.Parallel()

	r := validRecord()
	r.UserID = "user-123"

	tests := []struct {
		name   string
		userID string
		want   bool
	}{
		{"owner", "user-123", true},
		{"not owner", "user-456", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := r.IsOwnedBy(tt.userID); got != tt.want {
				t.Errorf("IsOwnedBy(%q) = %v, want %v", tt.userID, got, tt.want)
			}
		})
	}
}

func TestTreatmentRecord_IsDateOrCategoryChanged(t *testing.T) {
	t.Parallel()

	r := validRecord()
	originalID := r.TreatmentID
	originalDate := r.TreatmentDate

	tests := []struct {
		name        string
		treatmentID uuid.UUID
		date        time.Time
		want        bool
	}{
		{"no change", originalID, originalDate, false},
		{"treatment changed", uuid.New(), originalDate, true},
		{"date changed", originalID, originalDate.Add(24 * time.Hour), true},
		{"both changed", uuid.New(), originalDate.Add(24 * time.Hour), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := r.IsDateOrCategoryChanged(tt.treatmentID, tt.date); got != tt.want {
				t.Errorf("IsDateOrCategoryChanged() = %v, want %v", got, tt.want)
			}
		})
	}
}
