from fastapi import APIRouter

from app.agents.ask_data import AskDataAgent, AskDataRequest

router = APIRouter(tags=["askdata"])


@router.post("/query")
async def ask_query(req: AskDataRequest) -> dict:
    """电商问数：自然语言查询指标（当前为规则占位，后续接 NL2SQL / Agent）。"""
    agent = AskDataAgent()
    return await agent.query(req)
