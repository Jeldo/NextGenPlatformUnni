# Functional Design — Unit 3-B: Go Calendar Service

## 1. DB Schema

> 규칙:
> - 모든 테이블에 `seq` (BIGSERIAL) 운영 편의 필드 추가
> - UUID는 애플리케이션에서 생성/주입 (DB DEFAULT 사용 안 함)
> - enum CHECK 제약 없음 — 검증은 애플리케이션에서 수행 (확장성)
> - created_at, updated_at은 애플리케이션에서 주입

### treatment_records

```sql
CREATE TABLE treatment_records (
    seq               BIGSERIAL,
    id                UUID PRIMARY KEY,
    user_id           UUID NOT NULL,
    source            VARCHAR(10) NOT NULL,
    reservation_id    VARCHAR(255),
    category_id       UUID NOT NULL,
    treatment_id      UUID NOT NULL,
    dosage_type       VARCHAR(10),
    dosage_value      DECIMAL(10,2),
    treatment_date    DATE NOT NULL,
    hospital_name     VARCHAR(255) NOT NULL,
    hospital_location VARCHAR(255),
    doctor_name       VARCHAR(255),
    memo              TEXT,
    created_at        TIMESTAMPTZ NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL
);

-- 기간별 조회 (GET /api/records?from=&to=)
CREATE INDEX idx_records_user_date ON treatment_records (user_id, treatment_date);

-- 멱등성 체크 (Event Consumer — reservation_id 중복 방지)
CREATE UNIQUE INDEX idx_records_reservation_id ON treatment_records (reservation_id) WHERE reservation_id IS NOT NULL;

-- 통계 집계 (GET /api/statistics — user별 category 그룹핑)
CREATE INDEX idx_records_user_category ON treatment_records (user_id, category_id);

-- 리마인드 분기: 같은 treatment의 미래 기록 존재 여부 확인
CREATE INDEX idx_records_user_treatment_date ON treatment_records (user_id, treatment_id, treatment_date);
```

### scheduled_treatments

```sql
CREATE TABLE scheduled_treatments (
    seq             BIGSERIAL,
    id              UUID PRIMARY KEY,
    user_id         UUID NOT NULL,
    record_id       UUID NOT NULL REFERENCES treatment_records(id) ON DELETE CASCADE,
    category_id     UUID NOT NULL,
    treatment_id    UUID NOT NULL,
    scheduled_date  DATE NOT NULL,
    cycle_days      INTEGER NOT NULL,
    status          VARCHAR(10) NOT NULL DEFAULT 'PENDING',
    created_at      TIMESTAMPTZ NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL
);

-- 리마인드 배치 조회 (status=PENDING, scheduled_date <= today)
CREATE INDEX idx_schedules_due ON scheduled_treatments (status, scheduled_date) WHERE status = 'PENDING';

-- 사용자별 예정일 조회 (GET /api/schedules)
CREATE INDEX idx_schedules_user ON scheduled_treatments (user_id);

-- record 삭제 시 cascade 조회
CREATE INDEX idx_schedules_record ON scheduled_treatments (record_id);
```

---

## 2. Tech Stack 결정

| 항목 | 선택 | 근거 |
|------|------|------|
| Go 버전 | 최신 (1.24+) | 최신 기능 활용 |
| HTTP Router | `echo` (v4) | 미들웨어 생태계, 간결한 API |
| DB Driver | `pgx/v5` | 가장 빠른 PostgreSQL 드라이버 |
| Migration | `golang-migrate` | 파일 기반, 간단 |
| UUID | `google/uuid` (v7 생성) | UUIDv7 지원 |
| Logging | `log/slog` (표준 라이브러리) | 구조화 로깅 |
| Config | `os.Getenv` 직접 사용 | 단순, 외부 의존성 없음 |
| Testing | `testing` + Docker Compose DB | Classicist, 실제 DB 사용 |
| PBT | `pgregory.net/rapid` | Go PBT 표준 |

---

## 3. 비즈니스 규칙 상세

### CycleRule 조회 단위: treatment (시술명 단위)

> 주기는 **treatment 단위**로 관리됨 (category가 아님)

```
- GET /api/cycle-rules/{treatmentId} 로 조회
- 예: "사각턱 보톡스" = 90일, "이마 보톡스" = 60일 (같은 카테고리라도 다를 수 있음)
```

