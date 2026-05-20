import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { TreatmentCategory } from "@/types";

export function useCategories() {
  return useQuery({
    queryKey: ["categories"],
    queryFn: () => api<TreatmentCategory[]>("/api/treatment-data/categories"),
    staleTime: 600_000,
  });
}
