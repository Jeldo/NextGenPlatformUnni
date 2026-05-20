import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { WeeklyCalendarGrid } from "@/components/WeeklyCalendarGrid";
import type { TreatmentRecord, ScheduledTreatment } from "@/types";

const mockRecords: TreatmentRecord[] = [
  {
    id: "rec-1", user_id: "u1", source: "MANUAL", category_id: "c1",
    treatment_id: "t1", dosage_type: null, dosage_value: null,
    treatment_date: "2026-01-15T00:00:00Z", hospital_name: "강남언니의원",
    hospital_location: null, doctor_name: null, memo: null,
    created_at: "2026-01-15T09:00:00Z", updated_at: "2026-01-15T09:00:00Z",
  },
];

const mockSchedules: ScheduledTreatment[] = [
  {
    id: "sch-1", record_id: "rec-1", category_id: "c1", treatment_id: "t1",
    scheduled_date: "2026-04-15T00:00:00Z", cycle_days: 90, status: "PENDING",
  },
];

describe("WeeklyCalendarGrid", () => {
  it("주간 캘린더 grid를 렌더링한다", () => {
    render(<WeeklyCalendarGrid records={[]} schedules={[]} onDateSelect={vi.fn()} />);
    expect(screen.getByRole("grid", { name: "주간 캘린더" })).toBeInTheDocument();
  });

  it("요일 라벨을 표시한다", () => {
    render(<WeeklyCalendarGrid records={[]} schedules={[]} onDateSelect={vi.fn()} />);
    expect(screen.getByText("월")).toBeInTheDocument();
    expect(screen.getByText("일")).toBeInTheDocument();
  });

  it("이전/다음 주 버튼이 있다", () => {
    render(<WeeklyCalendarGrid records={[]} schedules={[]} onDateSelect={vi.fn()} />);
    expect(screen.getByLabelText("이전 주")).toBeInTheDocument();
    expect(screen.getByLabelText("다음 주")).toBeInTheDocument();
  });

  it("날짜 클릭 시 onDateSelect를 호출한다", () => {
    const onDateSelect = vi.fn();
    render(<WeeklyCalendarGrid records={mockRecords} schedules={mockSchedules} onDateSelect={onDateSelect} />);
    // 7개 날짜 버튼 중 하나 클릭
    const buttons = screen.getAllByRole("button").filter(b => b.getAttribute("aria-label")?.includes("월"));
    if (buttons.length > 0) buttons[0].click();
    expect(onDateSelect).toHaveBeenCalled();
  });
});
