# Application Design (Comprehensive) — 시술 관리 캘린더

## Architecture Overview

| 서비스 | 언어 | 아키텍처 | 포트 | 배포 | 책임 |
|--------|------|----------|------|------|------|
| **Calendar Web** | Next.js | App Router + React Query | 3000 | Vercel/독립 컨테이너 | SPA UI, 3단계 드롭다운, 캘린더 뷰, 구글 캘린더 연동 |
| Calendar API | Go | Hexagonal | 8080 | 독립 컨테이너 | 시술 기록 CRUD, 예정일, 통계, 드롭다운 프록시 |
| Event Consumer | Go | Hexagonal (공유) | — | 별도 컨테이너 | 예약 확정 이벤트 → 자동 기입 |
| Admin API | Python/FastAPI | Layered CRUD | 8081 | 독립 컨테이너 | 시술 데이터 + 추천 주기 관리 + AI 시술 정보 예측 (Bedrock) |

**공유 인프라**: PostgreSQL (단일 DB)
**UI 스택**: Tailwind CSS + HeroUI
**AWS AI 서비스**: Amazon Bedrock (Claude Opus) — 시술 카테고리/주기 자동 예측

---

## Key Design Decisions

| # | 결정 | 근거 |
|---|------|------|
| 1 | Hexagonal Architecture | 외부 의존성(SQS, 알림, Admin API) Port/Adapter 분리 → 모킹 용이, 테스트 격리 |
| 2 | Event Consumer 별도 프로세스 | API 서버와 독립 스케일링, 장애 격리, SQS 메시지 보존으로 복구 용이 |
| 3 | 코드베이스 공유 (Go 모듈) | Calendar API + Event Consumer가 동일 도메인 로직 재사용, 중복 제거 |
| 4 | 예정일 별도 테이블 저장 | 조회 성능 최적화, 리마인드 배치 쿼리 단순화, Status 상태 관리 |
| 5 | Circuit Breaker (Admin API 호출) | Admin API 장애 시 Calendar API graceful degradation |
| 6 | 멱등성 (Event Consumer) | 중복 이벤트 안전 처리, DLQ로 파싱 불가 메시지 격리 |

---

## State Transitions

### ScheduledTreatment Status

```
[PENDING] ──리마인드 발송 성공──→ [REMINDED]
[PENDING] ──새 시술 기록 생성──→ [COMPLETED]
[REMINDED] ──새 시술 기록 생성──→ [COMPLETED]
```

### Event Processing Status

```
[RECEIVED] → [VALIDATED] → [PROCESSING] → [COMPLETED] → ACK
                 ↓                ↓
              [SKIPPED]       [RETRY] → [DLQ]
              (surgery)       (3회 초과)
```

---

## Communication Summary

```
[예약 시스템] ──SQS──→ [Event Consumer] ──도메인 로직──→ [PostgreSQL]
                                                              ↑
[Calendar Web] ──REST──→ [Calendar API] ──────────────────→ │
(Next.js)                     │                              │
                              ├──HTTP (3s timeout)──→ [Admin API] ──→ │
                              │                           │
                              └──HTTP (5s timeout)──→ [알림 시스템]

[Calendar Web] ──Google Calendar API──→ [Google Calendar] (클라이언트 직접)

[Admin Client] ──REST──→ [Admin API] ──→ [PostgreSQL]
                              │
                              └──HTTPS (boto3)──→ [AWS Bedrock Claude Opus]
```

---

## Failure Handling Summary

| 장애 | Calendar API 동작 | 사용자 영향 |
|------|-------------------|------------|
| Admin API 다운 | 시술 기록 정상, 예정일 미생성 | 예정일 표시 안됨 (기록은 정상) |
| 알림 시스템 다운 | 배치 스킵, 다음 배치 재시도 | 리마인드 지연 |
| SQS 지연 | Event Consumer 대기 | 자동 기입 지연 |
| PostgreSQL 다운 | 전체 서비스 중단 | 서비스 이용 불가 |

---

## 상세 문서 참조 (서비스별)

| 서비스 | 파일 | 내용 |
|--------|------|------|
| **Calendar Service (Go)** | [calendar-service.md](calendar-service.md) | Hexagonal 구조, 도메인 모델, API 엔드포인트, 메서드 플로우, 에러 처리, Event Consumer |
| **Admin Service (FastAPI)** | [admin-service.md](admin-service.md) | CRUD 구조, 데이터 모델, API 엔드포인트, Cascade 규칙 |
| **Web Frontend (Next.js)** | [web-frontend.md](web-frontend.md) | 페이지/컴포넌트, React Query Hooks, 컴포넌트 플로우, 에러 처리 |

## 기존 문서 (통합 뷰)

- [components.md](components.md) — 전체 컴포넌트 목록
- [component-methods.md](component-methods.md) — 전체 메서드 상세
- [services.md](services.md) — 전체 서비스 정의
- [component-dependency.md](component-dependency.md) — 전체 의존성 매트릭스
