"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button, Input, Textarea } from "@heroui/react";
import { TreatmentDropdown } from "@/components/TreatmentDropdown";
import { DosageInput } from "@/components/DosageInput";
import { HospitalInput } from "@/components/HospitalInput";
import { useCreateRecord } from "@/hooks/useCreateRecord";

export default function AddRecordPage() {
  const router = useRouter();
  const { mutate: createRecord, isPending } = useCreateRecord();

  const [treatmentDate, setTreatmentDate] = useState("");
  const [hospitalName, setHospitalName] = useState("");
  const [categoryId, setCategoryId] = useState("");
  const [treatmentId, setTreatmentId] = useState("");
  const [dosageType, setDosageType] = useState("");
  const [dosageValue, setDosageValue] = useState("");
  const [memo, setMemo] = useState("");

  const isValid = treatmentDate && hospitalName && categoryId && treatmentId;

  const handleSubmit = () => {
    if (!isValid) return;
    createRecord(
      {
        treatment_date: new Date(treatmentDate).toISOString(),
        hospital_name: hospitalName,
        category_id: categoryId,
        treatment_id: treatmentId,
        dosage_type: dosageType || undefined,
        dosage_value: dosageValue || undefined,
        memo: memo || undefined,
      },
      { onSuccess: () => router.push("/calendar") },
    );
  };

  return (
    <main className="min-h-screen p-4 pb-20">
      <h1 className="text-xl font-bold mb-6">시술 추가</h1>

      <div className="flex flex-col gap-4">
        <Input
          type="date"
          label="시술 날짜"
          value={treatmentDate}
          onValueChange={setTreatmentDate}
          isRequired
          aria-label="시술 날짜"
        />

        <HospitalInput value={hospitalName} onChange={setHospitalName} />

        <TreatmentDropdown
          categoryId={categoryId}
          treatmentId={treatmentId}
          dosageType={dosageType}
          onCategoryChange={setCategoryId}
          onTreatmentChange={setTreatmentId}
          onDosageTypeChange={setDosageType}
        />

        {dosageType && (
          <DosageInput value={dosageValue} unit={dosageType} onChange={setDosageValue} />
        )}

        <Textarea
          label="메모 (선택)"
          placeholder="추가 메모를 입력하세요"
          value={memo}
          onValueChange={setMemo}
          aria-label="메모"
        />

        <Button
          color="primary"
          size="lg"
          className="w-full mt-4"
          isDisabled={!isValid}
          isLoading={isPending}
          onPress={handleSubmit}
        >
          저장
        </Button>
      </div>
    </main>
  );
}
