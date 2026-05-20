"use client";

import { Select, SelectItem } from "@heroui/react";
import { useCategories } from "@/hooks/useCategories";
import { useTreatments } from "@/hooks/useTreatments";
import { useDosageTypes } from "@/hooks/useDosageTypes";

export interface TreatmentDropdownProps {
  categoryId: string;
  treatmentId: string;
  dosageType: string;
  onCategoryChange: (id: string) => void;
  onTreatmentChange: (id: string) => void;
  onDosageTypeChange: (unit: string) => void;
}

export function TreatmentDropdown({
  categoryId,
  treatmentId,
  dosageType,
  onCategoryChange,
  onTreatmentChange,
  onDosageTypeChange,
}: TreatmentDropdownProps) {
  const { data: categories = [] } = useCategories();
  const { data: treatments = [] } = useTreatments(categoryId);
  const { data: dosageTypes = [] } = useDosageTypes(treatmentId);

  return (
    <div className="flex flex-col gap-2">
      <Select
        placeholder="카테고리 선택"
        selectedKeys={categoryId ? [categoryId] : []}
        onSelectionChange={(keys) => {
          const key = Array.from(keys)[0] as string;
          onCategoryChange(key);
          onTreatmentChange("");
          onDosageTypeChange("");
        }}
        aria-label="시술 카테고리"
        variant="bordered"
        size="md"
        classNames={{ trigger: "border border-black/20 rounded-lg" }}
      >
        {categories.map((cat) => (
          <SelectItem key={cat.id}>{cat.name}</SelectItem>
        ))}
      </Select>

      <Select
        placeholder="시술명 선택"
        isDisabled={!categoryId}
        selectedKeys={treatmentId ? [treatmentId] : []}
        onSelectionChange={(keys) => {
          const key = Array.from(keys)[0] as string;
          onTreatmentChange(key);
          onDosageTypeChange("");
        }}
        aria-label="시술명"
        variant="bordered"
        size="md"
        classNames={{ trigger: "border border-black/20 rounded-lg" }}
      >
        {treatments.map((t) => (
          <SelectItem key={t.id}>{t.name}</SelectItem>
        ))}
      </Select>

      {dosageTypes.length > 0 && (
        <Select
          placeholder="단위 선택"
          isDisabled={!treatmentId}
          selectedKeys={dosageType ? [dosageType] : []}
          onSelectionChange={(keys) => {
            const key = Array.from(keys)[0] as string;
            onDosageTypeChange(key);
          }}
          aria-label="용량 단위"
          variant="bordered"
          size="md"
          classNames={{ trigger: "border border-black/20 rounded-lg" }}
        >
          {dosageTypes.map((d) => (
            <SelectItem key={d.unit}>{d.unit}</SelectItem>
          ))}
        </Select>
      )}
    </div>
  );
}
