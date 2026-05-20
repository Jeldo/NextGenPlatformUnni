import os

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(title="Admin Service", version="0.1.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_methods=["*"],
    allow_headers=["*"],
)

# --- Dummy Data ---

CATEGORIES = [
    {"id": "019234ab-1111-7def-8000-000000000001", "name": "보톡스"},
    {"id": "019234ab-1111-7def-8000-000000000002", "name": "필러"},
    {"id": "019234ab-1111-7def-8000-000000000003", "name": "레이저"},
]

TREATMENTS = {
    "019234ab-1111-7def-8000-000000000001": [
        {"id": "019234ab-2222-7def-8000-000000000001", "category_id": "019234ab-1111-7def-8000-000000000001", "name": "사각턱"},
        {"id": "019234ab-2222-7def-8000-000000000002", "category_id": "019234ab-1111-7def-8000-000000000001", "name": "이마"},
        {"id": "019234ab-2222-7def-8000-000000000003", "category_id": "019234ab-1111-7def-8000-000000000001", "name": "스킨 보톡스"},
    ],
    "019234ab-1111-7def-8000-000000000002": [
        {"id": "019234ab-2222-7def-8000-000000000004", "category_id": "019234ab-1111-7def-8000-000000000002", "name": "입술"},
        {"id": "019234ab-2222-7def-8000-000000000005", "category_id": "019234ab-1111-7def-8000-000000000002", "name": "팔자주름"},
        {"id": "019234ab-2222-7def-8000-000000000006", "category_id": "019234ab-1111-7def-8000-000000000002", "name": "볼"},
    ],
    "019234ab-1111-7def-8000-000000000003": [
        {"id": "019234ab-2222-7def-8000-000000000007", "category_id": "019234ab-1111-7def-8000-000000000003", "name": "리프팅"},
        {"id": "019234ab-2222-7def-8000-000000000008", "category_id": "019234ab-1111-7def-8000-000000000003", "name": "토닝"},
        {"id": "019234ab-2222-7def-8000-000000000009", "category_id": "019234ab-1111-7def-8000-000000000003", "name": "제모"},
    ],
}

DOSAGE_TYPES = {
    "019234ab-2222-7def-8000-000000000001": [
        {"id": "019234ab-3333-7def-8000-000000000001", "treatment_id": "019234ab-2222-7def-8000-000000000001", "unit": "shot"},
        {"id": "019234ab-3333-7def-8000-000000000002", "treatment_id": "019234ab-2222-7def-8000-000000000001", "unit": "volume"},
    ],
    "019234ab-2222-7def-8000-000000000004": [
        {"id": "019234ab-3333-7def-8000-000000000003", "treatment_id": "019234ab-2222-7def-8000-000000000004", "unit": "volume"},
        {"id": "019234ab-3333-7def-8000-000000000004", "treatment_id": "019234ab-2222-7def-8000-000000000004", "unit": "vial"},
    ],
    "019234ab-2222-7def-8000-000000000007": [
        {"id": "019234ab-3333-7def-8000-000000000005", "treatment_id": "019234ab-2222-7def-8000-000000000007", "unit": "joule"},
        {"id": "019234ab-3333-7def-8000-000000000006", "treatment_id": "019234ab-2222-7def-8000-000000000007", "unit": "minute"},
    ],
}

CYCLE_RULES = [
    {"category_id": "019234ab-1111-7def-8000-000000000001", "cycle_days": 90, "description": "보톡스 추천 주기", "updated_at": "2026-01-01T00:00:00Z"},
    {"category_id": "019234ab-1111-7def-8000-000000000002", "cycle_days": 180, "description": "필러 추천 주기", "updated_at": "2026-01-01T00:00:00Z"},
    {"category_id": "019234ab-1111-7def-8000-000000000003", "cycle_days": 30, "description": "레이저 추천 주기", "updated_at": "2026-01-01T00:00:00Z"},
]


# --- Health ---

@app.get("/health")
def health():
    return {"status": "ok", "service": "admin-service", "version": "0.1.0"}


# --- Cycle Rules ---

@app.post("/api/cycle-rules", status_code=201)
def create_cycle_rule():
    return CYCLE_RULES[0]


@app.get("/api/cycle-rules")
def list_cycle_rules():
    return CYCLE_RULES


@app.get("/api/cycle-rules/{category_id}")
def get_cycle_rule(category_id: str):
    for rule in CYCLE_RULES:
        if rule["category_id"] == category_id:
            return rule
    return CYCLE_RULES[0]


@app.put("/api/cycle-rules/{category_id}")
def update_cycle_rule(category_id: str):
    return CYCLE_RULES[0]


@app.delete("/api/cycle-rules/{category_id}", status_code=204)
def delete_cycle_rule(category_id: str):
    return None


# --- Categories ---

@app.post("/api/categories", status_code=201)
def create_category():
    return CATEGORIES[0]


@app.get("/api/categories")
def list_categories():
    return CATEGORIES


@app.put("/api/categories/{id}")
def update_category(id: str):
    return CATEGORIES[0]


@app.delete("/api/categories/{id}", status_code=204)
def delete_category(id: str):
    return None


# --- Treatments ---

@app.post("/api/categories/{category_id}/treatments", status_code=201)
def create_treatment(category_id: str):
    return TREATMENTS.get(category_id, [{}])[0]


@app.get("/api/categories/{category_id}/treatments")
def list_treatments(category_id: str):
    return TREATMENTS.get(category_id, [])


@app.put("/api/treatments/{id}")
def update_treatment(id: str):
    return TREATMENTS["019234ab-1111-7def-8000-000000000001"][0]


@app.delete("/api/treatments/{id}", status_code=204)
def delete_treatment(id: str):
    return None


# --- Dosage Types ---

@app.post("/api/treatments/{treatment_id}/dosage-types", status_code=201)
def create_dosage_type(treatment_id: str):
    return DOSAGE_TYPES.get(treatment_id, [{}])[0]


@app.get("/api/treatments/{treatment_id}/dosage-types")
def list_dosage_types(treatment_id: str):
    return DOSAGE_TYPES.get(treatment_id, [])


@app.delete("/api/dosage-types/{id}", status_code=204)
def delete_dosage_type(id: str):
    return None


if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PORT", "8081"))
    uvicorn.run(app, host="0.0.0.0", port=port)
