from contextlib import asynccontextmanager
from collections.abc import AsyncIterator

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.routers import ai_suggest, cycle_rules, treatment_data


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncIterator[None]:
    yield


app = FastAPI(title="Admin Service", version="0.1.0", lifespan=lifespan)

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
