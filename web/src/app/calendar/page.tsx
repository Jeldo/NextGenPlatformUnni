"use client";

import { useState } from "react";
import { useRecords } from "@/hooks/useRecords";
import { useSchedules } from "@/hooks/useSchedules";
import { useStatistics } from "@/hooks/useStatistics";
import { useCompleteSchedule } from "@/hooks/useCompleteSchedule";
import { useDeleteSchedule } from "@/hooks/useDeleteSchedule";
import { WeeklyCalendarGrid } from "@/components/WeeklyCalendarGrid";
import { DateBottomSheet } from "@/components/DateBottomSheet";
import { FloatingAddButton } from "@/components/FloatingAddButton";
import { ScheduleConfirmModal } from "@/components/ScheduleConfirmModal";
import type { ScheduledTreatment } from "@/types";

export default function CalendarPage() {
  const { data: records = [] } = useRecords("2026-01-01", "2026-12-31");
  const { data: schedules = [] } = useSchedules();
  const { data: statistics = [] } = useStatistics();
  const completeSchedule = useCompleteSchedule();
  const deleteSchedule = useDeleteSchedule();

  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [selectedSchedule, setSelectedSchedule] = useState<ScheduledTreatment | null>(null);

  const dateRecords = records.filter((r) => r.treatment_date.slice(0, 10) === selectedDate);
  const dateSchedules = schedules.filter((s) => s.scheduled_date.slice(0, 10) === selectedDate);

  const handleScheduleClick = (schedule: ScheduledTreatment) => {
    setSelectedSchedule(schedule);
  };

  const handleComplete = (schedule: ScheduledTreatment) => {
    completeSchedule.mutate(schedule.id, { onSuccess: () => setSelectedSchedule(null) });
  };

  const handleDeleteSchedule = (schedule: ScheduledTreatment) => {
    deleteSchedule.mutate(schedule.id, { onSuccess: () => setSelectedSchedule(null) });
  };

  // 통계 계산
  const totalCount = statistics.reduce((sum, s) => sum + s.count, 0);
  const topCategory = statistics.length > 0 ? statistics.reduce((a, b) => a.count > b.count ? a : b) : null;
  const lastRecord = records.length > 0 ? records.reduce((a, b) => a.treatment_date > b.treatment_date ? a : b) : null;
  const nextSchedule = schedules.filter(s => s.status === "PENDING").sort((a, b) => a.scheduled_date.localeCompare(b.scheduled_date))[0];

  const now = new Date();
  const monthLabel = `${now.getFullYear()}년 ${now.getMonth() + 1}월`;

  const formatShortDate = (dateStr: string) => {
    const d = new Date(dateStr);
    return `${d.getMonth() + 1}월 ${d.getDate()}일`;
  };

  return (
    <main className="min-h-screen bg-[#f5f5f5] pb-20">
      {/* Month header */}
      <div className="bg-white px-5 pt-5 pb-3">
        <h1 className="text-xl font-bold text-black">{monthLabel}</h1>
      </div>

      {/* Weekly Calendar */}
      <div className="bg-white px-4 pb-4 border-b border-black/5">
        <WeeklyCalendarGrid
          records={records}
          schedules={schedules}
          onDateSelect={setSelectedDate}
        />
      </div>

      {/* 나의 시술 섹션 */}
      <div className="bg-white mt-2 px-5 py-5">
        <h2 className="text-lg font-bold text-black mb-4">나의 시술</h2>

        {/* 이번 달 / 누적 카드 */}
        <div className="grid grid-cols-2 gap-3 mb-4">
          <div className="bg-[#f5f5f5] rounded-xl p-4">
            <p className="text-xs text-gray-description mb-1">이번 달 시술</p>
            <p className="text-2xl font-bold text-black">{records.length}<span className="text-sm font-normal">회</span></p>
          </div>
          <div className="bg-[#f5f5f5] rounded-xl p-4">
            <p className="text-xs text-gray-description mb-1">누적 시술</p>
            <p className="text-2xl font-bold text-black">{totalCount}<span className="text-sm font-normal">회</span></p>
          </div>
        </div>

        {/* 요약 행 */}
        <div className="flex flex-col gap-0">
          {lastRecord && (
            <div className="flex justify-between py-3 border-b border-black/5">
              <span className="text-sm text-black">마지막 시술</span>
              <span className="text-sm text-gray-description">{formatShortDate(lastRecord.treatment_date)}</span>
            </div>
          )}
          {nextSchedule && (
            <div className="flex justify-between py-3 border-b border-black/5">
              <span className="text-sm text-black">다음 예정일</span>
              <span className="text-sm text-gray-description">{formatShortDate(nextSchedule.scheduled_date)}</span>
            </div>
          )}
          {topCategory && (
            <div className="flex justify-between py-3">
              <span className="text-sm text-black">가장 많은 시술</span>
              <span className="text-sm text-gray-description">{topCategory.category_name} ({topCategory.count}회)</span>
            </div>
          )}
        </div>
      </div>

      <FloatingAddButton />

      <DateBottomSheet
        date={selectedDate ?? ""}
        records={dateRecords}
        schedules={dateSchedules}
        isOpen={!!selectedDate}
        onClose={() => setSelectedDate(null)}
        onScheduleClick={handleScheduleClick}
      />

      <ScheduleConfirmModal
        schedule={selectedSchedule}
        isOpen={!!selectedSchedule}
        onClose={() => setSelectedSchedule(null)}
        onComplete={handleComplete}
        onDelete={handleDeleteSchedule}
        isLoading={completeSchedule.isPending || deleteSchedule.isPending}
      />
    </main>
  );
}
