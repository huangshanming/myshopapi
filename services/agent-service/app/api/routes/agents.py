from fastapi import APIRouter

from app.agents.shopping_guide import ShoppingGuideAgent, ShoppingGuidePlanRequest

router = APIRouter(tags=["agents"])


@router.post("/shopping-guide/plan")
async def shopping_guide_plan(req: ShoppingGuidePlanRequest) -> dict:
    """智能导购规划（当前为规则占位，后续接 LLM / RAG）。"""
    agent = ShoppingGuideAgent()
    return await agent.plan(req)
