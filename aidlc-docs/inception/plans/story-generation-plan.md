# Story Generation Plan — 시술 관리 캘린더

## Execution Checklist

- [x] 페르소나 정의
- [x] 유저 스토리 작성 (INVEST 기준)
- [x] Acceptance Criteria 작성
- [x] 페르소나-스토리 매핑

---

## Questions

아래 질문에 답변해주세요. `[Answer]:` 뒤에 선택지를 입력해주세요.

---

### Question 1
유저 스토리 분류 방식은 어떤 것을 선호하시나요?

A) User Journey 기반 — 사용자 플로우 순서대로 (시술 기록 → 조회 → 예정일 확인 → 알림 수신)
B) Feature 기반 — 기능 단위로 그룹핑 (자동 기입, 수동 CRUD, 예정일, 알림, 통계, 구글 캘린더)
C) Persona 기반 — 사용자 유형별 (고객 스토리, 관리자 스토리)
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

### Question 2
유저 스토리의 세분화 수준은 어느 정도가 적절한가요?

A) 큰 단위 (Epic 수준) — 예: "고객은 시술을 기록할 수 있다" (약 8~10개 스토리)
B) 중간 단위 — 예: "고객은 3단계 드롭다운으로 시술 정보를 선택할 수 있다" (약 15~20개 스토리)
C) 작은 단위 — 예: "고객은 시술 카테고리 드롭다운을 클릭하면 바텀시트를 볼 수 있다" (약 30개+ 스토리)
X) Other (please describe after [Answer]: tag below)

[Answer]: A와 C를 결합해서 큰 epic 안에 여러 스토리가 있는 방식.

---

### Question 3
관리자 페르소나의 범위는 어떻게 되나요?

A) 내부 운영팀 (강남언니 직원) — 시술 데이터와 추천 주기를 관리하는 담당자
B) 병원 관계자 — 자기 병원의 시술 정보를 관리하는 사람
C) 시스템 관리자 — 기술적 설정만 담당
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

### Question 4
고객 페르소나를 세분화할 필요가 있나요?

A) 단일 페르소나 — "시술을 주기적으로 받는 고객" 하나로 충분
B) 2개 페르소나 — "강남언니로 예약하는 고객" + "외부에서 시술받고 수동 기록하는 고객"
C) 3개 페르소나 — B + "처음 시술 캘린더를 사용하는 신규 고객"
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

### Question 5
Acceptance Criteria 형식은 어떤 것을 선호하시나요?

A) Given-When-Then (BDD 스타일) — 예: "Given 보톡스 시술 기록이 있을 때, When 90일이 경과하면, Then 리마인드 알림이 발송된다"
B) 체크리스트 형식 — 예: "✓ 시술 카테고리 선택 가능 ✓ 선택 시 바텀시트 표시"
C) 혼합 — 핵심 플로우는 Given-When-Then, 단순 UI는 체크리스트
X) Other (please describe after [Answer]: tag below)

[Answer]: A
