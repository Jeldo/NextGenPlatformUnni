# Application Design Plan — 시술 관리 캘린더

## Execution Checklist

- [x] 컴포넌트 식별 및 책임 정의
- [x] 컴포넌트 메서드 시그니처 정의
- [x] 서비스 레이어 설계
- [x] 컴포넌트 의존성 및 통신 패턴 정의
- [x] 통합 설계 문서 생성

---

## Questions

아래 질문에 답변해주세요. `[Answer]:` 뒤에 선택지를 입력해주세요.

---

### Question 1
Go 캘린더 서비스의 내부 아키텍처 패턴은?

A) Layered Architecture — Handler → Service → Repository 3계층
B) Hexagonal (Ports & Adapters) — 도메인 중심, 외부 의존성을 어댑터로 분리
C) Clean Architecture — Use Case 중심, 의존성 역전
X) Other (please describe after [Answer]: tag below)

[Answer]: B

---

### Question 2
FastAPI 관리 서비스의 구조는?

A) 단순 CRUD — Router → Service → Repository (최소 구조)
B) Domain-Driven — 도메인 모델 중심 설계
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

### Question 3
캘린더 서비스의 API 스타일은?

A) RESTful API (리소스 기반)
B) GraphQL
C) gRPC
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

### Question 4
이벤트 구독 처리와 API 서빙을 같은 프로세스에서 할까요?

A) 같은 프로세스 — HTTP 서버 + 이벤트 컨슈머가 하나의 바이너리
B) 분리 — API 서버와 이벤트 컨슈머를 별도 프로세스로 배포
X) Other (please describe after [Answer]: tag below)

[Answer]: B

---

### Question 5
시술 기록과 예정일의 데이터 관계는?

A) 예정일은 조회 시 실시간 계산 (별도 저장 안 함)
B) 예정일을 별도 테이블에 저장하고 시술 기록 변경 시 갱신
X) Other (please describe after [Answer]: tag below)

[Answer]: B
