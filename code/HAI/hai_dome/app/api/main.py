from fastapi import APIRouter

from app.api.routes import ai_dome

api_router = APIRouter()
api_router.include_router(ai_dome.router)

