import { test, expect } from "@playwright/test";

test.describe("캘린더 메인 페이지", () => {
  test("주간 캘린더와 나의 시술 섹션이 표시된다", async ({ page }) => {
    await page.goto("/calendar");
    await expect(page.getByRole("grid", { name: "주간 캘린더" })).toBeVisible();
    await expect(page.getByText("나의 시술")).toBeVisible();
    await expect(page.getByText("이번 달 시술")).toBeVisible();
    await expect(page.getByText("누적 시술")).toBeVisible();
  });

  test("이전/다음 주 이동이 동작한다", async ({ page }) => {
    await page.goto("/calendar");
    const weekLabel = page.locator("text=/\\d+월 \\d+일 — \\d+월 \\d+일/");
    const initialText = await weekLabel.textContent();
    await page.getByLabel("다음 주").click();
    await expect(weekLabel).not.toHaveText(initialText!);
    await page.getByLabel("이전 주").click();
    await expect(weekLabel).toHaveText(initialText!);
  });

  test("날짜 클릭 시 바텀시트가 열린다", async ({ page }) => {
    await page.goto("/calendar");
    const dateButtons = page.getByRole("button").filter({ hasText: /^\d+$/ });
    await dateButtons.first().click();
    await expect(page.getByRole("dialog")).toBeVisible();
    await expect(page.getByText("이 날짜에 시술 추가")).toBeVisible();
  });

  test("FAB 클릭 시 시술 추가 페이지로 이동한다", async ({ page }) => {
    await page.goto("/calendar");
    await page.getByLabel("시술 추가하기").click();
    await expect(page).toHaveURL(/\/calendar\/records\/new/);
  });
});

test.describe("시술 추가 페이지", () => {
  test("필수 필드 미입력 시 저장 버튼이 비활성화된다", async ({ page }) => {
    await page.goto("/calendar/records/new");
    await expect(page.getByRole("button", { name: "저장" })).toBeDisabled();
  });

  test("뒤로가기 버튼이 동작한다", async ({ page }) => {
    await page.goto("/calendar");
    await page.getByLabel("시술 추가하기").click();
    await page.getByRole("button", { name: "‹" }).click();
    await expect(page).toHaveURL(/\/calendar$/);
  });
});

test.describe("시술 상세 페이지", () => {
  const recordUrl = "/calendar/records/019234ab-5678-7def-8000-000000000001";

  test("시술 상세 정보가 표시된다", async ({ page }) => {
    await page.goto(recordUrl);
    await expect(page.getByText("시술 상세")).toBeVisible();
    await expect(page.getByRole("heading")).toContainText("강남언니의원");
    await expect(page.getByText(/1월 15일/)).toBeVisible();
  });

  test("수정하기 버튼 클릭 시 수정 모드로 전환된다", async ({ page }) => {
    await page.goto(recordUrl);
    await page.getByRole("button", { name: "수정하기" }).click();
    await expect(page.getByText("시술 정보 수정")).toBeVisible();
    await expect(page.getByLabel("날짜")).toBeVisible();
    await expect(page.getByLabel("병원명")).toBeVisible();
  });

  test("수정 모드에서 뒤로가기 시 상세 보기로 돌아온다", async ({ page }) => {
    await page.goto(recordUrl);
    await page.getByRole("button", { name: "수정하기" }).click();
    await expect(page.getByText("시술 정보 수정")).toBeVisible();
    await page.getByRole("button", { name: "‹" }).click();
    await expect(page.getByText("시술 상세")).toBeVisible();
  });

  test("··· 메뉴에서 삭제 클릭 시 삭제 모달이 표시된다", async ({ page }) => {
    await page.goto(recordUrl);
    await page.getByRole("button", { name: "···" }).click();
    await page.getByRole("button", { name: "삭제" }).first().click();
    await expect(page.getByText("시술 기록 삭제")).toBeVisible();
    await expect(page.getByText("삭제 후 복구할 수 없습니다")).toBeVisible();
  });

  test("구글 캘린더 등록 버튼이 표시된다", async ({ page }) => {
    await page.goto(recordUrl);
    await expect(page.getByLabel("구글 캘린더에 등록하기")).toBeVisible();
  });
});
