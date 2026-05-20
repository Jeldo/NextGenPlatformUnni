# Calendar Service (Go) — Comprehensive Design

> Hexagonal Architecture (4-Layer) | REST API | Port 8080 | Event Consumer (별도 프로세스)

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  1. PRESENTATION LAYER (Inbound Adapters)                    │
│  ┌──────────────┐  ┌──────────────────────┐                │
│  │ HTTP Handler  │  │ SQS Event Consumer   │                │
│  │ (REST API)    │  │ (별도 프로세스)        │                │
│  └──────┬───────┘  └──────────┬───────────┘                │
│         │                      │                             │
├─────────┼──────────────────────┼─────────────────────────────┤
│  2. APPLICATION LAYER (Use Cases / Orchestration)            │
│  ┌──────────────────────────────────────────────────┐       │
│  │  TreatmentRecordUseCase                          │       │
│  │  ScheduledTreatmentUseCase                       │       │
│  │  TreatmentStatisticsUseCase                      │       │
│  │  ReminderUseCase                                 │       │
│  └──────────────────────────────────────────────────┘       │
│         │                                                    │
├─────────┼────────────────────────────────────────────────────┤
│  3. DOMAIN LAYER (Models & Domain Services)                  │
│  ┌──────────────────────────────────────────────────┐       │
│  │  TreatmentRecord (model)                         │       │
│  │  ScheduledTreatment (model)                      │       │
│  │  TreatmentStat (model)                           │       │
│  │  CycleCalculator (domain service)                │       │
│  │  Ports (Repository interfaces, Client interfaces)│       │
│  └──────────────────────────────────────────────────┘       │
│         ^                                                    │
├─────────┼────────────────────────────────────────────────────┤
│  4. INFRASTRUCTURE LAYER (Outbound Adapters)                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Postgres     │  │ HTTP Cycle   │  │ HTTP Notif   │      │
│  │ Repositories │  │ Rule Client  │  │ Client       │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘

의존성 방향: Presentation → Application → Domain ← Infrastructure
```

---

## Components

### 1. Presentation Layer (Inbound Adapters)

| Component | 책임 |
|-----------|------|
| **HTTPHandler** | REST API 요청 수신, 입력 파싱/검증, Application Layer 호출, HTTP 응답 반환 |
| **SQSEventConsumer** | SQS 메시지 폴링, 이벤트 파싱, Application Layer 호출, ACK/NACK |
| **CronScheduler** | 별도 CronJob 프로세스 — 매일 ProcessRemindersCommand 트리거 |
| **MockEventHTTPAdapter** | 개발용 Mock 엔드포인트 (`/mock/events/reservation-fixed`) |

### 2. Application Layer (CQRS — Command/Query Handlers)

> CUD는 Command + CommandHandler, Read는 Query + QueryHandler로 분리

#### Commands (CUD)

| Command | Handler | 책임 |
|---------|---------|------|
| **CreateRecordCommand** | CreateRecordHandler | 시술 기록 생성 → 주기 조회 → 예정일 계산 |
| **UpdateRecordCommand** | UpdateRecordHandler | 시술 기록 수정 → 예정일 재계산 |
| **DeleteRecordCommand** | DeleteRecordHandler | 시술 기록 삭제 → 예정일 삭제 |
| **HandleReservationFixedCommand** | HandleReservationFixedHandler | 예약 확정 이벤트 → 자동 기입 |
| **CalculateScheduleCommand** | CalculateScheduleHandler | 예정일 계산/저장 |
| **CompleteScheduleCommand** | CompleteScheduleHandler | 예정일 완료 처리 (PENDING/REMINDED → COMPLETED) |
| **DeleteScheduleCommand** | DeleteScheduleHandler | 예정일 삭제 |
| **ProcessRemindersCommand** | ProcessRemindersHandler | 리마인드 배치 처리 |

#### Queries (Read)

| Query | Handler | 책임 |
|-------|---------|------|
| **ListRecordsQuery** | ListRecordsHandler | 기간별 시술 기록 목록 조회 |
| **GetRecordQuery** | GetRecordHandler | 시술 기록 단건 조회 |
| **ListSchedulesQuery** | ListSchedulesHandler | 예정일 목록 조회 |
| **GetScheduleQuery** | GetScheduleHandler | 예정일 단건 조회 |
| **GetStatisticsQuery** | GetStatisticsHandler | 시술 통계 집계 조회 |

### 3. Domain Layer (Rich Domain Models & Domain Services)

> Anemic Domain Model 금지 — 비즈니스 로직을 최대한 Domain Model에 포함

| Component | 유형 | 책임 |
|-----------|------|------|
| **TreatmentRecord** | Rich Model | 시술 기록 엔티티, 검증, 수정 적용 |
| **ScheduledTreatment** | Rich Model | 예정일 엔티티, 상태 전이 규칙, 리마인드 메시지 생성 |
| **TreatmentStat** | Value Object | 통계 값 객체 |
| **CycleCalculator** | Domain Service | 시술일 + 주기 → 예정일 계산 (순수 로직, PBT 대상) |
| **TreatmentRecordRepository** | Port (interface) | 시술 기록 영속성 인터페이스 (FindFutureByCategory 포함) |
| **ScheduledTreatmentRepository** | Port (interface) | 예정일 영속성 인터페이스 |
| **CycleRuleClient** | Port (interface) | 추천 주기 조회 인터페이스 |
| **NotificationClient** | Port (interface) | 알림 발송 인터페이스 |
| **EventSubscriber** | Port (interface) | 이벤트 수신 인터페이스 |

### 4. Infrastructure Layer (Outbound Adapters)

| Adapter | Port 구현 | 환경 |
|---------|-----------|------|
| **PostgresRecordRepo** | TreatmentRecordRepository | 공통 |
| **PostgresScheduleRepo** | ScheduledTreatmentRepository | 공통 |
| **HTTPCycleRuleClient** | CycleRuleClient | 공통 |
| **HTTPNotificationClient** | NotificationClient | Production |
| **MockNotificationClient** | NotificationClient | Development (로그만) |
| **SQSEventSubscriber** | EventSubscriber | Production |
| **InMemoryEventSubscriber** | EventSubscriber | Development |

---

## Data Structures

### Core Models

> enum은 string 리터럴 대신 typed const 사용

```go
// === Enum Types (typed const) ===

