import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { ScheduleConfirmModal } from "@/components/ScheduleConfirmModal";
import type { ScheduledTreatment } from "@/types";

const mockSchedule: ScheduledTreatment = {
  id: "sch-1", record_id: "rec-1", category_id: "c1", treatment_id: "t1",
  scheduled_date: "2026-04-15T00:00:00Z", cycle_days: 90, status: "PENDING",
};

describe("ScheduleConfirmModal", () => {
  const defaultProps = {
    schedule: mockSchedule,
    isOpen: true,
    onClose: vi.fn(),
    onComplete: vi.fn(),
    onDelete: vi.fn(),
  };

  it("모달 제목을 표시한다", () => {
    render(<ScheduleConfirmModal {...defaultProps} />);
    expect(screen.getByText("이 날짜에 시술을 받으셨나요?")).toBeInTheDocument();
  });

  it("예정 날짜를 표시한다", () => {
    render(<ScheduleConfirmModal {...defaultProps} />);
    expect(screen.getByText("2026-04-15 예정된 시술")).toBeInTheDocument();
  });

  it("'받았어요' 버튼 클릭 시 onComplete를 호출한다", () => {
    render(<ScheduleConfirmModal {...defaultProps} />);
    screen.getByText("받았어요").click();
    expect(defaultProps.onComplete).toHaveBeenCalledWith(mockSchedule);
  });

  it("'삭제' 버튼 클릭 시 onDelete를 호출한다", () => {
    render(<ScheduleConfirmModal {...defaultProps} />);
    screen.getByText("삭제").click();
    expect(defaultProps.onDelete).toHaveBeenCalledWith(mockSchedule);
  });

  it("schedule이 null이면 렌더링하지 않는다", () => {
    const { container } = render(<ScheduleConfirmModal {...defaultProps} schedule={null} />);
    expect(container.firstChild).toBeNull();
  });
});
