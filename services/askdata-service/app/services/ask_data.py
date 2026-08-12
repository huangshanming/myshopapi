from __future__ import annotations

from app.agent.ask_data import AskDataAgent
from app.api.schemas import AskDataRequest, AskDataResponse


class AskDataService:
    """问数业务编排：串联 agent / repository / clients。"""

    def __init__(self, agent: AskDataAgent | None = None) -> None:
        self.agent = agent or AskDataAgent()

    async def query(self, req: AskDataRequest) -> dict:
        result: AskDataResponse = await self.agent.run(req)
        return result.model_dump()
