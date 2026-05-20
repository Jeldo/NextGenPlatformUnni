package dto

import (
	"fmt"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
)

type RecordResponse struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	TreatmentID   string  `json:"treatment_id"`
	TreatmentName string  `json:"treatment_name"`
	CategoryID    string  `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	TreatmentDate string  `json:"treatment_date"`
	HospitalName  *string `json:"hospital_name,omitempty"`
	DosageValue   *string `json:"dosage_value,omitempty"`
	DosageUnit    *string `json:"dosage_unit,omitempty"`
	Source        string  `json:"source"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type ScheduleResponse struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	RecordID      string  `json:"record_id"`
	TreatmentID   string  `json:"treatment_id"`
	TreatmentName string  `json:"treatment_name"`
	ScheduledDate string  `json:"scheduled_date"`
	Status        string  `json:"status"`
	RemindedAt    *string `json:"reminded_at,omitempty"`
	CompletedAt   *string `json:"completed_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type StatisticsResponse struct {
	Stats []StatItem `json:"stats"`
}

type StatItem struct {
	TreatmentID   string `json:"treatment_id"`
	TreatmentName string `json:"treatment_name"`
	CategoryName  string `json:"category_name"`
	Count         int    `json:"count"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ToRecordResponse(r *model.TreatmentRecord) RecordResponse {
	resp := RecordResponse{
		ID:            r.ID.String(),
		UserID:        r.UserID,
		TreatmentID:   r.TreatmentID.String(),
		TreatmentName: r.TreatmentName,
		CategoryID:    r.CategoryID.String(),
		CategoryName:  r.CategoryName,
		TreatmentDate: r.TreatmentDate.Format(time.RFC3339),
		HospitalName:  r.HospitalName,
		Source:        string(r.Source),
		CreatedAt:     r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     r.UpdatedAt.Format(time.RFC3339),
	}
	if r.DosageValue != nil {
		v := formatFloat(*r.DosageValue)
		resp.DosageValue = &v
	}
	if r.DosageUnit != nil {
		u := string(*r.DosageUnit)
		resp.DosageUnit = &u
	}
	return resp
}

func ToScheduleResponse(s *model.ScheduledTreatment) ScheduleResponse {
	resp := ScheduleResponse{
		ID:            s.ID.String(),
		UserID:        s.UserID,
		RecordID:      s.RecordID.String(),
		TreatmentID:   s.TreatmentID.String(),
		TreatmentName: s.TreatmentName,
		ScheduledDate: s.ScheduledDate.Format(time.RFC3339),
		Status:        string(s.Status),
		CreatedAt:     s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     s.UpdatedAt.Format(time.RFC3339),
	}
	if s.RemindedAt != nil {
		t := s.RemindedAt.Format(time.RFC3339)
		resp.RemindedAt = &t
	}
	if s.CompletedAt != nil {
		t := s.CompletedAt.Format(time.RFC3339)
		resp.CompletedAt = &t
	}
	return resp
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%g", f)
}
