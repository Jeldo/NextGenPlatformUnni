from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from app.models import Category, CycleRule, DosageType, Treatment


class TreatmentDataRepository:
    def __init__(self, session: AsyncSession):
        self._session = session

    # --- Category ---

    async def list_categories(self) -> list[Category]:
        result = await self._session.execute(select(Category).order_by(Category.name))
        return list(result.scalars().all())

    async def get_category(self, category_id: str) -> Category | None:
        return await self._session.get(Category, category_id)

    async def get_category_by_name(self, name: str) -> Category | None:
        result = await self._session.execute(select(Category).where(Category.name == name))
        return result.scalar_one_or_none()

    async def create_category(self, category: Category) -> Category:
        self._session.add(category)
        await self._session.flush()
        return category

    async def delete_category(self, category: Category) -> None:
        await self._session.delete(category)

    # --- Treatment ---

    async def list_treatments(self, category_id: str) -> list[Treatment]:
        result = await self._session.execute(
            select(Treatment).where(Treatment.category_id == category_id).order_by(Treatment.name)
        )
        return list(result.scalars().all())

    async def get_treatment(self, treatment_id: str) -> Treatment | None:
        return await self._session.get(Treatment, treatment_id)

    async def get_treatment_by_name(self, category_id: str, name: str) -> Treatment | None:
        result = await self._session.execute(
            select(Treatment).where(Treatment.category_id == category_id, Treatment.name == name)
        )
        return result.scalar_one_or_none()

    async def create_treatment(self, treatment: Treatment) -> Treatment:
        self._session.add(treatment)
        await self._session.flush()
        return treatment

    async def delete_treatment(self, treatment: Treatment) -> None:
        await self._session.delete(treatment)

    # --- DosageType ---

    async def list_dosage_types(self, treatment_id: str) -> list[DosageType]:
        result = await self._session.execute(
            select(DosageType).where(DosageType.treatment_id == treatment_id)
        )
        return list(result.scalars().all())

    async def get_dosage_type(self, dosage_type_id: str) -> DosageType | None:
        return await self._session.get(DosageType, dosage_type_id)

    async def get_dosage_type_by_unit(self, treatment_id: str, unit: str) -> DosageType | None:
        result = await self._session.execute(
            select(DosageType).where(DosageType.treatment_id == treatment_id, DosageType.unit == unit)
        )
        return result.scalar_one_or_none()

    async def create_dosage_type(self, dosage_type: DosageType) -> DosageType:
        self._session.add(dosage_type)
        await self._session.flush()
        return dosage_type

    async def delete_dosage_type(self, dosage_type: DosageType) -> None:
        await self._session.delete(dosage_type)


class CycleRuleRepository:
    def __init__(self, session: AsyncSession):
        self._session = session

    async def list_rules(self) -> list[CycleRule]:
        result = await self._session.execute(select(CycleRule))
        return list(result.scalars().all())

    async def get_rule(self, treatment_id: str) -> CycleRule | None:
        return await self._session.get(CycleRule, treatment_id)

    async def create_rule(self, rule: CycleRule) -> CycleRule:
        self._session.add(rule)
        await self._session.flush()
        return rule

    async def delete_rule(self, rule: CycleRule) -> None:
        await self._session.delete(rule)
