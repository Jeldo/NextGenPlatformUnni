from pydantic import BaseModel, Field
from fastapi import APIRouter

from app.services.ai_suggest import suggest_treatment

router = APIRouter(prefix="/api/ai", tags=["ai"])


class SuggestRequest(BaseModel):
    treatment_name: str = Field(min_length=1, max_length=200)


class SuggestResponse(BaseModel):
    category: str
    cycle_days: int
    dosage_unit: str
    reasoning: str


@router.post("/suggest-treatment", response_model=SuggestResponse)
async def ai_suggest_treatment(req: SuggestRequest):
    result = await suggest_treatment(req.treatment_name)
    return SuggestResponse(**result)
