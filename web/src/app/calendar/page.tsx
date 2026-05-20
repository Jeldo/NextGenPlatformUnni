"use client";

import { useState } from "react";
import { useRecords } from "@/hooks/useRecords";
import { useSchedules } from "@/hooks/useSchedules";
import { useStatistics } from "@/hooks/useStatistics";
import { WeeklyCalendarGrid } from "@/components/WeeklyCalendarGrid";
import { DateBottomSheet } from "@/components/DateBottomSheet";
import { TreatmentStats } from "@/components/TreatmentStats";
import { FloatingAddButton } from "@/components/FloatingAddButton";
import type { ScheduledTreatment } from "@/types";

export default function CalendarPage() {
  const { data: records = [] } = useRecords("2026-01-01", "2026-12-31");
  const { data: schedules = [] } = useSchedules();
  const { data: statistics = [] } = useStatistics();

  const [selectedDate, setSelectedDate] = useState<string | null>(null);

  const dateRecords = records.filter((r) => r.treatment_date.slice(0, 10) === selectedDate);
  const dateSchedules = schedules.filter((s) => s.scheduled_date.slice(0, 10) === selectedDate);

  const handleScheduleClick = (_schedule: ScheduledTreatment) => {
    // Phase 4에서 ScheduleConfirmModal 연결
  };

  return (
    <main className="min-h-screen p-4 pb-20">
      <h1 className="text-xl font-bold mb-4">시술 캘린더</h1>

      <WeeklyCalendarGrid
        records={records}
        schedules={schedules}
        onDateSelect={setSelectedDate}
      />

      <TreatmentStats statistics={statistics} />

      <DateBottomSheet
        date={selectedDate ?? ""}
        records={dateRecords}
        schedules={dateSchedules}
        isOpen={!!selectedDate}
        onClose={() => setSelectedDate(null)}
        onScheduleClick={handleScheduleClick}
      />

      <FloatingAddButton />
    </main>
  );
}
