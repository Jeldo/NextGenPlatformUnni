from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession

from app.db.session import get_session
from app.repositories import CycleRuleRepository, TreatmentDataRepository
from app.schemas import CycleRuleCreate, CycleRuleResponse, CycleRuleUpdate
from app.services import CycleRuleService

router = APIRouter(prefix="/api/cycle-rules", tags=["cycle-rules"])


def _get_service(session: AsyncSession = Depends(get_session)) -> CycleRuleService:
    return CycleRuleService(CycleRuleRepository(session), TreatmentDataRepository(session))


@router.get("", response_model=list[CycleRuleResponse])
async def list_cycle_rules(svc: CycleRuleService = Depends(_get_service)):
    return await svc.list_rules()


@router.get("/{treatment_id}", response_model=CycleRuleResponse)
async def get_cycle_rule(treatment_id: str, svc: CycleRuleService = Depends(_get_service)):
    return await svc.get_rule(treatment_id)


@router.post("", response_model=CycleRuleResponse, status_code=201)
async def create_cycle_rule(data: CycleRuleCreate, svc: CycleRuleService = Depends(_get_service)):
    return await svc.create_rule(data)


@router.put("/{treatment_id}", response_model=CycleRuleResponse)
async def update_cycle_rule(treatment_id: str, data: CycleRuleUpdate, svc: CycleRuleService = Depends(_get_service)):
    return await svc.update_rule(treatment_id, data)


@router.delete("/{treatment_id}", status_code=204)
async def delete_cycle_rule(treatment_id: str, svc: CycleRuleService = Depends(_get_service)):
    await svc.delete_rule(treatment_id)
