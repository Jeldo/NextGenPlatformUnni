"use client";

import { useRecords } from "@/hooks/useRecords";
import { useSchedules } from "@/hooks/useSchedules";
import { useStatistics } from "@/hooks/useStatistics";

export default function CalendarPage() {
  const { data: records, isLoading: recordsLoading } = useRecords("2026-01-01", "2026-12-31");
  const { data: schedules, isLoading: schedulesLoading } = useSchedules();
  const { data: statistics, isLoading: statsLoading } = useStatistics();

  if (recordsLoading || schedulesLoading || statsLoading) {
    return <div className="p-4">Loading...</div>;
  }

  return (
    <main className="min-h-screen p-4">
      <h1 className="text-2xl font-bold text-primary">시술 캘린더</h1>

      {/* WeeklyCalendarGrid placeholder */}
      <section className="mt-4">
        <h2 className="text-lg font-semibold">시술 기록 ({records?.length ?? 0}건)</h2>
        <ul className="mt-2 space-y-2">
          {records?.map((r) => (
            <li key={r.id} className="border rounded p-3">
              <span className="font-medium">{r.hospital_name}</span>
              <span className="text-gray-description ml-2">{r.treatment_date.slice(0, 10)}</span>
            </li>
          ))}
        </ul>
      </section>

      {/* Schedules */}
      <section className="mt-4">
        <h2 className="text-lg font-semibold">예정일 ({schedules?.length ?? 0}건)</h2>
        <ul className="mt-2 space-y-2">
          {schedules?.map((s) => (
            <li key={s.id} className="border border-dashed rounded p-3">
              <span>{s.scheduled_date.slice(0, 10)}</span>
              <span className="text-gray-description ml-2">{s.status}</span>
            </li>
          ))}
        </ul>
      </section>

      {/* TreatmentStats */}
      <section className="mt-4">
        <h2 className="text-lg font-semibold">시술 통계</h2>
        <ul className="mt-2 flex gap-4">
          {statistics?.map((s) => (
            <li key={s.category_id} className="bg-gray-100 rounded px-3 py-1">
              {s.category_name} {s.count}회
            </li>
          ))}
        </ul>
      </section>
    </main>
  );
}
