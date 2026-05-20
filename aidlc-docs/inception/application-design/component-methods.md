# Component Methods (Comprehensive) — 시술 관리 캘린더

> 각 메서드의 동작 흐름, 입출력 데이터 구조, 에러 처리, 상태 전이를 상세 정의.

---

## Frontend — Component Methods & Flows

### TreatmentDropdown (3단계 드롭다운)

```
동작 흐름:
1. 컴포넌트 마운트 시 useCategories() 호출 → 카테고리 목록 로드
2. 카테고리 드롭다운 클릭 → CategoryBottomSheet 열림
3. 카테고리 선택 시:
   a. selectedCategoryId 상태 업데이트
   b. useTreatments(categoryId) 트리거 → 시술명 목록 로드
   c. 두 번째 드롭다운 활성화
4. 시술명 드롭다운 클릭 → TreatmentBottomSheet 열림
5. 시술명 선택 시:
   a. selectedTreatmentId 상태 업데이트
   b. useDosageTypes(treatmentId) 트리거 → 용량 단위 로드
   c. 세 번째 영역 활성화
6. 용량 입력: 숫자 인풋 + 단위 드롭다운

상태 관리:
- selectedCategoryId: string | null
- selectedTreatmentId: string | null
- dosageValue: number | null
- dosageUnit: string | null

에러 처리:
- API 로드 실패 → 에러 메시지 표시 + 재시도 버튼
- 부모 선택 없이 하위 드롭다운 비활성화 (disabled 상태)
```

### CalendarGrid

```
동작 흐름:
1. 현재 월 기준 useRecords(from, to) + useSchedules() 호출
2. 날짜별로 기록과 예정일 매핑
3. 각 날짜 셀에:
   - 확정 기록: RecordCard (실선 스타일)
   - 예정일: ScheduleCard (점선 스타일)
4. 월 변경 시: from/to 업데이트 → React Query 자동 refetch
5. 날짜 클릭 시: 해당 날짜의 기록 목록 표시

상태 관리:
- currentMonth: Date
- selectedDate: Date | null

에러 처리:
- 로딩 중: Skeleton UI 표시
- 실패 시: 에러 배너 + 재시도
```

### AddRecordPage

```
동작 흐름:
1. 날짜 선택 (DatePicker — 클릭 시 바텀시트로 월 캘린더 표시, 날짜 선택)
2. 병원명 입력 (HospitalInput — 자유 텍스트)
3. 시술 정보 입력 (TreatmentDropdown)
4. 메모 입력 (선택, TextArea)
5. [저장] 버튼 클릭:
   a. 입력값 클라이언트 검증 (날짜, 카테고리, 시술명 필수)
   b. useCreateRecord() mutation 호출
   c. 성공 시: 캘린더 페이지로 이동 + records 캐시 무효화
   d. 실패 시: 에러 토스트 표시

에러 처리:
- 필수값 미입력 → 필드별 인라인 에러 메시지
- API 실패 → 토스트 알림 + 폼 데이터 유지
```

### GoogleCalendarButton

```
동작 흐름:
1. 버튼 클릭
2. Google OAuth 팝업 (gapi.auth2 또는 Google Identity Services)
3. 인증 성공 시:
   a. access_token 획득
   b. 시술 기록/예정일 데이터로 Google Calendar Event 구성
   c. Google Calendar API POST /calendars/primary/events 호출
   d. 성공 → 성공 토스트
4. 인증 실패/취소 시: 에러 토스트

에러 처리:
- OAuth 팝업 차단 → 안내 메시지
- Google API 호출 실패 → 재시도 안내
- 토큰 만료 → 재인증 요청
```

### HospitalInput

