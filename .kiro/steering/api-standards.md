# API Standards

이 문서는 모든 서비스(Go, Python)의 REST API 설계 규칙을 정의합니다.

## URL 규칙

- 리소스명은 **복수형 명사** 사용: `/api/records`, `/api/categories`
- 계층 관계는 중첩 URL: `/api/categories/{id}/treatments`
- 동사는 URL에 사용하지 않음 (예외: `/api/schedules/{id}/complete` 같은 상태 전이 액션)
- URL은 **kebab-case**: `/api/cycle-rules`, `/api/treatment-data`
- 버전 prefix는 이 프로젝트에서는 생략 (단일 버전 MVP)

## HTTP 메서드

| 동작 | 메서드 | 응답 코드 |
|------|--------|-----------|
| 생성 | POST | 201 Created |
| 목록 조회 | GET | 200 OK |
| 단건 조회 | GET | 200 OK |
| 전체 수정 | PUT | 200 OK |
| 부분 수정 | PATCH | 200 OK |
| 삭제 | DELETE | 204 No Content |

## 요청/응답 형식

- Content-Type: `application/json`
- 필드명: **snake_case** (Go에서는 JSON 태그로 변환)
- 날짜/시간: ISO 8601 (`2026-01-15T09:00:00Z`)
- ID: UUIDv7 string (시간 순서 정렬 가능)

## 에러 응답 형식 (통일)

모든 서비스에서 동일한 에러 응답 구조를 사용합니다:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "treatment_date is required"
  }
}
```

| HTTP Status | 사용 상황 |
|-------------|----------|
| 400 | 입력값 검증 실패 |
| 404 | 리소스 미존재 |
| 409 | 중복/충돌 |
| 422 | 요청 형식 오류 (FastAPI 기본) |
| 500 | 서버 내부 오류 (상세 노출 금지) |
| 502 | 외부 서비스 호출 실패 |

## 페이지네이션

목록 API에서 필요 시:
- Query params: `?from=2026-01-01&to=2026-01-31` (기간 기반)
- 또는: `?page=1&size=20` (오프셋 기반)
- 응답에 총 개수 포함 불필요 (MVP)

## Health Check

모든 서비스는 `GET /health` 엔드포인트를 제공합니다:

```json
{
  "status": "ok",
  "service": "calendar-service",
  "version": "0.1.0"
}
```

## CORS

- 개발 환경: `localhost:3000` 허용
- 프로덕션: 화이트리스트 방식
- `allow_origins=["*"]` + `allow_credentials=True` 조합 금지
