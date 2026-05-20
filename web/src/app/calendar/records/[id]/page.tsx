"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Button, Input, Textarea } from "@heroui/react";
import { useRecord } from "@/hooks/useRecord";
import { useUpdateRecord } from "@/hooks/useUpdateRecord";
import { useDeleteRecord } from "@/hooks/useDeleteRecord";
import { GoogleCalendarButton } from "@/components/GoogleCalendarButton";

export default function RecordDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const { data: record, isLoading } = useRecord(id);
  const updateRecord = useUpdateRecord();
  const deleteRecord = useDeleteRecord();

  const [isEditing, setIsEditing] = useState(false);
  const [editDate, setEditDate] = useState("");
  const [editHospital, setEditHospital] = useState("");
  const [editMemo, setEditMemo] = useState("");

  if (isLoading) return <div className="p-4">Loading...</div>;
  if (!record) return <div className="p-4">기록을 찾을 수 없습니다.</div>;

  const startEdit = () => {
    setEditDate(record.treatment_date.slice(0, 10));
    setEditHospital(record.hospital_name);
    setEditMemo(record.memo ?? "");
    setIsEditing(true);
  };

  const handleSave = () => {
    updateRecord.mutate(
      { id, data: { treatment_date: new Date(editDate).toISOString(), hospital_name: editHospital, memo: editMemo || undefined } },
      { onSuccess: () => setIsEditing(false) },
    );
  };

  const handleDelete = () => {
    deleteRecord.mutate(id, { onSuccess: () => router.push("/calendar") });
  };

  if (isEditing) {
    return (
      <main className="min-h-screen p-4">
        <h1 className="text-xl font-bold mb-4">시술 수정</h1>
        <div className="flex flex-col gap-4">
          <Input type="date" label="시술 날짜" value={editDate} onValueChange={setEditDate} aria-label="시술 날짜" />
          <Input label="병원명" value={editHospital} onValueChange={setEditHospital} aria-label="병원명" />
          <Textarea label="메모" value={editMemo} onValueChange={setEditMemo} aria-label="메모" />
          <div className="flex gap-3">
            <Button variant="bordered" className="flex-1" onPress={() => setIsEditing(false)}>취소</Button>
            <Button color="primary" className="flex-1" isLoading={updateRecord.isPending} onPress={handleSave}>저장</Button>
          </div>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-screen p-4">
      <h1 className="text-xl font-bold mb-4">시술 상세</h1>

      <div className="space-y-3">
        <div className="flex justify-between">
          <span className="text-gray-description">날짜</span>
          <span>{record.treatment_date.slice(0, 10)}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-description">병원</span>
          <span>{record.hospital_name}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-description">소스</span>
          <span>{record.source === "AUTO" ? "자동" : "수동"}</span>
        </div>
        {record.dosage_type && (
          <div className="flex justify-between">
            <span className="text-gray-description">용량</span>
            <span>{record.dosage_value} {record.dosage_type}</span>
          </div>
        )}
        {record.memo && (
          <div className="flex justify-between">
            <span className="text-gray-description">메모</span>
            <span>{record.memo}</span>
          </div>
        )}
      </div>

      <div className="mt-6 flex gap-3">
        <Button variant="bordered" className="flex-1" onPress={startEdit}>수정</Button>
        <Button color="danger" variant="bordered" className="flex-1" isLoading={deleteRecord.isPending} onPress={handleDelete}>삭제</Button>
      </div>

      <div className="mt-4">
        <GoogleCalendarButton title={record.hospital_name} date={record.treatment_date} />
      </div>
    </main>
  );
}
