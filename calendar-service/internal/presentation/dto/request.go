package dto

type CreateRecordRequest struct {
	UserID        string  `json:"user_id"`
	TreatmentID   string  `json:"treatment_id"`
	TreatmentName string  `json:"treatment_name"`
	CategoryID    string  `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	TreatmentDate string  `json:"treatment_date"`
	HospitalName  *string `json:"hospital_name,omitempty"`
	DosageValue   *string `json:"dosage_value,omitempty"`
	DosageUnit    *string `json:"dosage_unit,omitempty"`
}

type UpdateRecordRequest struct {
	TreatmentID   string  `json:"treatment_id"`
	TreatmentName string  `json:"treatment_name"`
	CategoryID    string  `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	TreatmentDate string  `json:"treatment_date"`
	HospitalName  *string `json:"hospital_name,omitempty"`
	DosageValue   *string `json:"dosage_value,omitempty"`
	DosageUnit    *string `json:"dosage_unit,omitempty"`
}
