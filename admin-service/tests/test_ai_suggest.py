import json
from unittest.mock import patch

import pytest
from httpx import AsyncClient


@pytest.mark.asyncio
class TestAISuggest:
    async def test_suggest_treatment_success(self, client: AsyncClient):
        mock_response = {
            "category": "보톡스",
            "cycle_days": 168,
            "dosage_unit": "volume",
            "reasoning": "사각턱 보톡스는 보톡스 카테고리에 속합니다.",
        }
        with patch("app.services.ai_suggest._invoke_bedrock", return_value=mock_response):
            res = await client.post("/api/ai/suggest-treatment", json={"treatment_name": "사각턱 보톡스"})
        assert res.status_code == 200
        assert res.json()["category"] == "보톡스"
        assert res.json()["cycle_days"] == 168
        assert res.json()["dosage_unit"] == "volume"

    async def test_suggest_treatment_empty_name(self, client: AsyncClient):
        res = await client.post("/api/ai/suggest-treatment", json={"treatment_name": ""})
        assert res.status_code == 422

    async def test_suggest_treatment_bedrock_failure(self, client: AsyncClient):
        with patch("app.services.ai_suggest._invoke_bedrock", side_effect=RuntimeError("connection error")):
            res = await client.post("/api/ai/suggest-treatment", json={"treatment_name": "울쎄라"})
        assert res.status_code == 502
        assert res.json()["error"]["code"] == "BAD_GATEWAY"
