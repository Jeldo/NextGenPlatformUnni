# Technology Stack

## 서비스 구성

| 서비스 | 언어 | 프레임워크 | 포트 |
|--------|------|-----------|------|
| Calendar API | Go | net/http (TBD) | 8080 |
| Event Consumer | Go | (Calendar와 코드 공유) | — |
| Cron Scheduler | Go | (Calendar와 코드 공유) | — |
| Admin API | Python | FastAPI | 8081 |
| Web Frontend | TypeScript | Next.js (App Router) | 3000 |

## 인프라

| 항목 | 선택 |
|------|------|
| 클라우드 | AWS |
| 데이터베이스 | PostgreSQL 16 |
| 메시지 큐 | AWS SQS (프로덕션) / 인메모리 (개발) |
| 푸시 알림 | 기존 시스템 활용 (HTTP 호출) |
| 컨테이너 | Docker + Docker Compose (로컬) |
| 태스크 러너 | Taskfile |

## 프론트엔드 스택

| 항목 | 선택 |
|------|------|
| UI 라이브러리 | HeroUI v3 |
| 스타일링 | Tailwind CSS |
| 상태 관리 | React Query (TanStack Query) |
| Google 연동 | Google Identity Services + Calendar API |
| Color Theme | Primary `#F66336`, Black `#131517`, Gray `#697683` |

## 아키텍처 패턴

| 서비스 | 패턴 |
|--------|------|
| Calendar Service | Hexagonal 4-Layer + CQRS |
| Admin Service | Layered CRUD |
| Frontend | App Router + React Query hooks |

## 서비스 간 통신

| From → To | 방식 |
|-----------|------|
| Frontend → Calendar API | REST (HTTP) |
| Calendar API → Admin API | REST (내부 HTTP, 3s timeout, circuit breaker) |
| Calendar API → 알림 시스템 | REST (HTTP, 5s timeout) |
| 예약 시스템 → Event Consumer | AWS SQS |

## 개발 환경

- Docker Compose로 전체 서비스 로컬 실행
- 환경 변수: `.env.example` 참조
- 명령어: `task up`, `task down`, `task health`

## 상세 문서
- 설계: `aidlc-docs/inception/application-design/`
- Go 규칙: `.kiro/steering/go-conventions.md`
- Python 규칙: `.kiro/steering/python-conventions.md`
- Frontend 규칙: `.kiro/steering/frontend-conventions.md`
