"use client";

import { Input } from "@heroui/react";

export interface DosageInputProps {
  value: string;
  unit: string;
  onChange: (value: string) => void;
}

export function DosageInput({ value, unit, onChange }: DosageInputProps) {
  if (!unit) return null;

  return (
    <Input
      type="number"
      label="용량"
      placeholder="숫자 입력"
      value={value}
      onValueChange={onChange}
      endContent={<span className="text-sm text-gray-description">{unit}</span>}
      aria-label={`용량 (${unit})`}
    />
  );
}
