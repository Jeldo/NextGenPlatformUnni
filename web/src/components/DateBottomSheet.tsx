"use client";

import { useRouter } from "next/navigation";
import { RecordCard } from "@/components/RecordCard";
import { ScheduleCard } from "@/components/ScheduleCard";
import type { TreatmentRecord, ScheduledTreatment } from "@/types";

export interface DateBottomSheetProps {
  date: string;
  records: TreatmentRecord[];
  schedules: ScheduledTreatment[];
  isOpen: boolean;
  onClose: () => void;
  onScheduleClick: (schedule: ScheduledTreatment) => void;
}

export function DateBottomSheet({
  date,
  records,
  schedules,
  isOpen,
  onClose,
  onScheduleClick,
}: DateBottomSheetProps) {
  const router = useRouter();

  if (!isOpen) return null;

  const d = new Date(date);
  const days = ["일", "월", "화", "수", "목", "금", "토"];
  const label = `${d.getFullYear()}년 ${d.getMonth() + 1}월 ${d.getDate()}일 (${days[d.getDay()]})`;

  return (
    <div className="fixed inset-0 z-40" role="dialog" aria-label={`${label} 시술 목록`}>
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/40" onClick={onClose} />

      {/* Sheet */}
      <div className="absolute bottom-0 left-0 right-0 bg-white rounded-t-2xl max-h-[70vh] overflow-y-auto">
        {/* Handle */}
        <div className="w-10 h-1 bg-gray-300 rounded-full mx-auto mt-3" />

        {/* Header */}
        <div className="flex items-center justify-between px-5 pt-4 pb-3">
          <h2 className="text-lg font-bold text-black">{label}</h2>
          <button
            onClick={onClose}
            className="w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center text-gray-description"
            aria-label="닫기"
          >
            ✕
          </button>
        </div>

        <div className="px-5 pb-6">
          {/* 시술 추가 버튼 - 주황 점선 */}
          <button
            onClick={() => router.push(`/calendar/records/new?date=${date}`)}
            className="w-full py-3 border-2 border-dashed border-primary rounded-xl text-primary text-sm font-medium mb-4"
          >
            + 이 날짜에 시술 추가
          </button>

          {/* 카드 목록 */}
          <div className="flex flex-col gap-3">
            {records.map((r) => (
              <RecordCard
                key={r.id}
                record={r}
                onClick={() => router.push(`/calendar/records/${r.id}`)}
              />
            ))}
            {schedules.map((s) => (
              <ScheduleCard key={s.id} schedule={s} onClick={onScheduleClick} />
            ))}
            {records.length === 0 && schedules.length === 0 && (
              <p className="text-center text-gray-description py-4 text-sm">이 날짜에 시술 기록이 없습니다</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
