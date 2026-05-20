from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession

from app.db.session import get_session
from app.repositories import TreatmentDataRepository
from app.schemas import (
    CategoryCreate,
    CategoryResponse,
    CategoryUpdate,
    DosageTypeCreate,
    DosageTypeResponse,
    TreatmentCreate,
    TreatmentResponse,
    TreatmentUpdate,
)
from app.services import TreatmentDataService

router = APIRouter(prefix="/api", tags=["treatment-data"])


def _get_service(session: AsyncSession = Depends(get_session)) -> TreatmentDataService:
    return TreatmentDataService(TreatmentDataRepository(session))


# --- Categories ---

@router.get("/categories", response_model=list[CategoryResponse])
async def list_categories(svc: TreatmentDataService = Depends(_get_service)):
    return await svc.list_categories()


@router.post("/categories", response_model=CategoryResponse, status_code=201)
async def create_category(data: CategoryCreate, svc: TreatmentDataService = Depends(_get_service)):
    return await svc.create_category(data)


@router.put("/categories/{category_id}", response_model=CategoryResponse)
async def update_category(category_id: str, data: CategoryUpdate, svc: TreatmentDataService = Depends(_get_service)):
    return await svc.update_category(category_id, data)


@router.delete("/categories/{category_id}", status_code=204)
async def delete_category(category_id: str, svc: TreatmentDataService = Depends(_get_service)):
    await svc.delete_category(category_id)


# --- Treatments ---

@router.get("/categories/{category_id}/treatments", response_model=list[TreatmentResponse])
async def list_treatments(category_id: str, svc: TreatmentDataService = Depends(_get_service)):
    return await svc.list_treatments(category_id)


@router.post("/categories/{category_id}/treatments", response_model=TreatmentResponse, status_code=201)
async def create_treatment(category_id: str, data: TreatmentCreate, svc: TreatmentDataService = Depends(_get_service)):
    return await svc.create_treatment(category_id, data)


@router.put("/treatments/{treatment_id}", response_model=TreatmentResponse)
async def update_treatment(treatment_id: str, data: TreatmentUpdate, svc: TreatmentDataService = Depends(_get_service)):
    return await svc.update_treatment(treatment_id, data)


@router.delete("/treatments/{treatment_id}", status_code=204)
async def delete_treatment(treatment_id: str, svc: TreatmentDataService = Depends(_get_service)):
    await svc.delete_treatment(treatment_id)


# --- Dosage Types ---

@router.get("/treatments/{treatment_id}/dosage-types", response_model=list[DosageTypeResponse])
async def list_dosage_types(treatment_id: str, svc: TreatmentDataService = Depends(_get_service)):
    return await svc.list_dosage_types(treatment_id)


@router.post("/treatments/{treatment_id}/dosage-types", response_model=DosageTypeResponse, status_code=201)
async def create_dosage_type(treatment_id: str, data: DosageTypeCreate, svc: TreatmentDataService = Depends(_get_service)):
    return await svc.create_dosage_type(treatment_id, data)


@router.delete("/dosage-types/{dosage_type_id}", status_code=204)
async def delete_dosage_type(dosage_type_id: str, svc: TreatmentDataService = Depends(_get_service)):
    await svc.delete_dosage_type(dosage_type_id)
