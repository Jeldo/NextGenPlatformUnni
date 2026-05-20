"use client";

import { useRouter } from "next/navigation";
import { Button } from "@heroui/react";
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
  const label = `${d.getFullYear()}년 ${d.getMonth() + 1}월 ${d.getDate()}일`;

  return (
    <div className="fixed inset-0 z-40" role="dialog" aria-label={`${label} 시술 목록`}>
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/30" onClick={onClose} />

      {/* Sheet */}
      <div className="absolute bottom-0 left-0 right-0 bg-white rounded-t-2xl p-4 max-h-[70vh] overflow-y-auto">
        <div className="w-10 h-1 bg-gray-300 rounded-full mx-auto mb-4" />

        <h2 className="text-lg font-bold mb-3">{label}</h2>

        <Button
          variant="bordered"
          color="primary"
          className="w-full mb-4"
          onPress={() => router.push(`/calendar/records/new?date=${date}`)}
        >
          + 이 날짜에 시술 추가
        </Button>

        <div className="flex flex-col gap-2">
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
            <p className="text-center text-gray-description py-4">이 날짜에 시술 기록이 없습니다</p>
          )}
        </div>
      </div>
    </div>
  );
}
