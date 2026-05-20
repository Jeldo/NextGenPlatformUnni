"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useCategories } from "@/hooks/useCategories";
import { useTreatments } from "@/hooks/useTreatments";
import { useDosageTypes } from "@/hooks/useDosageTypes";
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

  const { data: categories = [] } = useCategories();
  const { data: treatments = [] } = useTreatments(categoryId);
  const { data: dosageTypes = [] } = useDosageTypes(treatmentId);

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
    <main className="min-h-screen bg-white flex flex-col relative">
      {/* Header */}
      <header className="flex items-center gap-3 px-5 py-4 border-b border-black/10 sticky top-0 bg-white z-10">
        <button onClick={() => router.back()} className="text-xl min-w-[44px] min-h-[44px] flex items-center">‹</button>
        <span className="text-[17px] font-semibold">시술 추가</span>
      </header>

      {/* Form */}
      <div className="flex-1 px-5 py-5 pb-[120px] flex flex-col gap-[18px]">
        {/* 날짜 */}
        <div className="flex flex-col gap-1.5">
          <label className="text-[13px] font-semibold text-black">날짜<span className="text-primary ml-0.5">*</span></label>
          <input
            type="date"
            value={treatmentDate}
            onChange={(e) => setTreatmentDate(e.target.value)}
            className="w-full px-3.5 py-3 border border-black/20 rounded-lg text-[15px] text-black bg-white"
            aria-label="시술 날짜"
          />
        </div>

        {/* 병원명 */}
        <div className="flex flex-col gap-1.5">
          <label className="text-[13px] font-semibold text-black">병원명<span className="text-primary ml-0.5">*</span></label>
          <input
            type="text"
            placeholder="병원명을 입력하세요"
            value={hospitalName}
            onChange={(e) => setHospitalName(e.target.value)}
            className="w-full px-3.5 py-3 border border-black/20 rounded-lg text-[15px] text-black bg-white placeholder:text-gray-300"
            aria-label="병원명"
          />
        </div>

        {/* 시술 정보 */}
        <div className="flex flex-col gap-1.5">
          <label className="text-[13px] font-semibold text-black">시술 정보<span className="text-primary ml-0.5">*</span></label>
          <div className="flex flex-col gap-2">
            <div className="relative">
              <select
                value={categoryId}
                onChange={(e) => { setCategoryId(e.target.value); setTreatmentId(""); setDosageType(""); }}
                className="w-full px-3.5 py-3 border border-black/20 rounded-lg text-[15px] text-black bg-white appearance-none pr-10"
                aria-label="시술 카테고리"
              >
                <option value="">카테고리 선택</option>
                {categories.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
              </select>
              <span className="absolute right-3.5 top-1/2 -translate-y-1/2 text-primary text-xs pointer-events-none">▼</span>
            </div>

            <div className="relative">
              <select
                value={treatmentId}
                onChange={(e) => { setTreatmentId(e.target.value); setDosageType(""); }}
                disabled={!categoryId}
                className="w-full px-3.5 py-3 border border-black/20 rounded-lg text-[15px] text-black bg-white appearance-none pr-10 disabled:bg-gray-50 disabled:text-gray-300"
                aria-label="시술명"
              >
                <option value="">시술명 선택</option>
                {treatments.map((t) => <option key={t.id} value={t.id}>{t.name}</option>)}
              </select>
              <span className="absolute right-3.5 top-1/2 -translate-y-1/2 text-primary text-xs pointer-events-none">▼</span>
            </div>

            {dosageTypes.length > 0 && (
              <div className="flex gap-2">
                <input
                  type="number"
                  placeholder="용량"
                  value={dosageValue}
                  onChange={(e) => setDosageValue(e.target.value)}
                  className="flex-1 px-3.5 py-3 border border-black/20 rounded-lg text-[15px] text-black bg-white"
                  aria-label="용량"
                />
                <div className="relative w-[110px]">
                  <select
                    value={dosageType}
                    onChange={(e) => setDosageType(e.target.value)}
                    className="w-full px-3.5 py-3 border border-black/20 rounded-lg text-[15px] text-black bg-white appearance-none pr-8"
                    aria-label="용량 단위"
                  >
                    <option value="">단위</option>
                    {dosageTypes.map((d) => <option key={d.unit} value={d.unit}>{d.unit}</option>)}
                  </select>
                  <span className="absolute right-3 top-1/2 -translate-y-1/2 text-primary text-xs pointer-events-none">▼</span>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* 메모 */}
        <div className="flex flex-col gap-1.5">
          <label className="text-[13px] font-semibold text-black">메모</label>
          <textarea
            placeholder="메모를 입력하세요 (선택)"
            value={memo}
            onChange={(e) => setMemo(e.target.value)}
            className="w-full px-3.5 py-3 border border-black/20 rounded-lg text-[15px] text-black bg-white placeholder:text-gray-300 resize-none h-20"
            aria-label="메모"
          />
        </div>
      </div>

      {/* 하단 고정 저장 버튼 */}
      <div className="fixed bottom-0 left-0 right-0 px-5 py-3 pb-7 bg-white border-t border-black/8">
        <button
          disabled={!isValid || isPending}
          onClick={handleSubmit}
          className="w-full py-3.5 bg-primary text-white text-base font-semibold rounded-xl disabled:opacity-50"
        >
          {isPending ? "저장 중..." : "저장"}
        </button>
      </div>
    </main>
  );
}
