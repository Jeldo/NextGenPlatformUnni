import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { HospitalInput } from "@/components/HospitalInput";

describe("HospitalInput", () => {
  it("라벨과 placeholder를 표시한다", () => {
    render(<HospitalInput value="" onChange={vi.fn()} />);
    expect(screen.getByLabelText("병원명")).toBeInTheDocument();
  });

  it("value를 표시한다", () => {
    render(<HospitalInput value="강남언니의원" onChange={vi.fn()} />);
    const input = screen.getByLabelText("병원명") as HTMLInputElement;
    expect(input.value).toBe("강남언니의원");
  });
});
