# Functional Design — Unit 1: 환경설정

## Unit Context

Unit 1은 모든 서비스의 프로젝트 초기 셋업과 로컬 개발 환경 구축을 다룹니다.

---

## 프로젝트 구조

### Go Calendar Service (Hexagonal 4-Layer)

```
calendar-service/
├── cmd/
│   ├── api/              # HTTP 서버 엔트리포인트
│   │   └── main.go
│   ├── consumer/         # SQS Event Consumer 엔트리포인트
│   │   └── main.go
│   └── cron/             # CronJob 엔트리포인트
│       └── main.go
├── internal/
│   ├── presentation/     # 1. Presentation Layer
│   │   ├── handler/      # HTTP handlers
│   │   ├── middleware/    # HTTP middleware
│   │   └── dto/          # Request/Response DTOs
│   ├── application/      # 2. Application Layer (CQRS)
│   │   ├── command/      # Commands + Handlers
│   │   └── query/        # Queries + Handlers
│   ├── domain/           # 3. Domain Layer
│   │   ├── model/        # Rich Domain Models
│   │   ├── service/      # Domain Services (CycleCalculator)
│   │   └── port/         # Port interfaces (Repository, Client)
│   └── infrastructure/   # 4. Infrastructure Layer
│       ├── postgres/     # PostgreSQL repositories
│       ├── http/         # HTTP clients (Admin API, Notification)
│       ├── sqs/          # SQS subscriber
│       └── mock/         # Mock adapters (dev)
├── migrations/           # DB 마이그레이션 파일
├── config/               # 설정 로드
├── go.mod
├── go.sum
├── Taskfile.yml
└── Dockerfile
```

### Next.js Frontend

```
web/
├── src/
│   ├── app/                    # App Router
│   │   ├── calendar/
│   │   │   ├── page.tsx        # CalendarPage
│   │   │   └── records/
│   │   │       ├── new/
│   │   │       │   └── page.tsx  # AddRecordPage
│   │   │       └── [id]/
│   │   │           └── page.tsx  # RecordDetailPage
│   │   ├── layout.tsx
│   │   └── globals.css
│   ├── components/
│   │   ├── calendar/           # WeeklyCalendarGrid, DateBottomSheet
│   │   ├── treatment/          # TreatmentDropdown, DosageInput, BottomSheets
│   │   ├── record/             # RecordCard, ScheduleCard
│   │   ├── common/             # FloatingAddButton, HospitalInput
│   │   └── modal/              # ScheduleConfirmModal
│   ├── hooks/                  # React Query hooks
│   ├── lib/                    # API client, utils
│   └── types/                  # TypeScript 타입 정의
├── public/
├── tailwind.config.ts
├── next.config.ts
├── package.json
├── tsconfig.json
└── Dockerfile
```

### FastAPI Admin Service

```
admin-service/
├── app/
│   ├── main.py               # FastAPI app 엔트리포인트
│   ├── config.py             # 설정
│   ├── routers/
│   │   ├── cycle_rules.py
│   │   └── treatment_data.py
│   ├── services/
│   │   ├── cycle_rule_service.py
│   │   └── treatment_data_service.py
│   ├── repositories/
│   │   ├── cycle_rule_repository.py
│   │   └── treatment_data_repository.py
│   ├── models/               # SQLAlchemy models
│   │   └── models.py
│   └── schemas/              # Pydantic schemas
│       └── schemas.py
├── migrations/               # Alembic
│   └── versions/
├── alembic.ini
├── pyproject.toml
├── Dockerfile
└── tests/
```

### 공통 인프라

```
(project root)/
├── docker-compose.yml        # PostgreSQL + 모든 서비스
├── .env.example
├── Taskfile.yml              # 전체 프로젝트 명령어
├── calendar-service/
├── web/
├── admin-service/
├── docs/
└── aidlc-docs/
```

---

## Docker Compose 구성

| 서비스 | 이미지 | 포트 | 역할 |
|--------|--------|------|------|
| postgres | postgres:16 | 5432 | 공유 DB |
| calendar-api | calendar-service (빌드) | 8080 | Calendar REST API |
| calendar-consumer | calendar-service (빌드) | — | Event Consumer |
| calendar-cron | calendar-service (빌드) | — | Reminder CronJob |
| admin-api | admin-service (빌드) | 8081 | Admin REST API |
| web | web (빌드) | 3000 | Next.js Frontend |

---

## 환경 변수

```env
# PostgreSQL
DATABASE_URL=postgres://user:password@localhost:5432/treatment_calendar

# Calendar Service
CALENDAR_PORT=8080
ADMIN_API_URL=http://admin-api:8081
NOTIFICATION_API_URL=http://localhost:9090  # mock
SQS_QUEUE_URL=mock://local
ENV=development

# Admin Service
ADMIN_PORT=8081
DATABASE_URL=postgres://user:password@localhost:5432/treatment_calendar

# Frontend
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_GOOGLE_CLIENT_ID=placeholder
```

---

## Health Check 엔드포인트

모든 서비스에 `GET /health` 구현:

```json
{
  "status": "ok",
  "service": "calendar-service",
  "version": "0.1.0"
}
```

---

## Tailwind + HeroUI Color Theme 설정

```typescript
// tailwind.config.ts
colors: {
  primary: {
    DEFAULT: '#F66336',
  },
  black: '#131517',
  gray: {
    description: '#697683',
  },
}
```

---

## 완료 기준

- [ ] `docker-compose up` 으로 모든 서비스 실행 가능
- [ ] 각 서비스 `GET /health` 응답 확인
- [ ] Frontend `http://localhost:3000` 접속 가능 (빈 페이지)
- [ ] PostgreSQL 연결 확인
- [ ] Color Theme 적용 확인 (Primary 색상 표시)
