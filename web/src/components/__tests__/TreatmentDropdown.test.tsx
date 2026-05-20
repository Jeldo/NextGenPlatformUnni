import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { TreatmentDropdown } from "@/components/TreatmentDropdown";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

vi.mock("@/hooks/useCategories", () => ({
  useCategories: () => ({
    data: [
      { id: "cat-1", name: "보톡스" },
      { id: "cat-2", name: "필러" },
    ],
  }),
}));

vi.mock("@/hooks/useTreatments", () => ({
  useTreatments: () => ({
    data: [{ id: "treat-1", category_id: "cat-1", name: "사각턱" }],
  }),
}));

vi.mock("@/hooks/useDosageTypes", () => ({
  useDosageTypes: () => ({
    data: [{ id: "dose-1", treatment_id: "treat-1", unit: "cc" }],
  }),
}));

function wrapper({ children }: { children: React.ReactNode }) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return <QueryClientProvider client={qc}>{children}</QueryClientProvider>;
}

describe("TreatmentDropdown", () => {
  const defaultProps = {
    categoryId: "",
    treatmentId: "",
    dosageType: "",
    onCategoryChange: vi.fn(),
    onTreatmentChange: vi.fn(),
    onDosageTypeChange: vi.fn(),
  };

  it("카테고리 셀렉트를 렌더링한다", () => {
    render(<TreatmentDropdown {...defaultProps} />, { wrapper });
    expect(screen.getAllByLabelText("시술 카테고리").length).toBeGreaterThan(0);
  });

  it("시술명 셀렉트를 렌더링한다", () => {
    render(<TreatmentDropdown {...defaultProps} />, { wrapper });
    expect(screen.getAllByLabelText("시술명").length).toBeGreaterThan(0);
  });

  it("용량 단위 셀렉트를 렌더링한다 (dosageTypes 있을 때)", () => {
    render(<TreatmentDropdown {...defaultProps} treatmentId="treat-1" />, { wrapper });
    expect(screen.getAllByLabelText("용량 단위").length).toBeGreaterThan(0);
  });
});
