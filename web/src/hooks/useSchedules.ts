import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ScheduledTreatment } from "@/types";

export function useSchedules() {
  return useQuery({
    queryKey: ["schedules"],
    queryFn: () => api<ScheduledTreatment[]>("/api/schedules"),
    staleTime: 30_000,
  });
}
