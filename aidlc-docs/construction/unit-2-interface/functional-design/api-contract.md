# Functional Design — Unit 2: 인터페이스 확정 + Dummy 통신 검증

## 목표

모든 API 엔드포인트의 계약(Request/Response)을 확정하고, dummy 데이터로 프론트↔백엔드 전체 통신을 검증한다.

---

## Calendar Service API Contract

### POST /api/records

**Request:**
```json
{
  "category_id": "019234ab-5678-7def-8000-000000000101",
  "treatment_id": "019234ab-5678-7def-8000-000000000201",
  "dosage_type": "volume",
  "dosage_value": "10.0",
  "treatment_date": "2026-01-15T00:00:00Z",
  "hospital_name": "강남언니의원",
  "hospital_location": "서울 강남구",
  "doctor_name": "김의사",
  "memo": "사각턱 보톡스 10cc"
}
```

**Response (201):**
```json
{
  "id": "019234ab-5678-7def-8000-000000000001",
  "user_id": "019234ab-5678-7def-8000-000000000901",
  "source": "MANUAL",
  "category_id": "019234ab-5678-7def-8000-000000000101",
  "treatment_id": "019234ab-5678-7def-8000-000000000201",
  "dosage_type": "volume",
  "dosage_value": "10.0",
  "treatment_date": "2026-01-15T00:00:00Z",
  "hospital_name": "강남언니의원",
  "hospital_location": "서울 강남구",
  "doctor_name": "김의사",
  "memo": "사각턱 보톡스 10cc",
  "created_at": "2026-01-15T09:00:00Z",
  "updated_at": "2026-01-15T09:00:00Z"
}
```

### GET /api/records?from=2026-01-01&to=2026-01-31

**Response (200):**
```json
[
  {
    "id": "019234ab-5678-7def-8000-000000000001",
    "user_id": "019234ab-5678-7def-8000-000000000901",
    "source": "AUTO",
    "category_id": "019234ab-5678-7def-8000-000000000101",
    "treatment_id": "019234ab-5678-7def-8000-000000000201",
    "dosage_type": "volume",
    "dosage_value": "10.0",
    "treatment_date": "2026-01-15T00:00:00Z",
    "hospital_name": "강남언니의원",
    "hospital_location": "서울 강남구",
    "doctor_name": null,
    "memo": null,
    "created_at": "2026-01-15T09:00:00Z",
    "updated_at": "2026-01-15T09:00:00Z"
  }
]
```

### GET /api/records/{id}

**Response (200):** 단건 Record (위와 동일 구조)

### PUT /api/records/{id}

**Request:** (partial — 변경할 필드만)
```json
{
  "treatment_date": "2026-01-16T00:00:00Z",
  "memo": "날짜 수정"
}
```

**Response (200):** 수정된 Record

### DELETE /api/records/{id}

**Response:** 204 No Content

### GET /api/schedules

**Response (200):**
```json
[
  {
    "id": "019234ab-5678-7def-8000-000000000501",
    "record_id": "019234ab-5678-7def-8000-000000000001",
    "category_id": "019234ab-5678-7def-8000-000000000101",
    "treatment_id": "019234ab-5678-7def-8000-000000000201",
    "scheduled_date": "2026-04-15T00:00:00Z",
    "cycle_days": 90,
    "status": "PENDING"
  }
]
```

### GET /api/schedules/{id}

**Response (200):** 단건 Schedule

### PATCH /api/schedules/{id}/complete

**Response (200):**
```json
{
  "id": "019234ab-5678-7def-8000-000000000501",
  "status": "COMPLETED",
  ...
}
```

### DELETE /api/schedules/{id}

**Response:** 204 No Content

### GET /api/statistics

**Response (200):**
```json
[
  { "category_id": "019234ab-5678-7def-8000-000000000101", "category_name": "보톡스", "count": 5 },
  { "category_id": "019234ab-5678-7def-8000-000000000102", "category_name": "필러", "count": 3 }
]
```

### GET /api/treatment-data/categories

**Response (200):**
```json
[
  { "id": "019234ab-5678-7def-8000-000000000101", "name": "보톡스" },
  { "id": "019234ab-5678-7def-8000-000000000102", "name": "필러" },
  { "id": "019234ab-5678-7def-8000-000000000103", "name": "레이저" }
]
```

### GET /api/treatment-data/categories/{id}/treatments

