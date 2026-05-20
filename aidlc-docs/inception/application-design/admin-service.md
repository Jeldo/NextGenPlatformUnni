# Admin Service (Python/FastAPI) — Comprehensive Design

> Layered CRUD | REST API | Port 8081

---

## Architecture

```
Router (FastAPI)
  ├── CycleRuleService → CycleRuleRepository → PostgreSQL
  └── TreatmentDataService → TreatmentDataRepository → PostgreSQL
```

---

## Components

| Component | 책임 |
|-----------|------|
| **Router** | API 엔드포인트 라우팅, 입력 검증 (Pydantic) |
| **CycleRuleService** | 추천 주기 CRUD 비즈니스 로직 |
| **TreatmentDataService** | 시술 마스터 데이터 CRUD (카테고리/시술명/용량유형) |
| **CycleRuleRepository** | 추천 주기 DB 접근 (SQLAlchemy) |
| **TreatmentDataRepository** | 시술 데이터 DB 접근 (SQLAlchemy) |

---

## Data Models

```python
# Pydantic Schemas (API)
class CycleRuleCreate(BaseModel):
    category_id: str
    cycle_days: int  # > 0
    description: str | None = None

class CycleRuleResponse(BaseModel):
    category_id: str
    cycle_days: int
    description: str | None
    updated_at: datetime

class CategoryCreate(BaseModel):
    name: str

class CategoryResponse(BaseModel):
    id: str
    name: str

class TreatmentCreate(BaseModel):
    name: str

class TreatmentResponse(BaseModel):
    id: str
    category_id: str
    name: str

class DosageTypeCreate(BaseModel):
    unit: str   # "shot", "minute", "volume", "vial", "joule"

class DosageTypeResponse(BaseModel):
    id: str
    treatment_id: str
    unit: str
```

---

## REST API Endpoints

### Cycle Rules (추천 주기)

| Method | Path | 설명 | 응답 |
|--------|------|------|------|
| POST | /api/cycle-rules | 추천 주기 생성 | 201 |
| GET | /api/cycle-rules | 전체 주기 목록 | 200 |
| GET | /api/cycle-rules/{categoryId} | 카테고리별 주기 조회 | 200 |
| PUT | /api/cycle-rules/{categoryId} | 주기 수정 | 200 |
| DELETE | /api/cycle-rules/{categoryId} | 주기 삭제 | 204 |

### Treatment Data (시술 마스터 데이터)

| Method | Path | 설명 | 응답 |
|--------|------|------|------|
| POST | /api/categories | 카테고리 생성 | 201 |
| GET | /api/categories | 카테고리 목록 | 200 |
| PUT | /api/categories/{id} | 카테고리 수정 | 200 |
| DELETE | /api/categories/{id} | 카테고리 삭제 (cascade) | 204 |
| POST | /api/categories/{id}/treatments | 시술명 생성 | 201 |
| GET | /api/categories/{id}/treatments | 시술명 목록 | 200 |
| PUT | /api/treatments/{id} | 시술명 수정 | 200 |
| DELETE | /api/treatments/{id} | 시술명 삭제 (cascade) | 204 |
| POST | /api/treatments/{id}/dosage-types | 용량 단위 생성 | 201 |
| GET | /api/treatments/{id}/dosage-types | 용량 단위 목록 | 200 |
| DELETE | /api/dosage-types/{id} | 용량 단위 삭제 | 204 |

---

## Method Flows

### CycleRuleService.create_cycle_rule

```
1. 입력 검증 (cycle_days > 0)
2. category_id로 카테고리 존재 확인
3. 기존 규칙 존재 확인 → 존재 시 409 Conflict
4. DB 저장
5. CycleRuleResponse 반환

에러: category 미존재→404, 이미 존재→409, cycle_days<=0→422
```

### CycleRuleService.update_cycle_rule

```
1. category_id로 기존 규칙 조회
2. 미존재 시 404
3. 필드 업데이트 → DB 저장
4. CycleRuleResponse 반환
```

### CycleRuleService.delete_cycle_rule

```
1. category_id로 기존 규칙 조회
2. 미존재 시 404
3. DB 삭제
4. 기존에 계산된 예정일에는 영향 없음 (이미 저장된 것은 유지)
```

### TreatmentDataService.create_category

```
1. 이름 중복 확인 → 중복 시 409
2. DB 저장
3. CategoryResponse 반환
```

### TreatmentDataService.create_treatment

```
1. category_id로 부모 카테고리 존재 확인 → 미존재 시 404
2. 같은 카테고리 내 이름 중복 확인 → 중복 시 409
3. DB 저장
4. TreatmentResponse 반환
```

### TreatmentDataService.create_dosage_type

```
1. treatment_id로 부모 시술 존재 확인 → 미존재 시 404
2. 같은 시술 내 type 중복 확인 → 중복 시 409
3. DB 저장
4. DosageTypeResponse 반환
```

---

## Cascade Delete Rules

| 삭제 대상 | Cascade 동작 | 영향 |
|-----------|-------------|------|
| Category 삭제 | 하위 Treatment + DosageType 모두 삭제 | 프론트 드롭다운에서 사라짐 |
| Treatment 삭제 | 하위 DosageType 모두 삭제 | 프론트 드롭다운에서 사라짐 |
| CycleRule 삭제 | 없음 | 기존 예정일 유지, 새 기록에 예정일 미생성 |

---

## Error Response Format

```json
{
  "detail": "Category with id 'xxx' not found"
}
```

| HTTP Status | 상황 |
|-------------|------|
| 201 | 생성 성공 |
| 200 | 조회/수정 성공 |
| 204 | 삭제 성공 |
| 404 | 리소스 미존재 |
| 409 | 중복 (Conflict) |
| 422 | 입력값 검증 실패 (Pydantic) |
| 500 | 서버 내부 오류 |

---

## Calendar Service와의 관계

- Calendar Service가 이 서비스의 API를 호출하여 주기/시술 데이터 조회
- 호출 경로: `GET /api/cycle-rules/{categoryId}`, `GET /api/categories`, etc.
- 이 서비스가 다운되면: Calendar Service는 예정일 미생성 (graceful degradation)
