import asyncio

from langgraph.runtime import Runtime

from app.agent.context import DataAgentContext
from app.agent.state import DataAgentState


async def correct_sql(state: DataAgentState, runtime: Runtime[DataAgentContext]):
    """根据校验错误信息校正 SQL。"""
    writer = runtime.stream_writer
    writer("校正 SQL")
    await asyncio.sleep(0.5)
