// Record
export interface TreatmentRecord {
  id: string;
  user_id: string;
  source: "AUTO" | "MANUAL";
  category_id: string;
  treatment_id: string;
  dosage_type: string | null;
  dosage_value: string | null;
  treatment_date: string;
  hospital_name: string;
  hospital_location: string | null;
  doctor_name: string | null;
  memo: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateRecordRequest {
  category_id: string;
  treatment_id: string;
  dosage_type?: string;
  dosage_value?: string;
  treatment_date: string;
  hospital_name: string;
  hospital_location?: string;
  doctor_name?: string;
  memo?: string;
}

export interface UpdateRecordRequest {
  category_id?: string;
  treatment_id?: string;
  dosage_type?: string;
  dosage_value?: string;
  treatment_date?: string;
  hospital_name?: string;
  hospital_location?: string;
  doctor_name?: string;
  memo?: string;
}

// Schedule
export interface ScheduledTreatment {
  id: string;
  record_id: string;
  category_id: string;
  treatment_id: string;
  scheduled_date: string;
  cycle_days: number;
  status: "PENDING" | "REMINDED" | "COMPLETED";
}

// Statistics
export interface TreatmentStat {
  category_id: string;
  category_name: string;
  count: number;
}

// Treatment Data
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
  unit: "shot" | "minute" | "volume" | "vial" | "joule";
}
