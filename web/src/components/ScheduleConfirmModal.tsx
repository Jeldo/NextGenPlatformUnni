"use client";

import { Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, Button } from "@heroui/react";
import type { ScheduledTreatment } from "@/types";

export interface ScheduleConfirmModalProps {
  schedule: ScheduledTreatment | null;
  isOpen: boolean;
  onClose: () => void;
  onComplete: (schedule: ScheduledTreatment) => void;
  onDelete: (schedule: ScheduledTreatment) => void;
  isLoading?: boolean;
}

export function ScheduleConfirmModal({
  schedule,
  isOpen,
  onClose,
  onComplete,
  onDelete,
  isLoading,
}: ScheduleConfirmModalProps) {
  if (!schedule) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} placement="center">
      <ModalContent>
        <ModalHeader>이 날짜에 시술을 받으셨나요?</ModalHeader>
        <ModalBody>
          <p className="text-sm text-gray-description">
            {schedule.scheduled_date.slice(0, 10)} 예정된 시술
          </p>
        </ModalBody>
        <ModalFooter>
          <Button
            variant="light"
            color="danger"
            onPress={() => onDelete(schedule)}
            isDisabled={isLoading}
          >
            삭제
          </Button>
          <Button
            color="primary"
            onPress={() => onComplete(schedule)}
            isLoading={isLoading}
          >
            받았어요
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}
