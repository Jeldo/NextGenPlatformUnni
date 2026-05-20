"use client";

import type { TreatmentRecord } from "@/types";

export interface RecordCardProps {
  record: TreatmentRecord;
  onClick?: (record: TreatmentRecord) => void;
}

export function RecordCard({ record, onClick }: RecordCardProps) {
  return (
    <button
      onClick={() => onClick?.(record)}
      className="w-full text-left border border-black/15 rounded-xl p-4 hover:bg-gray-50 transition-colors"
    >
      <p className="text-[15px] font-semibold text-black">
        {record.hospital_name}
      </p>
      <p className="text-[13px] text-gray-description mt-1">
        {record.treatment_date.slice(0, 10).replace(/-/g, ".")}
      </p>
      <span className="inline-block mt-2 text-[11px] px-2 py-0.5 bg-gray-100 rounded text-gray-description">
        확정
      </span>
    </button>
  );
}
