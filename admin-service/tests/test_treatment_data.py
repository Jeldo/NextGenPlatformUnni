import pytest
from httpx import AsyncClient


@pytest.mark.asyncio
class TestCategories:
    async def test_create_category(self, client: AsyncClient):
        res = await client.post("/api/categories", json={"name": "보톡스"})
        assert res.status_code == 201
        assert res.json()["name"] == "보톡스"
        assert "id" in res.json()

    async def test_create_duplicate_category(self, client: AsyncClient):
        await client.post("/api/categories", json={"name": "보톡스"})
        res = await client.post("/api/categories", json={"name": "보톡스"})
        assert res.status_code == 409
        assert res.json()["error"]["code"] == "CONFLICT"

    async def test_list_categories(self, client: AsyncClient):
        await client.post("/api/categories", json={"name": "보톡스"})
        await client.post("/api/categories", json={"name": "필러"})
        res = await client.get("/api/categories")
        assert res.status_code == 200
        assert len(res.json()) == 2

    async def test_update_category(self, client: AsyncClient):
        create = await client.post("/api/categories", json={"name": "보톡스"})
        cat_id = create.json()["id"]
        res = await client.put(f"/api/categories/{cat_id}", json={"name": "리프팅"})
        assert res.status_code == 200
        assert res.json()["name"] == "리프팅"

    async def test_delete_category(self, client: AsyncClient):
        create = await client.post("/api/categories", json={"name": "보톡스"})
        cat_id = create.json()["id"]
        res = await client.delete(f"/api/categories/{cat_id}")
        assert res.status_code == 204
        # Verify gone
        list_res = await client.get("/api/categories")
        assert len(list_res.json()) == 0

    async def test_delete_nonexistent_category(self, client: AsyncClient):
        res = await client.delete("/api/categories/nonexistent")
        assert res.status_code == 404
        assert res.json()["error"]["code"] == "NOT_FOUND"


@pytest.mark.asyncio
class TestTreatments:
    async def test_create_treatment(self, client: AsyncClient):
        cat = await client.post("/api/categories", json={"name": "보톡스"})
        cat_id = cat.json()["id"]
        res = await client.post(f"/api/categories/{cat_id}/treatments", json={"name": "사각턱"})
        assert res.status_code == 201
        assert res.json()["name"] == "사각턱"
        assert res.json()["category_id"] == cat_id

    async def test_create_treatment_nonexistent_category(self, client: AsyncClient):
        res = await client.post("/api/categories/fake-id/treatments", json={"name": "사각턱"})
        assert res.status_code == 404

    async def test_create_duplicate_treatment(self, client: AsyncClient):
        cat = await client.post("/api/categories", json={"name": "보톡스"})
        cat_id = cat.json()["id"]
        await client.post(f"/api/categories/{cat_id}/treatments", json={"name": "사각턱"})
        res = await client.post(f"/api/categories/{cat_id}/treatments", json={"name": "사각턱"})
        assert res.status_code == 409

    async def test_list_treatments(self, client: AsyncClient):
        cat = await client.post("/api/categories", json={"name": "보톡스"})
        cat_id = cat.json()["id"]
        await client.post(f"/api/categories/{cat_id}/treatments", json={"name": "사각턱"})
        await client.post(f"/api/categories/{cat_id}/treatments", json={"name": "이마"})
        res = await client.get(f"/api/categories/{cat_id}/treatments")
        assert res.status_code == 200
        assert len(res.json()) == 2

    async def test_cascade_delete(self, client: AsyncClient):
        cat = await client.post("/api/categories", json={"name": "보톡스"})
        cat_id = cat.json()["id"]
        await client.post(f"/api/categories/{cat_id}/treatments", json={"name": "사각턱"})
        # Delete category cascades to treatments
        await client.delete(f"/api/categories/{cat_id}")
        res = await client.get(f"/api/categories/{cat_id}/treatments")
        assert res.status_code == 404


@pytest.mark.asyncio
class TestDosageTypes:
    async def test_create_dosage_type(self, client: AsyncClient):
        cat = await client.post("/api/categories", json={"name": "보톡스"})
        treat = await client.post(f"/api/categories/{cat.json()['id']}/treatments", json={"name": "사각턱"})
        treat_id = treat.json()["id"]
        res = await client.post(f"/api/treatments/{treat_id}/dosage-types", json={"unit": "shot"})
        assert res.status_code == 201
        assert res.json()["unit"] == "shot"

    async def test_invalid_dosage_unit(self, client: AsyncClient):
        cat = await client.post("/api/categories", json={"name": "보톡스"})
        treat = await client.post(f"/api/categories/{cat.json()['id']}/treatments", json={"name": "사각턱"})
        treat_id = treat.json()["id"]
        res = await client.post(f"/api/treatments/{treat_id}/dosage-types", json={"unit": "invalid"})
        assert res.status_code == 422

    async def test_duplicate_dosage_type(self, client: AsyncClient):
        cat = await client.post("/api/categories", json={"name": "보톡스"})
        treat = await client.post(f"/api/categories/{cat.json()['id']}/treatments", json={"name": "사각턱"})
        treat_id = treat.json()["id"]
        await client.post(f"/api/treatments/{treat_id}/dosage-types", json={"unit": "shot"})
        res = await client.post(f"/api/treatments/{treat_id}/dosage-types", json={"unit": "shot"})
        assert res.status_code == 409


@pytest.mark.asyncio
class TestSeedData:
    """Verify seed_data.json is correctly loaded and queryable."""

    async def test_seed_categories_count(self, seeded_client: AsyncClient):
        res = await seeded_client.get("/api/categories")
        assert res.status_code == 200
        assert len(res.json()) == 6

    async def test_seed_category_names(self, seeded_client: AsyncClient):
        res = await seeded_client.get("/api/categories")
        names = {c["name"] for c in res.json()}
        assert "리프팅" in names
        assert "보톡스" in names
        assert "필러" in names

    async def test_seed_treatments_per_category(self, seeded_client: AsyncClient):
        res = await seeded_client.get("/api/categories")
        lifting_cat = next(c for c in res.json() if c["name"] == "리프팅")
        treats = await seeded_client.get(f"/api/categories/{lifting_cat['id']}/treatments")
        assert treats.status_code == 200
        assert len(treats.json()) == 5  # 울쎄라, 슈링크, 써마지, 올리지오, 인모드

    async def test_seed_dosage_types(self, seeded_client: AsyncClient):
        # 울쎄라 should have "shot" dosage type
        res = await seeded_client.get("/api/treatments/019e442a-cf92-7000-a4ec-61cc66c26af2/dosage-types")
        assert res.status_code == 200
        units = [d["unit"] for d in res.json()]
        assert "shot" in units

    async def test_seed_cycle_rules(self, seeded_client: AsyncClient):
        # 울쎄라 cycle rule = 336 days
        res = await seeded_client.get("/api/cycle-rules/019e442a-cf92-7000-a4ec-61cc66c26af2")
        assert res.status_code == 200
        assert res.json()["cycle_days"] == 336
