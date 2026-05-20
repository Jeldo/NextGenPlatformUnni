import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { DateBottomSheet } from "@/components/DateBottomSheet";
import type { TreatmentRecord, ScheduledTreatment } from "@/types";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

const mockRecord: TreatmentRecord = {
  id: "rec-1", user_id: "u1", source: "MANUAL", category_id: "c1",
  treatment_id: "t1", dosage_type: null, dosage_value: null,
  treatment_date: "2026-01-15T00:00:00Z", hospital_name: "강남언니의원",
  hospital_location: null, doctor_name: null, memo: null,
  created_at: "2026-01-15T09:00:00Z", updated_at: "2026-01-15T09:00:00Z",
};

const mockSchedule: ScheduledTreatment = {
  id: "sch-1", record_id: "rec-1", category_id: "c1", treatment_id: "t1",
  scheduled_date: "2026-01-15T00:00:00Z", cycle_days: 90, status: "PENDING",
};

describe("DateBottomSheet", () => {
  const defaultProps = {
    date: "2026-01-15",
    records: [mockRecord],
    schedules: [mockSchedule],
    isOpen: true,
    onClose: vi.fn(),
    onScheduleClick: vi.fn(),
  };

  it("isOpen=false이면 렌더링하지 않는다", () => {
    const { container } = render(<DateBottomSheet {...defaultProps} isOpen={false} />);
    expect(container.firstChild).toBeNull();
  });

  it("날짜 라벨을 표시한다", () => {
    render(<DateBottomSheet {...defaultProps} />);
    expect(screen.getByText("2026년 1월 15일")).toBeInTheDocument();
  });

  it("시술 추가 버튼을 표시한다", () => {
    render(<DateBottomSheet {...defaultProps} />);
    expect(screen.getByText("+ 이 날짜에 시술 추가")).toBeInTheDocument();
  });

  it("RecordCard를 표시한다", () => {
    render(<DateBottomSheet {...defaultProps} />);
    expect(screen.getByText("강남언니의원")).toBeInTheDocument();
  });

  it("ScheduleCard를 표시한다", () => {
    render(<DateBottomSheet {...defaultProps} />);
    expect(screen.getByText("예정")).toBeInTheDocument();
  });

  it("빈 상태 메시지를 표시한다", () => {
    render(<DateBottomSheet {...defaultProps} records={[]} schedules={[]} />);
    expect(screen.getByText("이 날짜에 시술 기록이 없습니다")).toBeInTheDocument();
  });
});
