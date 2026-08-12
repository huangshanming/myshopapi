from fastapi import APIRouter

from app.api.schemas import AskDataRequest
from app.services.ask_data import AskDataService

router = APIRouter(tags=["askdata"])


@router.post("/query")
async def ask_query(req: AskDataRequest) -> dict:
    """电商问数：自然语言查询指标（当前为规则占位）。"""
    return await AskDataService().query(req)
