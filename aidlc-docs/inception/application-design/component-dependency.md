# Component Dependencies (Comprehensive) — 시술 관리 캘린더

## Frontend — Dependencies

### Calendar Web → External

| Component | Depends On | 프로토콜 | 비고 |
|-----------|-----------|----------|------|
| React Query Hooks | Calendar API | HTTP (REST) | 모든 데이터 조회/변경 |
| GoogleCalendarButton | Google Calendar API | HTTP (OAuth 2.0) | 클라이언트 직접 호출 |
| HospitalInput | 로컬 캐시 (useRecords) | — | 기존 기록에서 병원명 추출 |

### Component 간 의존성

| Component | Depends On | 관계 |
|-----------|-----------|------|
| CalendarPage | WeeklyCalendarGrid, RecordCard, ScheduleCard | 렌더링 |
| WeeklyCalendarGrid | useRecords, useSchedules | 데이터 |
| AddRecordPage | TreatmentDropdown, HospitalInput, DatePicker | 폼 구성 |
| RecordDetailPage | useRecord, TreatmentDropdown, GoogleCalendarButton | 상세/수정 |
| TreatmentDropdown | CategoryBottomSheet, TreatmentBottomSheet, DosageInput | 계층 구조 |
| TreatmentDropdown | useCategories, useTreatments, useDosageTypes | 데이터 |
| StatisticsPage | useStatistics | 데이터 |

---

## Dependency Matrix

### Calendar API — Internal Dependencies

| Component | Depends On | 의존 방향 | 통신 방식 |
|-----------|-----------|-----------|-----------|
| HTTPHandler | TreatmentRecordService | Inbound → Domain | 직접 호출 |
| HTTPHandler | ScheduledTreatmentService | Inbound → Domain | 직접 호출 |
| HTTPHandler | TreatmentStatisticsService | Inbound → Domain | 직접 호출 |
| TreatmentRecordService | TreatmentRecordRepository (port) | Domain → Outbound | 인터페이스 |
| TreatmentRecordService | ScheduledTreatmentService | Domain → Domain | 직접 호출 |
| TreatmentRecordService | CycleRuleClient (port) | Domain → Outbound | 인터페이스 |
| ScheduledTreatmentService | ScheduledTreatmentRepository (port) | Domain → Outbound | 인터페이스 |
| ScheduledTreatmentService | CycleRuleClient (port) | Domain → Outbound | 인터페이스 |
| ReminderService | ScheduledTreatmentService | Domain → Domain | 직접 호출 |
| ReminderService | NotificationClient (port) | Domain → Outbound | 인터페이스 |

### Calendar API — External Dependencies

| Adapter | 외부 시스템 | 프로토콜 | Timeout | Retry |
|---------|-----------|----------|---------|-------|
| PostgresRecordRepo | PostgreSQL | TCP/SQL | 5s | 0 |
| PostgresScheduleRepo | PostgreSQL | TCP/SQL | 5s | 0 |
| HTTPCycleRuleClient | Admin API | HTTP | 3s | 1회 |
| HTTPNotificationClient | 알림 시스템 | HTTP | 5s | 1회 |

### Event Consumer — Dependencies

| Component | Depends On | 통신 방식 |
|-----------|-----------|-----------|
| EventConsumer | SQSEventSubscriber / InMemorySubscriber | 인터페이스 (Port) |
| ReservationEventHandler | TreatmentRecordService | 도메인 코드 공유 |
| SQSEventSubscriber | AWS SQS | AWS SDK |

### Admin API — Internal Dependencies

| Component | Depends On | 통신 방식 |
|-----------|-----------|-----------|
| Router | CycleRuleService | 직접 호출 |
| Router | TreatmentDataService | 직접 호출 |
| Router | AISuggestService | 직접 호출 |
| CycleRuleService | CycleRuleRepository | 직접 호출 |
| TreatmentDataService | TreatmentDataRepository | 직접 호출 |
| CycleRuleRepository | PostgreSQL | SQLAlchemy |
| TreatmentDataRepository | PostgreSQL | SQLAlchemy |
| AISuggestService | AWS Bedrock | boto3 (HTTP) |

### Admin API — External Dependencies

| Adapter | 외부 시스템 | 프로토콜 | Timeout | 비고 |
|---------|-----------|----------|---------|------|
| AISuggestService | AWS Bedrock (Claude Opus) | HTTPS (boto3) | 30s (default) | us-east-1, cross-region inference |

---

## Dependency Direction (Hexagonal 원칙)

```
┌─────────────────────────────────────────────────────────┐
│                    ADAPTERS (외부)                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │HTTP      │  │Postgres  │  │SQS       │             │
│  │Handler   │  │Repos     │  │Subscriber│             │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘             │
│       │              │              │                    │
├───────┼──────────────┼──────────────┼────────────────────┤
│       v              v              v    PORTS (인터페이스)│
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │Treatment │  │Record    │  │Event     │             │
│  │API Port  │  │Repo Port │  │Sub Port  │             │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘             │
│       │              │              │                    │
├───────┼──────────────┼──────────────┼────────────────────┤
│       v              ^              ^    DOMAIN (핵심)    │
│  ┌─────────────────────────────────────────────┐        │
│  │  TreatmentRecordService                     │        │
│  │  ScheduledTreatmentService                  │        │
│  │  TreatmentStatisticsService                 │        │
│  │  ReminderService                            │        │
│  └─────────────────────────────────────────────┘        │
└─────────────────────────────────────────────────────────┘

화살표 방향:
  Adapter → Port: Adapter가 Port를 구현
  Domain → Port: Domain이 Port 인터페이스를 사용 (의존성 역전)
  Inbound Adapter → Domain: Handler가 Service를 호출
```

---

## 환경별 Adapter 교체

| Port | Production Adapter | Development Adapter |
|------|--------------------|---------------------|
| EventSubscriber | SQSEventSubscriber | InMemoryEventSubscriber |
| CycleRuleClient | HTTPCycleRuleClient | HTTPCycleRuleClient (동일, Admin API 로컬 실행) |
| NotificationClient | HTTPNotificationClient | MockNotificationClient (로그만 출력) |
| TreatmentRecordRepository | PostgresRecordRepo | PostgresRecordRepo (동일, 로컬 DB) |

---

## Failure Isolation

| 장애 지점 | 영향 범위 | 대응 |
|-----------|----------|------|
| Admin API 다운 | 예정일 미생성 (시술 기록은 정상) | Circuit Breaker, graceful degradation |
| 알림 시스템 다운 | 리마인드 미발송 | 다음 배치에서 재시도 |
| SQS 지연 | 자동 기입 지연 | 메시지 보존 (SQS 기본 4일) |
| PostgreSQL 다운 | 전체 서비스 중단 | Health check + 자동 복구 |
| Event Consumer 다운 | 자동 기입 중단 | SQS 메시지 보존, 복구 후 처리 |
