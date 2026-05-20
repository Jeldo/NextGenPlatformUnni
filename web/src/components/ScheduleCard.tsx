"use client";

import { Card, CardBody } from "@heroui/react";
import type { ScheduledTreatment } from "@/types";

export interface ScheduleCardProps {
  schedule: ScheduledTreatment;
  onClick?: (schedule: ScheduledTreatment) => void;
}

export function ScheduleCard({ schedule, onClick }: ScheduleCardProps) {
  return (
    <Card
      isPressable={!!onClick}
      onPress={() => onClick?.(schedule)}
      className="border border-dashed border-primary/50 shadow-sm"
    >
      <CardBody className="p-3 gap-1">
        <div className="flex items-center justify-between">
          <span className="text-sm font-semibold text-primary">예정</span>
          <span className="text-xs text-gray-description">
            {schedule.cycle_days}일 주기
          </span>
        </div>
        <p className="text-sm text-gray-description">
          {schedule.scheduled_date.slice(0, 10)}
        </p>
      </CardBody>
    </Card>
  );
}
