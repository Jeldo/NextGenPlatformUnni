CREATE TABLE IF NOT EXISTS treatment_records (
    seq BIGSERIAL PRIMARY KEY,
    id UUID NOT NULL UNIQUE,
    user_id VARCHAR(255) NOT NULL,
    treatment_id UUID NOT NULL,
    treatment_name VARCHAR(255) NOT NULL,
    category_id UUID NOT NULL,
    category_name VARCHAR(255) NOT NULL,
    treatment_date TIMESTAMPTZ NOT NULL,
    hospital_name VARCHAR(255),
    dosage_value DOUBLE PRECISION,
    dosage_unit VARCHAR(50),
    source VARCHAR(20) NOT NULL,
    reservation_id VARCHAR(255) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_treatment_records_user_date ON treatment_records (user_id, treatment_date);
CREATE INDEX idx_treatment_records_reservation ON treatment_records (reservation_id) WHERE reservation_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS scheduled_treatments (
    seq BIGSERIAL PRIMARY KEY,
    id UUID NOT NULL UNIQUE,
    user_id VARCHAR(255) NOT NULL,
    record_id UUID NOT NULL REFERENCES treatment_records(id) ON DELETE CASCADE,
    treatment_id UUID NOT NULL,
    treatment_name VARCHAR(255) NOT NULL,
    scheduled_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL,
    reminded_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_scheduled_treatments_user_date ON scheduled_treatments (user_id, scheduled_date);
CREATE INDEX idx_scheduled_treatments_status_date ON scheduled_treatments (status, scheduled_date);
CREATE INDEX idx_scheduled_treatments_record ON scheduled_treatments (record_id);
