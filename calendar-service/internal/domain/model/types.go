package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RecordSource
type RecordSource string

const (
	SourceAuto   RecordSource = "AUTO"
	SourceManual RecordSource = "MANUAL"
)

// ScheduleStatus
type ScheduleStatus string

const (
	StatusPending   ScheduleStatus = "PENDING"
	StatusReminded  ScheduleStatus = "REMINDED"
	StatusCompleted ScheduleStatus = "COMPLETED"
)

// DosageUnit
type DosageUnit string

const (
	UnitShot   DosageUnit = "shot"
	UnitMinute DosageUnit = "minute"
	UnitVolume DosageUnit = "volume"
	UnitVial   DosageUnit = "vial"
	UnitJoule  DosageUnit = "joule"
)

// TreatmentType
type TreatmentType string

const (
	TypeProcedure TreatmentType = "PROCEDURE"
	TypeSurgery   TreatmentType = "SURGERY"
)

// ErrorCode
type ErrorCode string

const (
	ErrValidation   ErrorCode = "VALIDATION_ERROR"
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrConflict     ErrorCode = "CONFLICT"
	ErrForbidden    ErrorCode = "FORBIDDEN"
	ErrInternal     ErrorCode = "INTERNAL_ERROR"
	ErrExternalAPI  ErrorCode = "EXTERNAL_API_ERROR"
)

// AppError
type AppError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewAppError(code ErrorCode, message string, cause error) *AppError {
	return &AppError{Code: code, Message: message, Cause: cause}
}

// TreatmentStat
type TreatmentStat struct {
	TreatmentID   uuid.UUID
	TreatmentName string
	CategoryName  string
	Count         int
}

// ReminderMessage
type ReminderMessage struct {
	UserID        string
	TreatmentName string
	DaysSince     int
	ScheduleID    uuid.UUID
}

// ReservationFixedEvent
type ReservationFixedEvent struct {
	ReservationID string        `json:"reservation_id"`
	UserID        string        `json:"user_id"`
	TreatmentID   string        `json:"treatment_id"`
	TreatmentName string        `json:"treatment_name"`
	CategoryID    string        `json:"category_id"`
	CategoryName  string        `json:"category_name"`
	TreatmentType TreatmentType `json:"treatment_type"`
	TreatmentDate time.Time     `json:"treatment_date"`
	HospitalName  string        `json:"hospital_name"`
	DosageValue   *float64      `json:"dosage_value,omitempty"`
	DosageUnit    *DosageUnit   `json:"dosage_unit,omitempty"`
}
