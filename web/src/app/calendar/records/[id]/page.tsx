"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Button, Input, Textarea, Modal, ModalContent, ModalHeader, ModalBody, ModalFooter } from "@heroui/react";
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
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [showMenu, setShowMenu] = useState(false);
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

  const formatDate = (dateStr: string) => {
    const d = new Date(dateStr);
    const days = ["일", "월", "화", "수", "목", "금", "토"];
    return `${d.getFullYear()}년 ${d.getMonth() + 1}월 ${d.getDate()}일 (${days[d.getDay()]})`;
  };

  const title = [record.hospital_name, record.dosage_value ? `${record.dosage_value}${record.dosage_type}` : ""].filter(Boolean).join(" ");

  // 수정 모드
  if (isEditing) {
    return (
      <main className="min-h-screen bg-white flex flex-col">
        {/* Header */}
        <header className="flex items-center px-5 py-4 border-b border-black/10 sticky top-0 bg-white z-10">
          <button onClick={() => setIsEditing(false)} className="text-xl mr-3 min-w-[44px] min-h-[44px] flex items-center">‹</button>
          <span className="text-[17px] font-semibold">시술 정보 수정</span>
        </header>

        <div className="flex-1 px-5 py-5 flex flex-col gap-4">
          <div className="flex flex-col gap-1.5">
            <label className="text-[13px] font-semibold text-black">날짜</label>
            <Input type="date" value={editDate} onValueChange={setEditDate} aria-label="날짜" />
          </div>
          <div className="flex flex-col gap-1.5">
            <label className="text-[13px] font-semibold text-black">병원명</label>
            <Input value={editHospital} onValueChange={setEditHospital} aria-label="병원명" />
          </div>
          <div className="flex flex-col gap-1.5">
            <label className="text-[13px] font-semibold text-black">메모</label>
            <Textarea value={editMemo} onValueChange={setEditMemo} placeholder="메모를 입력하세요 (선택)" aria-label="메모" />
          </div>
          <Button color="primary" size="lg" className="w-full mt-4 font-semibold" isLoading={updateRecord.isPending} onPress={handleSave}>
            저장
          </Button>
        </div>
      </main>
    );
  }

  // 상세 보기 모드
  return (
    <main className="min-h-screen bg-white flex flex-col relative">
      {/* Header */}
      <header className="flex items-center justify-between px-5 py-4 border-b border-black/10 sticky top-0 bg-white z-10">
        <div className="flex items-center">
          <button onClick={() => router.push("/calendar")} className="text-xl mr-3 min-w-[44px] min-h-[44px] flex items-center">‹</button>
          <span className="text-[17px] font-semibold">시술 상세</span>
        </div>
        <div className="relative">
          <button onClick={() => setShowMenu(!showMenu)} className="text-xl px-2 py-1 min-w-[44px] min-h-[44px] flex items-center justify-center">···</button>
          {showMenu && (
            <div className="absolute right-0 top-10 bg-white border border-black/12 rounded-xl shadow-lg min-w-[120px] z-20 overflow-hidden">
              <button
                onClick={() => { setShowMenu(false); setShowDeleteModal(true); }}
                className="w-full text-left px-4 py-3 text-sm text-red-500 hover:bg-gray-50"
              >
                삭제
              </button>
            </div>
          )}
        </div>
      </header>

      {/* Content */}
      <div className="flex-1 pb-[100px]">
        {/* Title */}
        <div className="px-5 pt-5 pb-4">
          <h1 className="text-2xl font-bold text-black leading-tight">{title}</h1>
        </div>

        {/* Detail rows */}
        <div className="px-5">
          <div className="flex items-start gap-3.5 py-3.5 border-b border-black/7">
            <span className="text-lg flex-shrink-0 w-5 text-center">🏥</span>
            <span className="text-[15px] text-black">{record.hospital_name}</span>
          </div>
          <div className="flex items-start gap-3.5 py-3.5 border-b border-black/7">
            <span className="text-lg flex-shrink-0 w-5 text-center">📅</span>
            <span className="text-[15px] text-black">{formatDate(record.treatment_date)}</span>
          </div>
          <div className="flex items-start gap-3.5 py-3.5">
            <span className="text-lg flex-shrink-0 w-5 text-center">📝</span>
            <span className={`text-[15px] ${record.memo ? "text-black" : "text-gray-description"}`}>
              {record.memo || "메모 없음"}
            </span>
          </div>
        </div>
      </div>

      {/* Bottom fixed actions */}
      <div className="fixed bottom-0 left-0 right-0 bg-white border-t border-black/8 px-5 py-3 pb-7 flex gap-2.5">
        <GoogleCalendarButton title={title} date={record.treatment_date} />
        <Button color="primary" size="lg" className="flex-1 font-semibold" onPress={startEdit}>
          수정하기
        </Button>
      </div>

      {/* Delete Modal */}
      <Modal isOpen={showDeleteModal} onClose={() => setShowDeleteModal(false)} placement="center">
        <ModalContent>
          <ModalHeader className="justify-center">시술 기록 삭제</ModalHeader>
          <ModalBody className="text-center">
            <p className="text-sm text-gray-description">이 시술 기록을 삭제하시겠습니까?<br />삭제 후 복구할 수 없습니다.</p>
          </ModalBody>
          <ModalFooter>
            <Button variant="bordered" className="flex-1" onPress={() => setShowDeleteModal(false)}>취소</Button>
            <Button color="danger" className="flex-1 font-semibold" isLoading={deleteRecord.isPending} onPress={handleDelete}>삭제</Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </main>
  );
}