type RecordSource string
const (
    SourceAuto   RecordSource = "AUTO"
    SourceManual RecordSource = "MANUAL"
)

type ScheduleStatus string
const (
    StatusPending   ScheduleStatus = "PENDING"
    StatusReminded  ScheduleStatus = "REMINDED"
    StatusCompleted ScheduleStatus = "COMPLETED"
)

type DosageUnit string
const (
    DosageVolume DosageUnit = "volume"
)

type TreatmentType string
const (
    TypeSurgery   TreatmentType = "surgery"
    TypeTreatment TreatmentType = "treatment"
)

// === Rich Domain Models ===

type TreatmentRecord struct {
    ID              string
    UserID          string
    Source          RecordSource
    CategoryID      string
    TreatmentID     string
    DosageType      *DosageUnit
    DosageValue     *float64
    TreatmentDate   time.Time
    HospitalName    string
    HospitalLocation *string
    DoctorName      *string
    Memo            *string
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// --- TreatmentRecord Business Logic ---

// Validate 필수 필드 검증
func (r *TreatmentRecord) Validate() error

// ApplyUpdate 변경사항 적용, 변경된 필드 목록 반환 (날짜/카테고리 변경 감지용)
func (r *TreatmentRecord) ApplyUpdate(input UpdateRecordInput) (changedFields []string)

// IsDateOrCategoryChanged 예정일 재계산이 필요한 변경인지 판단
func (r *TreatmentRecord) IsDateOrCategoryChanged(changes []string) bool

// IsOwnedBy 소유권 검증
func (r *TreatmentRecord) IsOwnedBy(userID string) bool

// IsSurgery 수술 여부 판단 (이벤트 필터링용)
func (e *ReservationFixedEvent) IsSurgery() bool


type ScheduledTreatment struct {
    ID              string
    UserID          string
    RecordID        string
    CategoryID      string
    TreatmentID     string
    ScheduledDate   time.Time
    CycleDays       int
    Status          ScheduleStatus
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// --- ScheduledTreatment Business Logic ---

// MarkReminded 리마인드 발송 완료 상태 전이 (PENDING → REMINDED)
func (s *ScheduledTreatment) MarkReminded() error

// MarkCompleted 새 시술 기록 생성으로 완료 처리 (PENDING|REMINDED → COMPLETED)
func (s *ScheduledTreatment) MarkCompleted() error

// IsDue 리마인드 대상인지 판단
func (s *ScheduledTreatment) IsDue(today time.Time) bool

// IsOwnedBy 소유권 검증
func (s *ScheduledTreatment) IsOwnedBy(userID string) bool

// BuildReminderMessage 알림 메시지 생성 (도메인 로직)
func (s *ScheduledTreatment) BuildReminderMessage(treatmentName string, treatmentDate time.Time) ReminderMessage


type TreatmentStat struct {
    CategoryID   string
    CategoryName string
    Count        int
}

type ReminderMessage struct {
    UserID        string
    Title         string
    Body          string
    TreatmentName string
    DaysSince     int
}

type ReservationFixedEvent struct {
    ReservationID string
    UserID        string
    HospitalName  string
    TreatmentType TreatmentType
    CategoryID    string
    TreatmentID   string
    FixedDate     time.Time
    EventTime     time.Time
}
```

### HTTP Request/Response DTOs

> HTTP 레이어에서는 Request/Response suffix 사용
> 값 타입 = required, 포인터(*) = optional

```go
type CreateRecordRequest struct {
    // Required
    CategoryID    string    `json:"category_id" validate:"required"`
    TreatmentID   string    `json:"treatment_id" validate:"required"`
    TreatmentDate time.Time `json:"treatment_date" validate:"required"`
    HospitalName  string    `json:"hospital_name" validate:"required"`

    // Optional
    DosageType       *string  `json:"dosage_type,omitempty"`
    DosageValue      *float64 `json:"dosage_value,omitempty"`
    HospitalLocation *string  `json:"hospital_location,omitempty"`
    DoctorName       *string  `json:"doctor_name,omitempty"`
    Memo             *string  `json:"memo,omitempty"`
}

type UpdateRecordRequest struct {
    // All optional (partial update)
    CategoryID       *string    `json:"category_id,omitempty"`
    TreatmentID      *string    `json:"treatment_id,omitempty"`
    DosageType       *string    `json:"dosage_type,omitempty"`
    DosageValue      *float64   `json:"dosage_value,omitempty"`
    TreatmentDate    *time.Time `json:"treatment_date,omitempty"`
    HospitalName     *string    `json:"hospital_name,omitempty"`
    HospitalLocation *string    `json:"hospital_location,omitempty"`
    DoctorName       *string    `json:"doctor_name,omitempty"`
    Memo             *string    `json:"memo,omitempty"`
}

type RecordResponse struct {
    ID               string     `json:"id"`
    UserID           string     `json:"user_id"`
    Source           string     `json:"source"`
    CategoryID       string     `json:"category_id"`
    TreatmentID      string     `json:"treatment_id"`
    TreatmentDate    time.Time  `json:"treatment_date"`
    HospitalName     string     `json:"hospital_name"`
    DosageType       *string    `json:"dosage_type,omitempty"`
    DosageValue      *float64   `json:"dosage_value,omitempty"`
    HospitalLocation *string    `json:"hospital_location,omitempty"`
    DoctorName       *string    `json:"doctor_name,omitempty"`
    Memo             *string    `json:"memo,omitempty"`
    CreatedAt        time.Time  `json:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at"`
}

type ScheduleResponse struct {
    ID            string     `json:"id"`
    RecordID      string     `json:"record_id"`
    CategoryID    string     `json:"category_id"`
    TreatmentID   string     `json:"treatment_id"`
    ScheduledDate time.Time  `json:"scheduled_date"`
    CycleDays     int        `json:"cycle_days"`
    Status        string     `json:"status"`
}

type StatisticsResponse struct {
    CategoryID   string `json:"category_id"`
    CategoryName string `json:"category_name"`
    Count        int    `json:"count"`
}
```

### Error Types

```go
type AppError struct {
    Code    ErrorCode
    Message string
    Cause   error
}

type ErrorCode string
const (
    ErrNotFound          ErrorCode = "NOT_FOUND"
    ErrValidation        ErrorCode = "VALIDATION_ERROR"
    ErrCycleRuleNotFound ErrorCode = "CYCLE_RULE_NOT_FOUND"
    ErrExternalService   ErrorCode = "EXTERNAL_SERVICE_ERROR"
    ErrDatabase          ErrorCode = "DATABASE_ERROR"
    ErrUnauthorized      ErrorCode = "UNAUTHORIZED"
    ErrConflict          ErrorCode = "CONFLICT"
)
```

---

## REST API Endpoints

| Method | Path | 설명 | 응답 |
|--------|------|------|------|
| POST | /api/records | 시술 수동 추가 | 201 + Record |
| GET | /api/records?from=&to= | 기간별 시술 목록 | 200 + []Record |
| GET | /api/records/{id} | 시술 단건 조회 | 200 + Record |
| PUT | /api/records/{id} | 시술 수정 | 200 + Record |
| DELETE | /api/records/{id} | 시술 삭제 | 204 |
| GET | /api/schedules | 예정일 목록 | 200 + []Schedule |
| GET | /api/schedules/{id} | 예정일 단건 조회 | 200 + Schedule |
| PATCH | /api/schedules/{id}/complete | 예정일 완료 처리 | 200 + Schedule |
| DELETE | /api/schedules/{id} | 예정일 삭제 | 204 |
| GET | /api/statistics | 시술 통계 | 200 + []Stat |
| GET | /api/treatment-data/categories | 카테고리 목록 (프록시) | 200 |
| GET | /api/treatment-data/categories/{id}/treatments | 시술명 목록 (프록시) | 200 |
| GET | /api/treatment-data/treatments/{id}/dosage-types | 용량 단위 (프록시) | 200 |

### Mock Endpoints (개발 전용)

| Method | Path | 설명 |
|--------|------|------|
| POST | /mock/events/reservation-fixed | 예약 확정 이벤트 수동 발행 |

### Error Response Format

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "treatment_date is required"
  }
}
```

| HTTP Status | Error Code | 상황 |
|-------------|-----------|------|
| 400 | VALIDATION_ERROR | 입력값 검증 실패 |
| 404 | NOT_FOUND | 리소스 미존재 |
| 409 | CONFLICT | 중복 리소스 |
| 500 | DATABASE_ERROR | DB 오류 |
| 502 | EXTERNAL_SERVICE_ERROR | 외부 서비스 호출 실패 |

---

## Method Flows (by Layer)

### Presentation Layer → Application Layer

#### HTTPHandler.CreateRecord (POST /api/records)
```
1. HTTP 요청 파싱 (JSON → CreateRecordCommand)
2. 입력값 기본 검증 (필수 필드 존재 여부)
3. CreateRecordHandler.Handle(ctx, cmd) 호출
4. 성공 → 201 + JSON 응답
5. 에러 → ErrorCode에 따라 HTTP Status 매핑
```

#### SQSEventConsumer.Poll
```
1. SQS ReceiveMessage (long poll 20s)
2. 메시지 파싱 → HandleReservationFixedCommand
3. HandleReservationFixedHandler.Handle(ctx, cmd) 호출
4. 성공 → DeleteMessage (ACK)
5. 실패 → visibility timeout 후 재시도, 3회 초과 시 DLQ
```

### Application Layer (CQRS Handlers)

#### CreateRecordHandler.Handle(cmd CreateRecordCommand)
```
1. record 생성 → record.Validate()
2. TreatmentRecordRepository.Save(record)
3. CycleRuleClient.GetCycleRule(record.CategoryID)
   - 주기 존재 → CalculateScheduleHandler.Handle()
   - 주기 미존재 → 건너뜀 (에러 아님)
4. 생성된 TreatmentRecord 반환

에러: 검증실패→ErrValidation, DB실패→ErrDatabase, CycleRule실패→graceful skip
```

#### UpdateRecordHandler.Handle(cmd UpdateRecordCommand)
```
1. TreatmentRecordRepository.FindByID(recordID)
2. record.IsOwnedBy(requestUserID) → false면 ErrUnauthorized
3. changes := record.ApplyUpdate(cmd.Input)
4. TreatmentRecordRepository.Update(record)
5. record.IsDateOrCategoryChanged(changes) → true면:
   a. ScheduledTreatmentRepository.DeleteByRecordID()
   b. CycleRuleClient.GetCycleRule(record.CategoryID)
   c. CalculateScheduleHandler.Handle()
6. 수정된 Record 반환

에러: 미존재→ErrNotFound, 소유권→ErrUnauthorized, CycleRule실패→기존 예정일 유지
```

#### DeleteRecordHandler.Handle(cmd DeleteRecordCommand)
```
1. record := FindByID → record.IsOwnedBy() 검증
2. ScheduledTreatmentRepository.DeleteByRecordID() (실패 시 로그 후 계속)
3. TreatmentRecordRepository.Delete()

에러: 미존재→ErrNotFound, 소유권→ErrUnauthorized
```

#### HandleReservationFixedHandler.Handle(cmd HandleReservationFixedCommand)
```
1. cmd.Event.IsSurgery() → true면 무시 (로그)
2. 멱등성 체크 (reservationID로 기존 기록 조회)
3. CreateRecordCommand 구성 (Source=AUTO)
4. CreateRecordHandler.Handle() 호출

에러: 중복→스킵(성공), Create실패→ErrDatabase
```

#### CalculateScheduleHandler.Handle(cmd CalculateScheduleCommand)
```
1. scheduledDate := CycleCalculator.Calculate(treatmentDate, cycleDays)
2. schedule := ScheduledTreatment{Status: StatusPending, ScheduledDate: scheduledDate, ...}
3. ScheduledTreatmentRepository.Save(schedule)

에러: DB실패→ErrDatabase
```

#### CompleteScheduleHandler.Handle(cmd CompleteScheduleCommand)
```
1. schedule := ScheduledTreatmentRepository.FindByID(scheduleID)
2. schedule.IsOwnedBy(userID) 검증
3. schedule.MarkCompleted()
4. ScheduledTreatmentRepository.Update(schedule)
5. 완료된 Schedule 반환

에러: 미존재→ErrNotFound, 소유권→ErrUnauthorized, 이미 완료→ErrConflict
```

#### DeleteScheduleHandler.Handle(cmd DeleteScheduleCommand)
```
1. schedule := ScheduledTreatmentRepository.FindByID(scheduleID)
2. schedule.IsOwnedBy(userID) 검증
3. ScheduledTreatmentRepository.Delete(scheduleID)

에러: 미존재→ErrNotFound, 소유권→ErrUnauthorized
```

#### ProcessRemindersHandler.Handle(cmd ProcessRemindersCommand) (배치)
```
1. schedules := ScheduledTreatmentRepository.FindDue(today)
2. 각 schedule:
   a. 해당 카테고리의 미래 예약 존재 여부 확인 (TreatmentRecordRepository.FindFutureByCategory)
   b. 메시지 분기:
      - 미래 예약 있음 → "보톡스 예약이 n월n일에 있어요"
      - 미래 예약 없음 → "보톡스 맞은 지 N일이 됐어요. 다음 시술을 예약해볼까요?"
   c. NotificationClient.SendReminder(msg)
   d. 성공 → schedule.MarkReminded() → Repository.Update(schedule)
   e. 실패 → 로그, 다음 배치 재시도
```

#### ListRecordsHandler.Handle(query ListRecordsQuery)
```
1. 기간 검증 (from < to, 최대 1년)
2. TreatmentRecordRepository.FindByUserAndPeriod(userID, from, to)
3. 결과 반환

에러: 기간 유효하지 않음→ErrValidation
```

#### GetRecordHandler.Handle(query GetRecordQuery)
```
1. TreatmentRecordRepository.FindByID(recordID)
2. 소유권 검증
3. Record 반환

에러: 미존재→ErrNotFound, 소유권→ErrUnauthorized
```

#### ListSchedulesHandler.Handle(query ListSchedulesQuery)
```
1. ScheduledTreatmentRepository.FindByUser(userID)
2. 결과 반환
```

#### GetStatisticsHandler.Handle(query GetStatisticsQuery)
```
1. TreatmentRecordRepository.CountByUserGrouped(userID)
2. 카테고리별 집계 반환
```

### Domain Layer (Rich Model Methods)

#### TreatmentRecord.Validate
```go
func (r *TreatmentRecord) Validate() error {
    if r.CategoryID == "" { return ErrValidation("category_id required") }
    if r.TreatmentID == "" { return ErrValidation("treatment_id required") }
    if r.TreatmentDate.IsZero() { return ErrValidation("treatment_date required") }
    if r.HospitalName == "" { return ErrValidation("hospital_name required") }
    if r.DosageType != nil && *r.DosageType != DosageVolume {
        return ErrValidation("dosage_type must be volume")
    }
    return nil
}
```

#### TreatmentRecord.ApplyUpdate
```go
func (r *TreatmentRecord) ApplyUpdate(input UpdateRecordInput) []string {
    var changed []string
    if input.TreatmentDate != nil && !input.TreatmentDate.Equal(r.TreatmentDate) {
        r.TreatmentDate = *input.TreatmentDate
        changed = append(changed, "treatment_date")
    }
    if input.CategoryID != nil && *input.CategoryID != r.CategoryID {
        r.CategoryID = *input.CategoryID
        changed = append(changed, "category_id")
    }
    // ... 나머지 필드 적용
    r.UpdatedAt = time.Now()
    return changed
}
```

#### TreatmentRecord.IsDateOrCategoryChanged
```go
func (r *TreatmentRecord) IsDateOrCategoryChanged(changes []string) bool {
    for _, c := range changes {
        if c == "treatment_date" || c == "category_id" { return true }
    }
    return false
}
```

#### ScheduledTreatment.MarkReminded
```go
func (s *ScheduledTreatment) MarkReminded() error {
    if s.Status != StatusPending {
        return ErrValidation("can only remind PENDING schedules")
    }
    s.Status = StatusReminded
    s.UpdatedAt = time.Now()
    return nil
}
```

#### ScheduledTreatment.MarkCompleted
```go
func (s *ScheduledTreatment) MarkCompleted() error {
    if s.Status == StatusCompleted {
        return ErrConflict("already completed")
    }
    s.Status = StatusCompleted
    s.UpdatedAt = time.Now()
    return nil
}
```

#### ScheduledTreatment.BuildReminderMessage
```go
func (s *ScheduledTreatment) BuildReminderMessage(treatmentName string, treatmentDate time.Time) ReminderMessage {
    daysSince := int(time.Since(treatmentDate).Hours() / 24)
    return ReminderMessage{
        UserID:        s.UserID,
        Title:         "시술 리마인드",
        Body:          fmt.Sprintf("%s 맞은 지 %d일이 됐어요. 다음 시술을 예약해볼까요?", treatmentName, daysSince),
        TreatmentName: treatmentName,
        DaysSince:     daysSince,
    }
}
```

#### ReservationFixedEvent.IsSurgery
```go
func (e *ReservationFixedEvent) IsSurgery() bool {
    return e.TreatmentType == TypeSurgery
}
```

#### CycleCalculator.Calculate (Domain Service — PBT 대상)
```go
func (c *CycleCalculator) Calculate(treatmentDate time.Time, cycleDays int) time.Time {
    return treatmentDate.AddDate(0, 0, cycleDays)
}

// Properties (PBT):
// - 결과는 항상 treatmentDate보다 미래
// - cycleDays > 0이면 결과 != treatmentDate
// - Calculate(date, a) + b일 == Calculate(date, a+b)
```

#### ScheduledTreatment State Transitions
```
[PENDING] ──MarkReminded()──→ [REMINDED]
[PENDING] ──MarkCompleted()──→ [COMPLETED]
[REMINDED] ──MarkCompleted()──→ [COMPLETED]
```

---

## Event Consumer (별도 프로세스)

### Processing State Machine

```
[RECEIVED] → [VALIDATED] → [PROCESSING] → [COMPLETED] → ACK
                 ↓                ↓
              [SKIPPED]       [RETRY] → [DLQ]
              (surgery)       (3회 초과)
```

### Retry Policy

| 상황 | 처리 |
|------|------|
| 파싱 실패 | DLQ (재시도 불가) |
| DB 일시 오류 | 재시도 3회 (exponential backoff: 1s, 2s, 4s) |
| 중복 이벤트 | 스킵 + ACK (멱등성) |
| 3회 재시도 실패 | DLQ + 알림 |

---

## External Dependencies

| 대상 | 프로토콜 | Timeout | Retry | Circuit Breaker |
|------|----------|---------|-------|-----------------|
| Admin API | HTTP | 3s | 1회 | 5회 실패 → 30s open |
| 알림 시스템 | HTTP | 5s | 1회 | — |
| PostgreSQL | TCP/SQL | 5s | 0 | — |
| AWS SQS | AWS SDK | 20s (long poll) | 자동 | — |

---

## Failure Isolation

| 장애 | Calendar API 동작 | 사용자 영향 |
|------|-------------------|------------|
| Admin API 다운 | 기록 정상, 예정일 미생성 | 예정일 표시 안됨 |
| 알림 시스템 다운 | 배치 스킵, 다음 배치 재시도 | 리마인드 지연 |
| SQS 지연 | Event Consumer 대기 | 자동 기입 지연 |
| PostgreSQL 다운 | 전체 중단 | 서비스 이용 불가 |
