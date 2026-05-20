import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { StatisticsResponse } from "@/types";

export function useStatistics() {
  return useQuery({
    queryKey: ["statistics"],
    queryFn: () => api<StatisticsResponse>("/api/statistics"),
    staleTime: 300_000,
  });
}
