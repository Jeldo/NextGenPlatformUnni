import uuid
from datetime import datetime, timezone

from fastapi import HTTPException

from app.models import Category, CycleRule, DosageType, Treatment
from app.repositories import CycleRuleRepository, TreatmentDataRepository
from app.schemas import (
    CategoryCreate,
    CategoryUpdate,
    CycleRuleCreate,
    CycleRuleUpdate,
    DosageTypeCreate,
    TreatmentCreate,
    TreatmentUpdate,
)


def _new_id() -> str:
    return str(uuid.uuid4())


class TreatmentDataService:
    def __init__(self, repo: TreatmentDataRepository):
        self._repo = repo

    # --- Category ---

    async def list_categories(self) -> list[Category]:
        return await self._repo.list_categories()

    async def create_category(self, data: CategoryCreate) -> Category:
        if await self._repo.get_category_by_name(data.name):
            raise HTTPException(status_code=409, detail=f"Category '{data.name}' already exists")
        category = Category(id=_new_id(), name=data.name)
        return await self._repo.create_category(category)

    async def update_category(self, category_id: str, data: CategoryUpdate) -> Category:
        category = await self._repo.get_category(category_id)
        if not category:
            raise HTTPException(status_code=404, detail=f"Category '{category_id}' not found")
        existing = await self._repo.get_category_by_name(data.name)
        if existing and existing.id != category_id:
            raise HTTPException(status_code=409, detail=f"Category '{data.name}' already exists")
        category.name = data.name
        return category

    async def delete_category(self, category_id: str) -> None:
        category = await self._repo.get_category(category_id)
        if not category:
            raise HTTPException(status_code=404, detail=f"Category '{category_id}' not found")
        await self._repo.delete_category(category)

    # --- Treatment ---

    async def list_treatments(self, category_id: str) -> list[Treatment]:
        if not await self._repo.get_category(category_id):
            raise HTTPException(status_code=404, detail=f"Category '{category_id}' not found")
        return await self._repo.list_treatments(category_id)

    async def create_treatment(self, category_id: str, data: TreatmentCreate) -> Treatment:
        if not await self._repo.get_category(category_id):
            raise HTTPException(status_code=404, detail=f"Category '{category_id}' not found")
        if await self._repo.get_treatment_by_name(category_id, data.name):
            raise HTTPException(status_code=409, detail=f"Treatment '{data.name}' already exists in this category")
        treatment = Treatment(id=_new_id(), category_id=category_id, name=data.name)
        return await self._repo.create_treatment(treatment)

    async def update_treatment(self, treatment_id: str, data: TreatmentUpdate) -> Treatment:
        treatment = await self._repo.get_treatment(treatment_id)
        if not treatment:
            raise HTTPException(status_code=404, detail=f"Treatment '{treatment_id}' not found")
        existing = await self._repo.get_treatment_by_name(treatment.category_id, data.name)
        if existing and existing.id != treatment_id:
            raise HTTPException(status_code=409, detail=f"Treatment '{data.name}' already exists in this category")
        treatment.name = data.name
        return treatment

    async def delete_treatment(self, treatment_id: str) -> None:
        treatment = await self._repo.get_treatment(treatment_id)
        if not treatment:
            raise HTTPException(status_code=404, detail=f"Treatment '{treatment_id}' not found")
        await self._repo.delete_treatment(treatment)

    # --- DosageType ---

    async def list_dosage_types(self, treatment_id: str) -> list[DosageType]:
        if not await self._repo.get_treatment(treatment_id):
            raise HTTPException(status_code=404, detail=f"Treatment '{treatment_id}' not found")
        return await self._repo.list_dosage_types(treatment_id)

    async def create_dosage_type(self, treatment_id: str, data: DosageTypeCreate) -> DosageType:
        if not await self._repo.get_treatment(treatment_id):
            raise HTTPException(status_code=404, detail=f"Treatment '{treatment_id}' not found")
        if await self._repo.get_dosage_type_by_unit(treatment_id, data.unit):
            raise HTTPException(status_code=409, detail=f"DosageType '{data.unit}' already exists for this treatment")
        dosage_type = DosageType(id=_new_id(), treatment_id=treatment_id, unit=data.unit)
        return await self._repo.create_dosage_type(dosage_type)

    async def delete_dosage_type(self, dosage_type_id: str) -> None:
        dosage_type = await self._repo.get_dosage_type(dosage_type_id)
        if not dosage_type:
            raise HTTPException(status_code=404, detail=f"DosageType '{dosage_type_id}' not found")
        await self._repo.delete_dosage_type(dosage_type)


class CycleRuleService:
    def __init__(self, repo: CycleRuleRepository, treatment_repo: TreatmentDataRepository):
        self._repo = repo
        self._treatment_repo = treatment_repo

    async def list_rules(self) -> list[CycleRule]:
        return await self._repo.list_rules()

    async def get_rule(self, treatment_id: str) -> CycleRule:
        rule = await self._repo.get_rule(treatment_id)
        if not rule:
            raise HTTPException(status_code=404, detail=f"CycleRule for treatment '{treatment_id}' not found")
        return rule

    async def create_rule(self, data: CycleRuleCreate) -> CycleRule:
        if not await self._treatment_repo.get_treatment(data.treatment_id):
            raise HTTPException(status_code=404, detail=f"Treatment '{data.treatment_id}' not found")
        if await self._repo.get_rule(data.treatment_id):
            raise HTTPException(status_code=409, detail=f"CycleRule for treatment '{data.treatment_id}' already exists")
        rule = CycleRule(
            treatment_id=data.treatment_id,
            cycle_days=data.cycle_days,
            description=data.description,
            updated_at=datetime.now(timezone.utc),
        )
        return await self._repo.create_rule(rule)

    async def update_rule(self, treatment_id: str, data: CycleRuleUpdate) -> CycleRule:
        rule = await self._repo.get_rule(treatment_id)
        if not rule:
            raise HTTPException(status_code=404, detail=f"CycleRule for treatment '{treatment_id}' not found")
        rule.cycle_days = data.cycle_days
        if data.description is not None:
            rule.description = data.description
        rule.updated_at = datetime.now(timezone.utc)
        return rule

    async def delete_rule(self, treatment_id: str) -> None:
        rule = await self._repo.get_rule(treatment_id)
        if not rule:
            raise HTTPException(status_code=404, detail=f"CycleRule for treatment '{treatment_id}' not found")
        await self._repo.delete_rule(rule)
