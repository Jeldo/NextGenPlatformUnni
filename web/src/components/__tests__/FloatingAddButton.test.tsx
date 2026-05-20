import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { FloatingAddButton } from "@/components/FloatingAddButton";

const mockPush = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}));

describe("FloatingAddButton", () => {
  it("'시술 추가하기' 접근성 라벨을 가진다", () => {
    render(<FloatingAddButton />);
    expect(screen.getByLabelText("시술 추가하기")).toBeInTheDocument();
  });

  it("클릭 시 /calendar/records/new로 이동한다", () => {
    render(<FloatingAddButton />);
    screen.getByLabelText("시술 추가하기").click();
    expect(mockPush).toHaveBeenCalledWith("/calendar/records/new");
  });
});
