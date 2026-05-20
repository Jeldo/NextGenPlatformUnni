"use client";

import { useState, useMemo } from "react";
import type { TreatmentRecord, ScheduledTreatment } from "@/types";

export interface WeeklyCalendarGridProps {
  records: TreatmentRecord[];
  schedules: ScheduledTreatment[];
  onDateSelect: (date: string) => void;
}

const DAY_LABELS = ["일", "월", "화", "수", "목", "금", "토"];

function getWeekStart(date: Date): Date {
  const d = new Date(date);
  d.setDate(d.getDate() - d.getDay());
  d.setHours(0, 0, 0, 0);
  return d;
}

function formatDate(date: Date): string {
  return date.toISOString().slice(0, 10);
}

export function WeeklyCalendarGrid({ records, schedules, onDateSelect }: WeeklyCalendarGridProps) {
  const [weekStart, setWeekStart] = useState(() => getWeekStart(new Date()));

  const weekDates = useMemo(() => {
    return Array.from({ length: 7 }, (_, i) => {
      const d = new Date(weekStart);
      d.setDate(d.getDate() + i);
      return d;
    });
  }, [weekStart]);

  const recordsByDate = useMemo(() => {
    const map: Record<string, TreatmentRecord[]> = {};
    records.forEach((r) => {
      const key = r.treatment_date.slice(0, 10);
      (map[key] ??= []).push(r);
    });
    return map;
  }, [records]);

  const schedulesByDate = useMemo(() => {
    const map: Record<string, ScheduledTreatment[]> = {};
    schedules.forEach((s) => {
      const key = s.scheduled_date.slice(0, 10);
      (map[key] ??= []).push(s);
    });
    return map;
  }, [schedules]);

  const prevWeek = () => setWeekStart((prev) => { const d = new Date(prev); d.setDate(d.getDate() - 7); return d; });
  const nextWeek = () => setWeekStart((prev) => { const d = new Date(prev); d.setDate(d.getDate() + 7); return d; });

  const today = formatDate(new Date());

  return (
    <div className="w-full" role="grid" aria-label="주간 캘린더">
      {/* Header */}
      <div className="flex items-center justify-between py-2">
        <div className="flex items-center gap-1">
          <button onClick={prevWeek} className="p-2 min-w-[44px] min-h-[44px] text-sm" aria-label="이전 주">‹</button>
          <span className="text-sm font-medium">
            {weekDates[0].getMonth() + 1}월 {weekDates[0].getDate()}일 — {weekDates[6].getMonth() + 1}월 {weekDates[6].getDate()}일
          </span>
          <button onClick={nextWeek} className="p-2 min-w-[44px] min-h-[44px] text-sm" aria-label="다음 주">›</button>
        </div>
        {/* 범례 */}
        <div className="flex items-center gap-2 text-[11px] text-gray-description">
          <span className="flex items-center gap-1"><span className="w-2 h-2 bg-black rounded-sm" />확정</span>
          <span className="flex items-center gap-1"><span className="w-2 h-2 border border-dashed border-black rounded-sm" />예정</span>
        </div>
      </div>

      {/* Day labels */}
      <div className="grid grid-cols-7 text-center text-[11px] text-gray-description mb-1 font-medium">
        {DAY_LABELS.map((label) => (
          <span key={label}>{label}</span>
        ))}
      </div>

      {/* Date cells */}
      <div className="grid grid-cols-7 gap-px min-h-[100px]">
        {weekDates.map((date) => {
          const key = formatDate(date);
          const dayRecords = recordsByDate[key] || [];
          const daySchedules = schedulesByDate[key] || [];
          const isToday = key === today;

          return (
            <button
              key={key}
              onClick={() => onDateSelect(key)}
              className="flex flex-col items-center p-1 min-h-[90px] rounded-md hover:bg-gray-50"
              aria-label={`${date.getMonth() + 1}월 ${date.getDate()}일`}
            >
              <span className={`text-sm mb-1.5 w-7 h-7 flex items-center justify-center rounded-full ${
                isToday ? "bg-primary text-white font-bold" : "text-black"
              }`}>
                {date.getDate()}
              </span>
              <div className="flex flex-col gap-0.5 w-full overflow-hidden">
                {dayRecords.slice(0, 2).map((r) => (
                  <div key={r.id} className="text-[9px] bg-black text-white rounded px-1 py-0.5 truncate leading-tight">
                    {r.hospital_name.slice(0, 3)}…
                  </div>
                ))}
                {daySchedules.slice(0, 2).map((s) => (
                  <div key={s.id} className="text-[9px] border border-dashed border-black rounded px-1 py-0.5 truncate leading-tight text-black">
                    예정…
                  </div>
                ))}
              </div>
            </button>
          );
        })}
      </div>
    </div>
  );
}
