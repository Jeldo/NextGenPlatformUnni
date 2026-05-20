# Services (Comprehensive) — 시술 관리 캘린더

## Service Overview

```
+-------------------+       +-------------------+       +-------------------+
|  Calendar Web     |       |  Calendar API     |       |  Admin API        |
|  (Next.js)        |       |  (Go, REST)       |       |  (FastAPI)        |
|  Port: 3000       |       |  Port: 8080       |       |  Port: 8081       |
+-------------------+       +-------------------+       +-------------------+
        |                           |                           |
        +-------- REST (8080) ------+                           |
        |                           +------ HTTP (internal) ----+
        |                           |                           |
        |                           v                           v
        |               +-------------------+                   |
        |               |  Event Consumer   |                   |
        |               |  (Go, SQS)        |                   |
        |               +-------------------+                   |
        |                           |                           |
        v                           v                           v
+---------------------------------------------------------------+
|                      PostgreSQL (shared)                        |
+---------------------------------------------------------------+
        |
        v
+-------------------+
| Google Calendar   |
| API (client-side) |
+-------------------+
```

---

## Service 0: Calendar Web (Next.js)

| 항목 | 설명 |
|------|------|
| **프레임워크** | Next.js (App Router) |
| **상태 관리** | React Query (TanStack Query) |
| **UI** | Tailwind CSS + HeroUI |
| **포트** | 3000 |
| **배포** | Vercel 또는 독립 컨테이너 |

### Pages & Routes

| Route | Page | 설명 |
|-------|------|------|
| `/calendar` | CalendarPage | 메인 캘린더 뷰 (월별) |
| `/calendar/records/new` | AddRecordPage | 시술 수동 추가 |
| `/calendar/records/[id]` | RecordDetailPage | 시술 상세/수정/삭제 |
| `/calendar/statistics` | StatisticsPage | 시술 통계 |

### Client-Side Data Flow

```
[Page Component]
  → React Query Hook (useRecords, useSchedules, etc.)
    → fetch() to Calendar API (REST)
      → Response cached by React Query
        → UI re-render
```

### Mutation Flow (생성/수정/삭제)

```
[Form Submit]
  → Client-side validation
    → useMutation() → POST/PUT/DELETE to Calendar API
      → onSuccess: invalidateQueries(['records']) → 자동 refetch
      → onError: 에러 토스트 표시
```

### Google Calendar Integration (Client-Side)

```
[GoogleCalendarButton 클릭]
  → Google Identity Services (OAuth 2.0 팝업)
    → access_token 획득
      → Google Calendar API 직접 호출 (POST /calendars/primary/events)
        → 성공/실패 UI 피드백
```

### Error Handling Strategy

| 상황 | UI 처리 |
|------|---------|
| API 로딩 중 | Skeleton UI (HeroUI Skeleton) |
| API 실패 (네트워크) | 에러 배너 + 재시도 버튼 |
| API 실패 (4xx) | 토스트 알림 + 필드별 에러 메시지 |
| 빈 데이터 | Empty State 컴포넌트 |
| Google OAuth 실패 | 에러 토스트 + 재시도 안내 |

### Caching Strategy (React Query)

| 데이터 | staleTime | 이유 |
|--------|-----------|------|
| records (기간별) | 30s | 자주 변경될 수 있음 |
| schedules | 30s | 기록 변경 시 연동 |
| statistics | 5min | 자주 변경되지 않음 |
| categories | 10min | 관리자만 변경 |
| treatments | 10min | 관리자만 변경 |
| dosageTypes | 10min | 관리자만 변경 |

---

## Service 1: Calendar API Service

| 항목 | 설명 |
|------|------|
| **언어** | Go |
| **아키텍처** | Hexagonal (Ports & Adapters) |
| **API 스타일** | RESTful |
| **포트** | 8080 |
| **배포** | 독립 컨테이너/프로세스 |

### REST API Endpoints

| Method | Path | 설명 | 응답 |
|--------|------|------|------|
| POST | /api/records | 시술 수동 추가 | 201 + Record |
| GET | /api/records?from=&to= | 기간별 시술 목록 | 200 + []Record |
| GET | /api/records/{id} | 시술 단건 조회 | 200 + Record |
| PUT | /api/records/{id} | 시술 수정 | 200 + Record |
| DELETE | /api/records/{id} | 시술 삭제 | 204 |
| GET | /api/schedules | 예정일 목록 | 200 + []Schedule |
| GET | /api/statistics | 시술 통계 | 200 + []Stat |
| GET | /api/treatment-data/categories | 카테고리 목록 (드롭다운용) | 200 |
| GET | /api/treatment-data/categories/{id}/treatments | 시술명 목록 | 200 |
| GET | /api/treatment-data/treatments/{id}/dosage-types | 용량 단위 | 200 |

### Mock Endpoints (개발 환경 전용)

| Method | Path | 설명 |
|--------|------|------|
| POST | /mock/events/reservation-fixed | 예약 확정 이벤트 수동 발행 |

### Orchestration Patterns

**시술 기록 생성 (동기 체인):**
```
Handler.CreateRecord()
  → validate input
  → RecordService.CreateRecord()
      → RecordRepo.Save()
      → CycleRuleClient.GetCycleRule()  [실패 시 graceful skip]
      → ScheduleService.CalculateAndSave()  [주기 있을 때만]
          → ScheduleRepo.Save()
  → return 201
```

