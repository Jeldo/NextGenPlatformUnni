import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { ScheduleCard } from "@/components/ScheduleCard";
import type { ScheduledTreatment } from "@/types";

const mockSchedule: ScheduledTreatment = {
  id: "sch-1",
  record_id: "rec-1",
  category_id: "cat-1",
  treatment_id: "treat-1",
  scheduled_date: "2026-04-15T00:00:00Z",
  cycle_days: 90,
  status: "PENDING",
};

describe("ScheduleCard", () => {
  it("'예정' 라벨과 주기를 표시한다", () => {
    render(<ScheduleCard schedule={mockSchedule} />);
    expect(screen.getByText("예정")).toBeInTheDocument();
    expect(screen.getByText("90일 주기")).toBeInTheDocument();
  });

  it("예정 날짜를 표시한다", () => {
    render(<ScheduleCard schedule={mockSchedule} />);
    expect(screen.getByText("2026-04-15")).toBeInTheDocument();
  });

  it("onClick 호출 시 schedule을 전달한다", () => {
    const onClick = vi.fn();
    const { container } = render(<ScheduleCard schedule={mockSchedule} onClick={onClick} />);
    const button = container.querySelector("button");
    button?.click();
    expect(onClick).toHaveBeenCalledWith(mockSchedule);
  });
});
