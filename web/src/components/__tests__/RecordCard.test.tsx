import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { RecordCard } from "@/components/RecordCard";
import type { TreatmentRecord } from "@/types";

const mockRecord: TreatmentRecord = {
  id: "rec-1",
  user_id: "user-1",
  source: "MANUAL",
  category_id: "cat-1",
  treatment_id: "treat-1",
  dosage_type: "volume",
  dosage_value: "10.0",
  treatment_date: "2026-01-15T00:00:00Z",
  hospital_name: "강남언니의원",
  hospital_location: "서울 강남구",
  doctor_name: "김의사",
  memo: "사각턱 보톡스",
  created_at: "2026-01-15T09:00:00Z",
  updated_at: "2026-01-15T09:00:00Z",
};

describe("RecordCard", () => {
  it("병원명과 날짜를 표시한다", () => {
    render(<RecordCard record={mockRecord} />);
    expect(screen.getByText("강남언니의원")).toBeInTheDocument();
    expect(screen.getByText("2026-01-15")).toBeInTheDocument();
  });

  it("수동 기입 시 '수동' 라벨을 표시한다", () => {
    render(<RecordCard record={mockRecord} />);
    expect(screen.getByText("수동")).toBeInTheDocument();
  });

  it("자동 기입 시 '자동' 라벨을 표시한다", () => {
    render(<RecordCard record={{ ...mockRecord, source: "AUTO" }} />);
    expect(screen.getByText("자동")).toBeInTheDocument();
  });

  it("메모가 있으면 표시한다", () => {
    render(<RecordCard record={mockRecord} />);
    expect(screen.getByText("사각턱 보톡스")).toBeInTheDocument();
  });

  it("메모가 없으면 메모 영역이 없다", () => {
    render(<RecordCard record={{ ...mockRecord, memo: null }} />);
    expect(screen.queryByText("사각턱 보톡스")).not.toBeInTheDocument();
  });

  it("onClick 호출 시 record를 전달한다", async () => {
    const onClick = vi.fn();
    const { container } = render(<RecordCard record={mockRecord} onClick={onClick} />);
    const button = container.querySelector("button");
    button?.click();
    expect(onClick).toHaveBeenCalledWith(mockRecord);
  });
});
