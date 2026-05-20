# Web Frontend (Next.js) — Comprehensive Design

> Next.js App Router | React Query | Tailwind CSS + HeroUI | Port 3000

---

## Tech Stack

| 항목 | 선택 |
|------|------|
| Framework | Next.js (App Router) |
| State Management | React Query (TanStack Query) — 서버 상태 |
| UI Library | HeroUI (Tailwind CSS 기반) |
| Styling | Tailwind CSS |
| Color Theme | Primary: `#F66336`, Black: `#131517`, Gray (description): `#697683` |
| Google 연동 | Google Identity Services + Calendar API |
| 배포 | Vercel 또는 독립 컨테이너 |

---

## Pages & Routes

| Route | Page Component | 설명 |
|-------|---------------|------|
| `/calendar` | CalendarPage | 주간 캘린더 뷰 + 날짜 선택 바텀시트 + 하단 시술 통계 |
| `/calendar/records/new` | AddRecordPage | 시술 수동 추가 |
| `/calendar/records/[id]` | RecordDetailPage | 시술 상세/수정/삭제 + 구글 캘린더 등록 |

---

## Feature Components

| Component | 책임 |
|-----------|------|
| **WeeklyCalendarGrid** | 주간 캘린더 그리드 (날짜별 시술 dot 표시) |
| **DateBottomSheet** | 날짜 선택 시 올라오는 바텀시트 (해당 날짜 시술 목록) |
| **RecordCard** | 시술 기록 카드 (실선 스타일 — 확정 시술) |
| **ScheduleCard** | 예정일 카드 (점선 스타일 — 예정 시술) |
| **ScheduleConfirmModal** | 예정일 카드 클릭 시 "받았어요" / "삭제" 확인 모달 |
| **TreatmentStats** | 시술 통계 섹션 (CalendarPage 하단 인라인) |
| **TreatmentDropdown** | 3단계 드롭다운 (카테고리 → 시술명 → 용량) |
| **CategoryBottomSheet** | 카테고리 선택 바텀시트 |
| **TreatmentBottomSheet** | 시술명 선택 바텀시트 |
| **DosageInput** | 용량: 숫자인풋 + 단위 드롭다운 |
| **HospitalInput** | 병원명 자유 텍스트 입력 |
| **GoogleCalendarButton** | 상세페이지 하단 outline 버튼 "구글 캘린더에 등록하기" |
| **FloatingAddButton** | 우측 하단 "+" 플로팅 버튼 → AddRecordPage 이동 |
| **ReminderBanner** | 리마인드 알림 배너 |

---

## React Query Hooks

| Hook | API 호출 | staleTime | 비고 |
|------|----------|-----------|------|
| `useRecords(from, to)` | GET /api/records | 30s | 기간별 조회 |
| `useRecord(id)` | GET /api/records/{id} | 60s | 단건 조회 |
| `useSchedules()` | GET /api/schedules | 30s | 예정일 목록 |
| `useSchedule(id)` | GET /api/schedules/{id} | 60s | 예정일 단건 (구글 캘린더용) |
| `useStatistics()` | GET /api/statistics | 5min | 통계 |
| `useCategories()` | GET /api/treatment-data/categories | 10min | 드롭다운 1단계 |
| `useTreatments(categoryId)` | GET /api/treatment-data/categories/{id}/treatments | 10min | 드롭다운 2단계 |
| `useDosageTypes(treatmentId)` | GET /api/treatment-data/treatments/{id}/dosage-types | 10min | 드롭다운 3단계 |
| `useCreateRecord()` | POST /api/records | mutation | invalidates records, statistics |
| `useUpdateRecord()` | PUT /api/records/{id} | mutation | invalidates records, statistics |
| `useDeleteRecord()` | DELETE /api/records/{id} | mutation | invalidates records, statistics, schedules |
| `useCompleteSchedule()` | PATCH /api/schedules/{id}/complete | mutation | invalidates schedules, records |
| `useDeleteSchedule()` | DELETE /api/schedules/{id} | mutation | invalidates schedules |

---

