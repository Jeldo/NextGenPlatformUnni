# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: calendar.spec.ts >> 시술 상세 페이지 >> ··· 메뉴에서 삭제 클릭 시 삭제 모달이 표시된다
- Location: e2e/calendar.spec.ts:77:7

# Error details

```
Test timeout of 30000ms exceeded.
```

```
Error: locator.click: Test timeout of 30000ms exceeded.
Call log:
  - waiting for getByRole('button', { name: '···' })

```

# Test source

```ts
  1  | import { test, expect } from "@playwright/test";
  2  | 
  3  | test.describe("캘린더 메인 페이지", () => {
  4  |   test("주간 캘린더와 나의 시술 섹션이 표시된다", async ({ page }) => {
  5  |     await page.goto("/calendar");
  6  |     await expect(page.getByRole("grid", { name: "주간 캘린더" })).toBeVisible();
  7  |     await expect(page.getByText("나의 시술")).toBeVisible();
  8  |     await expect(page.getByText("이번 달 시술")).toBeVisible();
  9  |     await expect(page.getByText("누적 시술")).toBeVisible();
  10 |   });
  11 | 
  12 |   test("이전/다음 주 이동이 동작한다", async ({ page }) => {
  13 |     await page.goto("/calendar");
  14 |     const weekLabel = page.locator("text=/\\d+월 \\d+일 — \\d+월 \\d+일/");
  15 |     const initialText = await weekLabel.textContent();
  16 |     await page.getByLabel("다음 주").click();
  17 |     await expect(weekLabel).not.toHaveText(initialText!);
  18 |     await page.getByLabel("이전 주").click();
  19 |     await expect(weekLabel).toHaveText(initialText!);
  20 |   });
  21 | 
  22 |   test("날짜 클릭 시 바텀시트가 열린다", async ({ page }) => {
  23 |     await page.goto("/calendar");
  24 |     const dateButtons = page.getByRole("button").filter({ hasText: /^\d+$/ });
  25 |     await dateButtons.first().click();
  26 |     await expect(page.getByRole("dialog")).toBeVisible();
  27 |     await expect(page.getByText("이 날짜에 시술 추가")).toBeVisible();
  28 |   });
  29 | 
  30 |   test("FAB 클릭 시 시술 추가 페이지로 이동한다", async ({ page }) => {
  31 |     await page.goto("/calendar");
  32 |     await page.getByLabel("시술 추가하기").click();
  33 |     await expect(page).toHaveURL(/\/calendar\/records\/new/);
  34 |   });
  35 | });
  36 | 
  37 | test.describe("시술 추가 페이지", () => {
  38 |   test("필수 필드 미입력 시 저장 버튼이 비활성화된다", async ({ page }) => {
  39 |     await page.goto("/calendar/records/new");
  40 |     await expect(page.getByRole("button", { name: "저장" })).toBeDisabled();
  41 |   });
  42 | 
  43 |   test("뒤로가기 버튼이 동작한다", async ({ page }) => {
  44 |     await page.goto("/calendar");
  45 |     await page.getByLabel("시술 추가하기").click();
  46 |     await page.getByRole("button", { name: "‹" }).click();
  47 |     await expect(page).toHaveURL(/\/calendar$/);
  48 |   });
  49 | });
  50 | 
  51 | test.describe("시술 상세 페이지", () => {
  52 |   const recordUrl = "/calendar/records/019234ab-5678-7def-8000-000000000001";
  53 | 
  54 |   test("시술 상세 정보가 표시된다", async ({ page }) => {
  55 |     await page.goto(recordUrl);
  56 |     await expect(page.getByText("시술 상세")).toBeVisible();
  57 |     await expect(page.getByRole("heading")).toContainText("강남언니의원");
  58 |     await expect(page.getByText(/1월 15일/)).toBeVisible();
  59 |   });
  60 | 
  61 |   test("수정하기 버튼 클릭 시 수정 모드로 전환된다", async ({ page }) => {
  62 |     await page.goto(recordUrl);
  63 |     await page.getByRole("button", { name: "수정하기" }).click();
  64 |     await expect(page.getByText("시술 정보 수정")).toBeVisible();
  65 |     await expect(page.getByLabel("날짜")).toBeVisible();
  66 |     await expect(page.getByLabel("병원명")).toBeVisible();
  67 |   });
  68 | 
  69 |   test("수정 모드에서 뒤로가기 시 상세 보기로 돌아온다", async ({ page }) => {
  70 |     await page.goto(recordUrl);
  71 |     await page.getByRole("button", { name: "수정하기" }).click();
  72 |     await expect(page.getByText("시술 정보 수정")).toBeVisible();
  73 |     await page.getByRole("button", { name: "‹" }).click();
  74 |     await expect(page.getByText("시술 상세")).toBeVisible();
  75 |   });
  76 | 
  77 |   test("··· 메뉴에서 삭제 클릭 시 삭제 모달이 표시된다", async ({ page }) => {
  78 |     await page.goto(recordUrl);
> 79 |     await page.getByRole("button", { name: "···" }).click();
     |                                                     ^ Error: locator.click: Test timeout of 30000ms exceeded.
  80 |     await page.getByRole("button", { name: "삭제" }).first().click();
  81 |     await expect(page.getByText("시술 기록 삭제")).toBeVisible();
  82 |     await expect(page.getByText("삭제 후 복구할 수 없습니다")).toBeVisible();
  83 |   });
  84 | 
  85 |   test("구글 캘린더 등록 버튼이 표시된다", async ({ page }) => {
  86 |     await page.goto(recordUrl);
  87 |     await expect(page.getByLabel("구글 캘린더에 등록하기")).toBeVisible();
  88 |   });
  89 | });
  90 | 
```