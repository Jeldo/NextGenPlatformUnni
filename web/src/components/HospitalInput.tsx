"use client";

import { Input } from "@heroui/react";

export interface HospitalInputProps {
  value: string;
  onChange: (value: string) => void;
}

export function HospitalInput({ value, onChange }: HospitalInputProps) {
  return (
    <Input
      label="병원명"
      placeholder="병원명을 입력하세요"
      value={value}
      onValueChange={onChange}
      isRequired
      aria-label="병원명"
    />
  );
}
