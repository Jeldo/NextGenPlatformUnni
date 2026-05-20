# Go Code Conventions

이 문서는 Go 서비스의 코딩 규칙을 정의합니다. Google Go Style Guide를 기반으로 합니다.

## 프로젝트 구조 (Hexagonal 4-Layer)

```
cmd/
  api/main.go           # HTTP 서버 엔트리포인트
  consumer/main.go      # 이벤트 컨슈머 엔트리포인트
  cron/main.go          # 배치 작업 엔트리포인트
internal/
  presentation/         # Inbound Adapters (HTTP handlers, DTO)
  application/          # Use Cases (CQRS: command/, query/)
  domain/               # Models, Services, Ports
  infrastructure/       # Outbound Adapters (DB, HTTP clients)
config/                 # 설정 로드
migrations/             # DB 마이그레이션
```

## 네이밍

### 패키지
- 소문자, 단수형: `model`, `handler`, `command`
- 패키지명과 디렉토리명 일치
- `util`, `common`, `helper` 같은 범용 이름 금지

### 타입/인터페이스
- 인터페이스: 동사 + er 접미사 (`Reader`, `Validator`) 또는 역할 명사 (`Repository`, `Client`)
- 구조체: 명사 (`TreatmentRecord`, `CreateRecordHandler`)
- 인터페이스는 사용하는 쪽에서 정의 (Port 패키지 예외)

### 함수/메서드
- 동사로 시작: `Validate()`, `Calculate()`, `Save()`
- Getter에 `Get` prefix 불필요: `record.ID()` (not `record.GetID()`)
- Bool 반환: `Is`, `Has`, `Can` prefix: `IsOwnedBy()`, `IsSurgery()`

### 변수
- camelCase
- 짧은 스코프에서는 짧은 이름 허용: `r` (record), `s` (schedule), `ctx` (context)
- 긴 스코프에서는 설명적 이름: `treatmentRecord`, `scheduledDate`

### 상수/Enum
- typed const 사용 (string 리터럴 금지)
```go
type RecordSource string
const (
    SourceAuto   RecordSource = "AUTO"
    SourceManual RecordSource = "MANUAL"
)
```

## 에러 처리

- 에러는 즉시 처리하거나 반환. 무시 금지.
- 커스텀 에러 타입 사용:
```go
type AppError struct {
    Code    ErrorCode
    Message string
    Cause   error
}
```
- `errors.Is()`, `errors.As()`로 에러 체인 검사
- 에러 메시지는 소문자로 시작, 마침표 없음
- sentinel 에러보다 typed 에러 선호

## Context

- 모든 public 함수의 첫 번째 인자는 `context.Context`
- context에 비즈니스 데이터 저장 금지 (request ID, trace ID만 허용)
- context 취소를 존중: 장기 작업에서 `ctx.Done()` 체크

## 의존성 주입

- 생성자 함수로 의존성 주입:
```go
func NewCreateRecordHandler(
    repo port.TreatmentRecordRepository,
    cycleClient port.CycleRuleClient,
    scheduleHandler *CalculateScheduleHandler,
) *CreateRecordHandler {
    return &CreateRecordHandler{repo: repo, cycleClient: cycleClient, scheduleHandler: scheduleHandler}
}
```
- DI 프레임워크 사용하지 않음 (수동 와이어링)
- 인터페이스는 Port 패키지에 정의, 구현은 Infrastructure에 위치

## Rich Domain Model 규칙

- 도메인 모델에 비즈니스 로직 포함 (Anemic Model 금지)
- 검증은 모델 자체에서: `record.Validate()`
- 상태 전이는 모델 메서드로: `schedule.MarkReminded()`
- 불변식(invariant)은 모델이 보장

## CQRS 패턴

- Command: 상태 변경 (Create, Update, Delete)
- Query: 읽기 전용 (List, Get)
- **네이밍 규칙**:
  - Command struct: `XxxCommand` (예: `CreateRecordCommand`)
  - Command handler: `XxxCommandHandler` (예: `CreateRecordCommandHandler`)
  - Query struct: `XxxQuery` (예: `ListRecordsQuery`)
  - Query handler: `XxxQueryHandler` (예: `ListRecordsQueryHandler`)
- Handle 메서드는 항상 Command/Query struct를 받음 (primitive 파라미터 금지)
- Handler는 하나의 Command/Query만 처리
- Handler 간 직접 호출 허용 (같은 레이어)

## HTTP Handler 규칙

- Handler는 얇게 유지: 파싱 → 핸들러 호출 → 응답 변환
- 비즈니스 로직을 Handler에 넣지 않음
- Request/Response DTO는 `presentation/dto/` 패키지에 정의
- JSON 태그 필수: `json:"field_name"`
- validate 태그로 필수 필드 표시: `validate:"required"`

## 동시성

- goroutine 시작 시 반드시 종료 메커니즘 확보 (context 취소, done 채널)
- 공유 상태 접근 시 mutex 또는 채널 사용
- `sync.WaitGroup`으로 goroutine 완료 대기

## 로깅

- 구조화 로깅 사용 (slog 또는 zerolog)
- 로그 레벨: DEBUG, INFO, WARN, ERROR
- 에러 로그에는 context 정보 포함 (userID, recordID 등)
- 민감 정보 로깅 금지

## 안티 패턴 (피할 것)

- `init()` 함수에서 부작용 발생
- 글로벌 변수로 상태 관리
- `interface{}` / `any` 남용
- panic으로 에러 처리
- 테스트에서 `time.Sleep()` 사용
- 순환 import