## Component Flows

### CalendarPage (메인 뷰)

```
구조:
┌─────────────────────────┐
│  WeeklyCalendarGrid     │
│  (주간 뷰, 날짜별 dot)  │
├─────────────────────────┤
│  TreatmentStats         │
│  (시술 통계 인라인)      │
├─────────────────────────┤
│                    [+]  │  ← FloatingAddButton (우측 하단)
└─────────────────────────┘

동작:
1. WeeklyCalendarGrid: 주 단위로 날짜 표시, 각 날짜 셀에 시술 카드 표시 (구글 캘린더 스타일)
   - 카드 내용: "병원명 / 시술명"
   - width 초과 시 ellipsis 처리
   - 확정 시술: 실선 카드, 예정 시술: 점선 카드
2. 날짜 클릭 → DateBottomSheet 올라옴
3. 하단: TreatmentStats (useStatistics)
4. 우측 하단: FloatingAddButton "+" → /calendar/records/new
```

### WeeklyCalendarGrid

```
1. 현재 주 기준 useRecords(weekStart, weekEnd) + useSchedules()
2. 각 날짜 셀에 시술 카드 표시 (구글 캘린더 스타일):
   - 확정 시술: 실선 카드 "병원명 / 시술명" (overflow → ellipsis)
   - 예정 시술: 점선 카드 "시술명 (예정)" (overflow → ellipsis)
3. 주 스와이프 → 이전/다음 주 이동 → 자동 refetch
4. 날짜 클릭 → DateBottomSheet 열림 (selectedDate 설정)

상태: currentWeekStart, selectedDate
```

### DateBottomSheet

```
날짜 선택 시 올라오는 바텀시트

구조:
┌─────────────────────────────┐
│  2026년 1월 15일 (수)        │
├─────────────────────────────┤
│  [+ 이 날짜에 시술 추가]     │  ← 추가 버튼
├─────────────────────────────┤
│  ┌─────────────────────┐   │
│  │ 보톡스 - 사각턱      │   │  ← RecordCard (실선) → 클릭 시 상세페이지
│  │ 강남언니의원         │   │
│  └─────────────────────┘   │
│  ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┐   │
│  │ 필러 - 예정          │   │  ← ScheduleCard (점선) → 클릭 시 ScheduleConfirmModal
│  └ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┘   │
└─────────────────────────────┘

동작:
1. selectedDate 기준으로 해당 날짜의 records + schedules 필터링
2. [+ 이 날짜에 시술 추가] → /calendar/records/new?date={selectedDate}
3. RecordCard(실선) 클릭 → /calendar/records/[id]
4. ScheduleCard(점선) 클릭 → ScheduleConfirmModal 열림
```

### ScheduleConfirmModal

```
예정일 카드 클릭 시 표시되는 모달

┌─────────────────────────┐
│  이 날짜에 시술을         │
│  받으셨나요?             │
├─────────────────────────┤
│  [삭제]      [받았어요]  │
└─────────────────────────┘

동작:
- [받았어요] 클릭:
  1. 예정일(점선) → 확정 기록(실선)으로 변환
  2. useCreateRecord() 호출 (해당 예정일의 시술 정보로 Record 생성)
  3. useCompleteSchedule() 호출 → PATCH /api/schedules/{id}/complete
  4. records + schedules 캐시 무효화
  
- [삭제] 클릭:
  1. useDeleteSchedule() 호출 → DELETE /api/schedules/{id}
  2. schedules 캐시 무효화
```

### RecordDetailPage (상세 페이지)

```
1. useRecord(id) → 시술 기록 상세 조회
2. 시술 정보 표시 (카테고리, 시술명, 용량, 날짜, 병원명)
3. 액션 버튼:
   - [수정] → 수정 모드 전환
   - [삭제] → 확인 모달 → useDeleteRecord()
4. 하단: GoogleCalendarButton
   - outline 스타일 버튼
   - 라벨: "구글 캘린더에 등록하기"

구조:
┌─────────────────────┐
│   시술 상세 정보     │
├─────────────────────┤
│   [수정] [삭제]     │
├─────────────────────┤
│   [구글 캘린더에     │
│    등록하기] (outline)│
└─────────────────────┘
```

