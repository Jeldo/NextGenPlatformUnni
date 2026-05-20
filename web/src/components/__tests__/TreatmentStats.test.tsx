import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { TreatmentStats } from "@/components/TreatmentStats";
import type { TreatmentStat } from "@/types";

const mockStats: TreatmentStat[] = [
  { category_id: "cat-1", category_name: "보톡스", count: 5 },
  { category_id: "cat-2", category_name: "필러", count: 3 },
];

describe("TreatmentStats", () => {
  it("시술별 횟수를 표시한다", () => {
    render(<TreatmentStats statistics={mockStats} />);
    expect(screen.getByText("보톡스 5회")).toBeInTheDocument();
    expect(screen.getByText("필러 3회")).toBeInTheDocument();
  });

  it("빈 배열이면 아무것도 렌더링하지 않는다", () => {
    const { container } = render(<TreatmentStats statistics={[]} />);
    expect(container.firstChild).toBeNull();
  });

  it("접근성: list role을 가진다", () => {
    render(<TreatmentStats statistics={mockStats} />);
    expect(screen.getByRole("list", { name: "시술 통계" })).toBeInTheDocument();
  });
});
