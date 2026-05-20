"""Seed the database with initial treatment data from seed_data.json."""
import asyncio
import json
from datetime import datetime, timezone
from pathlib import Path

from sqlalchemy import select, text
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

from app.config import settings
from app.models import Base, Category, CycleRule, DosageType, Treatment

SEED_FILE = Path(__file__).parent.parent / "seed_data.json"


async def seed(session: AsyncSession) -> None:
    # Skip if data already exists
    result = await session.execute(select(Category).limit(1))
    if result.scalar_one_or_none():
        print("Seed data already exists, skipping.")
        return

    data = json.loads(SEED_FILE.read_text())

    for cat in data["categories"]:
        session.add(Category(id=cat["id"], name=cat["name"]))

    for t in data["treatments"]:
        session.add(Treatment(id=t["id"], category_id=t["category_id"], name=t["name"]))

    for dt in data["dosage_types"]:
        session.add(DosageType(id=dt["id"], treatment_id=dt["treatment_id"], unit=dt["unit"]))

    now = datetime.now(timezone.utc)
    for cr in data["cycle_rules"]:
        session.add(CycleRule(
            treatment_id=cr["treatment_id"],
            cycle_days=cr["cycle_days"],
            description=cr["description"],
            updated_at=now,
        ))

    await session.commit()
    print(f"Seeded: {len(data['categories'])} categories, {len(data['treatments'])} treatments, "
          f"{len(data['dosage_types'])} dosage types, {len(data['cycle_rules'])} cycle rules.")


async def main() -> None:
    engine = create_async_engine(settings.database_url)
    async_session = async_sessionmaker(engine, class_=AsyncSession, expire_on_commit=False)
    async with async_session() as session:
        await seed(session)
    await engine.dispose()


if __name__ == "__main__":
    asyncio.run(main())
