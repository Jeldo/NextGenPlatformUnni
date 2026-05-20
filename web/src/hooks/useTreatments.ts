import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { Treatment } from "@/types";

export function useTreatments(categoryId: string | null) {
  return useQuery({
    queryKey: ["treatments", categoryId],
    queryFn: () => api<Treatment[]>(`/api/treatment-data/categories/${categoryId}/treatments`),
    staleTime: 600_000,
    enabled: !!categoryId,
  });
}
