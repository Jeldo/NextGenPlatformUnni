from contextlib import asynccontextmanager
from collections.abc import AsyncIterator

from fastapi import FastAPI, HTTPException, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from app.routers import ai_suggest, cycle_rules, treatment_data


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncIterator[None]:
    yield


app = FastAPI(title="Admin Service", version="0.1.0", lifespan=lifespan)


@app.exception_handler(HTTPException)
async def http_exception_handler(request: Request, exc: HTTPException) -> JSONResponse:
    code_map = {400: "BAD_REQUEST", 404: "NOT_FOUND", 409: "CONFLICT", 422: "VALIDATION_ERROR", 502: "BAD_GATEWAY"}
    return JSONResponse(
        status_code=exc.status_code,
        content={"error": {"code": code_map.get(exc.status_code, "INTERNAL_ERROR"), "message": exc.detail}},
    )

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(cycle_rules.router)
app.include_router(treatment_data.router)
app.include_router(ai_suggest.router)


@app.get("/health")
async def health():
    return {"status": "ok", "service": "admin-service", "version": "0.1.0"}


if __name__ == "__main__":
    import uvicorn
    from app.config import settings

    uvicorn.run(app, host="0.0.0.0", port=settings.port)