**시술 기록 수정 (조건부 재계산):**
```
Handler.UpdateRecord()
  → RecordRepo.FindByID()
  → verify ownership
  → apply changes
  → RecordRepo.Update()
  → IF date or category changed:
      → ScheduleRepo.DeleteByRecordID()
      → CycleRuleClient.GetCycleRule()
      → ScheduleService.CalculateAndSave()
  → return 200
```

**시술 기록 삭제 (cascade):**
```
Handler.DeleteRecord()
  → RecordRepo.FindByID()
  → verify ownership
  → ScheduleRepo.DeleteByRecordID()  [실패 시 로그 후 계속]
  → RecordRepo.Delete()
  → return 204
```

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

## Service 2: Event Consumer Service

| 항목 | 설명 |
|------|------|
| **언어** | Go |
| **배포** | 별도 컨테이너/프로세스 |
| **입력** | AWS SQS (프로덕션) / 인메모리 채널 (개발) |
| **코드 공유** | Calendar API와 같은 Go 모듈 |

### Processing Flow

```
SQS.ReceiveMessage()
  → parse ReservationFixedEvent
  → IF event.TreatmentType == "surgery": ACK & skip
  → idempotency check (reservationID 중복 확인)
  → TreatmentRecordService.CreateRecord(source=AUTO)
  → ACK message
```

### Retry & Error Policy

| 상황 | 처리 |
|------|------|
| 파싱 실패 | DLQ로 이동 (재시도 불가) |
| DB 일시 오류 | 재시도 3회 (exponential backoff) |
| 중복 이벤트 | 스킵 + ACK (멱등성) |
| 3회 재시도 실패 | DLQ로 이동 + 알림 |

### State Transitions (이벤트 처리)

```
[RECEIVED] → 파싱 성공 → [VALIDATED]
[VALIDATED] → surgery 필터 → [SKIPPED] → ACK
[VALIDATED] → treatment → [PROCESSING]
[PROCESSING] → 중복 → [DUPLICATE] → ACK
[PROCESSING] → 저장 성공 → [COMPLETED] → ACK
[PROCESSING] → 저장 실패 → [RETRY] → 3회 초과 → [DLQ]
```

---

## Service 3: Admin API Service

| 항목 | 설명 |
|------|------|
| **언어** | Python (FastAPI) |
| **아키텍처** | Layered CRUD |
| **포트** | 8081 |
| **배포** | 독립 컨테이너/프로세스 |

### REST API Endpoints

| Method | Path | 설명 |
|--------|------|------|
| POST | /api/cycle-rules | 추천 주기 생성 |
| GET | /api/cycle-rules | 전체 주기 목록 |
| GET | /api/cycle-rules/{categoryId} | 카테고리별 주기 조회 |
| PUT | /api/cycle-rules/{categoryId} | 주기 수정 |
| DELETE | /api/cycle-rules/{categoryId} | 주기 삭제 |
| POST | /api/categories | 시술 카테고리 생성 |
| GET | /api/categories | 카테고리 목록 |
| PUT | /api/categories/{id} | 카테고리 수정 |
| DELETE | /api/categories/{id} | 카테고리 삭제 |
| POST | /api/categories/{id}/treatments | 시술명 생성 |
| GET | /api/categories/{id}/treatments | 시술명 목록 |
| PUT | /api/treatments/{id} | 시술명 수정 |
| DELETE | /api/treatments/{id} | 시술명 삭제 |
| POST | /api/treatments/{id}/dosage-types | 용량 단위 생성 |
| GET | /api/treatments/{id}/dosage-types | 용량 단위 목록 |
| DELETE | /api/dosage-types/{id} | 용량 단위 삭제 |

### Cascade Delete Rules

| 삭제 대상 | Cascade 동작 |
|-----------|-------------|
| Category 삭제 | 하위 Treatment + DosageType 모두 삭제 |
| Treatment 삭제 | 하위 DosageType 모두 삭제 |
| CycleRule 삭제 | 기존 예정일에는 영향 없음 (이미 계산된 것은 유지) |

---

## Service 간 통신 상세

### Calendar API → Admin API

| 호출 | 목적 | 실패 시 |
|------|------|---------|
| GET /api/cycle-rules/{categoryId} | 예정일 계산용 주기 조회 | graceful skip (예정일 미생성) |
| GET /api/categories | 드롭다운 데이터 프록시 | 502 반환 |
| GET /api/categories/{id}/treatments | 드롭다운 데이터 프록시 | 502 반환 |
| GET /api/treatments/{id}/dosage-types | 드롭다운 데이터 프록시 | 502 반환 |

### Timeout & Circuit Breaker

| 설정 | 값 |
|------|-----|
| HTTP Timeout | 3초 |
| Circuit Breaker 임계값 | 5회 연속 실패 |
| Circuit Open 유지 시간 | 30초 |
| Retry | 1회 (즉시) |

---

## Reminder Batch Job

| 항목 | 설명 |
|------|------|
| **트리거** | 매일 09:00 (Cron/CloudWatch Events) |
| **실행 위치** | Calendar API 프로세스 내 또는 별도 Lambda |
| **처리 대상** | Status=PENDING, ScheduledDate <= today |

### Batch Flow

```
1. ListDueReminders(today)
2. For each schedule:
   a. Build ReminderMessage
   b. NotificationClient.SendReminder()
   c. IF success: update Status → REMINDED
   d. IF fail: log, skip (next batch retry)
3. Log summary (total, success, failed)
```
