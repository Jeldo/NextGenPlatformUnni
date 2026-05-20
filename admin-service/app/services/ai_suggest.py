import asyncio
import json

import boto3
from fastapi import HTTPException

from app.config import settings

_client = boto3.client("bedrock-runtime", region_name=settings.aws_region)

SYSTEM_PROMPT = """당신은 미용 시술 전문가입니다. 시술명을 입력받으면 다음을 예측해주세요:

1. category: 시술이 속하는 카테고리 (다음 중 하나: 리프팅, 스킨부스터, 레이저 (색소/홍조), 여드름·모공, 보톡스, 필러)
2. cycle_days: 추천 재시술 주기 (일 단위 정수)
3. dosage_unit: 용량 단위 (다음 중 하나: shot, minute, volume, vial, unit, joule)
4. reasoning: 판단 근거 (한국어, 1~2문장)

반드시 아래 JSON 형식으로만 응답하세요:
{"category": "...", "cycle_days": 0, "dosage_unit": "...", "reasoning": "..."}"""


def _invoke_bedrock(treatment_name: str) -> dict:
    response = _client.invoke_model(
        modelId=settings.bedrock_model_id,
        body=json.dumps({
            "anthropic_version": "bedrock-2023-05-31",
            "max_tokens": 256,
            "messages": [{"role": "user", "content": f"시술명: {treatment_name}"}],
            "system": SYSTEM_PROMPT,
        }),
    )
    body = json.loads(response["body"].read())
    return json.loads(body["content"][0]["text"])


async def suggest_treatment(treatment_name: str) -> dict:
    try:
        return await asyncio.to_thread(_invoke_bedrock, treatment_name)
    except json.JSONDecodeError:
        raise HTTPException(status_code=502, detail="AI 응답 파싱 실패")
    except Exception as e:
        raise HTTPException(status_code=502, detail=f"Bedrock 호출 실패: {str(e)}")
