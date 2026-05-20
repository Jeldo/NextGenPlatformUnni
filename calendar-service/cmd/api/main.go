package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /health", healthHandler)

	// Records
	mux.HandleFunc("POST /api/records", createRecordHandler)
	mux.HandleFunc("GET /api/records", listRecordsHandler)
	mux.HandleFunc("GET /api/records/{id}", getRecordHandler)
	mux.HandleFunc("PUT /api/records/{id}", updateRecordHandler)
	mux.HandleFunc("DELETE /api/records/{id}", deleteRecordHandler)

	// Schedules
	mux.HandleFunc("GET /api/schedules", listSchedulesHandler)
	mux.HandleFunc("GET /api/schedules/{id}", getScheduleHandler)
	mux.HandleFunc("PATCH /api/schedules/{id}/complete", completeScheduleHandler)
	mux.HandleFunc("DELETE /api/schedules/{id}", deleteScheduleHandler)

	// Statistics
	mux.HandleFunc("GET /api/statistics", getStatisticsHandler)

	// Treatment Data (proxy)
	mux.HandleFunc("GET /api/treatment-data/categories", listCategoriesHandler)
	mux.HandleFunc("GET /api/treatment-data/categories/{id}/treatments", listTreatmentsHandler)
	mux.HandleFunc("GET /api/treatment-data/treatments/{id}/dosage-types", listDosageTypesHandler)

	// Mock
	mux.HandleFunc("POST /mock/events/reservation-fixed", mockEventHandler)

	log.Printf("Calendar API starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(mux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok", "service": "calendar-service", "version": "0.1.0"})
}

// --- Dummy Data ---

var dummyRecords = []map[string]any{
	{
		"id": "019234ab-5678-7def-8000-000000000001", "user_id": "019234ab-0000-7def-8000-000000000001",
		"source": "AUTO", "category_id": "019234ab-1111-7def-8000-000000000001",
		"treatment_id": "019234ab-2222-7def-8000-000000000001",
		"dosage_type": "volume", "dosage_value": "10.0",
		"treatment_date": "2026-01-15T00:00:00Z", "hospital_name": "강남언니의원",
		"hospital_location": "서울 강남구", "doctor_name": nil, "memo": nil,
		"created_at": "2026-01-15T09:00:00Z", "updated_at": "2026-01-15T09:00:00Z",
	},
	{
		"id": "019234ab-5678-7def-8000-000000000002", "user_id": "019234ab-0000-7def-8000-000000000001",
		"source": "MANUAL", "category_id": "019234ab-1111-7def-8000-000000000002",
		"treatment_id": "019234ab-2222-7def-8000-000000000004",
		"dosage_type": "volume", "dosage_value": "1.0",
		"treatment_date": "2026-01-20T00:00:00Z", "hospital_name": "청담피부과",
		"hospital_location": "서울 강남구", "doctor_name": "박의사", "memo": "입술 필러",
		"created_at": "2026-01-20T10:00:00Z", "updated_at": "2026-01-20T10:00:00Z",
	},
	{
		"id": "019234ab-5678-7def-8000-000000000003", "user_id": "019234ab-0000-7def-8000-000000000001",
		"source": "AUTO", "category_id": "019234ab-1111-7def-8000-000000000003",
		"treatment_id": "019234ab-2222-7def-8000-000000000007",
		"dosage_type": "joule", "dosage_value": "15.0",
		"treatment_date": "2026-02-01T00:00:00Z", "hospital_name": "강남언니의원",
		"hospital_location": nil, "doctor_name": nil, "memo": nil,
		"created_at": "2026-02-01T09:00:00Z", "updated_at": "2026-02-01T09:00:00Z",
	},
}

var dummySchedules = []map[string]any{
	{
		"id": "019234ab-9999-7def-8000-000000000001", "record_id": "019234ab-5678-7def-8000-000000000001",
		"category_id": "019234ab-1111-7def-8000-000000000001", "treatment_id": "019234ab-2222-7def-8000-000000000001",
		"scheduled_date": "2026-04-15T00:00:00Z", "cycle_days": 90, "status": "PENDING",
	},
	{
		"id": "019234ab-9999-7def-8000-000000000002", "record_id": "019234ab-5678-7def-8000-000000000002",
		"category_id": "019234ab-1111-7def-8000-000000000002", "treatment_id": "019234ab-2222-7def-8000-000000000004",
		"scheduled_date": "2026-07-20T00:00:00Z", "cycle_days": 180, "status": "PENDING",
	},
}