```
단순 텍스트 인풋 (자유 입력)
- placeholder: "병원명을 입력하세요"
- 자동완성 없음

상태 관리:
- inputValue: string
```
```

---

## Data Structures

### Core Domain Models

```go
type TreatmentRecord struct {
    ID              string
    UserID          string
    Source          RecordSource       // AUTO | MANUAL
    CategoryID      string
    TreatmentID     string
    DosageType      *string            // "shot" | "minute" | "volume" | "vial" | "joule" | nil
    DosageValue     *float64
    TreatmentDate   time.Time
    HospitalName    string
    HospitalLocation *string
    DoctorName      *string
    Memo            *string
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type RecordSource string
const (
    SourceAuto   RecordSource = "AUTO"
    SourceManual RecordSource = "MANUAL"
)

type ScheduledTreatment struct {
    ID              string
    UserID          string
    RecordID        string             // 기반이 된 시술 기록 ID
    CategoryID      string
    TreatmentID     string
    ScheduledDate   time.Time          // 마지막 시술일 + 추천 주기
    CycleDays       int                // 적용된 주기 (일)
    Status          ScheduleStatus     // PENDING | REMINDED | COMPLETED
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type ScheduleStatus string
const (
    StatusPending   ScheduleStatus = "PENDING"
    StatusReminded  ScheduleStatus = "REMINDED"
    StatusCompleted ScheduleStatus = "COMPLETED"
)

type CycleRule struct {
    CategoryID  string
    CycleDays   int
    Description string
    UpdatedAt   time.Time
}

type TreatmentCategory struct {
    ID   string
    Name string
}

type Treatment struct {
    ID         string
    CategoryID string
    Name       string
}

type DosageType struct {
    ID          string
    TreatmentID string
    Type        string  // "shot", "minute", "volume", "vial", "joule"
    Unit        string  // "cc", "회" 등
}

type TreatmentStat struct {
    CategoryID   string
    CategoryName string
    Count        int
}

type ReminderMessage struct {
    UserID        string
    Title         string
    Body          string
    TreatmentName string
    DaysSince     int
}

type ReservationFixedEvent struct {
    ReservationID string
    UserID        string
    HospitalName  string
    TreatmentType string    // "surgery" | "treatment"
    CategoryID    string
    TreatmentID   string
    FixedDate     time.Time
    EventTime     time.Time
}
```

### Input/Output DTOs

```go
type CreateRecordInput struct {
    UserID          string
    Source          RecordSource
    CategoryID      string
    TreatmentID     string
    DosageType      *string
    DosageValue     *string   // API: string, 내부 도메인: float64로 변환
    TreatmentDate   time.Time
    HospitalName    string
    HospitalLocation *string
    DoctorName      *string
    Memo            *string
}

type UpdateRecordInput struct {
    CategoryID      *string
    TreatmentID     *string
    DosageType      *string
    DosageValue     *string   // API: string, 내부 도메인: float64로 변환
    TreatmentDate   *time.Time
    HospitalName    *string
    HospitalLocation *string
    DoctorName      *string
    Memo            *string
}
```

### Error Types

```go
type AppError struct {
    Code    ErrorCode
    Message string
    Cause   error
}

type ErrorCode string
const (
    ErrNotFound         ErrorCode = "NOT_FOUND"
    ErrValidation       ErrorCode = "VALIDATION_ERROR"
    ErrCycleRuleNotFound ErrorCode = "CYCLE_RULE_NOT_FOUND"
    ErrExternalService  ErrorCode = "EXTERNAL_SERVICE_ERROR"
    ErrDatabase         ErrorCode = "DATABASE_ERROR"
    ErrUnauthorized     ErrorCode = "UNAUTHORIZED"
    ErrConflict         ErrorCode = "CONFLICT"
)
```

---

## Calendar Service — Domain Methods

### TreatmentRecordService

#### CreateRecord

```
동작 흐름:
1. 입력값 검증 (CategoryID, TreatmentID, TreatmentDate 필수)
2. TreatmentRecordRepository.Save() 호출
3. CycleRuleClient.GetCycleRule(categoryID) 호출
   - 주기 존재 시: ScheduledTreatmentService.CalculateAndSave() 호출
   - 주기 미존재 시: 예정일 생성 건너뜀 (에러 아님)
4. 생성된 TreatmentRecord 반환

에러 처리:
- 입력값 누락 → ErrValidation 반환
- DB 저장 실패 → ErrDatabase 반환
- CycleRuleClient 실패 → 로그 기록, 예정일 없이 기록만 생성 (graceful degradation)
```

#### UpdateRecord

```
동작 흐름:
1. TreatmentRecordRepository.FindByID(recordID) 조회
2. 소유권 검증 (record.UserID == 요청자)
3. 변경 필드 적용
4. TreatmentRecordRepository.Update() 호출
5. 날짜 또는 카테고리 변경 시:
   a. ScheduledTreatmentRepository.DeleteByRecordID(recordID)
   b. CycleRuleClient.GetCycleRule(newCategoryID) 호출
   c. ScheduledTreatmentService.CalculateAndSave() 호출
6. 수정된 TreatmentRecord 반환

에러 처리:
- 기록 미존재 → ErrNotFound 반환
- 소유권 불일치 → ErrUnauthorized 반환
- DB 업데이트 실패 → ErrDatabase 반환
- CycleRuleClient 실패 → 로그 기록, 기존 예정일 유지
```

#### DeleteRecord

```
동작 흐름:
1. TreatmentRecordRepository.FindByID(recordID) 조회
2. 소유권 검증
3. ScheduledTreatmentRepository.DeleteByRecordID(recordID)
4. TreatmentRecordRepository.Delete(recordID)
5. 성공 반환

에러 처리:
- 기록 미존재 → ErrNotFound 반환
- 소유권 불일치 → ErrUnauthorized 반환
- 예정일 삭제 실패 → 로그 기록 후 기록 삭제 계속 진행
```

#### ListRecords

```
동작 흐름:
1. 기간 검증 (from < to, 최대 1년)
2. TreatmentRecordRepository.FindByUserAndPeriod(userID, from, to)
3. 결과 반환 (빈 배열 가능)

에러 처리:
- 기간 유효하지 않음 → ErrValidation 반환
```

---

### ScheduledTreatmentService

#### CalculateAndSave

```
동작 흐름:
1. CycleRuleClient.GetCycleRule(record.CategoryID) 호출
2. scheduledDate = record.TreatmentDate + cycleDays
3. ScheduledTreatment 생성 (Status=PENDING)
4. ScheduledTreatmentRepository.Save() 호출
5. 생성된 ScheduledTreatment 반환

에러 처리:
- CycleRule 미존재 → ErrCycleRuleNotFound 반환 (호출자가 graceful 처리)
- DB 저장 실패 → ErrDatabase 반환
```

#### ListDueReminders

```
동작 흐름:
1. targetDate 기준으로 Status=PENDING이고 ScheduledDate <= targetDate인 항목 조회
2. 결과 반환

용도: 리마인드 배치 작업에서 호출
```

---

### TreatmentStatisticsService

#### GetStatsByUser

```
동작 흐름:
1. TreatmentRecordRepository.CountByUserGrouped(userID) 호출
2. 카테고리별 집계 결과 반환

에러 처리:
- DB 조회 실패 → ErrDatabase 반환
- 기록 없음 → 빈 배열 반환 (에러 아님)
```

---

### ReminderService (배치)

#### ProcessDueReminders

```
동작 흐름:
1. ScheduledTreatmentService.ListDueReminders(today) 호출
2. 각 예정일에 대해:
   a. 시술명 조회 (CycleRuleClient 또는 캐시)
   b. 경과 일수 계산 (today - 원본 시술일)
   c. ReminderMessage 구성
   d. NotificationClient.SendReminder() 호출
   e. 성공 시: Status를 REMINDED로 업데이트
   f. 실패 시: 로그 기록, 다음 항목 계속 처리
3. 처리 결과 요약 반환

에러 처리:
- NotificationClient 실패 → 개별 항목 스킵, 다음 배치에서 재시도
- DB 업데이트 실패 → 로그 기록, 중복 알림 방지 로직 필요
```

---

## Event Consumer — Methods

### ReservationEventHandler.HandleReservationFixed

```
동작 흐름:
1. 이벤트 파싱 및 검증
2. TreatmentType 필터링:
   - "surgery" → 무시 (로그만 기록)
   - "treatment" → 계속 진행
3. CreateRecordInput 구성:
   - Source = AUTO
   - TreatmentDate = event.FixedDate
   - HospitalName = event.HospitalName
   - CategoryID = event.CategoryID
   - TreatmentID = event.TreatmentID
4. TreatmentRecordService.CreateRecord() 호출
5. 성공 시: 메시지 ACK
6. 실패 시: 재시도 정책에 따라 처리

에러 처리:
- 이벤트 파싱 실패 → Dead Letter Queue로 이동
- 중복 이벤트 (같은 ReservationID) → 멱등성 체크, 이미 존재하면 스킵
- CreateRecord 실패 → 재시도 (최대 3회), 이후 DLQ
```

---

## Admin Service (FastAPI) — Methods

### CycleRuleService

#### create_cycle_rule

```
동작 흐름:
1. 입력 검증 (category_id 존재 여부, cycle_days > 0)
2. 기존 규칙 존재 확인 → 존재 시 ErrConflict
3. DB 저장
4. 생성된 CycleRule 반환

에러 처리:
- category_id 미존재 → 404
- 이미 존재 → 409 Conflict
- cycle_days <= 0 → 422 Validation Error
```

### TreatmentDataService

#### create_category / create_treatment / create_dosage_type

```
동작 흐름 (공통):
1. 입력 검증
2. 중복 확인 (같은 이름)
3. 부모 존재 확인 (treatment → category 존재?, dosage_type → treatment 존재?)
4. DB 저장
5. 생성된 엔티티 반환

에러 처리:
- 부모 미존재 → 404
- 이름 중복 → 409 Conflict
- 입력 누락 → 422 Validation Error
```

---

## Interaction Sequences

### Sequence 1: 예약 확정 → 자동 기입

```
예약시스템 → SQS → EventConsumer → TreatmentRecordService → PostgreSQL
                                         ↓
                                   CycleRuleClient → Admin API → PostgreSQL
                                         ↓
                                   ScheduledTreatmentService → PostgreSQL
```

```
1. [예약시스템] 예약 확정 이벤트 발행 → SQS
2. [EventConsumer] SQS 메시지 수신
3. [EventConsumer] TreatmentType == "surgery" 필터링
4. [EventConsumer] TreatmentRecordService.CreateRecord() 호출
5. [TreatmentRecordService] DB에 시술 기록 저장
6. [TreatmentRecordService] CycleRuleClient.GetCycleRule() 호출
7. [CycleRuleClient] Admin API HTTP GET /api/cycle-rules/{categoryId}
8. [TreatmentRecordService] 예정일 계산 (시술일 + 주기)
9. [ScheduledTreatmentService] DB에 예정일 저장
10. [EventConsumer] 메시지 ACK
```

### Sequence 2: 시술 수동 추가

```
1. [Client] POST /api/records (CreateRecordInput)
2. [Handler] 입력 검증
3. [TreatmentRecordService] CreateRecord() 호출
4. [Repository] DB 저장
5. [CycleRuleClient] 주기 조회
6. [ScheduledTreatmentService] 예정일 계산 및 저장
7. [Handler] 201 Created + TreatmentRecord 반환
```

### Sequence 3: 시술 정보 수정

```
1. [Client] PUT /api/records/{id} (UpdateRecordInput)
2. [Handler] 입력 검증
3. [TreatmentRecordService] UpdateRecord() 호출
4. [Repository] 기존 기록 조회
5. [Service] 소유권 검증
6. [Repository] 기록 업데이트
7. [Service] 날짜/카테고리 변경 감지 시:
   a. 기존 예정일 삭제
   b. 새 주기 조회
   c. 새 예정일 계산 및 저장
8. [Handler] 200 OK + 수정된 TreatmentRecord 반환
```

### Sequence 4: 리마인드 알림 배치

```
1. [Scheduler/Cron] ReminderService.ProcessDueReminders(today) 트리거
2. [ScheduledTreatmentService] Status=PENDING, ScheduledDate <= today 조회
3. 각 항목에 대해:
   a. 알림 메시지 구성
   b. NotificationClient.SendReminder() 호출
   c. Status → REMINDED 업데이트
4. 처리 완료 로그
```

### Sequence 6: 구글 캘린더 내보내기

```
1. [Client] 구글 캘린더 내보내기 버튼 클릭
2. [Client] Google OAuth 인증 (클라이언트 사이드)
3. [Client] GET /api/records/{id} 또는 GET /api/schedules/{id} 로 일정 데이터 조회
4. [Client] Google Calendar API 직접 호출 (이벤트 생성)
5. [Client] 성공/실패 UI 표시

※ 서버는 데이터 제공만, 구글 API 호출은 클라이언트 책임
```
