"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button, Input, Textarea } from "@heroui/react";
import { TreatmentDropdown } from "@/components/TreatmentDropdown";
import { DosageInput } from "@/components/DosageInput";
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
    <main className="min-h-screen bg-white flex flex-col">
      {/* Header */}
      <header className="flex items-center px-5 py-4 border-b border-black/10 sticky top-0 bg-white z-10">
        <button onClick={() => router.back()} className="text-xl mr-3 min-w-[44px] min-h-[44px] flex items-center">‹</button>
        <span className="text-[17px] font-semibold">시술 추가</span>
      </header>

      {/* Form */}
      <div className="flex-1 px-5 py-5 flex flex-col gap-5">
        {/* 날짜 */}
        <div className="flex flex-col gap-1.5">
          <label className="text-[13px] font-semibold text-black">
            날짜 <span className="text-red-500">*</span>
          </label>
          <Input
            type="date"
            value={treatmentDate}
            onValueChange={setTreatmentDate}
            isRequired
            aria-label="시술 날짜"
          />
        </div>

        {/* 병원명 */}
        <div className="flex flex-col gap-1.5">
          <label className="text-[13px] font-semibold text-black">
            병원명 <span className="text-red-500">*</span>
          </label>
          <Input
            placeholder="병원명을 입력하세요"
            value={hospitalName}
            onValueChange={setHospitalName}
            isRequired
            aria-label="병원명"
          />
        </div>

        {/* 시술 정보 */}
        <div className="flex flex-col gap-1.5">
          <label className="text-[13px] font-semibold text-black">
            시술 정보 <span className="text-red-500">*</span>
          </label>
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
        </div>

        {/* 메모 */}
        <div className="flex flex-col gap-1.5">
          <label className="text-[13px] font-semibold text-black">메모</label>
          <Textarea
            placeholder="메모를 입력하세요 (선택)"
            value={memo}
            onValueChange={setMemo}
            aria-label="메모"
            minRows={3}
          />
        </div>

        {/* 저장 */}
        <Button
          color="primary"
          size="lg"
          className="w-full mt-2 font-semibold"
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
