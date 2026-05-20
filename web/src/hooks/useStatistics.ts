import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { TreatmentStat } from "@/types";

export function useStatistics() {
  return useQuery({
    queryKey: ["statistics"],
    queryFn: () => api<TreatmentStat[]>("/api/statistics"),
    staleTime: 300_000,
  });
}
