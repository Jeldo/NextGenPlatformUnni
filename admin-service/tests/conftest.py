import asyncio
import json
from collections.abc import AsyncIterator
from pathlib import Path

import pytest
import pytest_asyncio
from httpx import ASGITransport, AsyncClient
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

from app.db.session import get_session
from app.main import app
from app.models import Base

TEST_DB_URL = "sqlite+aiosqlite:///./test.db"

engine = create_async_engine(TEST_DB_URL, echo=False)
TestSessionLocal = async_sessionmaker(engine, class_=AsyncSession, expire_on_commit=False)


async def override_get_session() -> AsyncIterator[AsyncSession]:
    async with TestSessionLocal() as session:
        try:
            yield session
            await session.commit()
        except Exception:
            await session.rollback()
            raise


app.dependency_overrides[get_session] = override_get_session


@pytest_asyncio.fixture(autouse=True)
async def setup_db():
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    yield
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.drop_all)


@pytest_asyncio.fixture
async def client() -> AsyncIterator[AsyncClient]:
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as ac:
        yield ac


@pytest_asyncio.fixture
async def seeded_client(client: AsyncClient) -> AsyncClient:
    """Client with seed_data.json pre-loaded into DB."""
    seed_file = Path(__file__).parent.parent / "seed_data.json"
    data = json.loads(seed_file.read_text())

    for cat in data["categories"]:
        await client.post("/api/categories", json={"name": cat["name"]})

    # We need to use the IDs from seed_data, so insert directly via session
    async with TestSessionLocal() as session:
        await session.execute(
            Base.metadata.tables["categories"].delete()
        )
        await session.execute(
            Base.metadata.tables["treatments"].delete()
        )
        await session.commit()

        from app.models import Category, Treatment, DosageType, CycleRule
        from datetime import datetime, timezone

        for cat in data["categories"]:
            session.add(Category(id=cat["id"], name=cat["name"]))
        await session.flush()

        for t in data["treatments"]:
            session.add(Treatment(id=t["id"], category_id=t["category_id"], name=t["name"]))
        await session.flush()

        for dt in data["dosage_types"]:
            session.add(DosageType(id=dt["id"], treatment_id=dt["treatment_id"], unit=dt["unit"]))
        await session.flush()

        for cr in data["cycle_rules"]:
            session.add(CycleRule(
                treatment_id=cr["treatment_id"],
                cycle_days=cr["cycle_days"],
                description=cr["description"],
                updated_at=datetime.now(timezone.utc),
            ))
        await session.commit()

    return client
