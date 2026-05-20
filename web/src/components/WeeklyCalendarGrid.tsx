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

  const prevWeek = () => {
    const d = new Date(weekStart);
    d.setDate(d.getDate() - 7);
    setWeekStart(d);
  };

  const nextWeek = () => {
    const d = new Date(weekStart);
    d.setDate(d.getDate() + 7);
    setWeekStart(d);
  };

  const today = formatDate(new Date());

  return (
    <div className="w-full" role="grid" aria-label="주간 캘린더">
      {/* Header: 주 이동 */}
      <div className="flex items-center justify-between mb-2">
        <button onClick={prevWeek} className="p-2 min-w-[44px] min-h-[44px]" aria-label="이전 주">
          ‹
        </button>
        <span className="text-sm font-semibold">
          {weekDates[0].getMonth() + 1}월 {weekDates[0].getDate()}일 ~ {weekDates[6].getMonth() + 1}월 {weekDates[6].getDate()}일
        </span>
        <button onClick={nextWeek} className="p-2 min-w-[44px] min-h-[44px]" aria-label="다음 주">
          ›
        </button>
      </div>

      {/* Day labels */}
      <div className="grid grid-cols-7 text-center text-xs text-gray-description mb-1">
        {DAY_LABELS.map((label) => (
          <span key={label}>{label}</span>
        ))}
      </div>

      {/* Date cells */}
      <div className="grid grid-cols-7 gap-px min-h-[120px]">
        {weekDates.map((date) => {
          const key = formatDate(date);
          const dayRecords = recordsByDate[key] || [];
          const daySchedules = schedulesByDate[key] || [];
          const isToday = key === today;

          return (
            <button
              key={key}
              onClick={() => onDateSelect(key)}
              className={`flex flex-col items-center p-1 min-h-[100px] rounded-md ${
                isToday ? "bg-primary/10" : "hover:bg-gray-50"
              }`}
              aria-label={`${date.getMonth() + 1}월 ${date.getDate()}일`}
            >
              <span className={`text-xs mb-1 ${isToday ? "text-primary font-bold" : ""}`}>
                {date.getDate()}
              </span>
              <div className="flex flex-col gap-0.5 w-full overflow-hidden">
                {dayRecords.slice(0, 2).map((r) => (
                  <div key={r.id} className="text-[10px] bg-primary/10 border border-primary/30 rounded px-0.5 truncate">
                    {r.hospital_name}
                  </div>
                ))}
                {daySchedules.slice(0, 2).map((s) => (
                  <div key={s.id} className="text-[10px] border border-dashed border-primary/50 rounded px-0.5 truncate text-primary">
                    예정
                  </div>
                ))}
                {(dayRecords.length + daySchedules.length > 2) && (
                  <span className="text-[9px] text-gray-description">+{dayRecords.length + daySchedules.length - 2}</span>
                )}
              </div>
            </button>
          );
        })}
      </div>
    </div>
  );
}
