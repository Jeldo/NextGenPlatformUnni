import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { ScheduleCard } from "@/components/ScheduleCard";
import type { ScheduledTreatment } from "@/types";

const mockSchedule: ScheduledTreatment = {
  id: "sch-1", record_id: "rec-1", category_id: "c1", treatment_id: "t1",
  scheduled_date: "2026-04-15T00:00:00Z", cycle_days: 90, status: "PENDING",
};

describe("ScheduleCard", () => {
  it("'예정' 텍스트를 표시한다", () => {
    render(<ScheduleCard schedule={mockSchedule} />);
    expect(screen.getByText("예정")).toBeInTheDocument();
  });

  it("'추천 주기 기반' 설명을 표시한다", () => {
    render(<ScheduleCard schedule={mockSchedule} />);
    expect(screen.getByText("추천 주기 기반")).toBeInTheDocument();
  });

  it("'예정 — 탭하여 확인' 뱃지를 표시한다", () => {
    render(<ScheduleCard schedule={mockSchedule} />);
    expect(screen.getByText("예정 — 탭하여 확인")).toBeInTheDocument();
  });

  it("클릭 시 onClick을 호출한다", () => {
    const onClick = vi.fn();
    render(<ScheduleCard schedule={mockSchedule} onClick={onClick} />);
    screen.getByText("예정 — 탭하여 확인").click();
    expect(onClick).toHaveBeenCalledWith(mockSchedule);
  });
});
