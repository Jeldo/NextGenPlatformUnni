import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { DosageType } from "@/types";

export function useDosageTypes(treatmentId: string | null) {
  return useQuery({
    queryKey: ["dosageTypes", treatmentId],
    queryFn: () => api<DosageType[]>(`/api/treatment-data/treatments/${treatmentId}/dosage-types`),
    staleTime: 600_000,
    enabled: !!treatmentId,
  });
}
