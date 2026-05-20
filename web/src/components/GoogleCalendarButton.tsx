"use client";

import { Button } from "@heroui/react";

export interface GoogleCalendarButtonProps {
  title: string;
  date: string;
  isSchedule?: boolean;
}

export function GoogleCalendarButton({ title, date, isSchedule }: GoogleCalendarButtonProps) {
  const handleExport = () => {
    const eventTitle = isSchedule ? `(예정) ${title}` : title;
    const dateStr = date.slice(0, 10).replace(/-/g, "");
    const url = `https://calendar.google.com/calendar/render?action=TEMPLATE&text=${encodeURIComponent(eventTitle)}&dates=${dateStr}/${dateStr}`;
    window.open(url, "_blank");
  };

  return (
    <Button
      variant="bordered"
      size="lg"
      className="flex-1"
      onPress={handleExport}
      aria-label="구글 캘린더에 등록하기"
    >
      구글 캘린더 등록
    </Button>
  );
}
