# Frontend (Next.js/TypeScript) Code Conventions

이 문서는 Next.js 프론트엔드의 코딩 규칙을 정의합니다.

## UI 컴포넌트 라이브러리

- **HeroUI (v3)를 필수로 사용**한다. 자체 컴포넌트를 만들기 전에 반드시 HeroUI에 해당 컴포넌트가 있는지 확인한다.
- 컴포넌트를 사용할 때 **반드시 heroui-react MCP 도구(`list_components`, `get_component_docs`)를 호출**하여 정확한 API, props, 사용법을 확인한 후 구현한다.
- **절대 컴포넌트 API를 유추하지 않는다.** MCP 조회 결과가 없으면 공식 문서를 추가로 확인한다.
- HeroUI에 없는 UI가 필요한 경우, HeroUI의 기본 컴포넌트를 조합하여 구현한다.

## 설계 참조

- 코드 구현 시 **반드시 inception 단계 산출물(`aidlc-docs/inception/`)을 참조**하여 기능 요구사항과 동작 흐름을 정확히 반영한다.
- 특히 `web-frontend.md`의 Component Flows, React Query Hooks, Error Handling Strategy를 따른다.
- 요구사항에 명시되지 않은 동작을 임의로 추가하지 않는다.

## 기술 스택 (변경 금지)

| 항목 | 선택 |
|------|------|
| Framework | Next.js (App Router) |
| State Management | React Query (TanStack Query) — 서버 상태 전용 |
| Styling | Tailwind CSS + HeroUI |
| Color Theme | Primary `#F66336`, Black `#131517`, Gray `#697683` |
| Google 연동 | Google Identity Services + Calendar API (클라이언트 사이드) |

## 모바일 웹뷰 디자인

이 앱은 **모바일 앱 내 웹뷰**로 동작합니다:
- 터치 타겟 최소 44x44px 유지
- 데스크톱 반응형 breakpoint 불필요 — 모바일 뷰포트 기준으로만 설계
- Safe area inset 고려 (`env(safe-area-inset-*)`)
- 스크롤 동작은 네이티브 스크롤 사용, 커스텀 스크롤바 금지

## 프로젝트 구조

```
src/
├── app/                        # App Router (pages)
│   └── calendar/
│       ├── page.tsx            # CalendarPage
│       └── records/
│           ├── new/page.tsx    # AddRecordPage
│           └── [id]/page.tsx   # RecordDetailPage
├── components/                 # UI 컴포넌트 (flat)
│   ├── WeeklyCalendarGrid.tsx
│   ├── DateBottomSheet.tsx
│   ├── RecordCard.tsx
│   ├── ScheduleCard.tsx
│   ├── ScheduleConfirmModal.tsx
│   ├── TreatmentDropdown.tsx
│   ├── CategoryBottomSheet.tsx
│   ├── TreatmentBottomSheet.tsx
│   ├── DosageInput.tsx
│   ├── HospitalInput.tsx
│   ├── GoogleCalendarButton.tsx
│   ├── FloatingAddButton.tsx
│   ├── TreatmentStats.tsx
│   └── ReminderBanner.tsx
├── hooks/                      # React Query hooks (파일당 1 hook)
│   ├── useRecords.ts
│   ├── useRecord.ts
│   ├── useSchedules.ts
│   ├── useStatistics.ts
│   ├── useCategories.ts
│   ├── useTreatments.ts
│   ├── useDosageTypes.ts
│   ├── useCreateRecord.ts
│   ├── useUpdateRecord.ts
│   ├── useDeleteRecord.ts
│   ├── useCompleteSchedule.ts
│   └── useDeleteSchedule.ts
├── services/
│   └── api.ts                  # API base client (fetch wrapper)
└── types/
    └── index.ts                # 공통 타입 정의
```

## 네이밍

- 컴포넌트: PascalCase (`WeeklyCalendarGrid.tsx`)
- Hook: camelCase, `use` prefix (`useRecords.ts`)
- 타입: PascalCase (`TreatmentRecord`, `CreateRecordRequest`)
- 파일명: 컴포넌트는 PascalCase, 나머지는 camelCase
- Props 타입: `{ComponentName}Props`

## 상태 관리

- **서버 상태**: React Query만 사용 (useState로 API 데이터 관리 금지)
- **UI 상태**: React useState/useReducer (최소한으로)
- **폼 상태**: 컴포넌트 로컬 state 또는 React Hook Form
- 전역 상태 라이브러리 사용 금지 (이 프로젝트에서는 불필요)

## React Query 규칙

- Hook 파일당 1개 hook만 export
- Query key는 배열 형태: `['records', { from, to }]`
- Mutation 성공 시 관련 query invalidate
- staleTime은 데이터 특성에 맞게 설정 (변경 빈도 기준)
- 에러/로딩 상태를 항상 처리

## 에러/로딩/빈 상태 처리

모든 비동기 UI는 3가지 상태를 명시적으로 처리합니다:

| 상태 | 처리 | HeroUI 컴포넌트 |
|------|------|-----------------|
| 로딩 | Skeleton UI | `<Skeleton>` |
| 에러 (네트워크) | 에러 배너 + 재시도 | `<Card>` + `<Button>` |
| 에러 (4xx) | Toast + 인라인 에러 | Toast + `<Input isInvalid>` |
| 빈 데이터 | Empty State | Custom component |

## TypeScript 규칙

- `any` 사용 금지 — `unknown` + 타입 가드 사용
- API 응답 타입은 `types/index.ts`에 중앙 정의
- Props는 interface로 정의 (type alias 대신)
- Optional chaining과 nullish coalescing 적극 활용
- 타입 단언(`as`) 최소화 — 타입 가드 선호

## 안티 패턴 (피할 것)

- `useEffect`로 데이터 fetching (React Query 사용)
- Props drilling 3단계 이상 (컴포넌트 분리 또는 composition)
- 인라인 스타일 사용 (Tailwind 클래스 사용)
- `console.log` 커밋 (개발 중에만 사용)
- 하드코딩된 API URL (환경변수 사용)
- 컴포넌트 내부에서 직접 fetch 호출 (hook으로 분리)
