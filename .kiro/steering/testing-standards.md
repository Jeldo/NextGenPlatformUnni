# Testing Standards

이 문서는 모든 서비스의 테스트 전략과 규칙을 정의합니다.

## 공통 원칙

- **Classicist Testing** 방식을 따른다: 실제 객체를 우선 사용하고, 외부 의존성만 모킹한다.
- 모든 기능 구현 시 **테스트 코드를 함께 작성**한다 (같은 PR에 포함).
- 테스트는 구현 세부사항이 아닌 **동작(behavior)**을 검증한다.

## 모킹 규칙

### 모킹 대상 (외부 의존성만)
- 외부 HTTP API 호출 (Admin API, 알림 시스템)
- 메시지 큐 (SQS)
- 외부 인증 (Google OAuth)

### 모킹 금지 (실제 객체 사용)
- 도메인 모델
- 서비스/유즈케이스 레이어
- 리포지토리 (테스트 DB 사용)

## Property-Based Testing (PBT)

순수 함수에 PBT를 적용합니다:
- Go: `rapid` 또는 `gopter` 라이브러리
- Python: `hypothesis` 라이브러리

### PBT 대상
- 날짜 계산 로직 (CycleCalculator)
- 데이터 변환/직렬화
- 검증 로직 (Validate 메서드)

### PBT 속성 예시
```
- Calculate(date, days) 결과는 항상 date보다 미래
- cycleDays > 0이면 결과 != date
- Validate()가 통과한 모델은 항상 저장 가능
```

## Go 테스트

### 구조
```
internal/
  domain/model/
    treatment_record_test.go    # 도메인 모델 단위 테스트
  domain/service/
    cycle_calculator_test.go    # PBT 포함
  application/command/
    create_record_handler_test.go
  infrastructure/postgres/
    record_repository_test.go   # 테스트 DB 사용
```

### 규칙
- 테스트 파일은 `_test.go` suffix
- 테이블 드리븐 테스트 사용
- 테스트 헬퍼는 `testutil` 패키지에 분리
- DB 테스트는 `testcontainers-go` 또는 Docker Compose 테스트 DB 사용
- `t.Parallel()` 적극 활용

### 예시
```go
func TestTreatmentRecord_Validate(t *testing.T) {
    tests := []struct {
        name    string
        record  TreatmentRecord
        wantErr bool
    }{
        {"valid record", validRecord(), false},
        {"missing category", recordWithout("category"), true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.record.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Python (FastAPI) 테스트

### 구조
```
tests/
  test_cycle_rules.py           # 라우트 통합 테스트
  test_treatment_data.py
  test_services/
    test_cycle_rule_service.py  # 서비스 단위 테스트
```

### 규칙
- `pytest` + `httpx.AsyncClient` 사용
- DB는 테스트 전용 트랜잭션 롤백 또는 `testcontainers`
- 외부 의존성은 `app.dependency_overrides`로 교체
- 통합 테스트: 라우트 → 서비스 → DB 1개 흐름 검증
- 단위 테스트: 서비스/리포지토리 집중

### 예시
```python
@pytest.fixture
async def client(app: FastAPI) -> AsyncIterator[AsyncClient]:
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as ac:
        yield ac

async def test_create_cycle_rule(client: AsyncClient):
    response = await client.post("/api/cycle-rules", json={
        "category_id": "cat-1",
        "cycle_days": 90,
    })
    assert response.status_code == 201
```

## Frontend (Next.js) 테스트

### 구조
```
__tests__/
  components/
    WeeklyCalendarGrid.test.tsx
    TreatmentDropdown.test.tsx
  hooks/
    useRecords.test.ts
  pages/
    CalendarPage.test.tsx
```

### 규칙
- React Testing Library + Jest (또는 Vitest)
- 컴포넌트: 렌더링 + 사용자 인터랙션 테스트
- Hook: React Query 훅 단위 테스트 (`@tanstack/react-query` test utils)
- 페이지: 주요 시나리오별 통합 테스트
- API 모킹: `msw` (Mock Service Worker)

### 예시
```tsx
test('RecordCard displays treatment name', () => {
  render(<RecordCard record={mockRecord} />);
  expect(screen.getByText('보톡스 - 사각턱')).toBeInTheDocument();
});
```

## 커버리지 기준

| 레이어 | 최소 커버리지 |
|--------|-------------|
| Domain Model | 90% |
| Application (Handler) | 80% |
| Infrastructure | 70% (통합 테스트) |
| Frontend Component | 80% |
| Frontend Hook | 90% |
