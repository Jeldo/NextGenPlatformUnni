from datetime import datetime

from pydantic import BaseModel, ConfigDict, Field


# --- Category ---

class CategoryCreate(BaseModel):
    name: str = Field(min_length=1, max_length=100)


class CategoryUpdate(BaseModel):
    name: str = Field(min_length=1, max_length=100)


class CategoryResponse(BaseModel):
    id: str
    name: str
    model_config = ConfigDict(from_attributes=True)


# --- Treatment ---

class TreatmentCreate(BaseModel):
    name: str = Field(min_length=1, max_length=200)


class TreatmentUpdate(BaseModel):
    name: str = Field(min_length=1, max_length=200)


class TreatmentResponse(BaseModel):
    id: str
    category_id: str
    name: str
    model_config = ConfigDict(from_attributes=True)


# --- DosageType ---

class DosageTypeCreate(BaseModel):
    unit: str = Field(pattern=r"^(shot|minute|volume|vial|unit|joule)$")


class DosageTypeResponse(BaseModel):
    id: str
    treatment_id: str
    unit: str
    model_config = ConfigDict(from_attributes=True)


# --- CycleRule ---

class CycleRuleCreate(BaseModel):
    treatment_id: str
    cycle_days: int = Field(gt=0)
    description: str | None = None


class CycleRuleUpdate(BaseModel):
    cycle_days: int = Field(gt=0)
    description: str | None = None


class CycleRuleResponse(BaseModel):
    treatment_id: str
    cycle_days: int
    description: str | None
    updated_at: datetime
    model_config = ConfigDict(from_attributes=True)
