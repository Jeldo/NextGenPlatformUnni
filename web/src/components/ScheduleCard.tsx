"use client";

import type { ScheduledTreatment } from "@/types";

export interface ScheduleCardProps {
  schedule: ScheduledTreatment;
  onClick?: (schedule: ScheduledTreatment) => void;
}

export function ScheduleCard({ schedule, onClick }: ScheduleCardProps) {
  return (
    <button
      onClick={() => onClick?.(schedule)}
      className="w-full text-left border border-dashed border-black/20 rounded-xl p-4 hover:bg-gray-50 transition-colors"
    >
      <p className="text-[15px] font-semibold text-black">
        예정
      </p>
      <p className="text-[13px] text-gray-description mt-1">
        추천 주기 기반
      </p>
      <span className="inline-block mt-2 text-[11px] px-2 py-0.5 bg-gray-100 rounded text-gray-description">
        예정 — 탭하여 확인
      </span>
    </button>
  );
}
