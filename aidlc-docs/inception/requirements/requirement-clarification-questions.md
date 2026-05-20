# Requirements Clarification Questions

이전 답변에서 몇 가지 확인이 필요한 부분이 있습니다.

---

## Property-Based Testing (PBT) 설명

**PBT란?** 일반 단위 테스트는 "입력 A → 결과 B" 같은 구체적 케이스를 검증하지만, PBT는 "어떤 입력이든 이 속성(property)을 만족해야 한다"를 검증합니다.

**예시 (시술 캘린더 맥락):**
- 일반 테스트: `보톡스 시술일 2026-01-01 + 주기 90일 = 예정일 2026-04-01`
- PBT: `어떤 시술이든, 시술일 + 추천주기 = 예정일이어야 하고, 예정일은 항상 시술일보다 미래여야 한다`

**이 프로젝트에서 PBT가 유용한 영역:**
- 추천 주기 계산 로직 (날짜 연산의 정확성)
- 시술 기록 CRUD의 데이터 무결성

**적용 수준:**
- Yes: 모든 비즈니스 로직에 PBT 적용 (엄격)
- Partial: 날짜 계산, 데이터 변환 같은 순수 함수에만 적용 (실용적)
- No: 일반 단위 테스트만 사용 (가장 간단)

---

## Clarification Question 1
위 설명을 참고하여, PBT를 어떻게 적용할까요?

A) Partial — 날짜 계산, 주기 로직 등 순수 함수에만 PBT 적용 (권장)
B) No — 일반 단위 테스트만 사용
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

## Clarification Question 2
보안 관련: "최소한의 보안"이란 구체적으로 어떤 수준을 의미하나요?

A) 입력값 검증(validation) + SQL injection 방지 수준
B) A + API rate limiting + 기본적인 에러 핸들링 (민감 정보 노출 방지)
C) B + HTTPS 강제 + 환경변수로 시크릿 관리
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

## Clarification Question 3
구글 캘린더 연동: 클라이언트 호출 방식을 선택하셨는데, 추가 검토 결과를 공유드립니다.

**방식 비교:**
| 방식 | 장점 | 단점 |
|------|------|------|
| B) 클라이언트 직접 호출 | 서버 부하 없음, 구현 간단 | 사용자가 Google OAuth 인증 필요, 토큰 관리가 클라이언트에 의존 |
| C) .ics 파일 다운로드 | OAuth 불필요, 모든 캘린더 앱 호환, 가장 간단 | 단방향만 가능(내보내기만), 자동 동기화 불가 |

PRD 요구사항이 "단건 내보내기"이므로, .ics 방식이 가장 적합할 수 있습니다. 어떤 방식을 선택하시겠어요?

A) 클라이언트에서 Google Calendar API 직접 호출 (Google OAuth 필요)
B) .ics 파일 다운로드 방식 (OAuth 불필요, 범용 캘린더 호환)
C) 둘 다 지원 (기본은 .ics, Google 로그인 시 API 호출)
X) Other (please describe after [Answer]: tag below)

[Answer]: A
