"use client";

import { Card, CardBody } from "@heroui/react";
import type { TreatmentRecord } from "@/types";

export interface RecordCardProps {
  record: TreatmentRecord;
  onClick?: (record: TreatmentRecord) => void;
}

export function RecordCard({ record, onClick }: RecordCardProps) {
  return (
    <Card
      isPressable={!!onClick}
      onPress={() => onClick?.(record)}
      className="border border-solid border-gray-200 shadow-sm"
    >
      <CardBody className="p-3 gap-1">
        <div className="flex items-center justify-between">
          <span className="text-sm font-semibold text-black">
            {record.hospital_name}
          </span>
          <span className="text-xs text-gray-description">
            {record.source === "AUTO" ? "자동" : "수동"}
          </span>
        </div>
        <p className="text-sm text-gray-description">
          {record.treatment_date.slice(0, 10)}
        </p>
        {record.memo && (
          <p className="text-xs text-gray-description truncate">{record.memo}</p>
        )}
      </CardBody>
    </Card>
  );
}