var dummyCategories = []map[string]any{
	{"id": "019234ab-1111-7def-8000-000000000001", "name": "보톡스"},
	{"id": "019234ab-1111-7def-8000-000000000002", "name": "필러"},
	{"id": "019234ab-1111-7def-8000-000000000003", "name": "레이저"},
}

var dummyTreatments = map[string][]map[string]any{
	"019234ab-1111-7def-8000-000000000001": {
		{"id": "019234ab-2222-7def-8000-000000000001", "category_id": "019234ab-1111-7def-8000-000000000001", "name": "사각턱"},
		{"id": "019234ab-2222-7def-8000-000000000002", "category_id": "019234ab-1111-7def-8000-000000000001", "name": "이마"},
		{"id": "019234ab-2222-7def-8000-000000000003", "category_id": "019234ab-1111-7def-8000-000000000001", "name": "스킨 보톡스"},
	},
	"019234ab-1111-7def-8000-000000000002": {
		{"id": "019234ab-2222-7def-8000-000000000004", "category_id": "019234ab-1111-7def-8000-000000000002", "name": "입술"},
		{"id": "019234ab-2222-7def-8000-000000000005", "category_id": "019234ab-1111-7def-8000-000000000002", "name": "팔자주름"},
		{"id": "019234ab-2222-7def-8000-000000000006", "category_id": "019234ab-1111-7def-8000-000000000002", "name": "볼"},
	},
	"019234ab-1111-7def-8000-000000000003": {
		{"id": "019234ab-2222-7def-8000-000000000007", "category_id": "019234ab-1111-7def-8000-000000000003", "name": "리프팅"},
		{"id": "019234ab-2222-7def-8000-000000000008", "category_id": "019234ab-1111-7def-8000-000000000003", "name": "토닝"},
		{"id": "019234ab-2222-7def-8000-000000000009", "category_id": "019234ab-1111-7def-8000-000000000003", "name": "제모"},
	},
}

var dummyDosageTypes = []map[string]any{
	{"id": "019234ab-3333-7def-8000-000000000001", "treatment_id": "019234ab-2222-7def-8000-000000000001", "unit": "shot"},
	{"id": "019234ab-3333-7def-8000-000000000002", "treatment_id": "019234ab-2222-7def-8000-000000000001", "unit": "volume"},
}

// --- Handlers ---

func createRecordHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 201, dummyRecords[0])
}

func listRecordsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, dummyRecords)
}

func getRecordHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, dummyRecords[0])
}

func updateRecordHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, dummyRecords[0])
}

func deleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}

func listSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, dummySchedules)
}

func getScheduleHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, dummySchedules[0])
}

func completeScheduleHandler(w http.ResponseWriter, r *http.Request) {
	completed := dummySchedules[0]
	completed["status"] = "COMPLETED"
	writeJSON(w, 200, completed)
}

func deleteScheduleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}

func getStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, []map[string]any{
		{"category_id": "019234ab-1111-7def-8000-000000000001", "category_name": "보톡스", "count": 3},
		{"category_id": "019234ab-1111-7def-8000-000000000002", "category_name": "필러", "count": 2},
		{"category_id": "019234ab-1111-7def-8000-000000000003", "category_name": "레이저", "count": 1},
	})
}

func listCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, dummyCategories)
}

func listTreatmentsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	treatments, ok := dummyTreatments[id]
	if !ok {
		treatments = []map[string]any{}
	}
	writeJSON(w, 200, treatments)
}

func listDosageTypesHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, dummyDosageTypes)
}

func mockEventHandler(w http.ResponseWriter, r *http.Request) {
	_ = strings.NewReader("")
	writeJSON(w, 202, map[string]string{"message": "event received"})
}
