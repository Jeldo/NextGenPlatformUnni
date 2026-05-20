# Units of Work — 시술 관리 캘린더

## Unit 구조 (순차 → 병렬)

```
Unit 1: 환경설정
    ↓
Unit 2: 인터페이스 확정 + Dummy 통신 검증
    ↓
Unit 3: 병렬 개발
  ├── Unit 3-A: Frontend (Next.js)
  ├── Unit 3-B: Backend (Go Calendar Service + Event Consumer + CronJob)
  └── Unit 3-C: Admin Service (FastAPI) + 시술 데이터
```

---

## Unit 1: 환경설정

> 프론트엔드 + 백엔드 프로젝트 초기 셋업을 동시에 진행

| 항목 | 내용 |
|------|------|
| **목표** | 모든 서비스의 프로젝트 구조 생성, 로컬 개발 환경 구축 |
| **선행 조건** | 없음 |
| **완료 기준** | 각 서비스가 로컬에서 빈 상태로 실행 가능 |

### 포함 작업:
- [ ] Go 프로젝트 초기화 (go mod, 4-layer 디렉토리 구조)
- [ ] Next.js 프로젝트 초기화 (App Router, Tailwind, HeroUI, React Query)
- [ ] FastAPI 프로젝트 초기화 (pyproject.toml, SQLAlchemy)
- [ ] PostgreSQL 로컬 환경 (Docker Compose)
- [ ] 공통 설정: .env, docker-compose.yml, Makefile/Taskfile
- [ ] Color Theme 설정 (#F66336, #131517, #697683)
- [ ] 각 서비스 Health Check 엔드포인트 (`GET /health`)

---

## Unit 2: 인터페이스 확정 + Dummy 통신 검증

> API 계약을 확정하고, dummy 데이터로 프론트↔백엔드 전체 통신 검증

| 항목 | 내용 |
|------|------|
| **목표** | 모든 API 엔드포인트 스텁 구현 + 프론트에서 호출 확인 |
| **선행 조건** | Unit 1 완료 |
| **완료 기준** | 프론트에서 모든 API 호출이 dummy 응답을 받아 UI에 표시됨 |

### 포함 작업:
- [ ] Calendar Service: 모든 REST API 스텁 구현 (하드코딩 dummy 응답)
  - POST/GET/PUT/DELETE /api/records
  - GET /api/schedules, GET /api/schedules/{id}
  - PATCH /api/schedules/{id}/complete, DELETE /api/schedules/{id}
  - GET /api/statistics
  - GET /api/treatment-data/* (프록시 스텁)
  - POST /mock/events/reservation-fixed
- [ ] Admin Service: 모든 REST API 스텁 구현 (dummy 응답)
  - /api/cycle-rules CRUD
  - /api/categories, /api/treatments, /api/dosage-types CRUD
- [ ] Frontend: 모든 React Query Hook 연결 + 기본 페이지 레이아웃
  - CalendarPage (WeeklyCalendarGrid 스켈레톤 + TreatmentStats)
  - AddRecordPage (폼 레이아웃)
  - RecordDetailPage (상세 레이아웃)
- [ ] 전체 통신 E2E 확인 (프론트 → Calendar API → Admin API)
- [ ] Request/Response DTO 확정 (JSON 스키마)

---

## Unit 3: 병렬 개발

> 인터페이스가 확정되었으므로 프론트/백엔드/데이터를 독립적으로 개발

### Unit 3-A: Frontend (Next.js)

| 항목 | 내용 |
|------|------|
| **목표** | 모든 UI 컴포넌트 + 인터랙션 구현 |
| **선행 조건** | Unit 2 완료 (API 계약 확정) |
| **완료 기준** | 모든 유저 스토리의 프론트엔드 동작 완료 |

#### 포함 작업:
- [ ] WeeklyCalendarGrid (주간 뷰, 카드 형태, 스와이프)
- [ ] DateBottomSheet (날짜별 시술 목록)
- [ ] ScheduleConfirmModal ([삭제] [받았어요])
- [ ] RecordCard / ScheduleCard (실선/점선)
- [ ] TreatmentDropdown (3단계 + 바텀시트)
- [ ] DosageInput (숫자 + 단위 드롭다운)
- [ ] AddRecordPage (DatePicker 바텀시트 월캘린더 + 폼)
- [ ] RecordDetailPage (상세 + 수정 + 삭제)
- [ ] GoogleCalendarButton (OAuth + "(예정)" prefix)
- [ ] TreatmentStats (통계 인라인)
- [ ] FloatingAddButton
- [ ] Error Handling (Skeleton, Toast, Empty State)

### Unit 3-B: Backend (Go Calendar Service)

| 항목 | 내용 |
|------|------|
| **목표** | Calendar Service 전체 비즈니스 로직 구현 |
| **선행 조건** | Unit 2 완료 (API 계약 확정) |
| **완료 기준** | 모든 Command/Query Handler 동작 + 테스트 통과 |

#### 포함 작업:
- [ ] Domain Layer: TreatmentRecord, ScheduledTreatment (Rich Model + Validate, ApplyUpdate 등)
- [ ] Domain Layer: CycleCalculator (+ PBT 테스트)
- [ ] Application Layer: CreateRecordCommandHandler, UpdateRecordCommandHandler, DeleteRecordCommandHandler
- [ ] Application Layer: HandleReservationFixedCommandHandler (멱등성)
- [ ] Application Layer: CalculateScheduleCommandHandler, CompleteScheduleCommandHandler, DeleteScheduleCommandHandler
- [ ] Application Layer: ProcessRemindersCommandHandler (예약 유무 분기)
- [ ] Application Layer: Query Handlers (List/Get Records, Schedules, Statistics)
- [ ] Infrastructure Layer: PostgreSQL Repositories
- [ ] Infrastructure Layer: HTTPCycleRuleClient (Admin API 호출 + Circuit Breaker)
- [ ] Infrastructure Layer: HTTPNotificationClient / MockNotificationClient
- [ ] Infrastructure Layer: SQSEventSubscriber / InMemoryEventSubscriber
- [ ] Presentation Layer: HTTP Handlers (Request 파싱 → Handler 호출 → Response)
- [ ] Presentation Layer: SQS Event Consumer (별도 프로세스)
- [ ] Presentation Layer: CronScheduler (별도 프로세스)
- [ ] Mock 엔드포인트 (/mock/events/reservation-fixed)
- [ ] DB 마이그레이션 (treatment_records, scheduled_treatments 테이블)
- [ ] 단위 테스트 (Classicist + PBT)

### Unit 3-C: Admin Service (FastAPI) + 시술 데이터

| 항목 | 내용 |
|------|------|
| **목표** | 관리자 API 구현 + 초기 시술 마스터 데이터 구축 |
| **선행 조건** | Unit 2 완료 (API 계약 확정) |
| **완료 기준** | 모든 CRUD 동작 + 시술 데이터 시딩 완료 |

#### 포함 작업:
- [x] CycleRuleService + Repository (CRUD + 검증)
- [x] TreatmentDataService + Repository (Category/Treatment/DosageType CRUD)
- [x] Cascade Delete 구현
- [x] DB 마이그레이션 (categories, treatments, dosage_types, cycle_rules 테이블)
- [x] 초기 시술 마스터 데이터 시딩 (보톡스, 필러, 레이저 등 카테고리 + 시술명 + 단위)
- [x] 초기 추천 주기 데이터 시딩 (카테고리별 기본 주기)
- [x] AI Suggest 엔드포인트 (AWS Bedrock Claude Opus — 시술 카테고리/주기/단위 예측)
- [ ] 단위 테스트

---

## Unit 의존성

```
[Unit 1] ──→ [Unit 2] ──→ [Unit 3-A]
                      ──→ [Unit 3-B]
                      ──→ [Unit 3-C]
```

- Unit 1 → Unit 2: 순차 (환경 필요)
- Unit 2 → Unit 3-*: 순차 (인터페이스 확정 필요)
- Unit 3-A, 3-B, 3-C: **병렬** (독립 개발 가능)
