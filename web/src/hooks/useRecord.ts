import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { TreatmentRecord } from "@/types";

export function useRecord(id: string) {
  return useQuery({
    queryKey: ["records", id],
    queryFn: () => api<TreatmentRecord>(`/api/records/${id}`),
    staleTime: 60_000,
    enabled: !!id,
  });
}
