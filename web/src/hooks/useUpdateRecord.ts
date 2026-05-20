import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { TreatmentRecord, UpdateRecordRequest } from "@/types";

export function useUpdateRecord(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: UpdateRecordRequest) =>
      api<TreatmentRecord>(`/api/records/${id}`, { method: "PUT", body: JSON.stringify(data) }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["records"] });
      queryClient.invalidateQueries({ queryKey: ["statistics"] });
    },
  });
}
