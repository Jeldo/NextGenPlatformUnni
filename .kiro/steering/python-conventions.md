# Python (FastAPI) Code Conventions

이 문서는 FastAPI 기반 Admin Service의 코딩 규칙을 정의합니다.

## 버전과 의존성

- **FastAPI 0.115+**, **Pydantic v2**, **Uvicorn** (개발) / **Gunicorn + Uvicorn workers** (프로덕션)
- 비동기 HTTP 클라이언트: `httpx` (`requests` 대신)
- ORM: `SQLAlchemy 2.x` (async)
- 마이그레이션: `Alembic`
- 테스트: `pytest` + `httpx.AsyncClient`

## 프로젝트 구조 (계층형 — 작은 서비스)

```
app/
├── main.py                 # FastAPI 앱 생성 및 라우터 등록
├── config.py               # Settings (pydantic-settings)
├── routers/                # APIRouter (도메인별)
│   ├── cycle_rules.py
│   └── treatment_data.py
├── services/               # 비즈니스 로직
├── repositories/           # 데이터 액세스
├── models/                 # SQLAlchemy ORM 엔티티
├── schemas/                # Pydantic 모델 (요청/응답 DTO)
└── db/
    └── session.py          # DB 세션 관리
```

- 라우터는 얇게 유지, 비즈니스 로직은 서비스에 위치
- 도메인 간 직접 import 피하기

## 비동기 (async) 가이드

- 라우트 핸들러는 `async def`로 작성
- 핸들러 안에서 **블로킹 호출 금지** (불가피하면 `await asyncio.to_thread(...)`)
- DB 드라이버는 async 드라이버(`asyncpg`) 사용
- CPU 바운드 작업은 백그라운드 워커로 이관

## 의존성 주입 (Depends)

- 횡단 관심사(세션, 설정)는 `Depends`로 주입
- 재사용 가능한 의존성은 별도 파일에 정의
- 테스트에서 `app.dependency_overrides`로 교체

```python
async def get_session() -> AsyncIterator[AsyncSession]:
    async with AsyncSessionLocal() as session:
        try:
            yield session
            await session.commit()
        except Exception:
            await session.rollback()
            raise
```

## Pydantic 스키마 규칙

- 요청/응답 DTO는 ORM 모델과 분리
- 패턴:
  - `XxxCreate`: 생성 요청
  - `XxxUpdate`: 부분 수정 (필드 모두 `Optional`)
  - `XxxResponse`: 응답 모델 (`model_config = ConfigDict(from_attributes=True)`)
- `response_model`을 모든 엔드포인트에 명시
- 민감 필드는 응답 스키마에 포함 금지

## 설정 관리

- `pydantic-settings`의 `BaseSettings` 사용
- `.env` 파일은 로컬 개발 전용
- 비밀값을 코드/리포지토리에 커밋 금지

```python
class Settings(BaseSettings):
    database_url: str
    port: int = 8081

    model_config = SettingsConfigDict(env_file=".env")
```

## 에러 처리

- 비즈니스 예외는 도메인 예외 클래스로 정의
- 글로벌 `exception_handler`로 HTTP 응답 매핑
- 에러 응답 형식 통일:
```json
{"error": {"code": "NOT_FOUND", "message": "Category not found"}}
```
- `HTTPException(detail=...)`에 내부 구현 정보 노출 금지

## 데이터베이스 패턴

- 세션은 의존성으로 주입, 요청 단위 트랜잭션
- 리포지토리에서 ORM 객체 다루기, 서비스는 도메인 로직만
- 라우터에서 `commit()` 직접 호출 금지
- N+1 쿼리 방지: `selectinload` / `joinedload` 사용

## 라이프사이클

- `lifespan` async context manager 사용 (`@app.on_event` 금지)
- DB 풀, HTTP 클라이언트는 `lifespan`에서 생성/정리

```python
@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncIterator[None]:
    # startup
    yield
    # shutdown

app = FastAPI(lifespan=lifespan)
```

## 네이밍

- 함수/변수: `snake_case`
- 클래스: `PascalCase`
- 상수: `UPPER_SNAKE_CASE`
- private: `_prefix`
- 파일명: `snake_case.py`

## 안티 패턴 (피할 것)

- 라우터에 비즈니스 로직과 DB 호출 직접 작성
- 동기 ORM 호출을 `async def` 핸들러에서 사용
- 글로벌 변수로 DB 세션/HTTP 클라이언트 보유
- `Exception`을 잡아 200 OK로 변환
- 응답 모델 없이 `dict` 그대로 반환
- `allow_origins=["*"]` + `allow_credentials=True` 조합
- 마이그레이션 없이 `create_all()`로 운영 스키마 변경
