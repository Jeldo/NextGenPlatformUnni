# Components — 시술 관리 캘린더

## Frontend: Calendar Web App (Next.js)

| 항목 | 설명 |
|------|------|
| **프레임워크** | Next.js (App Router) |
| **상태 관리** | React Query (TanStack Query) — 서버 상태 |
| **UI** | Tailwind CSS + HeroUI |
| **Google 연동** | Google Calendar API (클라이언트 사이드) |

### Page Components

| Component | 경로 | 책임 |
|-----------|------|------|
| **CalendarPage** | `/calendar` | 월별 캘린더 뷰, 시술 기록 + 예정일 표시 |
| **RecordDetailPage** | `/calendar/records/[id]` | 시술 기록 상세 조회/수정/삭제 |
| **AddRecordPage** | `/calendar/records/new` | 시술 수동 추가 |
| **StatisticsPage** | `/calendar/statistics` | 시술 통계 |

### Feature Components

| Component | 책임 |
|-----------|------|
| **TreatmentDropdown** | 3단계 드롭다운 (카테고리 → 시술명 → 용량) |
| **CategoryBottomSheet** | 카테고리 선택 바텀시트 |
| **TreatmentBottomSheet** | 시술명 선택 바텀시트 |
| **DosageInput** | 용량 숫자 인풋 + 단위 드롭다운 |
| **CalendarGrid** | 월별 캘린더 그리드 (확정 일정 + 예정일 점선 구분) |
| **RecordCard** | 시술 기록 카드 (목록 아이템) |
| **ScheduleCard** | 예정일 카드 (점선 스타일) |
| **HospitalInput** | 병원명 자유 텍스트 + 자동완성 |
| **GoogleCalendarButton** | 구글 캘린더 내보내기 버튼 (OAuth 트리거) |
| **ReminderBanner** | 리마인드 알림 배너 (앱 내 표시) |

### Hooks (React Query)

| Hook | API 호출 | 캐싱 전략 |
|------|----------|-----------|
| `useRecords(from, to)` | GET /api/records | staleTime: 30s |
| `useRecord(id)` | GET /api/records/{id} | staleTime: 60s |
| `useSchedules()` | GET /api/schedules | staleTime: 30s |
| `useStatistics()` | GET /api/statistics | staleTime: 5min |
| `useCategories()` | GET /api/treatment-data/categories | staleTime: 10min |
| `useTreatments(categoryId)` | GET /api/treatment-data/categories/{id}/treatments | staleTime: 10min |
| `useDosageTypes(treatmentId)` | GET /api/treatment-data/treatments/{id}/dosage-types | staleTime: 10min |
| `useCreateRecord()` | POST /api/records | mutation, invalidates records |
| `useUpdateRecord()` | PUT /api/records/{id} | mutation, invalidates records |
| `useDeleteRecord()` | DELETE /api/records/{id} | mutation, invalidates records |

---

## Service 1: Calendar Service (Go, Hexagonal Architecture)

### Domain Layer (Core)

| Component | 책임 |
|-----------|------|
| **TreatmentRecord** | 시술 기록 도메인 모델 및 비즈니스 규칙 |
| **ScheduledTreatment** | 다음 시술 예정일 도메인 모델 |
| **TreatmentStatistics** | 시술 통계 집계 로직 |

### Port Layer (Interfaces)

| Port | 방향 | 책임 |
|------|------|------|
| **TreatmentRecordRepository** | Outbound | 시술 기록 영속성 |
| **ScheduledTreatmentRepository** | Outbound | 예정일 영속성 |
| **CycleRuleClient** | Outbound | 추천 주기 데이터 조회 (FastAPI 호출) |
| **NotificationClient** | Outbound | 푸시 알림 발송 요청 |
| **EventSubscriber** | Inbound | 예약 확정 이벤트 수신 |
| **TreatmentAPI** | Inbound | REST API 요청 처리 |

### Adapter Layer (Implementations)

| Adapter | Port 구현 | 기술 |
|---------|-----------|------|
| **PostgresRecordRepo** | TreatmentRecordRepository | PostgreSQL |
| **PostgresScheduleRepo** | ScheduledTreatmentRepository | PostgreSQL |
| **HTTPCycleRuleClient** | CycleRuleClient | HTTP (FastAPI 호출) |
| **HTTPNotificationClient** | NotificationClient | HTTP (기존 알림 시스템) |
| **SQSEventSubscriber** | EventSubscriber | AWS SQS (프로덕션) |
| **InMemoryEventSubscriber** | EventSubscriber | 인메모리 채널 (개발) |
| **MockEventHTTPAdapter** | — | Mock HTTP 엔드포인트 (개발) |
| **HTTPHandler** | TreatmentAPI | net/http 또는 chi/echo |

---

## Service 2: Admin Service (FastAPI, Layered CRUD)

| Component | 책임 |
|-----------|------|
| **Router** | API 엔드포인트 라우팅 |
| **CycleRuleService** | 추천 주기 CRUD 비즈니스 로직 |
| **TreatmentDataService** | 시술 마스터 데이터 CRUD |
| **CycleRuleRepository** | 추천 주기 DB 접근 |
| **TreatmentDataRepository** | 시술 데이터 DB 접근 |

---

## Service 3: Event Consumer (Go, 별도 프로세스)

| Component | 책임 |
|-----------|------|
| **EventConsumer** | SQS 메시지 폴링 및 처리 |
| **ReservationEventHandler** | 예약 확정 이벤트 → 시술 기록 생성 |

> Event Consumer는 Calendar Service와 동일한 도메인/포트를 공유하되, 별도 바이너리로 배포됨.
