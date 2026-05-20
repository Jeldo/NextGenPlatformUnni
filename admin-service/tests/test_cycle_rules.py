import pytest
from httpx import AsyncClient


@pytest.mark.asyncio
class TestCycleRules:
    async def _create_treatment(self, client: AsyncClient) -> str:
        cat = await client.post("/api/categories", json={"name": "보톡스"})
        treat = await client.post(f"/api/categories/{cat.json()['id']}/treatments", json={"name": "사각턱"})
        return treat.json()["id"]

    async def test_create_cycle_rule(self, client: AsyncClient):
        treat_id = await self._create_treatment(client)
        res = await client.post("/api/cycle-rules", json={
            "treatment_id": treat_id, "cycle_days": 90, "description": "테스트 주기"
        })
        assert res.status_code == 201
        assert res.json()["cycle_days"] == 90
        assert res.json()["treatment_id"] == treat_id

    async def test_create_cycle_rule_invalid_days(self, client: AsyncClient):
        treat_id = await self._create_treatment(client)
        res = await client.post("/api/cycle-rules", json={
            "treatment_id": treat_id, "cycle_days": 0
        })
        assert res.status_code == 422

    async def test_create_cycle_rule_nonexistent_treatment(self, client: AsyncClient):
        res = await client.post("/api/cycle-rules", json={
            "treatment_id": "fake-id", "cycle_days": 90
        })
        assert res.status_code == 404

    async def test_create_duplicate_cycle_rule(self, client: AsyncClient):
        treat_id = await self._create_treatment(client)
        await client.post("/api/cycle-rules", json={"treatment_id": treat_id, "cycle_days": 90})
        res = await client.post("/api/cycle-rules", json={"treatment_id": treat_id, "cycle_days": 180})
        assert res.status_code == 409

    async def test_list_cycle_rules(self, client: AsyncClient):
        treat_id = await self._create_treatment(client)
        await client.post("/api/cycle-rules", json={"treatment_id": treat_id, "cycle_days": 90})
        res = await client.get("/api/cycle-rules")
        assert res.status_code == 200
        assert len(res.json()) == 1

    async def test_get_cycle_rule(self, client: AsyncClient):
        treat_id = await self._create_treatment(client)
        await client.post("/api/cycle-rules", json={"treatment_id": treat_id, "cycle_days": 90})
        res = await client.get(f"/api/cycle-rules/{treat_id}")
        assert res.status_code == 200
        assert res.json()["cycle_days"] == 90

    async def test_get_nonexistent_cycle_rule(self, client: AsyncClient):
        res = await client.get("/api/cycle-rules/fake-id")
        assert res.status_code == 404

    async def test_update_cycle_rule(self, client: AsyncClient):
        treat_id = await self._create_treatment(client)
        await client.post("/api/cycle-rules", json={"treatment_id": treat_id, "cycle_days": 90})
        res = await client.put(f"/api/cycle-rules/{treat_id}", json={"cycle_days": 180, "description": "변경"})
        assert res.status_code == 200
        assert res.json()["cycle_days"] == 180

    async def test_delete_cycle_rule(self, client: AsyncClient):
        treat_id = await self._create_treatment(client)
        await client.post("/api/cycle-rules", json={"treatment_id": treat_id, "cycle_days": 90})
        res = await client.delete(f"/api/cycle-rules/{treat_id}")
        assert res.status_code == 204
        # Verify gone
        get_res = await client.get(f"/api/cycle-rules/{treat_id}")
        assert get_res.status_code == 404
