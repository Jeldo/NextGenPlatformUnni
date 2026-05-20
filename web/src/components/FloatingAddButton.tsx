"use client";

import { Button } from "@heroui/react";
import { useRouter } from "next/navigation";

export function FloatingAddButton() {
  const router = useRouter();

  return (
    <Button
      isIconOnly
      color="primary"
      size="lg"
      className="fixed bottom-6 right-4 z-50 w-14 h-14 rounded-full shadow-lg"
      aria-label="시술 추가하기"
      onPress={() => router.push("/calendar/records/new")}
    >
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
        <path d="M12 5v14M5 12h14" />
      </svg>
    </Button>
  );
}