### 시술 기록 생성 시 예정일 계산

```
1. Record 저장 성공
2. Admin API에서 CycleRule 조회 (GET /api/cycle-rules/{treatmentId})
   - 200 + cycle_days → 예정일 = treatment_date + cycle_days
   - 404 (규칙 없음) → 예정일 미생성 (정상)
   - 5xx / timeout → 로그 기록, 예정일 미생성 (graceful degradation)
3. 예정일 저장 (status = PENDING)
```

### 시술 기록 수정 시 예정일 재계산

```
1. 날짜 또는 treatment 변경 감지
2. 변경된 경우에만:
   a. 기존 예정일 삭제 (record_id 기준)
   b. 새 CycleRule 조회 (treatmentId 기준)
   c. 새 예정일 계산 및 저장
3. 변경되지 않은 경우: 예정일 유지
```

### 멱등성 (Event Consumer)

```
1. ReservationFixedEvent 수신
2. reservation_id로 기존 기록 조회 (UNIQUE INDEX)
3. 이미 존재 → 스킵 + ACK (중복 이벤트)
4. 미존재 → 정상 처리
```

### 리마인드 배치 규칙

```
1. 대상: status = 'PENDING' AND scheduled_date <= today
2. 각 대상에 대해:
   a. 같은 treatment의 미래 기록 존재 여부 확인
      - 있음 → "보톡스(사각턱) 예약이 n월n일에 있어요"
      - 없음 → "보톡스(사각턱) 맞은 지 N일이 됐어요. 다음 시술을 예약해볼까요?"
   b. NotificationClient.SendReminder()
   c. 성공 → status = 'REMINDED'
   d. 실패 → 로그, 다음 배치 재시도
```

### Schedule Complete 규칙

```
1. PATCH /api/schedules/{id}/complete 호출
2. status가 PENDING 또는 REMINDED인 경우만 → COMPLETED
3. 이미 COMPLETED → 409 Conflict
```

---

## 4. Repository Interface 상세

```go
type TreatmentRecordRepository interface {
    Save(ctx context.Context, record *model.TreatmentRecord) error
    FindByID(ctx context.Context, id string) (*model.TreatmentRecord, error)
    FindByUserAndPeriod(ctx context.Context, userID string, from, to time.Time) ([]model.TreatmentRecord, error)
    FindByReservationID(ctx context.Context, reservationID string) (*model.TreatmentRecord, error)
    FindFutureByTreatment(ctx context.Context, userID, treatmentID string, after time.Time) (*model.TreatmentRecord, error)
    Update(ctx context.Context, record *model.TreatmentRecord) error
    Delete(ctx context.Context, id string) error
    CountByUserGrouped(ctx context.Context, userID string) ([]model.TreatmentStat, error)
}

type ScheduledTreatmentRepository interface {
    Save(ctx context.Context, schedule *model.ScheduledTreatment) error
    FindByID(ctx context.Context, id string) (*model.ScheduledTreatment, error)
    FindByUser(ctx context.Context, userID string) ([]model.ScheduledTreatment, error)
    FindDue(ctx context.Context, date time.Time) ([]model.ScheduledTreatment, error)
    DeleteByRecordID(ctx context.Context, recordID string) error
    Delete(ctx context.Context, id string) error
    Update(ctx context.Context, schedule *model.ScheduledTreatment) error
}
```

---

## 5. 완료 기준

- [ ] Domain Layer: TreatmentRecord, ScheduledTreatment (Rich Model + 모든 메서드)
- [ ] Domain Layer: CycleCalculator + PBT 테스트
- [ ] Application Layer: 모든 Command/Query Handlers
- [ ] Infrastructure Layer: PostgreSQL Repositories (pgx)
- [ ] Infrastructure Layer: HTTPCycleRuleClient (Circuit Breaker)
- [ ] Infrastructure Layer: MockNotificationClient
- [ ] Infrastructure Layer: InMemoryEventSubscriber
- [ ] Presentation Layer: HTTP Handlers (echo)
- [ ] Presentation Layer: SQS Event Consumer (별도 main)
- [ ] Presentation Layer: CronScheduler (별도 main)
- [ ] DB Migration 파일
- [ ] 단위 테스트 (Domain + Application) — Classicist
- [ ] 통합 테스트 (Repository — Docker Compose DB)