### AddRecordPage (시술 수동 추가)

```
1. 날짜 선택 (DatePicker — 클릭 시 바텀시트로 월 캘린더 표시, URL query에서 date 있으면 프리필)
2. 병원명 입력 (HospitalInput — 단순 텍스트)
3. 시술 정보 (TreatmentDropdown)
4. 메모 (선택, TextArea)
5. [저장]:
   a. 클라이언트 검증 (날짜, 카테고리, 시술명 필수)
   b. useCreateRecord() mutation
   c. 성공 → /calendar 이동 + records/statistics 캐시 무효화
   d. 실패 → 에러 토스트 + 폼 유지

에러: 필수값 미입력→인라인 에러, API 실패→토스트
```

### TreatmentDropdown (3단계 드롭다운)

```
1. 마운트 → useCategories() → 카테고리 목록 로드
2. 카테고리 드롭다운 클릭 → CategoryBottomSheet 열림
3. 카테고리 선택 → useTreatments(categoryId) → 두 번째 활성화
4. 시술명 드롭다운 클릭 → TreatmentBottomSheet 열림
5. 시술명 선택 → useDosageTypes(treatmentId) → DosageInput 활성화
6. DosageInput: [숫자 인풋] [단위 드롭다운]

상태: selectedCategoryId, selectedTreatmentId, volumeValue, volumeUnit
에러: API 실패 → 에러 + 재시도, 부모 미선택 → 하위 disabled
```

### GoogleCalendarButton

```
위치: RecordDetailPage 하단
스타일: outline 버튼, 라벨 "구글 캘린더에 등록하기"

동작:
1. 버튼 클릭 → Google OAuth 팝업
2. 인증 성공 → access_token
3. Event 제목 구성:
   - 확정 시술 (Record): "시술명 - 병원명"
   - 예정 시술 (Schedule): "(예정) 시술명"
4. Google Calendar API POST → 성공/실패 토스트
```

---

## Data Flow Patterns

### Query (조회)
```
[Page] → Hook → fetch() → Calendar API → Cache → UI
```

### Mutation (생성/수정/삭제)
```
[Form] → Validation → useMutation() → API → onSuccess: invalidate → refetch
```

### Google Calendar (클라이언트 직접)
```
[Button] → OAuth → access_token → Google Calendar API → 토스트
```

---

## Error Handling Strategy

| 상황 | UI 처리 | HeroUI 컴포넌트 |
|------|---------|-----------------|
| API 로딩 중 | Skeleton UI | `<Skeleton>` |
| API 실패 (네트워크) | 에러 배너 + 재시도 | `<Card>` + `<Button>` |
| API 실패 (4xx) | 토스트 + 필드 에러 | Toast + inline error |
| 빈 데이터 | Empty State | Custom component |
| Google OAuth 실패 | 에러 토스트 | Toast |
| 필수값 미입력 | 인라인 에러 메시지 | `<Input isInvalid>` |

---

## Component Dependencies

| Component | Depends On |
|-----------|-----------|
| CalendarPage | WeeklyCalendarGrid, DateBottomSheet, TreatmentStats, FloatingAddButton |
| WeeklyCalendarGrid | useRecords, useSchedules |
| DateBottomSheet | RecordCard, ScheduleCard, ScheduleConfirmModal |
| ScheduleConfirmModal | useCreateRecord, useCompleteSchedule, useDeleteSchedule |
| AddRecordPage | TreatmentDropdown, HospitalInput, DatePicker, useCreateRecord |
| RecordDetailPage | useRecord, TreatmentDropdown, GoogleCalendarButton, useUpdateRecord, useDeleteRecord |
| TreatmentStats | useStatistics |
| TreatmentDropdown | CategoryBottomSheet, TreatmentBottomSheet, DosageInput, useCategories, useTreatments, useDosageTypes |
| GoogleCalendarButton | Google Identity Services SDK, useRecord or useSchedule |
