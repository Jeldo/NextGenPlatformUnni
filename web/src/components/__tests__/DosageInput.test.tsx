import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { DosageInput } from "@/components/DosageInput";

describe("DosageInput", () => {
  it("unit이 있으면 렌더링한다", () => {
    render(<DosageInput value="10" unit="cc" onChange={vi.fn()} />);
    expect(screen.getByLabelText("용량 (cc)")).toBeInTheDocument();
    expect(screen.getByText("cc")).toBeInTheDocument();
  });

  it("unit이 빈 문자열이면 렌더링하지 않는다", () => {
    const { container } = render(<DosageInput value="" unit="" onChange={vi.fn()} />);
    expect(container.firstChild).toBeNull();
  });

  it("value를 표시한다", () => {
    render(<DosageInput value="10" unit="cc" onChange={vi.fn()} />);
    const input = screen.getByLabelText("용량 (cc)") as HTMLInputElement;
    expect(input.value).toBe("10");
  });
});
