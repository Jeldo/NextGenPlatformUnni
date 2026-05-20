// Record
export interface TreatmentRecord {
  id: string;
  user_id: string;
  source: "AUTO" | "MANUAL";
  category_id: string;
  category_name: string;
  treatment_id: string;
  treatment_name: string;
  treatment_date: string;
  hospital_name?: string;
  dosage_value?: string;
  dosage_unit?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateRecordRequest {
  user_id: string;
  category_id: string;
  category_name: string;
  treatment_id: string;
  treatment_name: string;
  treatment_date: string;
  hospital_name?: string;
  dosage_value?: string;
  dosage_unit?: string;
}

export interface UpdateRecordRequest {
  category_id?: string;
  category_name?: string;
  treatment_id?: string;
  treatment_name?: string;
  treatment_date?: string;
  hospital_name?: string;
  dosage_value?: string;
  dosage_unit?: string;
}

// Schedule
export interface ScheduledTreatment {
  id: string;
  user_id: string;
  record_id: string;
  treatment_id: string;
  treatment_name: string;
  scheduled_date: string;
  status: "PENDING" | "REMINDED" | "COMPLETED";
  reminded_at?: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
}

// Statistics
export interface StatItem {
  treatment_id: string;
  treatment_name: string;
  category_name: string;
  count: number;
}

export interface StatisticsResponse {
  stats: StatItem[];
}

// Treatment Data (Admin)
export interface TreatmentCategory {
  id: string;
  name: string;
}

export interface Treatment {
  id: string;
  category_id: string;
  name: string;
}

export interface DosageType {
  id: string;
  treatment_id: string;
  unit: string;
}
