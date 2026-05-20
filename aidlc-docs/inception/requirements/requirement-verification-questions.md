# Requirements Verification Questions

아래 질문에 답변해주세요. 각 질문의 `[Answer]:` 태그 뒤에 선택지 알파벳을 입력해주세요.
제공된 선택지가 맞지 않으면 마지막 옵션(Other)을 선택하고 설명을 추가해주세요.

---

## Question 1
이 프로젝트의 플랫폼은 무엇인가요?

A) 모바일 앱 (iOS/Android) — 네이티브
B) 모바일 앱 — Flutter/React Native 등 크로스플랫폼
C) 웹 애플리케이션 (SPA)
D) 백엔드 API 서버만 (프론트엔드는 별도 팀에서 개발)
E) 풀스택 (백엔드 + 프론트엔드 모두)
X) Other (please describe after [Answer]: tag below)

[Answer]: C

---

## Question 2
백엔드 기술 스택 선호도가 있나요?

A) Kotlin + Spring Boot
B) Java + Spring Boot
C) Node.js (TypeScript)
D) Python (FastAPI/Django)
E) Go
X) Other (please describe after [Answer]: tag below)

[Answer]: E

---

## Question 3
데이터베이스 선호도가 있나요?

A) 관계형 DB (PostgreSQL)
B) 관계형 DB (MySQL)
C) NoSQL (MongoDB)
D) NoSQL (DynamoDB)
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

## Question 4
인프라/배포 환경은 어떻게 되나요?

A) AWS 클라우드
B) GCP 클라우드
C) On-premise / 자체 서버
D) 아직 미정 (추천 받고 싶음)
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

## Question 5
"예약 확정(FIXED) 이벤트"를 캘린더에 자동 기입하는 방식은 어떻게 구현하나요?

A) 기존 강남언니 예약 시스템에서 이벤트/메시지를 발행하고, 이 서비스가 구독하는 방식 (이벤트 기반)
B) 기존 예약 시스템 DB를 직접 조회하는 방식 (폴링/배치)
C) 기존 예약 시스템이 이 서비스의 API를 호출하는 방식 (동기 호출)
D) 아직 미정 (추천 받고 싶음)
X) Other (please describe after [Answer]: tag below)

[Answer]: A. 지금은 당장 이벤트를 구독할 수 있는 상황은 아님. 개발시에는 이벤트에 대한 모킹이 필요함.

---

## Question 6
푸시 알림 인프라는 이미 존재하나요?

A) 예 — 기존 푸시 알림 시스템이 있고, 해당 시스템에 알림 요청만 보내면 됨
B) 아니오 — 푸시 알림 시스템도 새로 구축해야 함
C) 부분적으로 있음 — 기존 시스템이 있지만 이 프로젝트에 맞게 확장 필요
X) Other (please describe after [Answer]: tag below)

[Answer]: A. 캘린더 서비스 개발시에는 알림 발송 요청 및 수신에 대한 모킹이 필요함.

---

## Question 7
시술 추천 주기 룰베이스 데이터는 어떻게 관리하나요?

A) 시술 카테고리별 고정 주기를 DB에 저장 (관리자가 수정 가능)
B) 설정 파일(config)로 관리
C) 별도 관리자 페이지에서 CRUD
D) 아직 미정 (추천 받고 싶음)
X) Other (please describe after [Answer]: tag below)

[Answer]: A. 관리자는 API를 통해서 데이터 수정 가능.

---

## Question 8
이 서비스의 예상 사용자 규모는 어느 정도인가요?

A) 소규모 (일일 활성 사용자 1만 이하)
B) 중규모 (일일 활성 사용자 1만~10만)
C) 대규모 (일일 활성 사용자 10만 이상)
D) 아직 미정
X) Other (please describe after [Answer]: tag below)

[Answer]: B

---

## Question 9
인증/인가는 어떻게 처리하나요?

A) 기존 강남언니 인증 시스템 연동 (토큰 기반)
B) 별도 인증 시스템 구축
C) 아직 미정 (추천 받고 싶음)
X) Other (please describe after [Answer]: tag below)

[Answer]: X. 인증은 이번 범위에서 제외.

---

## Question 10
구글 캘린더 연동 방식은?

A) Google Calendar API를 직접 호출 (서버 사이드)
B) 클라이언트에서 Google Calendar API 호출
C) .ics 파일 다운로드 방식 (범용 캘린더 호환)
D) 아직 미정 (추천 받고 싶음)
X) Other (please describe after [Answer]: tag below)

[Answer]: B, 어떤 방식이 보안, 안정성 측면에서 좋은지 검토하고 추천해줘.

---

## Question 11: Security Extensions
이 프로젝트에 보안 확장 규칙을 적용할까요?

A) Yes — 모든 보안 규칙을 blocking constraint로 적용 (프로덕션 수준 애플리케이션에 권장)
B) No — 보안 규칙 건너뛰기 (PoC, 프로토타입, 실험적 프로젝트에 적합)
X) Other (please describe after [Answer]: tag below)

[Answer]: B, 최소한의 보안은 적용하고 싶음.

---

## Question 12: Property-Based Testing Extension
이 프로젝트에 Property-Based Testing(PBT) 규칙을 적용할까요?

A) Yes — 모든 PBT 규칙을 blocking constraint로 적용 (비즈니스 로직, 데이터 변환, 직렬화, 상태 관리 컴포넌트가 있는 프로젝트에 권장)
B) Partial — 순수 함수와 직렬화 round-trip에만 PBT 규칙 적용
C) No — PBT 규칙 건너뛰기 (단순 CRUD, UI 전용, 얇은 통합 레이어에 적합)
X) Other (please describe after [Answer]: tag below)

[Answer]: X. 이게 무슨 질문인지 친절하게 설명해줘.
