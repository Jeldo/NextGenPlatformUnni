import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { RecordCard } from "@/components/RecordCard";
import type { TreatmentRecord } from "@/types";

const mockRecord: TreatmentRecord = {
  id: "rec-1", user_id: "user-1", source: "MANUAL",
  category_id: "cat-1", treatment_id: "treat-1",
  dosage_type: "volume", dosage_value: "10.0",
  treatment_date: "2026-01-15T00:00:00Z", hospital_name: "강남언니의원",
  hospital_location: "서울 강남구", doctor_name: "김의사",
  memo: "사각턱 보톡스", created_at: "2026-01-15T09:00:00Z", updated_at: "2026-01-15T09:00:00Z",
};

describe("RecordCard", () => {
  it("병원명을 표시한다", () => {
    render(<RecordCard record={mockRecord} />);
    expect(screen.getByText("강남언니의원")).toBeInTheDocument();
  });

  it("날짜를 표시한다", () => {
    render(<RecordCard record={mockRecord} />);
    expect(screen.getByText("2026.01.15")).toBeInTheDocument();
  });

  it("'확정' 뱃지를 표시한다", () => {
    render(<RecordCard record={mockRecord} />);
    expect(screen.getByText("확정")).toBeInTheDocument();
  });

  it("클릭 시 onClick을 호출한다", () => {
    const onClick = vi.fn();
    render(<RecordCard record={mockRecord} onClick={onClick} />);
    screen.getByText("강남언니의원").click();
    expect(onClick).toHaveBeenCalledWith(mockRecord);
  });
});
