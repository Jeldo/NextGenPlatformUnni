# Project Structure

## 모노레포 구조

```
NextGenPlatformUnni/
├── calendar-service/         # Go Calendar API + Event Consumer + Cron
├── admin-service/            # Python FastAPI Admin API
├── web/                      # Next.js Frontend
├── docker-compose.yml        # 로컬 개발 환경
├── Taskfile.yml              # 프로젝트 명령어
├── .env.example              # 환경 변수 템플릿
├── docs/                     # PRD 등 기획 문서
└── aidlc-docs/               # AI-DLC 설계 산출물 (코드 아님)
```

## Calendar Service (Go)

```
calendar-service/
├── cmd/
│   ├── api/main.go           # HTTP 서버
│   ├── consumer/main.go      # SQS Event Consumer
│   └── cron/main.go          # Reminder CronJob
├── internal/
│   ├── presentation/         # HTTP handlers, DTOs
│   ├── application/          # CQRS (command/, query/)
│   ├── domain/               # Rich models, services, ports
│   └── infrastructure/       # DB repos, HTTP clients, SQS
├── migrations/
├── config/
├── go.mod
└── Dockerfile
```

## Admin Service (Python)

```
admin-service/
├── app/
│   ├── main.py
│   ├── config.py
│   ├── routers/              # cycle_rules.py, treatment_data.py
│   ├── services/
│   ├── repositories/
│   ├── models/               # SQLAlchemy
│   └── schemas/              # Pydantic
├── migrations/               # Alembic
├── pyproject.toml
└── Dockerfile
```

## Web Frontend (Next.js)

```
web/
├── src/
│   ├── app/                  # App Router pages
│   ├── components/           # UI 컴포넌트 (flat)
│   ├── hooks/                # React Query hooks (파일당 1개)
│   ├── lib/                  # API client, utils
│   └── types/                # TypeScript 타입
├── tailwind.config.ts
├── package.json
└── Dockerfile
```

## 규칙

- **애플리케이션 코드**: 프로젝트 루트의 각 서비스 디렉토리에 위치
- **설계 문서**: `aidlc-docs/` 에만 위치 (코드와 분리)
- **Steering 파일**: `.kiro/steering/` (Kiro 컨텍스트)
- **환경 설정**: `.env.example` (커밋), `.env` (gitignore)
