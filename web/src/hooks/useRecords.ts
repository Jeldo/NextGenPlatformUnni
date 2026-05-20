import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { TreatmentRecord } from "@/types";

export function useRecords(from: string, to: string) {
  return useQuery({
    queryKey: ["records", { from, to }],
    queryFn: () => api<TreatmentRecord[]>(`/api/records?from=${from}&to=${to}`),
    staleTime: 30_000,
  });
}
