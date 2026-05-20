"use client";

import { useParams, useRouter } from "next/navigation";
import { useRecord } from "@/hooks/useRecord";
import { useDeleteRecord } from "@/hooks/useDeleteRecord";

export default function RecordDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const { data: record, isLoading } = useRecord(id);
  const deleteRecord = useDeleteRecord();

  if (isLoading) return <div className="p-4">Loading...</div>;
  if (!record) return <div className="p-4">기록을 찾을 수 없습니다.</div>;

  const handleDelete = () => {
    deleteRecord.mutate(id, { onSuccess: () => router.push("/calendar") });
  };

  return (
    <main className="min-h-screen p-4">
      <h1 className="text-xl font-bold">시술 상세</h1>

      <div className="mt-4 space-y-3">
        <div><span className="text-gray-description">날짜:</span> {record.treatment_date.slice(0, 10)}</div>
        <div><span className="text-gray-description">병원:</span> {record.hospital_name}</div>
        <div><span className="text-gray-description">소스:</span> {record.source}</div>
        {record.dosage_type && (
          <div><span className="text-gray-description">용량:</span> {record.dosage_value} {record.dosage_type}</div>
        )}
        {record.memo && <div><span className="text-gray-description">메모:</span> {record.memo}</div>}
      </div>

      <div className="mt-6 flex gap-3">
        <button className="border rounded p-2 flex-1">수정</button>
        <button onClick={handleDelete} className="border border-red-500 text-red-500 rounded p-2 flex-1">삭제</button>
      </div>

      <button className="mt-4 border rounded p-3 w-full text-primary border-primary">
        구글 캘린더에 등록하기
      </button>
    </main>
  );
}
