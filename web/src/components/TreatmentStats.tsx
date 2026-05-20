"use client";

import { Chip } from "@heroui/react";
import type { TreatmentStat } from "@/types";

export interface TreatmentStatsProps {
  statistics: TreatmentStat[];
}

export function TreatmentStats({ statistics }: TreatmentStatsProps) {
  if (statistics.length === 0) return null;

  return (
    <div className="flex gap-2 overflow-x-auto py-2" role="list" aria-label="시술 통계">
      {statistics.map((stat) => (
        <Chip key={stat.category_id} variant="flat" size="sm" role="listitem">
          {stat.category_name} {stat.count}회
        </Chip>
      ))}
    </div>
  );
}