**Response (200):**
```json
[
  { "id": "019234ab-5678-7def-8000-000000000201", "category_id": "019234ab-5678-7def-8000-000000000101", "name": "사각턱" },
  { "id": "019234ab-5678-7def-8000-000000000202", "category_id": "019234ab-5678-7def-8000-000000000101", "name": "이마" },
  { "id": "019234ab-5678-7def-8000-000000000203", "category_id": "019234ab-5678-7def-8000-000000000101", "name": "스킨 보톡스" }
]
```

### GET /api/treatment-data/treatments/{id}/dosage-types

**Response (200):**
```json
[
  { "id": "019234ab-5678-7def-8000-000000000301", "treatment_id": "019234ab-5678-7def-8000-000000000201", "unit": "cc" },
  { "id": "019234ab-5678-7def-8000-000000000302", "treatment_id": "019234ab-5678-7def-8000-000000000201", "unit": "ml" }
]
```

### POST /mock/events/reservation-fixed

**Request:**
```json
{
  "reservation_id": "019234ab-5678-7def-8000-000000000801",
  "user_id": "019234ab-5678-7def-8000-000000000901",
  "hospital_name": "강남언니의원",
  "treatment_type": "treatment",
  "category_id": "019234ab-5678-7def-8000-000000000101",
  "treatment_id": "019234ab-5678-7def-8000-000000000201",
  "fixed_date": "2026-01-20T00:00:00Z"
}
```

**Response (202):**
```json
{ "message": "event received" }
```

---

## Admin Service API Contract

### POST /api/cycle-rules

**Request:**
```json
{ "category_id": "019234ab-5678-7def-8000-000000000101", "cycle_days": 90, "description": "보톡스 추천 주기" }
```

**Response (201):**
```json
{ "category_id": "019234ab-5678-7def-8000-000000000101", "cycle_days": 90, "description": "보톡스 추천 주기", "updated_at": "2026-01-01T00:00:00Z" }
```

### GET /api/cycle-rules

**Response (200):**
```json
[
  { "category_id": "019234ab-5678-7def-8000-000000000101", "cycle_days": 90, "description": "보톡스 추천 주기", "updated_at": "..." },
  { "category_id": "019234ab-5678-7def-8000-000000000102", "cycle_days": 180, "description": "필러 추천 주기", "updated_at": "..." }
]
```

### GET /api/cycle-rules/{categoryId}

**Response (200):** 단건 CycleRule

### PUT /api/cycle-rules/{categoryId}

**Request:** `{ "cycle_days": 120 }`
**Response (200):** 수정된 CycleRule

### DELETE /api/cycle-rules/{categoryId}

**Response:** 204 No Content

### Category/Treatment/DosageType CRUD

위 Calendar Service 프록시 응답과 동일한 구조. Admin Service가 원본 데이터 소유.

---

## Frontend 연결 계획

### 페이지별 API 호출 매핑

| 페이지 | 호출 API | Hook |
|--------|----------|------|
| CalendarPage | GET /api/records, GET /api/schedules, GET /api/statistics | useRecords, useSchedules, useStatistics |
| AddRecordPage | GET /api/treatment-data/*, POST /api/records | useCategories, useTreatments, useDosageTypes, useCreateRecord |
| RecordDetailPage | GET /api/records/{id}, PUT, DELETE | useRecord, useUpdateRecord, useDeleteRecord |
| DateBottomSheet | (CalendarPage 데이터 재사용) | — |
| ScheduleConfirmModal | PATCH /api/schedules/{id}/complete, DELETE /api/schedules/{id} | useCompleteSchedule, useDeleteSchedule |

### Dummy 데이터 시나리오

스텁 API가 반환할 dummy 데이터:
- 사용자 1명 (`019234ab-5678-7def-8000-000000000901`)
- 시술 기록 3건 (보톡스, 필러, 레이저 각 1건, 날짜 분산)
- 예정일 2건 (보톡스 90일 후, 필러 180일 후)
- 카테고리 3개, 시술명 카테고리당 3개, 용량 단위 시술당 2개
- 통계: 보톡스 3회, 필러 2회, 레이저 1회

---

## 완료 기준

- [ ] Calendar Service: 모든 엔드포인트가 dummy JSON 반환
- [ ] Admin Service: 모든 엔드포인트가 dummy JSON 반환
- [ ] Frontend: 모든 페이지에서 API 호출 → 데이터 표시 확인
- [ ] 프론트 → Calendar API → Admin API (프록시) 전체 체인 동작
- [ ] Request/Response 형식이 이 문서와 일치
