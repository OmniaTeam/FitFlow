from fastapi import APIRouter
from app.api.v1.endpoints import (
    yandex_gpt,
    fitness_planner,
)

api_router = APIRouter()
api_router.include_router(yandex_gpt.router, prefix="/api/gpt/yandex-gpt", tags=["yandex-gpt"])
api_router.include_router(fitness_planner.router, prefix="/api/gpt/fitness", tags=["fitness"])
