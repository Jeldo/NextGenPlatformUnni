"use client";

import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useCategories } from "@/hooks/useCategories";
import { useTreatments } from "@/hooks/useTreatments";
import { useDosageTypes } from "@/hooks/useDosageTypes";
import { useCreateRecord } from "@/hooks/useCreateRecord";

export default function AddRecordPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const prefillDate = searchParams.get("date") || "";

  const [date, setDate] = useState(prefillDate);
  const [hospitalName, setHospitalName] = useState("");
  const [categoryId, setCategoryId] = useState<string | null>(null);
  const [treatmentId, setTreatmentId] = useState<string | null>(null);
  const [dosageType, setDosageType] = useState<string | null>(null);
  const [dosageValue, setDosageValue] = useState("");
  const [memo, setMemo] = useState("");

  const { data: categories } = useCategories();
  const { data: treatments } = useTreatments(categoryId);
  const { data: dosageTypes } = useDosageTypes(treatmentId);
  const createRecord = useCreateRecord();

  const handleSubmit = () => {
    if (!date || !categoryId || !treatmentId || !hospitalName) return;
    createRecord.mutate(
      {
        category_id: categoryId,
        treatment_id: treatmentId,
        dosage_type: dosageType || undefined,
        dosage_value: dosageValue || undefined,
        treatment_date: new Date(date).toISOString(),
        hospital_name: hospitalName,
        memo: memo || undefined,
      },
      { onSuccess: () => router.push("/calendar") }
    );
  };

  return (
    <main className="min-h-screen p-4">
      <h1 className="text-xl font-bold">시술 추가</h1>

      <div className="mt-4 space-y-4">
        <input type="date" value={date} onChange={(e) => setDate(e.target.value)} className="border rounded p-2 w-full" />
        <input placeholder="병원명" value={hospitalName} onChange={(e) => setHospitalName(e.target.value)} className="border rounded p-2 w-full" />

        {/* Category dropdown */}
        <select value={categoryId || ""} onChange={(e) => { setCategoryId(e.target.value); setTreatmentId(null); }} className="border rounded p-2 w-full">
          <option value="">카테고리 선택</option>
          {categories?.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
        </select>

        {/* Treatment dropdown */}
        <select value={treatmentId || ""} onChange={(e) => setTreatmentId(e.target.value)} disabled={!categoryId} className="border rounded p-2 w-full">
          <option value="">시술명 선택</option>
          {treatments?.map((t) => <option key={t.id} value={t.id}>{t.name}</option>)}
        </select>

        {/* Dosage */}
        <div className="flex gap-2">
          <select value={dosageType || ""} onChange={(e) => setDosageType(e.target.value)} disabled={!treatmentId} className="border rounded p-2">
            <option value="">단위</option>
            {dosageTypes?.map((d) => <option key={d.id} value={d.unit}>{d.unit}</option>)}
          </select>
          <input type="number" placeholder="용량" value={dosageValue} onChange={(e) => setDosageValue(e.target.value)} className="border rounded p-2 flex-1" />
        </div>

        <textarea placeholder="메모 (선택)" value={memo} onChange={(e) => setMemo(e.target.value)} className="border rounded p-2 w-full" />

        <button onClick={handleSubmit} disabled={createRecord.isPending} className="bg-primary text-white rounded p-3 w-full font-semibold">
          {createRecord.isPending ? "저장 중..." : "저장"}
        </button>
      </div>
    </main>
  );
}
